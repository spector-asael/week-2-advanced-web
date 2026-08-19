package handler

import (
	"log/slog"
	"time"

	"github.com/spector-asael/week-1/internal/data"
)

type ServerConfig struct {
	Port        int
	Environment string
	ReportDelay time.Duration
	DB          struct {
		DSN string
	}
}

type ApplicationDependencies struct {
	Config ServerConfig
	Logger *slog.Logger
	Models data.Models
}
