package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func Init(env string) {
	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

func Info(msg string, args ...interface{}) {
	Logger.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	Logger.Error(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	Logger.Debug(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	Logger.Warn(msg, args...)
}