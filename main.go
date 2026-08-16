package main

import (
	"davidc/todo-api/database"
	"davidc/todo-api/models"
	"davidc/todo-api/query"
	"davidc/todo-api/services"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

var (
	taskDbLock  = sync.RWMutex{}
	tasksDb     = make(map[uint64]models.Task)
	taskIdCount atomic.Uint64
)

const (
	TaskNotExistError = "Task does not exist."
	IdParamError      = "Incorrect id format."
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", ping)
	router.GET("/tasks", getAllTasks)
	router.POST("/tasks", postCreateTask)
	router.PUT("/tasks/:id", putUpdateTask)
	router.DELETE("/tasks/:id", deleteTask)
	return router
}

func setupStorage() {
	fileStore, details := services.ConfigStorage()
	log.Println("About to Init DB")
	err := database.InitDb(fileStore, details)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	setupStorage()
	router := setupRouter()
	router.Run(":8080")
}

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

// Get all Tasks
func getAllTasks(ctx *gin.Context) {
	allTasks, err := query.SelectAllTodos()
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, map[string]string{"err": "Internal Server Error"})
		return
	}
	log.Println(allTasks)
	ctx.IndentedJSON(http.StatusOK, allTasks)
}

// Create a task
func postCreateTask(ctx *gin.Context) {
	var newTask models.TaskRequest
	if err := ctx.BindJSON(&newTask); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, map[string]string{"err": "Internal Server Error"})
		return
	}
	err := query.CreateTask(newTask)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, map[string]string{"err": "Internal Server Error"})
		return
	}
	ctx.Status(http.StatusCreated)
}

// Update a Task
func putUpdateTask(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.String(http.StatusBadRequest, IdParamError)
		return
	}
	var updateTask models.TaskRequest
	if err := ctx.BindJSON(&updateTask); err != nil {
		return
	}
	taskDbLock.RLock()
	val, ok := tasksDb[id]
	taskDbLock.RUnlock()
	if !ok {
		ctx.String(http.StatusNotFound, TaskNotExistError)
		return
	}
	if val.Description != updateTask.Description {
		val.Description = updateTask.Description
	}
	if val.Completed != updateTask.Completed {
		val.Completed = updateTask.Completed
	}
	taskDbLock.Lock()
	tasksDb[id] = val
	taskDbLock.Unlock()
	ctx.IndentedJSON(http.StatusOK, val)
}

// Delete a Task
func deleteTask(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.String(http.StatusBadRequest, IdParamError)
		return
	}
	if _, ok := tasksDb[id]; !ok {
		ctx.String(http.StatusNotFound, TaskNotExistError)
		return
	}
	taskDbLock.Lock()
	delete(tasksDb, id)
	taskDbLock.Unlock()
	ctx.Status(http.StatusNoContent)
}
