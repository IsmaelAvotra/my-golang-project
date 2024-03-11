package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"

	"net/http"
	"os"
	"time"

	"github.com/IsmaelAvotra/pkg/database"
	"github.com/IsmaelAvotra/pkg/models"
	"github.com/IsmaelAvotra/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	statusInternalServerError = http.StatusInternalServerError
	statusBadRequest          = http.StatusBadRequest
	statusOK                  = http.StatusOK
	statusUnauthorized        = http.StatusUnauthorized
	statusForbidden           = http.StatusForbidden
)

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func LoginHandler(c *gin.Context) {
	incomingUser := models.LoginUser{}
	dbUser := models.User{}

	if err := c.ShouldBindJSON(&incomingUser); err != nil {
		utils.ErrorResponse(c, statusBadRequest, err.Error())
		return
	}
	filter := bson.M{"email": incomingUser.Email}
	err := database.DB.Collection("users").FindOne(c, filter).Decode(&dbUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ErrorResponse(c, statusUnauthorized, err.Error())
			return
		}
		utils.ErrorResponse(c, statusInternalServerError, err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(incomingUser.Password)); err != nil {
		utils.ErrorResponse(c, statusUnauthorized, "email or password is incorrect")
		return
	}

	token, err := GenerateToken(dbUser.Email, dbUser.Role)

	if err != nil {
		utils.ErrorResponse(c, statusInternalServerError, err.Error())
		return
	}

	c.JSON(statusOK, gin.H{"token": token})

}

func GenerateToken(email string, role string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute).Unix()

	claims := &Claims{
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JwtKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ExtractTokenFromRequest(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
		return "", errors.New("invalid authorization header")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	return token, nil
}

func IsAdmin(c *gin.Context) (bool, error) {
	token, err := ExtractTokenFromRequest(c)
	if err != nil {
		return false, err
	}
	claims, err := ValidateJWTToken(token)

	role := claims["role"]
	if role != "admin" {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func MyProtectedAdminEndpoint(c *gin.Context) {
	isAdmin, err := IsAdmin(c)
	if err != nil {
		utils.ErrorResponse(c, statusUnauthorized, err.Error())
		return
	}
	if !isAdmin {
		utils.ErrorResponse(c, statusForbidden, "You are not authorized to access this resource")
		return
	}
	c.JSON(statusOK, gin.H{"message": "Your are authorized to access this resource"})
}

func ValidateJWTToken(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch ve.Errors {
			case jwt.ValidationErrorMalformed:
				return nil, errors.New("malformed token")
			case jwt.ValidationErrorUnverifiable:
				return nil, errors.New("token could not be verified")
			case jwt.ValidationErrorSignatureInvalid:
				return nil, errors.New("invalid token signature")
			case jwt.ValidationErrorExpired:
				return nil, errors.New("expired token")
			default:
				return nil, errors.New("invalid token")
			}
		} else {
			return nil, err
		}
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func RegisterHandler(c *gin.Context) {
	var userToCreate models.User

	if err := c.ShouldBindJSON(&userToCreate); err != nil {
		utils.ErrorResponse(c, statusBadRequest, err.Error())
		return
	}

	existingUserByEmail, err := database.GetUserByEmail(userToCreate.Email)

	if err != nil {
		utils.ErrorResponse(c, statusInternalServerError, "Error checking email uniqueness")
	}

	if existingUserByEmail != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Email already exists")
		return
	}

	existingUserByUsername, err := database.GetUserByUsername(userToCreate.Username)
	if err != nil {
		utils.ErrorResponse(c, statusInternalServerError, "Error checking username uniqueness")
	}

	if existingUserByUsername != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Username already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(userToCreate.Password)
	if err != nil {
		utils.ErrorResponse(c, statusInternalServerError, err.Error())
		return
	}

	userToCreate.Role = "normal"

	newUser := models.User{
		Username:  userToCreate.Username,
		Email:     userToCreate.Email,
		Password:  hashedPassword,
		Role:      userToCreate.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertResult, err := database.DB.Collection("users").InsertOne(c, newUser)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Could not save user")
		return
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Invalid inserted ID")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful", "id": insertedID.Hex()})
}

func GenerateRandomKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic("Failed to generate random key: " + err.Error())
	}

	return base64.StdEncoding.EncodeToString(key)
}
