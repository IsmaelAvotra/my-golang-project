package auth

import (
	"errors"
	"os"

	"net/http"
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

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

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

	accessToken, refreshToken, err := GenerateTokens(dbUser.Email, dbUser.Role)

	if err != nil {
		utils.ErrorResponse(c, statusInternalServerError, err.Error())
		return
	}

	c.JSON(statusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func GenerateTokens(email, role string) (string, string, error) {
	accessTokenClaims := &Claims{
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   "access_token",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(JwtKey)
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   "refresh_token",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(JwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func RefreshTokenHandler(c *gin.Context) {
	refreshToken := c.Request.Header.Get("refresh_token")
	if refreshToken == "" {
		utils.ErrorResponse(c, statusUnauthorized, "refresh token is required")
		return
	}

	refreshClaims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		utils.ErrorResponse(c, statusUnauthorized, "invalid refresh token")
		return
	}

	newAccessToken, newRefreshToken, err := GenerateTokens(refreshClaims.Email, refreshClaims.Role)
	if err != nil {
		utils.ErrorResponse(c, statusUnauthorized, "failed to generate new tokens")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func MyProtectedAdminEndpoint(c *gin.Context) {
	c.JSON(statusOK, gin.H{"message": "Your are authorized to access this resource"})
}

func RegisterHandler(c *gin.Context) {
	var userToCreate models.User

	if err := c.ShouldBindJSON(&userToCreate); err != nil {
		utils.ErrorResponse(c, statusBadRequest, err.Error())
		return
	}

	existingUserByEmail, err := database.GetUserByEmail(userToCreate.Email)

	if err != nil {
		if err == errors.New("DB nil") {
			utils.ErrorResponse(c, statusInternalServerError, "Database connection error nil")
		} else {
			utils.ErrorResponse(c, statusInternalServerError, err.Error())
		}
		return
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
