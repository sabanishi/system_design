package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	kw := ctx.Query("kw")

	//終了状態を取得
	isDoneStr := ctx.Query("is_done")
	isDoneExist := true

	//isDoneをbool型に変換
	isDone, err := strconv.ParseBool(isDoneStr)
	if err != nil {
		fmt.Println("is_done is not bool")
		isDoneExist = false
	}

	// Get tasks in DB
	var tasks []database.Task
	switch {
	case kw != "" && isDoneExist:
		err = db.Select(&tasks, "SELECT * FROM tasks WHERE title LIKE ? AND is_done = ?", "%"+kw+"%", isDone)
	case kw != "":
		err = db.Select(&tasks, "SELECT * FROM tasks WHERE title LIKE ?", "%"+kw+"%")
	case isDoneExist:
		err = db.Select(&tasks, "SELECT * FROM tasks WHERE is_done = ?", isDone)
	default:
		err = db.Select(&tasks, "SELECT * FROM tasks")
	}
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Render task
	ctx.HTML(http.StatusOK, "task.html", task)
}

func NewTaskForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "form_new_task.html", gin.H{"Title": "Task registration"})
}

func RegisterTask(ctx *gin.Context) {
	title, exist := ctx.GetPostForm("title")
	if !exist {
		Error(http.StatusBadRequest, "title is not exist")(ctx)
		return
	}
	description, exist := ctx.GetPostForm("description")
	if !exist {
		Error(http.StatusBadRequest, "description is not exist")(ctx)
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	result, err := db.Exec("INSERT INTO tasks (title,description) VALUES (?,?)", title, description)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	path := "/list"
	if id, err := result.LastInsertId(); err == nil {
		path = fmt.Sprintf("/task/%d", id)
	}
	ctx.Redirect(http.StatusFound, path)
}

func EditTaskForm(ctx *gin.Context) {
	//IDを取得
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//Taskの取得
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id)
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	ctx.HTML(http.StatusOK, "form_edit_task.html", gin.H{"Title": fmt.Sprintf("Edit task %d", task.ID), "Task": task})
}

func UpdateTask(ctx *gin.Context) {
	//POSTされたデータを取得
	title, exist := ctx.GetPostForm("title")
	if !exist {
		Error(http.StatusBadRequest, "title is not exist")(ctx)
		return
	}
	description, exist := ctx.GetPostForm("description")
	if !exist {
		Error(http.StatusBadRequest, "description is not exist")(ctx)
		return
	}
	isDoneStr, exist := ctx.GetPostForm("is_done")
	if !exist {
		Error(http.StatusBadRequest, "is_done is not exist")(ctx)
		return
	}
	//isDoneをbool型に変換
	isDone, err := strconv.ParseBool(isDoneStr)
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	//IDを取得
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	//DB接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//Taskの更新
	_, err = db.Exec("UPDATE tasks SET title=?, description=?, is_done=? WHERE id=?", title, description, isDone, id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//リダイレクト処理
	path := fmt.Sprintf("/task/%d", id)
	ctx.Redirect(http.StatusFound, path)
}

func DeleteTask(ctx *gin.Context) {
	//IDを取得
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	//DB接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//Taskの削除
	_, err = db.Exec("DELETE FROM tasks WHERE id=?", id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//リダイレクト処理
	path := "/list"
	ctx.Redirect(http.StatusFound, path)
}
