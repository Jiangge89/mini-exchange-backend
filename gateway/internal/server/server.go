package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Jiangge89/mini-exchange-backend/gateway/internal/cache"
	"github.com/Jiangge89/mini-exchange-backend/gateway/internal/config"
	"github.com/Jiangge89/mini-exchange-backend/gateway/internal/metrics"
	"github.com/Jiangge89/mini-exchange-backend/gateway/internal/ratelimit"
	"github.com/Jiangge89/mini-exchange-backend/gateway/internal/upstream"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type JSONRPCRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
	ID      json.RawMessage   `json:"id"`
}

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
	ID      json.RawMessage `json:"id"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func New(cfg *config.Config, selector upstream.Selector, cache cache.Cache, limiter ratelimit.RateLimiter) *http.Server {
	mux := http.NewServeMux()
	h := &handler{
		selector: selector,
		cache:    cache,
		limiter:  limiter,
	}
	mux.HandleFunc("POST /", h.handleRPC)
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /health", handleHealth)

	return &http.Server{
		Addr:         cfg.Addr(),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

type handler struct {
	selector upstream.Selector
	cache    cache.Cache
	limiter  ratelimit.RateLimiter
}

func (h *handler) handleRPC(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	apiKey := r.Header.Get("X-API-Key")
	if apiKey != "" {
		allowed, _, err := h.limiter.Allow(r.Context(), apiKey)
		if err != nil {
			slog.Error("rate limiter error", "error", err)
		}
		if !allowed {
			metrics.RequestsTotal.WithLabelValues("rate_limited").Inc()
			writeJSONError(w, http.StatusTooManyRequests, -32000, "rate limit exceeded")
			return
		}
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 10*1024*1024))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, -32700, "failed to read request body")
		return
	}

	isBatch := len(body) > 0 && body[0] == '['
	if isBatch {
		h.handleBatch(w, r, body, start)
	} else {
		h.handleSingle(w, r, body, start)
	}
}

func (h *handler) handleSingle(w http.ResponseWriter, r *http.Request, body []byte, start time.Time) {
	var req JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, -32700, "parse error")
		return
	}

	resp := h.processRequest(r, &req)

	metrics.RequestsTotal.WithLabelValues(req.Method).Inc()
	metrics.RequestDuration.WithLabelValues(req.Method).Observe(time.Since(start).Seconds())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *handler) handleBatch(w http.ResponseWriter, r *http.Request, body []byte, start time.Time) {
	var reqs []JSONRPCRequest
	if err := json.Unmarshal(body, &reqs); err != nil {
		writeJSONError(w, http.StatusBadRequest, -32700, "parse error")
		return
	}

	responses := make([]JSONRPCResponse, len(reqs))
	for i, req := range reqs {
		responses[i] = h.processRequest(r, &req)
		metrics.RequestsTotal.WithLabelValues(req.Method).Inc()
	}
	metrics.RequestDuration.WithLabelValues("batch").Observe(time.Since(start).Seconds())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *handler) processRequest(r *http.Request, req *JSONRPCRequest) JSONRPCResponse {
	ctx := r.Context()

	shouldCache, ttl := h.cache.ShouldCache(req.Method, req.Params)
	if shouldCache {
		cacheKey := cache.BuildKey(req.Method, req.Params)
		if cached, ok, err := h.cache.Get(ctx, cacheKey); err == nil && ok {
			metrics.CacheHits.Inc()
			var resp JSONRPCResponse
			if err := json.Unmarshal(cached, &resp); err == nil {
				resp.ID = req.ID
				return resp
			}
		}
		metrics.CacheMisses.Inc()

		resp := h.forwardToUpstream(ctx, req)
		if resp.Error == nil {
			if respBytes, err := json.Marshal(resp); err == nil {
				h.cache.Set(ctx, cacheKey, respBytes, ttl)
			}
		}
		return resp
	}

	return h.forwardToUpstream(ctx, req)
}

func (h *handler) forwardToUpstream(ctx interface{ Done() <-chan struct{} }, req *JSONRPCRequest) JSONRPCResponse {
	reqBody, _ := json.Marshal(req)
	respBody, err := h.selector.Forward(ctx, reqBody)
	if err != nil {
		metrics.UpstreamErrors.Inc()
		return JSONRPCResponse{
			JSONRPC: "2.0",
			Error:   &JSONRPCError{Code: -32603, Message: "upstream error: " + err.Error()},
			ID:      req.ID,
		}
	}

	var resp JSONRPCResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			Error:   &JSONRPCError{Code: -32603, Message: "invalid upstream response"},
			ID:      req.ID,
		}
	}
	resp.ID = req.ID
	return resp
}

func writeJSONError(w http.ResponseWriter, httpStatus, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Error:   &JSONRPCError{Code: code, Message: msg},
		ID:      json.RawMessage("null"),
	}
	json.NewEncoder(w).Encode(resp)
}
