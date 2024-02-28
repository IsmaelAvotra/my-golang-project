package utils

import (
	"errors"
	"strings"

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

func RemoveAccents(input string) string {
	replacements := map[rune]rune{
		'À': 'A', 'Á': 'A', 'Â': 'A', 'Ã': 'A', 'Ä': 'A', 'Å': 'A',
		'à': 'a', 'á': 'a', 'â': 'a', 'ã': 'a', 'ä': 'a', 'å': 'a',
		'È': 'E', 'É': 'E', 'Ê': 'E', 'Ë': 'E',
		'è': 'e', 'é': 'e', 'ê': 'e', 'ë': 'e',
		'Ì': 'I', 'Í': 'I', 'Î': 'I', 'Ï': 'I',
		'ì': 'i', 'í': 'i', 'î': 'i', 'ï': 'i',
		'Ò': 'O', 'Ó': 'O', 'Ô': 'O', 'Õ': 'O', 'Ö': 'O',
		'ò': 'o', 'ó': 'o', 'ô': 'o', 'õ': 'o', 'ö': 'o',
		'Ù': 'U', 'Ú': 'U', 'Û': 'U', 'Ü': 'U',
		'ù': 'u', 'ú': 'u', 'û': 'u', 'ü': 'u',
		'Ç': 'C', 'ç': 'c',
	}
	var output strings.Builder
	for _, char := range input {
		if replacement, ok := replacements[char]; ok {
			output.WriteRune(replacement)
		} else {
			output.WriteRune(char)
		}
	}
	return output.String()
}
