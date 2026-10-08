package server

import (
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/s444v/spots/internal/config"
)

type Server struct {
	log   *slog.Logger
	webFS fs.FS
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

type healthzResp struct {
	Status string `json:"status"`
}

func New(cfg config.Config, webFS fs.FS, log *slog.Logger) *http.Server {
	s := &Server{log: log, webFS: webFS}

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           s.Logging(s.Recover(s.routes())),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealthz)
	mux.Handle("GET /", http.FileServerFS(s.webFS))
	return mux
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	h := healthzResp{Status: "ok"}
	s.writeJSON(w, http.StatusOK, h)
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *Server) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		duration := time.Since(start)
		s.log.Info("http request", "method", r.Method, "path", r.URL.Path, "status", rec.status, "duration_micros", duration.Microseconds())
	})
}

func (s *Server) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.log.Error("panic", "error", err, slog.String("debug", string(debug.Stack())))
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		s.log.Error("failed to encode JSON response",
			slog.Int("status", status),
			slog.Any("error", err),
		)
	}

}
