package middleware

import (
	"errors"
	"time"

	"github.com/IsmaelAvotra/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := ExtractTokenFromRequest(c)
		if err != nil {
			utils.ErrorResponse(c, StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				switch ve.Errors {
				case jwt.ValidationErrorMalformed:
					utils.ErrorResponse(c, StatusUnauthorized, "malformed token")
				case jwt.ValidationErrorUnverifiable:
					utils.ErrorResponse(c, StatusUnauthorized, "token could not be verified")
				case jwt.ValidationErrorSignatureInvalid:
					utils.ErrorResponse(c, StatusUnauthorized, "invalid token signature")

				case jwt.ValidationErrorExpired:
					refreshToken := c.Request.Header.Get("refresh_token")
					if refreshToken == "" {
						utils.ErrorResponse(c, StatusUnauthorized, "expired token, refresh token is required")
						c.Abort()
						return
					}

					refreshClaims, err := ValidateRefreshToken(refreshToken)
					if err != nil {
						utils.ErrorResponse(c, StatusUnauthorized, "invalid refresh token")
						c.Abort()
						return
					}

					newAccessToken, _, err := GenerateTokens(refreshClaims.Email, refreshClaims.Role)
					if err != nil {
						utils.ErrorResponse(c, StatusUnauthorized, "failed to generate new access token")
						c.Abort()
						return
					}

					c.Request.Header.Set("Authorization", "Bearer "+newAccessToken)
					accessToken = newAccessToken

				default:
					utils.ErrorResponse(c, StatusUnauthorized, "invalid token")
				}
			} else {
				utils.ErrorResponse(c, StatusUnauthorized, err.Error())
			}
			c.Abort()
			return
		}

		if !parsedToken.Valid {
			utils.ErrorResponse(c, StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
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
