package mcp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sundance/backend/services/mcp-gateway/internal/tools"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ServerTransport string

const (
	ServerTransportHTTP  ServerTransport = "http"
	ServerTransportStdio ServerTransport = "stdio"
)

type ServerOpts struct {
	name      string
	version   string
	transport ServerTransport
	logger    *slog.Logger
	port      *int
	registry  tools.Registry
}

type ServerOptsFn func(*ServerOpts)

func ServerWithLogger(logger *slog.Logger) ServerOptsFn {
	return func(so *ServerOpts) {
		so.logger = logger
	}
}

func ServerWithName(name string) ServerOptsFn {
	return func(so *ServerOpts) {
		so.name = name
	}
}

func ServerWithPort(port *int) ServerOptsFn {
	return func(so *ServerOpts) {
		so.port = port
	}
}

func ServerWithToolRegistry(registry tools.Registry) ServerOptsFn {
	return func(so *ServerOpts) {
		so.registry = registry
	}
}

func ServerWithTransport(transport ServerTransport) ServerOptsFn {
	return func(so *ServerOpts) {
		so.transport = transport
	}
}

func ServerWithVersion(version string) ServerOptsFn {
	return func(so *ServerOpts) {
		so.version = version
	}
}

type Server struct {
	server    *mcp.Server
	transport ServerTransport
	logger    *slog.Logger
	port      *int
}

func NewServer(opts ...ServerOptsFn) *Server {
	var o ServerOpts

	for _, opt := range opts {
		opt(&o)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    o.name,
		Version: o.version,
	}, nil)

	return &Server{
		server:    server,
		transport: o.transport,
		port:      o.port,
		logger:    o.logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverErrChan := make(chan error, 1)

	go func() {
		var run func(context.Context) error

		switch s.transport {
		case ServerTransportStdio:
			run = s.stdio
		case ServerTransportHTTP:
			run = s.http
		}

		if err := run(signalCtx); err != nil {
			serverErrChan <- err
		}

		close(serverErrChan)
	}()

	select {
	case err := <-serverErrChan:
		if err == nil || errors.Is(err, mcp.ErrConnectionClosed) {
			return nil
		}

		return err
	case <-signalCtx.Done():
		s.logger.Info("application received shutdown signal, waiting for server cleanup...")

		if err := <-serverErrChan; err != nil && !errors.Is(err, mcp.ErrConnectionClosed) {
			return err
		}
	}

	return nil
}

func (s *Server) stdio(ctx context.Context) error {
	return s.server.Run(ctx, &mcp.StdioTransport{})
}

func (s *Server) http(ctx context.Context) error {
	if s.port == nil {
		return errors.New("http transport requires a port to listen on but got nil")
	}

	handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		return s.server
	}, nil)

	mux := chi.NewRouter()

	mux.Use(httplog.RequestLogger(s.logger, &httplog.Options{
		Schema: httplog.SchemaOTEL,
	}))
	mux.Mount("/mcp", handler)

	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", *s.port),
		Handler:     mux,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
	}

	s.logger.Info(fmt.Sprintf("application listening on :%d", *s.port))

	return server.ListenAndServe()
}
