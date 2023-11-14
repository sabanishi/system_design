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

	// parse page given as a parameter
	page, err := strconv.Atoi(ctx.Param("page"))
	if err != nil {
		page = 1
	}

	kw := ctx.Query("kw")

	//終了状態を取得
	isDoneStr := ctx.Query("is_done")
	isDoneExist := true

	//isDoneをbool型に変換
	isDone, err := strconv.ParseBool(isDoneStr)
	if err != nil {
		isDoneExist = false
	}

	// Get tasks in DB
	var tasks []database.Task
	query :=
		"SELECT id, title, created_at, is_done, description " +
			"FROM tasks " +
			"INNER JOIN ownership ON task_id = id " +
			"WHERE owner_id = ?"
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

	//表示するページを制限する
	var taskPage []database.Task
	var hasBeforePage = page > 1
	var hasAfterPage = len(tasks) > page*10
	for i := 0; i < len(tasks); i++ {
		if i >= (page-1)*10 && i < page*10 {
			taskPage = append(taskPage, tasks[i])
		}
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{
		"Title":         "Task list",
		"Page":          page,
		"HasBeforePage": hasBeforePage,
		"HasAfterPage":  hasAfterPage,
		"BeforePage":    page - 1,
		"AfterPage":     page + 1,
		"Tasks":         taskPage})
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
		"SELECT id, title, created_at, is_done, description,priority,deadline " +
			"FROM tasks " +
			"INNER JOIN ownership ON task_id = id " +
			"WHERE owner_id = ?"
	err = db.Get(&task, query+" AND tasks.id = ?", userID, id) // Use DB#Get for one entry
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	var owners []database.User
	err = db.Select(&owners,
		"SELECT id,name "+
			"FROM ownership "+
			"INNER JOIN users ON owner_id = id "+
			"WHERE task_id = ?", id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//優先度の表示
	var priority string
	switch task.Priority {
	case 0:
		priority = "低"
	case 1:
		priority = "中"
	case 2:
		priority = "高"
	}

	fmt.Println(task.Deadline)

	// Render task
	ctx.HTML(http.StatusOK, "task.html",
		gin.H{"Title": task.Title,
			"ID":          task.ID,
			"CreatedAt":   task.CreatedAt,
			"IsDone":      task.IsDone,
			"Description": task.Description,
			"Priority":    priority,
			"Deadline":    task.Deadline,
			"Owners":      owners})
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
	priority, exist := ctx.GetPostForm("priority")
	if !exist {
		Error(http.StatusBadRequest, "priority is not exist")(ctx)
		return
	}
	deadline, exist := ctx.GetPostForm("deadline")
	if !exist {
		Error(http.StatusBadRequest, "deadline is not exist")(ctx)
		return
	}

	//DB接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	tx := db.MustBegin()
	result, err := tx.Exec("INSERT INTO tasks (title,description,priority,deadline) VALUES (?,?,?,?)",
		title, description, priority, deadline)
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
	_, err = tx.Exec("INSERT INTO ownership (owner_id, task_id) VALUES (?,?)", userID, taskID)
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

	fmt.Println(task.Deadline)

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
	priority, exist := ctx.GetPostForm("priority")
	if !exist {
		Error(http.StatusBadRequest, "priority is not exist")(ctx)
		return
	}
	deadline, exist := ctx.GetPostForm("deadline")
	if !exist {
		Error(http.StatusBadRequest, "deadline is not exist")(ctx)
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
	_, err = db.Exec("UPDATE tasks SET title=?, description=?, is_done=?, priority=?,deadline=? "+
		"WHERE id=?",
		title, description, isDone, priority, deadline, id)
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

func AddOwnerForm(ctx *gin.Context) {
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

	//Taskの取得
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id)
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	ctx.HTML(http.StatusOK, "form_add_owner.html", gin.H{"Title": fmt.Sprintf("Add owner to task %d", task.ID), "Task": task})
}

func AddOwner(ctx *gin.Context) {
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
	//Taskの取得
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id)
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	//POSTされたデータを取得
	username := ctx.PostForm("username")
	if username == "" {
		ctx.HTML(http.StatusBadRequest, "form_add_owner.html",
			gin.H{"Title": fmt.Sprintf("Add owner to task %d", task.ID),
				"Task":  task,
				"Error": "Usernane is not provided"})
		return
	}

	//指定されたユーザーが存在するか確認する
	var duplicate int
	err = db.Get(&duplicate, "SELECT COUNT(*) FROM users WHERE name=?", username)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	if duplicate == 0 {
		//ユーザーが存在しない場合
		ctx.HTML(http.StatusBadRequest, "form_add_owner.html",
			gin.H{"Title": fmt.Sprintf("Add owner to task %d", task.ID),
				"Task":  task,
				"Error": "This user is not exist"})
		return
	}

	//既にアクセス権が付与されているか確認する
	var exist int
	err = db.Get(&exist,
		"SELECT COUNT(*) "+
			"FROM ownership WHERE owner_id=(SELECT id FROM users WHERE name=?) AND task_id=?", username, id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	if exist > 0 {
		//既にアクセス権が付与されている場合
		ctx.HTML(http.StatusBadRequest, "form_add_owner.html",
			gin.H{"Title": fmt.Sprintf("Add owner to task %d", task.ID),
				"Task":  task,
				"Error": "This user already has access to this task"})
		return
	}

	//ユーザー名からユーザーIDを取得する
	var userID int
	err = db.Get(&userID, "SELECT id FROM users WHERE name=?", username)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//Ownershipの追加
	_, err = db.Exec("INSERT INTO ownership (owner_id, task_id) VALUES (?,?)", userID, id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	//リダイレクト処理
	path := fmt.Sprintf("/task/%d", id)
	ctx.Redirect(http.StatusFound, path)
}
