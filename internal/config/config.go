package config

import (
	"fmt"
	"net"
	"os"
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
	_, _, err := net.SplitHostPort(val)
	if err != nil {
		return config, fmt.Errorf("invalid ADDR format: %w", err)
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
