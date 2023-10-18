package stateless

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Start(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "start.html", gin.H{"Target": "/stateless/start"})
}

func NameForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "name-form.html", gin.H{"Target": "/stateless/name"})
}

func BirthdayForm(ctx *gin.Context) {
	name, exist := ctx.GetPostForm("name")
	if !isExist(exist, "name", ctx) {
		return
	}
	ctx.HTML(http.StatusOK, "stateless-birthday-form.html", gin.H{"Name": name})
}

func MessageForm(ctx *gin.Context) {
	name, exist := ctx.GetPostForm("name")
	if !isExist(exist, "name", ctx) {
		return
	}

	birthday, exist := ctx.GetPostForm("birthday")
	if !isExist(exist, "birthday", ctx) {
		return
	}

	ctx.HTML(http.StatusOK, "stateless-message-form.html", gin.H{"Name": name, "Birthday": birthday})
}

func Confirm(ctx *gin.Context) {
	name, exist := ctx.GetPostForm("name")
	if !isExist(exist, "name", ctx) {
		return
	}

	birthday, exist := ctx.GetPostForm("birthday")
	if !isExist(exist, "birthday", ctx) {
		return
	}

	message, exist := ctx.GetPostForm("message")
	if !isExist(exist, "message", ctx) {
		return
	}

	ctx.HTML(http.StatusOK, "stateless-confirm.html", gin.H{"Name": name, "Birthday": birthday, "Message": message})
}

func isExist(exist bool, message string, ctx *gin.Context) bool {
	if !exist {
		ctx.String(http.StatusBadRequest, "parameter '%s' is not provided", message)
		return false
	}
	return true
}
