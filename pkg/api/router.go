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

		v1.GET("/universities", handlers.GetFilteredUniversitiesHandler)
		v1.GET("/universities/:univId", handlers.GetUniversityHandler)
		v1.DELETE("/universities/:univId", handlers.DeleteUniversityHandler)
		v1.PATCH("/universities/:univId", handlers.UpdateUniversityHandler)
		v1.POST("/create-university", handlers.CreateUniverity)

		v1.POST("/universities/create-program", handlers.CreateProgramHandler)
		v1.GET("/universities/programs", handlers.GetProgramsHandler)
		v1.GET("/universities/programs/:programId", handlers.GetProgramHandler)
		v1.PATCH("/universities/programs/:programId", handlers.UpdateProgramHandler)
		v1.DELETE("/universities/programs/:programId", handlers.DeleteProgramHandler)

	}
	return r
}
