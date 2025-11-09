package _http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	//http2 "github.com/kaffein/goffy_shopper/pkg/routers/http"
	"github.com/rs/zerolog/log"
)

const (
	// HTTP server timeout constants
	DefaultReadTimeout       = 10 * time.Second
	DefaultWriteTimeout      = 10 * time.Second
	DefaultIdleTimeout       = 120 * time.Second
	DefaultReadHeaderTimeout = 5 * time.Second
	DefaultShutdownTimeout   = 5 * time.Second
)

type Server struct {
	*gin.Engine
	httpSrv *http.Server
	host    string
	port    int
}

type Option func(s *Server)

func NewServer(opts ...Option) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithEngine(engine *gin.Engine) Option {
	return func(s *Server) {
		s.Engine = engine
	}
}

func (s *Server) Start(_ context.Context) error {
	addr := net.JoinHostPort(s.host, strconv.Itoa(s.port))
	log.Info().Str("address", addr).Msg("Starting HTTP server")

	NewRouterHTTP(s.Engine).ConfigRouterHTTP()

	s.httpSrv = &http.Server{
		Addr:              addr,
		Handler:           s.Engine,
		ReadTimeout:       DefaultReadTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		IdleTimeout:       DefaultIdleTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
	}

	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpSrv == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, DefaultShutdownTimeout)
	defer cancel()

	if err := s.httpSrv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Error shutting down HTTP server")
		return err
	}

	log.Logger.Warn().Str("server", "HTTP").Msg("HTTP server stopped gracefully")
	return nil

}
