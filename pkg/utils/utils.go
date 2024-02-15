package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// Hash the password
const costHash = 14

func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", errors.New("password length must be at least 8 characters")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), costHash)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
