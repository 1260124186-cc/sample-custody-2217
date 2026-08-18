package config

import (
	"net"
	"os"
	"strings"
)

type Config struct {
	Address string
}

func Load() Config {
	address := strings.TrimSpace(os.Getenv("SAMPLE_CUSTODY_ADDR"))
	if address == "" {
		address = "127.0.0.1:8090"
	}
	return Config{Address: address}
}

func (c Config) ListenAddress() string {
	if _, _, err := net.SplitHostPort(c.Address); err != nil {
		return "127.0.0.1:8090"
	}
	return c.Address
}
