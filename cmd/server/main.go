package main

import (
	"goserv/internal/db"
	"goserv/internal/server"
	"goserv/pkg/config"
	"goserv/pkg/templates"
	"log"
	"net/http"
)

func run() error {
	cfg := config.Load()

	conn, err := db.NewDB(cfg.DB)
	if err != nil {
		return err
	}
	defer conn.Close()

	tmplCache, err := templates.LoadTemplates("tmpl/*.html")
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	srv := server.NewServer(conn, tmplCache)

	addr := ":" + cfg.Port
	log.Printf("Starting server on %s...", addr)

	return http.ListenAndServe(addr, srv)
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
