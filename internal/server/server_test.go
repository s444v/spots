package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/s444v/spots/internal/config"
)

func TestHandler(t *testing.T) {
	webFS := fstest.MapFS{
		"index.html": {Data: []byte("<h1>spots</h1>")},
	}
	cfg := config.Config{Port: ":9999"}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	h := New(cfg, webFS, log).Handler

	tests := []struct {
		name      string
		method    string
		path      string
		wantCode  int
		wantType  string // префикс Content-Type, пусто = не проверять
		wantBody  string // подстрока, пусто = не проверять
		checkBody func(t *testing.T, body []byte)
	}{
		{
			name:     "healthz",
			method:   http.MethodGet,
			path:     "/healthz",
			wantCode: http.StatusOK,
			wantType: "application/json",
			checkBody: func(t *testing.T, body []byte) {
				var got map[string]string
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatalf("body is not valid JSON: %v (body=%q)", err, body)
				}
				if got["status"] != "ok" {
					t.Errorf(`status field = %q, want "ok"`, got["status"])
				}
			},
		},
		{
			name:     "index",
			method:   http.MethodGet,
			path:     "/",
			wantCode: http.StatusOK,
			wantType: "text/html",
			wantBody: "<h1>spots</h1>",
		},
		{"missing file", http.MethodGet, "/nope.js", http.StatusNotFound, "", "", nil},
		{"post not allowed", http.MethodPost, "/healthz", http.StatusMethodNotAllowed, "", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantCode)
			}
			if tt.wantType != "" {
				if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, tt.wantType) {
					t.Errorf("Content-Type = %q, want prefix %q", ct, tt.wantType)
				}
			}
			if tt.wantBody != "" && !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("body = %q, want it to contain %q", rec.Body.String(), tt.wantBody)
			}
			if tt.checkBody != nil {
				tt.checkBody(t, rec.Body.Bytes())
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	cfg := config.Config{Port: ":9999"} // без t.Setenv
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := New(cfg, fstest.MapFS{}, log)

	if srv.Addr != ":9999" {
		t.Errorf("Addr = %q, want :9999", srv.Addr)
	}
	if srv.Handler == nil {
		t.Error("Handler is nil")
	}
}
