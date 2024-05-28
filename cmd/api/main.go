package main

import (
	"net/http"

	"github.com/sgatu/ezmail/cmd/api/server"
	"github.com/sgatu/ezmail/internal/http/handlers"
)

func main() {
	server := server.NewServer()
	handlers.SetupRoutes(server)
	http.ListenAndServe(":3000", server)
}
