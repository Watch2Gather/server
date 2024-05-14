package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
	"google.golang.org/grpc"

	"github.com/Watch2Gather/server/cmd/movie/config"
	"github.com/Watch2Gather/server/cmd/movie/socket"
	"github.com/Watch2Gather/server/internal/movie/app"
	sharedkernel "github.com/Watch2Gather/server/internal/pkg/shared_kernel"
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

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load config", err)
	}

	slog.Info("Init app", "name", cfg.Name, "version", cfg.Version)

	log := slog.New(
		slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{Level: logger.ConvertLogLevel(cfg.Level)}))
	slog.SetDefault(log)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(sharedkernel.TokenInterceptor),
	)

	httpMux := http.NewServeMux()

	go func() {
		defer grpcServer.GracefulStop()
		<-ctx.Done()
	}()

	_, cleanup, err := app.InitApp(cfg, postgres.DBConnString(cfg.DsnURL), grpcServer)
	if err != nil {
		slog.Error("failed init app", err)
		cancel()
	}

	go _grpc(ctx, cfg, grpcServer, cancel)
	go _http(ctx, cfg, httpMux, cancel)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case v := <-quit:
		grpcServer.Stop()
		cleanup()
		slog.Info("signal.Notify", v)
	case done := <-ctx.Done():
		grpcServer.Stop()
		cleanup()
		slog.Info("ctx.Done", "app done", done)
	}
}

func _grpc(ctx context.Context, cfg *config.Config, server *grpc.Server, cancel func()) {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	network := "tcp"

	l, err := net.Listen(network, address)
	if err != nil {
		slog.Error("Failed to listen to address", err, "network", network, "address", address)
		cancel()
		<-ctx.Done()
	}

	slog.Info("Start grpc server...", "address", address)

	defer func() {
		if err1 := l.Close(); err != nil {
			slog.Error("failed to close grpc", err1, "network", network, "address", address)
			<-ctx.Done()
		}
	}()

	err = server.Serve(l)
	if err != nil {
		slog.Error("Failed to start gRPC server", err, "network", network, "address", address)
		cancel()
		<-ctx.Done()
	}
}

func _http(ctx context.Context, cfg *config.Config, mux *http.ServeMux, cancel func()) {
	go socket.NewWsHandler(mux)
	address := fmt.Sprintf("%s:%d", cfg.WSHost, cfg.WSPort)
	network := "tcp"

	l, err := net.Listen(network, address)
	if err != nil {
		slog.Error("Failed to listen to address", err, "network", network, "address", address)
		cancel()
		<-ctx.Done()
	}

	slog.Info("Start http server...", "address", address)

	defer func() {
		if err1 := l.Close(); err != nil {
			slog.Error("failed to close http", err1, "network", network, "address", address)
			<-ctx.Done()
		}
	}()

	err = http.Serve(l, mux)
	if err != nil {
		slog.Error("Failed to start http server", err, "network", network, "address", address)
		cancel()
		<-ctx.Done()
	}
}
