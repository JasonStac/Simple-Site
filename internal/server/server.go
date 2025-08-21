package server

import (
	"context"
	"goserv/ent/gen"
	"goserv/internal/database"
	pRepo "goserv/internal/domain/posts/repository"
	sRepo "goserv/internal/domain/sessions/repository"
	uRepo "goserv/internal/domain/users/repository"
	"goserv/pkg/config"
	"goserv/pkg/templates"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	cfg config.Config

	tmplCache *template.Template

	ent *gen.Client

	user    uRepo.User
	session sRepo.Session
	post    pRepo.Post

	router *chi.Mux

	httpServer *http.Server
}

func NewServer() *Server {
	server := &Server{
		cfg:    config.Load(),
		router: chi.NewRouter(),
	}

	server.initDB()
	server.initTemplates()
	server.initDomain()

	return server
}

func (s *Server) initDB() {
	conn, err := database.NewDB(s.cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v\n", err)
	}

	s.ent = conn
}

func (s *Server) initTemplates() {
	tmplCache, err := templates.LoadTemplates("tmpl/*.html")
	if err != nil {
		log.Fatalf("Failed to load templates: %v\n", err)
	}
	s.tmplCache = tmplCache
}

func (s *Server) Config() config.Config {
	return s.cfg
}

func (s *Server) Run() {
	s.httpServer = &http.Server{
		Addr:              s.cfg.Host + ":" + s.cfg.Port,
		Handler:           s.router,
		ReadHeaderTimeout: s.cfg.ReadHeaderTimeout,
	}

	go func() {
		start(s)
	}()

	_ = gracefulShutdown(context.Background(), s)
}

func start(s *Server) {
	log.Printf("Starting server on %s...", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) closeResources() {
	s.ent.Close()
}

func gracefulShutdown(ctx context.Context, s *Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down...")

	ctx, shutdown := context.WithTimeout(ctx, s.Config().GracefulTimeout*time.Second)
	defer shutdown()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
	s.closeResources()

	return nil
}
