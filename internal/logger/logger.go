package logger

import (
	"log/slog"
	"os"
)

func Init() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler).With(
		"version", "1.2.0",
		"env", "production",
		"app", "spots",
	)
	return logger
}
