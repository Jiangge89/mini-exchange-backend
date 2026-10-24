package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var cfg *Config
var loadCfgOnce sync.Once

type UpstreamConfig struct {
	Name string
	URL  string
}

type Config struct {
	Port      string
	RedisURL  string
	Upstreams []UpstreamConfig
	ChainID   uint64
}

func Load() (*Config, error) {
	port := envOr("GATEWAY_PORT", "8080")
	redisURL := envOr("REDIS_URL", "redis://localhost:6379")

	var upstreams []UpstreamConfig

	if url := os.Getenv("GETH_URL"); url != "" {
		upstreams = append(upstreams, UpstreamConfig{Name: "geth", URL: url})
	}
	if url := os.Getenv("ALCHEMY_HOODI_URL"); url != "" {
		upstreams = append(upstreams, UpstreamConfig{Name: "alchemy", URL: url})
	}
	if url := os.Getenv("PUBLIC_RPC_URL"); url != "" {
		upstreams = append(upstreams, UpstreamConfig{Name: "public", URL: url})
	}

	if len(upstreams) == 0 {
		return nil, fmt.Errorf("no upstreams configured; set at least one of GETH_URL, ALCHEMY_HOODI_URL, PUBLIC_RPC_URL")
	}

	chainIDStr := os.Getenv("CHAIN_ID")
	if chainIDStr == "" {
		return nil, fmt.Errorf("CHAIN_ID env variable is required")
	}
	chainID, err := strconv.ParseUint(chainIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse chain_id: %s err: %v", chainIDStr, err)
	}

	cfg = &Config{
		Port:      port,
		RedisURL:  redisURL,
		Upstreams: upstreams,
		ChainID:   chainID,
	}

	return cfg, nil
}

func (c *Config) Addr() string {
	port := c.Port
	if !strings.Contains(port, ":") {
		port = ":" + port
	}
	return port
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
