package server

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/s444v/spots/internal/config"
	"github.com/s444v/spots/internal/logger"
)

func NewServer(cfg config.Config, webFS fs.FS) *http.Server {
	return &http.Server{
		Handler:      NewHandler(webFS),
		Addr:         cfg.ADDR,
		ReadTimeout:  cfg.SHUTDOWN_TIMEOUT,
		WriteTimeout: cfg.SHUTDOWN_TIMEOUT,
		IdleTimeout:  cfg.SHUTDOWN_TIMEOUT,
	}
}

func NewHandler(webFS fs.FS) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)

	mux.Handle("GET /", http.FileServerFS(webFS))
	return mux
}

type healthzResp struct {
	Status string `json:"status"`
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	h := healthzResp{Status: "ok"}
	writeJSON(w, http.StatusOK, h)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(w, r)
		logger.Log.Info("logHandlers", "Method", r.Method, "Path", r.URL.Path, "Status", rec.status, "Duration", time.Since(start))
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Error("panic", "error", err, "debug", debug.Stack())
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
