package logger

import (
	"log/slog"
	"os"
)

// Log is the global logger instance
var Log *slog.Logger

// Init initializes the global logger with the specified verbosity level
func Init(verbose bool) {
	level := slog.LevelWarn
	if verbose {
		level = slog.LevelDebug
	}

	Log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
}

// InitDefault initializes the logger with default settings (warn level)
func InitDefault() {
	Init(false)
}
