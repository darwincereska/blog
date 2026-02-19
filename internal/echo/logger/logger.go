package logger

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

// NewCharmSlog returns a standard *slog.Logger powered by Charmbracelet
func NewCharmSlog() *slog.Logger {
	// 1. Initialize Charmbracelet
	options := log.Options{
		ReportTimestamp: true,
		ReportCaller:    false,
		Level:           log.DebugLevel,
		Prefix: "ECHO",
	}
	handler := log.NewWithOptions(os.Stderr, options)

	// 2. Return as *slog.Logger
	// Charmbracelet's Logger implements the slog.Handler interface
	return slog.New(handler)
}
