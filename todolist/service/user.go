package service

import (
	"crypto/sha256"
	"net/http"
	database "todolist.go/db"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

func NewUserForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "new_user_form.html", gin.H{"Title": "Register user"})
}

func RegisterUser(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	passwordConfirm := ctx.PostForm("password_confirm")
	switch {
	case username == "":
		//ユーザー名が入力されていない場合
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Usernane is not provided", "Username": username})
		return
	case password == "":
		//パスワードが入力されていない場合
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Password is not provided", "Password": password})
		return
	case password != passwordConfirm:
		//パスワードが一致しない場合
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Password does not match", "Username": username, "Password": password})
		return
	case utf8.RuneCountInString(password) < 8:
		//パスワードが8文字未満の場合
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Password must be at least 8 characters", "Username": username, "Password": password})
		return
	}

	//DB接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//重複チェック
	var duplicate int
	err = db.Get(&duplicate, "SELECT COUNT(*) FROM users WHERE name=?", username)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	if duplicate > 0 {
		//重複がある場合、入力画面にエラー表示状態で戻る
		ctx.HTML(http.StatusBadRequest, "new_user_form.html",
			gin.H{"Title": "Register user", "Error": "Username is already taken", "Username": username, "Password": password})
		return
	}

	//DBにユーザー情報を保存する
	result, err := db.Exec("INSERT INTO users(name,password) VALUES(?,?)", username, hash(password))
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//保存状態を確認する
	id, _ := result.LastInsertId()
	var user database.User
	err = db.Get(&user, "SELECT id, name, password FROM users WHERE id=?", id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func hash(pw string) []byte {
	const salt = "todolist.go#"
	h := sha256.New()
	h.Write([]byte(salt))
	h.Write([]byte(pw))
	return h.Sum(nil)
}
