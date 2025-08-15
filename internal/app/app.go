package app

import (
	"context"
	"github.com/Sanchir01/go-shortener-auth/internal/config"
	"github.com/Sanchir01/go-shortener-auth/pkg/logger"
	"log/slog"
)

type App struct {
	Cfg           *config.Config
	Log           *slog.Logger
	CloseLoggerFN func()
}

func New(ctx context.Context) (*App, error) {
	cfg := config.InitConfig()
	lg, cleanup := slogpretty.SetupAsyncLogger("development")
	lg.Info("initializing app")
	app := &App{
		Cfg:           cfg,
		Log:           lg,
		CloseLoggerFN: cleanup,
	}
	return app, nil
}
