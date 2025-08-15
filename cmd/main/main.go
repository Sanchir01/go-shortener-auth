package main

import (
	"context"
	"fmt"
	"github.com/Sanchir01/go-shortener-auth/internal/feature/auth"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sanchir01/go-shortener-auth/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Signal(syscall.SIGTERM), syscall.SIGINT)
	defer cancel()
	application, err := app.New(ctx)
	defer application.CloseLoggerFN()
	if err != nil {
		panic(err)
	}
	fmt.Println(application.Cfg)
	application.GRPCApp.MustStart(func(server *grpc.Server) {
		auth.NewServer(server)
	})
	<-ctx.Done()
}
