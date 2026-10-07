package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadAddr(t *testing.T) {
	tests := []struct {
		name     string
		port     string
		unset    bool
		wantPort string
		wantErr  bool
	}{
		{name: "default", unset: true, wantPort: ":8080"},
		{name: "custom", port: "9090", wantPort: ":9090"},
		{name: "invalid", port: "abc", wantErr: true},
		{name: "range", port: "70000", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", tt.port)
			if tt.unset {
				os.Unsetenv("PORT")
			}
			cfg, err := Load()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.Port != tt.wantPort {
				t.Errorf("Port = %s, want %s", cfg.Port, tt.wantPort)
			}
		})
	}
}

func TestLoadShutdown(t *testing.T) {
	tests := []struct {
		name        string
		timeout     string
		unset       bool
		wantTimeout time.Duration
		wantErr     bool
	}{
		{name: "default", unset: true, wantTimeout: 10 * time.Second},
		{name: "custom", timeout: "5s", wantTimeout: 5 * time.Second},
		{name: "invalid", timeout: "324234", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SHUTDOWN_TIMEOUT", tt.timeout)
			if tt.unset {
				os.Unsetenv("SHUTDOWN_TIMEOUT")
			}
			cfg, err := Load()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.ShutdownTimeout != tt.wantTimeout {
				t.Errorf("timeout = %s, want %s", cfg.ShutdownTimeout, tt.wantTimeout)
			}
		})
	}
}
