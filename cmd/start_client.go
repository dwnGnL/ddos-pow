package cmd

import (
	"context"
	"fmt"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/internal/handler/client"
	"github.com/dwnGnL/ddos-pow/lib/goerrors"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dwnGnL/ddos-pow/internal/application"
	"github.com/dwnGnL/ddos-pow/internal/service"
)

const (
	gracefulStop = 5 * time.Second
)

func StartClient(cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()
	s, err := buildClient(cfg)
	if err != nil {
		return fmt.Errorf("build service err:%w", err)
	}

	gracefulStopWithCtx := client.SetupHandlers(s, cfg)

	var group errgroup.Group

	group.Go(func() error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		goerrors.Log().Debug("wait for Ctrl-C")
		<-sigCh
		goerrors.Log().Debug("Ctrl-C signal")
		cancelCtx()
		shutdownCtx, shutdownCtxFunc := context.WithDeadline(ctx, time.Now().Add(gracefulStop))
		defer shutdownCtxFunc()

		_ = gracefulStopWithCtx(shutdownCtx)
		return nil
	})

	if err := group.Wait(); err != nil {
		goerrors.Log().WithError(err).Error("Stopping service with error")
	}
	return nil
}

func buildClient(conf *config.Config) (application.Core, error) {
	return service.New(conf), nil
}
