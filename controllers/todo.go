package controllers

import (
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"tododemo.com/m/models"
)

func HomePage(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Todo App",
	})
}

func GetTodos(context *gin.Context) {
	var todos []models.Todos
	if err := models.DB.Find(&todos).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	context.JSON(http.StatusOK, gin.H{"data": todos})
}

func GetTodo(context *gin.Context) {
	var todo models.Todos
	id := context.Param("id")
	if err := models.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Todo not found!"})
	}
	context.JSON(http.StatusOK, gin.H{"data": todo})
}

func AddTodo(context *gin.Context) {
	var newTodo models.Todos
	if err := context.BindJSON(&newTodo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	newTodo.Completed = false
	currentTime := time.Now()
	newTodo.AddedAt = currentTime.Unix()
	newTodo.UpdatedAt = currentTime.Unix()
	if err := models.DB.Create(&newTodo).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	context.JSON(http.StatusCreated, gin.H{"data": newTodo})
}

func ToggleTodoStatus(context *gin.Context) {
	var todo models.Todos
	id := context.Param("id")
	if err := models.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Todo not found!"})
		return
	}
	todo.Completed = !todo.Completed
	todo.UpdatedAt = time.Now().Unix()
	models.DB.Save(&todo)
	context.JSON(http.StatusOK, gin.H{"data": todo})
}

func DeleteTodo(context *gin.Context) {
	var todo models.Todos
	id := context.Param("id")
	if err := models.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Todo not found!"})
		return
	}
	models.DB.Delete(&todo)
	context.JSON(http.StatusOK, gin.H{"data": true})
}
