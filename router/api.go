package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/controllers"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/database"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/middlewares"
)

func Api() {
	database.ConDB()

	// User Endpoint
	r := gin.Default()
	r.POST("/users/register", controllers.Register)
	r.POST("/users/login", controllers.Login)

	r.Use(middlewares.JWTMiddleware())

	r.PUT("/users/:userId", controllers.Update)
	r.POST("/users/logout", controllers.Logout)
	r.DELETE("/users/:userId", controllers.Delete)

	// Photos Endpoint
	r.POST("/photos", controllers.Upload)
	r.GET("/photos", controllers.Show)
	r.PUT("/photos/:photoId", controllers.Upload)
	r.DELETE("/photos/:photoId", controllers.Remove)

	r.Run(":8080")
}
