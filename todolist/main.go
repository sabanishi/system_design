package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"todolist.go/db"
	"todolist.go/service"
)

const port = 8000

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize Gin engine
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")

	//セッションの準備
	store := cookie.NewStore([]byte("m-secret"))
	engine.Use(sessions.Sessions("user-session", store))

	// routing
	engine.Static("/assets", "./assets")
	engine.GET("/", service.Home)

	engine.GET("/list", service.LoginCheck, service.TaskList)
	engine.GET("/list/:page", service.LoginCheck, service.TaskList)

	//ログアウト
	engine.POST("/logout", service.LoginCheck, service.Logout)
	//退会
	engine.POST("/delete", service.LoginCheck, service.DeleteUser)

	taskGroup := engine.Group("/task")
	taskGroup.Use(service.LoginCheck)
	{
		taskGroup.GET("/:id", service.ShowTask) // ":id" is a parameter

		//タスクの追加/編集/削除
		taskGroup.GET("/new", service.NewTaskForm)
		taskGroup.POST("/new", service.RegisterTask)
		taskGroup.GET("/edit/:id", service.EditTaskForm)
		taskGroup.POST("/edit/:id", service.UpdateTask)
		taskGroup.GET("/delete/:id", service.DeleteTask)
		//タスクへのアクセス権付与
		taskGroup.GET("/add_owner/:id", service.AddOwnerForm)
		taskGroup.POST("add_owner/:id", service.AddOwner)
	}

	//ユーザー登録
	engine.GET("/user/new", service.NewUserForm)
	engine.POST("/user/new", service.RegisterUser)

	engine.GET("/login", service.LoginForm)
	engine.POST("/login", service.Login)

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
