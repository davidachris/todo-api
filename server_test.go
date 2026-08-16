package main

import (
	"bytes"
	"davidc/todo-api/database"
	"davidc/todo-api/models"
	"davidc/todo-api/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func routerAndHttpTest() (*gin.Engine, *httptest.ResponseRecorder) {
	return setupRouter(), httptest.NewRecorder()
}

func initTestDB(t *testing.T) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "todos.db")
	err := database.InitDb(services.NewNoopStore(), &services.S3Details{FileName: path})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPingRoute(t *testing.T) {
	router, w := routerAndHttpTest()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestCreateTask(t *testing.T) {
	initTestDB(t)
	newTask := models.TaskRequest{
		Description: "Hello",
		Completed:   0,
	}
	newTaskJson, _ := json.Marshal(newTask)
	router, w := routerAndHttpTest()
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(newTaskJson))
	req.Header.Add("Content-Type", "application/json")
	defer req.Body.Close()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}
