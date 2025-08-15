package grpcapp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rookie-ninja/rk-boot/v2"
	"github.com/rookie-ninja/rk-grpc/v2/boot"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	log          *slog.Logger
	boot         *rkboot.Boot
	registerFunc func(*grpc.Server)
}

func New(log *slog.Logger) *GRPCApp {
	return &GRPCApp{
		log: log,
	}
}

func (a *GRPCApp) MustStart(registerFunc func(*grpc.Server)) {
	if err := a.Start(registerFunc); err != nil {
		panic(err)
	}
}

func (a *GRPCApp) RegisterServices(registerFunc func(*grpc.Server)) {
	a.registerFunc = registerFunc
}

func (a *GRPCApp) Start(registerFunc func(*grpc.Server)) error {
	const op = "grpcapp.App.Start"
	log := a.log.With(slog.String("op", op))

	boot := rkboot.NewBoot(rkboot.WithBootConfigPath("boot.yaml", nil))
	a.boot = boot

	grpcEntry := rkgrpc.GetGrpcEntry("auth-service")
	if grpcEntry == nil {
		return fmt.Errorf("failed to get gRPC entry 'auth-service' from boot.yaml")
	}

	grpcEntry.AddRegFuncGrpc(func(server *grpc.Server) {
		registerFunc(server)
		log.Info("Services registered")
	})

	log.Info("Starting gRPC server from boot.yaml config")
	boot.Bootstrap(context.Background())
	log.Info("gRPC server started", "port", grpcEntry.Port)

	return nil
}

func (a *GRPCApp) Stop() {
	if a.boot != nil {
		a.boot.Shutdown(context.Background())
		a.log.Info("gRPC server stopped")
	}
}
