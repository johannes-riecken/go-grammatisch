package main

import (
	"github.com/gin-gonic/gin"
	"github.com/johannes-riecken/go-grammatisch/pkg/go-grammatisch/server"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("pkg/templates/*")
	server.AddRoutes(r)
	_ = r.Run()
}
