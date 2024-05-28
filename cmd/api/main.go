package main

import (
	"net/http"

	"github.com/sgatu/ezmail/cmd/api/server"
)

func main() {
	server := server.NewServer()
	http.ListenAndServe(":3000", server)
}
