package main

import (
	"github.com/johannes-riecken/go-grammatisch/pkg/go-grammatisch/server"
	"log"
	"net/http"
)

func mainAux() error {
	server.AddRoutes()
	return http.ListenAndServe(":8080", nil)
}

func main() {
	if err := mainAux(); err != nil {
		log.Fatal(err)
	}
}
