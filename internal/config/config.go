package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ADDR             string
	SHUTDOWN_TIMEOUT time.Duration
}

func Load() (Config, error) {
	var config Config
	val, ok := os.LookupEnv("ADDR")
	if !ok {
		val = "8080"
	}
	port, err := strconv.Atoi(val)
	if err != nil {
		return config, fmt.Errorf("invalid PORT %q: %w", val, err)
	}
	if port < 1 || port > 65535 {
		return config, fmt.Errorf("PORT %d out of range 1-65535", port)
	}
	config.ADDR = val
	val, ok = os.LookupEnv("SHUTDOWN_TIMEOUT")
	if !ok {
		val = "10s"
	}
	sTimeout, err := time.ParseDuration(val)
	if err != nil {
		return config, err
	}
	config.SHUTDOWN_TIMEOUT = sTimeout
	return config, nil
}
