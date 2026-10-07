package logger

import (
	"log/slog"
	"os"
)

func Init() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler).With(
		"app", "spots",
	)
	return logger
}
