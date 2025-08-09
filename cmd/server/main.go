package main

import (
	"goserv/internal/server"
)

func run() {
	srv := server.NewServer()
	srv.Run()
}

func main() {
	run()
}
