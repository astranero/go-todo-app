package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTodos(c *gin.Context){
	var todos []Todo 
	DB.Find(%todos)
	c.JSON(http.StatusOK, todos)
}


func CreateTodo(c *gin.Context){
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	DB.Create(&todo)
	c.JSON(http.statusOK, todo)
}