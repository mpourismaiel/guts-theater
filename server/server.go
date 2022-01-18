package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				s.logger.Fatal("graceful shutdown timed out... forcing exit.")
			}
			s.logger.Info("gracefully shutdown")
		}()

		err := s.server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	fields := []zapcore.Field{
		zap.String("address", s.server.Addr),
	}
	s.logger.Info("Server running on", fields...)
	if err = s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
		fields := []zapcore.Field{
			zap.String("address", s.server.Addr),
			zap.String("error", err.Error()),
		}
		s.logger.Error("http.Server.Serve failed on address", fields...)
	}
	<-serverCtx.Done()
	return nil
}

func (s *Server) Router() *chi.Mux {
	return s.router
}
