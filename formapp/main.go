package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// config
const port = 8000

func main() {
	// initialize Gin engine
	engine := gin.Default()

	engine.LoadHTMLGlob("templates/*.html")

	// routing
	engine.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world.")
	})
	engine.GET("/bye", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world.")
	})
	engine.GET("/hello.jp", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "bye.")
	})

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}

func rootHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "hello.html", nil)
}
