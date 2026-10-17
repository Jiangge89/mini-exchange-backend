package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	Port                 string
	GatewayURL           string
	TokenContractAddress common.Address
	PostgresHost         string
	PostgresPort         string
	PostgresUser         string
	PostgresPassword     string
	PostgresDB           string
}

func Load() (*Config, error) {
	port := envOr("INDEXER_PORT", "8081")
	gatewayURL := envOr("GETH_URL", "http://localhost:8080")

	tokenAddr := os.Getenv("TOKEN_CONTRACT_ADDRESS")
	if tokenAddr == "" {
		return nil, fmt.Errorf("TOKEN_CONTRACT_ADDRESS is required")
	}
	if !common.IsHexAddress(tokenAddr) {
		return nil, fmt.Errorf("TOKEN_CONTRACT_ADDRESS is not a valid address: %s", tokenAddr)
	}

	return &Config{
		Port:                 port,
		GatewayURL:           gatewayURL,
		TokenContractAddress: common.HexToAddress(tokenAddr),
		PostgresHost:         envOr("POSTGRES_HOST", "localhost"),
		PostgresPort:         envOr("POSTGRES_PORT", "5432"),
		PostgresUser:         envOr("POSTGRES_USER", "miniexchange"),
		PostgresPassword:     envOr("POSTGRES_PASSWORD", "miniexchange"),
		PostgresDB:           envOr("POSTGRES_DB", "miniexchange"),
	}, nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresDB)
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
