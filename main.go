package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"tododemo.com/m/controllers"
	"tododemo.com/m/middleware"
	"tododemo.com/m/models"
)

func main() {
	fmt.Print("Starting RESTful API server...\n")
	err := godotenv.Load()
	if err != nil {
		fmt.Print("Error loading .env file, proceeding with environment variables\n")
	}
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
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
	http.HandleFunc("/oauth/google/login", controllers.HandleGoogleLogin)
	router.GET("/oauth/google/callback", controllers.HandleGoogleCallback)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Route not found", "code": "PAGE_NOT_FOUND"})
	})

	router.Run("localhost:9090")
}
