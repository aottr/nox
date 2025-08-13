package logging

import (
	"log/slog"
	"os"
	"strings"
	"sync"
)

type Logger = *slog.Logger

var (
	once     sync.Once
	logger   *slog.Logger
	logLevel = &slog.LevelVar{}
)

func Get() Logger {
	Init()
	return logger
}

func Init() {
	once.Do(func() {
		logLevel.Set(slog.LevelInfo)

		h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
		logger = slog.New(h)
	})
}

func SetLevel(s string) {
	switch strings.ToLower(s) {
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "info":
		logLevel.Set(slog.LevelInfo)
	case "warn", "warning":
		logLevel.Set(slog.LevelWarn)
	case "error":
		logLevel.Set(slog.LevelError)
	default:
		logLevel.Set(slog.LevelInfo)
	}
}
