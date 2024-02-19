package api

import (
	"github.com/IsmaelAvotra/pkg/auth"
	"github.com/IsmaelAvotra/pkg/handlers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", auth.LoginHandler)
		v1.POST("/register", auth.RegisterHandler)

		v1.GET("/users", handlers.GetUsersHandler)
		v1.GET("/users/:id", handlers.GetUserHandler)
		v1.DELETE("/users/:id", handlers.DeleteUserHandler)
		v1.PATCH("/users/:id", handlers.UpdateUserHandler)

		v1.GET("/universities", handlers.GetUniversitiesHandler)
		v1.GET("/events/:eventId", handlers.GetEventById)
		v1.DELETE("/events/:eventId", handlers.DeleteEvent)
		v1.PATCH("/events/:eventId", handlers.UpdateEvent)
		v1.POST("/create-university", handlers.CreateUniverity)
	}

	return r
}
