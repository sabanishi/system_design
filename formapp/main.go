package main

import (
	"fmt"
	"net/http"

	"formapp.go/service"
	"formapp.go/service/stateless"
	"github.com/gin-gonic/gin"
)

// config
const port = 8000

func main() {
	// initialize Gin engine
	engine := gin.Default()

	engine.LoadHTMLGlob("templates/*.html")

	engine.GET("/", service.RootHandler)

	engine.GET("/stateless/start", stateless.Start)
	engine.POST("/stateless/start", stateless.NameForm)
	engine.POST("/stateless/name", stateless.BirthdayForm)
	engine.POST("/stateless/birthday", stateless.MessageForm)
	engine.POST("/stateless/message", stateless.Confirm)
	engine.POST("/stateless/confirm", stateless.Start)

	/* 3章の範囲
	// routing
	engine.GET("/hello", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "hello.html", nil)
	})
	engine.GET("/bye", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "bye.html", nil)
	})
	engine.GET("/hello.jp", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "hello_jp.html", nil)
	})
	engine.GET("/name-form", nameFormHandler)
	engine.POST("/register-name", registerNameHandler)
	*/

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}

func notImplemented(ctx *gin.Context) {
	msg := fmt.Sprintf("%s to %s is not implemented yet", ctx.Request.Method, ctx.Request.URL)
	ctx.String(http.StatusNotImplemented, msg)
}

func rootHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "hello.html", nil)
}

func nameFormHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "name_form.html", nil)
}

func registerNameHandler(ctx *gin.Context) {
	name, _ := ctx.GetPostForm("name")
	ctx.HTML(http.StatusOK, "result.html", gin.H{"Name": name})
}
