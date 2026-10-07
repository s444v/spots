package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	var config Config

	val, ok := os.LookupEnv("ADDR")
	if !ok {
		val = ":8080"
	}
	_, p, err := net.SplitHostPort(val)
	if err != nil {
		return config, fmt.Errorf("invalid ADDR format: %w", err)
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return config, fmt.Errorf("failed to conv port: %w", err)
	}
	if port < 1 || port > 65535 {
		return config, fmt.Errorf("port out of range: %d", port)
	}
	config.Addr = val
	val, ok = os.LookupEnv("SHUTDOWN_TIMEOUT")
	if !ok {
		val = "10s"
	}
	sTimeout, err := time.ParseDuration(val)
	if err != nil {
		return config, err
	}
	config.ShutdownTimeout = sTimeout
	return config, nil
}
