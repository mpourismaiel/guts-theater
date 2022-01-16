package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	router *chi.Mux
	server *http.Server
	port   string
}

func New(port string) (*Server, error) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
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
		port:   port,
	}, nil
}

func (s *Server) Serve() error {
	s.server = &http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", s.port),
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
			log.Fatalf("http.Server.Serve failed on address %s: %v", s.server.Addr, err)
		}
		wg.Done()
	}()

	log.Println("Server running on:", s.server.Addr)
	wg.Wait()
	return nil
}

func (s *Server) Router() *chi.Mux {
	return s.router
}
