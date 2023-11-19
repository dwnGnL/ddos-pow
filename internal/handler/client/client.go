package client

import (
	"context"
	"fmt"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/internal/application"
	"log"
	"net/http"
)

type GracefulStopFuncWithCtx func(ctx context.Context) error

func SetupHandlers(core application.Core, cfg *config.Config) GracefulStopFuncWithCtx {
	mux := http.NewServeMux()

	handlerRoutes := application.WithApp(core, mux)

	handler := newHandler(cfg)

	mux.HandleFunc("/request-challenge", handler.RequestChallenge)
	mux.HandleFunc("/request-resource", handler.RequestResource)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Client.Port),
		Handler: handlerRoutes,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start srv: %v", err)
		}
	}()

	return srv.Shutdown
}
