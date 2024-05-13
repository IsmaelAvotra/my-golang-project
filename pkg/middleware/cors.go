package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
		MaxAge: 12 * time.Hour,
	})
}