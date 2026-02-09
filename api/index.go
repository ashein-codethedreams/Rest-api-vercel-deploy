package handler // Can also be 'package main', but Vercel prefers this structure inside /api

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{ID: "1", Item: "Learn Go", Completed: false},
	{ID: "2", Item: "Build a REST API", Completed: false},
	{ID: "3", Item: "Become a Fullstack Developer", Completed: false},
}

var (
	router *gin.Engine
	once   sync.Once
)

// InitRouter initializes the Gin engine once
func InitRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/todos", addTodo)
	r.GET("/todos", getTodos)
	r.GET("/todos/:id", getTodo)
	r.PATCH("/todos/:id", toogleTodoStatus)
	r.DELETE("/todos/:id", deleteTodoByID)
	return r
}

// Handler is the entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		router = InitRouter()
	})
	router.ServeHTTP(w, r)
}

func addTodo(context *gin.Context) {
	var newTodo todo
	if err := context.BindJSON(&newTodo); err != nil {
		return
	}
	todos = append(todos, newTodo)
	context.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

func getTodoByID(id string) (*todo, error) {
	for i, t := range todos {
		if t.ID == id {
			return &todos[i], nil
		}
	}
	return nil, errors.New("todo not found")
}

func getTodo(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoByID(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found!"})
		return
	}
	context.IndentedJSON(http.StatusOK, todo)
}

func toogleTodoStatus(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoByID(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found!"})
		return
	}
	todo.Completed = !todo.Completed
	context.IndentedJSON(http.StatusOK, todo)
}

func deleteTodoByID(context *gin.Context) {
	id := context.Param("id")
	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			context.IndentedJSON(http.StatusOK, gin.H{"message": "Todo deleted successfully!"})
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found!"})
}
