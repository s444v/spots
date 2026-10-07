package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/s444v/spots/internal/config"
	"github.com/s444v/spots/internal/logger"
)

func TestHandler(t *testing.T) {
	webFS := fstest.MapFS{
		"index.html": {Data: []byte("<h1>spots</h1>")},
	}
	cfg := config.Config{Port: ":9999"}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	h := New(cfg, webFS, log).Handler

	tests := []struct {
		name     string
		method   string
		path     string
		wantCode int
		wantBody string
	}{
		{"healthz", http.MethodGet, "/healthz", http.StatusOK, `"status":"ok"`},
		{"index", http.MethodGet, "/", http.StatusOK, "<h1>spots</h1>"},
		{"missing file", http.MethodGet, "/nope.js", http.StatusNotFound, ""},
		{"post not allowed", http.MethodPost, "/healthz", http.StatusMethodNotAllowed, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantCode)
			}
			if tt.wantBody != "" && !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("body = %q, want it to contain %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	cfg := config.Config{Port: ":9999"} // без t.Setenv
	log := logger.Init()
	srv := New(cfg, fstest.MapFS{}, log)

	if srv.Addr != ":9999" {
		t.Errorf("Addr = %q, want :9999", srv.Addr)
	}
	if srv.Handler == nil {
		t.Error("Handler is nil")
	}
}
