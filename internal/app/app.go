package app

import (
	"context"
	"github.com/Sanchir01/go-shortener-auth/internal/config"
	grpcapp "github.com/Sanchir01/go-shortener-auth/internal/grpc"
	"github.com/Sanchir01/go-shortener-auth/pkg/logger"
	"log/slog"
)

type App struct {
	Cfg           *config.Config
	Log           *slog.Logger
	CloseLoggerFN func()
	GRPCApp       *grpcapp.GRPCApp
}

func New(ctx context.Context) (*App, error) {
	cfg := config.InitConfig()
	lg, cleanup := slogpretty.SetupAsyncLogger("development")
	lg.Info("initializing app")

	gRPC := grpcapp.New(lg)
	app := &App{
		Cfg:           cfg,
		Log:           lg,
		CloseLoggerFN: cleanup,
		GRPCApp:       gRPC,
	}
	return app, nil
}
