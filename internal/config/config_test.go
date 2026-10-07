package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadAddr(t *testing.T) {
	tests := []struct {
		name     string
		addr     string
		unset    bool
		wantAddr string
		wantErr  bool
	}{
		{name: "default", unset: true, wantAddr: ":8080"},
		{name: "custom", addr: ":9090", wantAddr: ":9090"},
		{name: "invalid", addr: "abc", wantErr: true},
		{name: "range", addr: "70000", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ADDR", tt.addr)
			if tt.unset {
				os.Unsetenv("ADDR")
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
			if cfg.Addr != tt.wantAddr {
				t.Errorf("Addr = %s, want %s", cfg.Addr, tt.wantAddr)
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
