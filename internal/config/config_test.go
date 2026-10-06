package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		port     string
		unset    bool
		wantPort string
		wantErr  bool
	}{
		{name: "default", unset: true, wantPort: "8080"},
		{name: "custom", port: "9090", wantPort: "9090"},
		{name: "invalid", port: "abc", wantErr: true},
		{name: "range", port: "70000", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ADDR", tt.port)
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
			if cfg.ADDR != tt.wantPort {
				t.Errorf("Port = %s, want %s", cfg.ADDR, tt.wantPort)
			}
		})
	}
}
