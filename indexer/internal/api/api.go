package api

import (
	"encoding/json"
	"log/slog"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Jiangge89/mini-exchange-backend/indexer/internal/store"
)

func NewServer(addr string, s store.Store) *http.Server {
	mux := http.NewServeMux()
	h := &handler{store: s}

	mux.HandleFunc("GET /api/transfers/{address}", h.getTransfers)
	mux.HandleFunc("GET /api/balance/{address}", h.getBalance)
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	return &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

type handler struct {
	store store.Store
}

type transferResponse struct {
	BlockNumber       uint64 `json:"block_number"`
	TxHash            string `json:"tx_hash"`
	LogIndex          uint   `json:"log_index"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	ConfirmationLevel string `json:"confirmation_level"`
}

func (h *handler) getTransfers(w http.ResponseWriter, r *http.Request) {
	addr := r.PathValue("address")
	if !common.IsHexAddress(addr) {
		writeError(w, http.StatusBadRequest, "invalid address")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	transfers, err := h.store.GetTransfers(r.Context(), common.HexToAddress(addr), limit, offset)
	if err != nil {
		slog.Error("failed to get transfers", "error", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := make([]transferResponse, len(transfers))
	for i, t := range transfers {
		resp[i] = transferResponse{
			BlockNumber:       t.BlockNumber,
			TxHash:            t.TxHash.Hex(),
			LogIndex:          t.LogIndex,
			From:              t.From.Hex(),
			To:                t.To.Hex(),
			Value:             t.Value.String(),
			ConfirmationLevel: string(t.ConfirmationLevel),
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

func (h *handler) getBalance(w http.ResponseWriter, r *http.Request) {
	addr := r.PathValue("address")
	if !common.IsHexAddress(addr) {
		writeError(w, http.StatusBadRequest, "invalid address")
		return
	}

	bal, err := h.store.GetBalance(r.Context(), common.HexToAddress(addr))
	if err != nil {
		slog.Error("failed to get balance", "error", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if bal == nil {
		bal = big.NewInt(0)
	}

	writeJSON(w, http.StatusOK, balanceResponse{
		Address: addr,
		Balance: bal.String(),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
