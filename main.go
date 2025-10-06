package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"tododemo.com/m/controllers"
	"tododemo.com/m/middleware"
	"tododemo.com/m/models"
)

func main() {
	fmt.Print("Starting RESTful API server...\n")
	router := gin.Default()
	router.SetTrustedProxies([]string{"localhost"})

	models.ConnectDatabase()
	router.GET("/", controllers.HomePage)
	privateRouter := router.Group("/api").Use(middleware.AuthJwt)
	{
		privateRouter.GET("/todos", controllers.GetTodos)
		privateRouter.GET("/todos/:id", controllers.GetTodo)
		privateRouter.PATCH("/todos/:id", controllers.ToggleTodoStatus)
		privateRouter.POST("/todos", controllers.AddTodo)
		privateRouter.DELETE("/todos/:id", controllers.DeleteTodo)
	}
	router.POST("/login", middleware.Login)
	router.POST("/user/register", controllers.RegisterUser)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Route not found", "code": "PAGE_NOT_FOUND"})
	})

	router.Run("localhost:9090")
}
