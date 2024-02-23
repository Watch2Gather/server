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
	"github.com/Watch2Gather/server/pkg/logger"
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
			os.Stdout, &slog.HandlerOptions{Level: logger.ConvertLogLevel(cfg.Log.Level)}))
	slog.SetDefault(log)

	server := grpc.NewServer()

	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()

	// TODO add app init with postgres
	// cleanup := prepareApp(ctx, cancel, cfg, server)
	cleanup := func() {}

	// gRPC server
	address := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
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
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
}
