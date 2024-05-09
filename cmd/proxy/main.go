package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Watch2Gather/server/cmd/proxy/config"
	"github.com/Watch2Gather/server/pkg/logger"
	"github.com/Watch2Gather/server/proto/gen"
)

func newGateway(
	ctx context.Context,
	cfg *config.Config,
	opts []gwruntime.ServeMuxOption,
) (http.Handler, error) {
	userEndpoint := fmt.Sprintf("%s:%d", cfg.UserHost, cfg.UserPort)
	roomEndpoint := fmt.Sprintf("%s:%d", cfg.RoomHost, cfg.RoomPort)
	movieEndpoint := fmt.Sprintf("%s:%d", cfg.MovieHost, cfg.MoviePort)

	mux := gwruntime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	var err error

	err = gen.RegisterUserServiceHandlerFromEndpoint(ctx, mux, userEndpoint, dialOpts)
	if err != nil {
		slog.Error("User service", "error", err)
		return nil, err
	}

	err = gen.RegisterRoomServiceHandlerFromEndpoint(ctx, mux, roomEndpoint, dialOpts)
	if err != nil {
		slog.Error("Room service", "error", err)
		return nil, err
	}

	err = gen.RegisterMovieServiceHandlerFromEndpoint(ctx, mux, movieEndpoint, dialOpts)
	if err != nil {
		slog.Error("Movie service", "error", err)
		return nil, err
	}

	return mux, nil
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)

				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	headers := []string{"*"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))

	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))

	slog.Info("preflight request", "http_path", r.URL.Path)
}

func withLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Run request", "http_method", r.Method, "http_host", r.Host, "http_url", r.URL)

		h.ServeHTTP(w, r)
	})
}

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg, err := config.NewCofig()
	if err != nil {
		glog.Fatalf("Config error: %s", err)
	}

	log := slog.New(
		slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{Level: logger.ConvertLogLevel(cfg.Level)}))
	slog.SetDefault(log)

	mux := http.NewServeMux()

	gw, err := newGateway(ctx, cfg, nil)
	if err != nil {
		slog.Error("failed to create a new gateway", err)
	}

	mux.Handle("/", gw)

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: allowCORS(withLogger(wsproxy.WebsocketProxy(mux))),
		// Handler: allowCORS(withLogger(mux)),
	}

	go func() {
		<-ctx.Done()
		slog.Info("shutting down the http server")

		if err := s.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown the server", err)
		}
	}()

	slog.Info("start listening...", "address", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))

	if err := s.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to listen and serve", err)
	}
}
