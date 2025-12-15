package logger

import (
	"context"
	"os"
	"sync"

	"log/slog"
)

var (
	defaultLogger *slog.Logger
	once          sync.Once
)

// Config de Logger
type Config struct {
	ServiceName string
	Environment string // "dev" | "prod"
	Level       string // "debug", "info", "warn", "error"
}

// InitLogger inițializează logger-ul global și returnează instanța
func InitLogger(cfg Config) *slog.Logger {
	once.Do(func() {
		var handler slog.Handler

		// Nivelul de logare
		var level slog.Level
		switch cfg.Level {
		case "debug":
			level = slog.LevelDebug
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}

		opts := &slog.HandlerOptions{
			Level: level,
			// Adăugăm sursa (fișier/linie) doar în develop
			AddSource: cfg.Environment == "development",
		}

		if cfg.Environment == "development" {
			// Text Handler pentru citire umană ușoară local
			handler = slog.NewTextHandler(os.Stdout, opts)
		} else {
			// JSON Handler pentru producție (Parsing automat)
			handler = slog.NewJSONHandler(os.Stdout, opts)
		}

		// Adăugăm atribute globale (Service Name)
		defaultLogger = slog.New(handler).With(
			slog.String("service", cfg.ServiceName),
			slog.String("env", cfg.Environment),
		)

		// Setăm ca logger default global pentru librăriile standard
		slog.SetDefault(defaultLogger)
	})

	return defaultLogger
}

// Get returnează logger-ul configurat. Dacă nu e init, returnează unul basic.
func Get() *slog.Logger {
	if defaultLogger == nil {
		return slog.Default()
	}
	return defaultLogger
}

// LogWithContext extrage TraceID din context dacă există (pentru tracing distribuit)
func For(ctx context.Context) *slog.Logger {
	// Aici am putea extrage TraceID din context dacă folosim OpenTelemetry
	// Momentan returnăm logger-ul standard
	return Get()
}
