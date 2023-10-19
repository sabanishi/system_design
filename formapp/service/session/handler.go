package session

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Start(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "start.html", gin.H{"Target": "/session/start"})
}

func NameForm(ctx *gin.Context) {
	session, err := NewSession()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Fail to create new session")
		return
	}
	ctx.SetCookie("userid", session.ID(), 600, "/session", "localhost:8000", false, false)
	ctx.HTML(http.StatusOK, "name-form.html", gin.H{"Target": "/session/name"})
}

func BirthdayForm(ctx *gin.Context) {
	id, err := ctx.Cookie("userid")
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid access")
	}
	session := Session{id}
	name, exist := ctx.GetPostForm("name")
	if !exist {
		ctx.String(http.StatusBadRequest, "parameter 'name' is not exist")
	}
	state, _ := session.GetState()
	state.Name = name
	session.SetState(state)
	ctx.HTML(http.StatusOK, "session-birthday-form.html", nil)
}

func MessageForm(ctx *gin.Context) {
	id, err := ctx.Cookie("userid")
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid access")
	}
	session := Session{id}

	birthday, exist := ctx.GetPostForm("birthday")
	if !exist {
		ctx.String(http.StatusBadRequest, "parameter 'birthday' is not exist")
	}

	state, _ := session.GetState()
	state.Birthday = birthday
	session.SetState(state)
	ctx.HTML(http.StatusOK, "session-message-form.html", nil)
}

func Confirm(ctx *gin.Context) {
	id, err := ctx.Cookie("userid")
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid access")
	}
	session := Session{id}

	message, exist := ctx.GetPostForm("message")
	if !exist {
		ctx.String(http.StatusBadRequest, "parameter 'message' is not exist")
	}

	state, _ := session.GetState()
	state.Message = message
	session.SetState(state)

	ctx.HTML(http.StatusOK, "confirm.html",
		gin.H{"Target": "/session/confirm/",
			"Name":     state.Name,
			"Birthday": state.Birthday,
			"Message":  state.Message})
}
