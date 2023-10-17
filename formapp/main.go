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

    // routing
    engine.GET("/", rootHandler)

    // start server
    engine.Run(fmt.Sprintf(":%d", port))
}

func rootHandler(ctx *gin.Context) {
    ctx.String(http.StatusOK, "Hello world.")
}
