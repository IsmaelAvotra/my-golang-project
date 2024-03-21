package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth() gin.HandlerFunc {
	hashedSecretKey := hashSecretKey([]byte(os.Getenv("API_SECRET_KEY")))

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		hashedApiKey := hashSecretKey([]byte(apiKey))

		if hashedApiKey == hashedSecretKey {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
		}
	}
}

func hashSecretKey(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}