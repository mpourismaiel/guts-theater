package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mpourismaiel/guts-theater/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

// creates new http server along with chi router, can be used to create multiple
// servers and makes it possible to create microservices if required
// also registers essential middlewares such as cleanpath, recoverer and prometheus calls
func New(conf *config.Config, logger *zap.Logger) (*Server, error) {
	r := chi.NewRouter()

	r.Use(chizap.New(logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)
	r.Use(patternHandler)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Handle("/metrics", promhttp.Handler())

	return &Server{
		router: r,
		addr:   conf.Address,
		port:   conf.Port,
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

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				s.logger.Fatal("graceful shutdown timed out... forcing exit.")
			}
			s.logger.Info("gracefully shutdown")
			cancel()
		}()

		err := s.server.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Fatal(err.Error())
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
