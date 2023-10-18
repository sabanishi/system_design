package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RootHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "hello.html", nil)
}
