package main

import (
	"errors"
	"net/http"

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

			context.IndentedJSON(http.StatusOK, gin.H{
				"message": "Todo deleted successfully!",
			})
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{
		"message": "Todo not found!",
	})
}

func main() {
	router := gin.Default()
	router.POST("/todos", addTodo)
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodo)
	router.PATCH("/todos/:id", toogleTodoStatus)
	router.DELETE("/todos/:id", deleteTodoByID)
	router.Run("localhost:8080")
}
