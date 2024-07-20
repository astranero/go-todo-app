package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"fmt"
	"log"
	"time"
)

func main(){
	router := gin.Default()

	randomString := uuid.New().String()
	log.Println("Application Started")

	ticker := time.NewTicker(5 * time.Second)
	for t := range ticker.C {
		fmt.Printf("%s: %s\n", t.Format(time.RFC3339), randomString)
	}

	InitDatabase()

	router.GET("/todos", GetTodos)
	router.POST("/todos", PostTodo)
	router.GET("/todos/:id", GetTodoByID)
	router.PUT("/todos/", UpdateTodo)
	router.DELETE("/todos/:id", DeleteTodo)

	router.Run(":8080")
}