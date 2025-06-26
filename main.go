package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goserv/handlers"
	"goserv/server"
)

func main() {
	port := 8080
	serv := server.NewServer(port)

	serv.Use(server.LoggingMiddleware)
	serv.Use(server.RecoveryMiddleware)

	serv.Router.GET("/", handlers.HomeHandler)
	serv.Router.GET("/time", handlers.TimeHandler)
	serv.Router.GET("/users", handlers.UsersHandler)
	serv.Router.POST("/users", handlers.CreateUserHandler)
	serv.Router.NotFound(handlers.NotFoundHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := serv.Run(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	fmt.Printf("Server is running on http://localhost:%d", port)

	<-stop

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := serv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}
