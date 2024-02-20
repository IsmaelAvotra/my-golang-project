package handlers

import (
	"net/http"

	"github.com/IsmaelAvotra/pkg/database"
	"github.com/IsmaelAvotra/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	StatusNotFound            = http.StatusNotFound
	StatusInternalServerError = http.StatusInternalServerError
	StatusOK                  = http.StatusOK
	StatusBadRequest          = http.StatusBadRequest
)

func GetUsersHandler(c *gin.Context) {
	users, err := database.GetAllUsers()
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}

	c.JSON(StatusOK, users)
}

func GetUserHandler(c *gin.Context) {
	userId := c.Param("id")

	user, err := database.GetUserByID(userId)

	if err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}

	if user == nil {
		utils.ErrorResponse(c, StatusNotFound, "user not found")
		return
	}

	c.JSON(StatusOK, user)
}

func DeleteUserHandler(c *gin.Context) {
	userId := c.Param("id")

	err := database.DeleteUser(userId)

	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}
	c.JSON(StatusOK, gin.H{"message": "user deleted with success"})
}

func UpdateUserHandler(c *gin.Context) {
	userId := c.Param("id")

	var update bson.M

	if err := c.BindJSON(&update); err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}

	err := database.UpdateUser(userId, update)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}

	c.JSON(StatusOK, gin.H{"message": "user updated successfully"})

}
