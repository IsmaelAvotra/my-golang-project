package auth

import (
	"crypto/rand"
	"encoding/base64"

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
	statusBadRequest          = http.StatusBadRequest
	statusInternalServerError = http.StatusInternalServerError
	statusOK                  = http.StatusOK
	statusUnauthorized        = http.StatusUnauthorized
)

type Claims struct {
	Email string `json:"email"`
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

	token, err := GenerateToken(dbUser.Email)

	if err != nil {
		utils.ErrorResponse(c, statusInternalServerError, err.Error())
		return
	}

	c.JSON(statusOK, gin.H{"token": token})

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

func GenerateToken(email string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute).Unix()

	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime,
		Issuer:    email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JwtKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRandomKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic("Failed to generate random key: " + err.Error())
	}

	return base64.StdEncoding.EncodeToString(key)
}
