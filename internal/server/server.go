package server

import (
	"net/http"
	"time"
)

func CreateServer() *http.Server {
	return &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 10 * time.Second,
	}
}
