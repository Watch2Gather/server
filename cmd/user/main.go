package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
	"google.golang.org/grpc"

	"github.com/Watch2Gather/server/cmd/user/config"
	sharedkernel "github.com/Watch2Gather/server/internal/pkg/shared_kernel"
	"github.com/Watch2Gather/server/internal/user/app"
	"github.com/Watch2Gather/server/pkg/logger"
	"github.com/Watch2Gather/server/pkg/postgres"
)

func main() {
	// set GOMAXPROCS
	_, err := maxprocs.Set()
	if err != nil {
		slog.Error("failed to set maxprocs")
	}

	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.NewCofig()
	if err != nil {
		slog.Error("failed to load config", err)
	}

	slog.Info("Init app", "name", cfg.Name, "version", cfg.Version)

	log := slog.New(
		slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{Level: logger.ConvertLogLevel(cfg.Level)}))
	slog.SetDefault(log)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(sharedkernel.TokenInterceptor),
	)

	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()

	_, cleanup, err := app.InitApp(cfg, postgres.DBConnString(cfg.DsnURL), server)
	if err != nil {
		slog.Error("failed init app", err)
		cancel()
	}

	// gRPC server
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	network := "tcp"

	l, err := net.Listen(network, address)
	if err != nil {
		slog.Error("Failed to listen to address", err, "network", network, "address", address)
		cancel()
		<-ctx.Done()
	}

	slog.Info("Start server...", "address", address)

	defer func() {
		if err1 := l.Close(); err != nil {
			slog.Error("failed to close", err1, "network", network, "address", address)
			<-ctx.Done()
		}
	}()

	err = server.Serve(l)
	if err != nil {
		slog.Error("Failed to start gRPC server", err, "network", network, "address", address)
		cancel()
		<-ctx.Done()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case v := <-quit:
		// server.Stop()
		cleanup()
		slog.Info("signal.Notify", v)
	case done := <-ctx.Done():
		// server.Stop()
		cleanup()
		slog.Info("ctx.Done", "app done", done)
	}
}
