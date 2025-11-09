// pkg/app_server.go
package pkg

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type AppServer struct {
	name     string
	version  string
	adapters []Server
}

type Option func(a *AppServer)

func WithName(name string) Option {
	return func(a *AppServer) {
		if a.name == "" {
			a.name = "default-server"
		}
		a.name = name
	}
}

func WithVersion(version string) Option {
	return func(a *AppServer) {
		if version == "" {
			version = "0.0.1"
		}
		a.version = version
	}
}

func NewServer(opts ...Option) *AppServer {
	a := &AppServer{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithAdapter(adapters ...Server) Option {
	return func(a *AppServer) {
		a.adapters = append(a.adapters, adapters...)
	}
}

func (a *AppServer) Run() error {
	startTime := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info().
		Str("name", a.name).
		Str("version", a.version).
		Msg("Starting server")

	var (
		wg      sync.WaitGroup
		errChan = make(chan error, len(a.adapters))
	)

	for _, adapter := range a.adapters {
		wg.Add(1)
		go func(ad Server) {
			defer wg.Done()
			if err := ad.Start(ctx); err != nil {
				select {
				case errChan <- err:
				default:
					log.Warn().Msg("error channel full, dropping error")
				}
			}
		}(adapter)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var runErr error

	select {
	case err := <-errChan:
		runErr = fmt.Errorf("adapter start failed: %w", err)
		log.Error().Err(err).Msg("Adapter start failed — initiating shutdown...")
		cancel()

	case sig := <-quit:
		log.Warn().Msgf("Received signal: %s — shutting down...", sig)
		cancel()

	case <-ctx.Done():
		log.Warn().Msg("Context canceled — shutting down...")
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	for _, adapter := range a.adapters {
		if err := adapter.Stop(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Adapter stop failed")
		}
	}

	wg.Wait()
	close(errChan)

	duration := time.Since(startTime)
	time.Sleep(200 * time.Millisecond)

	log.Warn().
		Str("server", a.name).
		Float64("shutdown_time_sec", duration.Seconds()).
		Msg("Shutdown complete")

	log.Warn().
		Int("goroutines", runtime.NumGoroutine()).
		Msg("Final runtime stats")

	return runErr
}
