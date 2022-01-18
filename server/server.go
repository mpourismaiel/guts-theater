package server

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/chizap"
)

type Server struct {
	router *chi.Mux
	server *http.Server
	port   string
	addr   string
	logger *zap.Logger
}

func New(address string, port string, logger *zap.Logger) (*Server, error) {
	r := chi.NewRouter()

	r.Use(chizap.New(logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	return &Server{
		router: r,
		addr:   address,
		port:   port,
		logger: logger,
	}, nil
}

func (s *Server) Serve() error {
	s.server = &http.Server{
		Addr:    net.JoinHostPort(s.addr, s.port),
		Handler: s.router,
	}

	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("net.Listen failed on address %s: %v", s.server.Addr, err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err = s.server.Serve(listener); err != nil {
			fields := []zapcore.Field{
				zap.String("address", s.server.Addr),
				zap.String("error", err.Error()),
			}
			s.logger.Error("http.Server.Serve failed on address", fields...)
		}
		wg.Done()
	}()

	fields := []zapcore.Field{
		zap.String("address", s.server.Addr),
	}
	s.logger.Info("Server running on", fields...)
	wg.Wait()
	return nil
}

func (s *Server) Router() *chi.Mux {
	return s.router
}
