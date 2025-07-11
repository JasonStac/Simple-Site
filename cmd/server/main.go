package main

import (
	"context"
	"goserv/internal/db"
	"goserv/internal/handlers"
	"goserv/pkg/config"
	"goserv/pkg/templates"
	"log"
	"net/http"
)

func run(ctx context.Context) error {
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

	mux := http.NewServeMux()
	contentHandler := handlers.NewContentHandler(conn, tmplCache)
	homeHandler := handlers.NewHomeHandler(tmplCache)
	mux.HandleFunc("/", homeHandler.ServeHTTP)
	mux.HandleFunc("/view", contentHandler.HandleViewContent)
	mux.HandleFunc("/add", contentHandler.HandleAddContent)

	mux.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("./content"))))

	addr := ":" + cfg.Port
	log.Printf("Starting server on %s...", addr)

	return http.ListenAndServe(addr, mux)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
