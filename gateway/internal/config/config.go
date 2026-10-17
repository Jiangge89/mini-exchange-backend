package config

import (
	"fmt"
	"os"
	"strings"
)

type UpstreamConfig struct {
	Name string
	URL  string
}

type Config struct {
	Port      string
	RedisURL  string
	Upstreams []UpstreamConfig
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

	return &Config{
		Port:      port,
		RedisURL:  redisURL,
		Upstreams: upstreams,
	}, nil
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
