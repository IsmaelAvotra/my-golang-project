package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/IsmaelAvotra/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	StatusNotFound            = http.StatusNotFound
	StatusInternalServerError = http.StatusInternalServerError
	StatusOK                  = http.StatusOK
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusForbidden           = http.StatusForbidden
)

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func ExtractTokenFromRequest(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
		return "", errors.New("invalid authorization header")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	return token, nil
}

func IsAdmin(c *gin.Context) (bool, error) {
	value, ok := c.Get("claims")
	if !ok {
		return false, errors.New("unauthorized: missing claims")
	}

	claims, ok := value.(jwt.MapClaims)
	if !ok {
		return false, errors.New("unauthorized: invalid claims format")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return false, errors.New("unauthorized: invalid claims format (missing 'role' key)")
	}

	return role == "admin", nil
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, err := IsAdmin(c)
		if err != nil {
			utils.ErrorResponse(c, StatusUnauthorized, "You are not authorized to access this resource")
			c.Abort()
			return
		}
		if !isAdmin {
			utils.ErrorResponse(c, StatusForbidden, "Forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
