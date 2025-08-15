package main

import (
	"context"
	"fmt"
	"github.com/Sanchir01/go-shortener-auth/internal/app"
	"os"
	"os/signal"
	"syscall"
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

	<-ctx.Done()

}
