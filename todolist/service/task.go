package service

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	userID := sessions.Default(ctx).Get("user")

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
	query := "SELECT id, title, created_at, is_done, description FROM tasks INNER JOIN ownership ON task_id = id WHERE user_id = ?"
	switch {
	case kw != "" && isDoneExist:
		err = db.Select(&tasks, query+"AND title LIKE ? AND is_done = ?", userID, "%"+kw+"%", isDone)
	case kw != "":
		err = db.Select(&tasks, query+"AND title LIKE ?", userID, "%"+kw+"%")
	case isDoneExist:
		err = db.Select(&tasks, query+"AND is_done = ?", userID, isDone)
	default:
		err = db.Select(&tasks, query, userID)
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
	userID := sessions.Default(ctx).Get("user")

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

	query :=
		"SELECT id, title, created_at, is_done, description " +
			"FROM tasks " +
			"INNER JOIN ownership ON task_id = id " +
			"WHERE user_id = ?"
	err = db.Get(&task, query+" AND tasks.id = ?", userID, id) // Use DB#Get for one entry
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
	userID := sessions.Default(ctx).Get("user")

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

	//DB接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	tx := db.MustBegin()
	result, err := tx.Exec("INSERT INTO tasks (title,description) VALUES (?,?)", title, description)
	if err != nil {
		tx.Rollback()
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	taskID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	_, err = tx.Exec("INSERT INTO ownership (user_id, task_id) VALUES (?,?)", userID, taskID)
	if err != nil {
		tx.Rollback()
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	tx.Commit()

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/task/%d", taskID))
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
