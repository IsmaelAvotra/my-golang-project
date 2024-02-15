package handlers

import (
	"net/http"

	"github.com/IsmaelAvotra/pkg/database"
	"github.com/IsmaelAvotra/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	statusNotFound            = http.StatusNotFound
	statusInternalServerError = http.StatusInternalServerError
	statusOK                  = http.StatusOK
	statusBadRequest          = http.StatusBadRequest
)

func GetUsersHandler(c *gin.Context) {
	users, err := database.GetAllUsers()
	if err != nil {
		utils.ErrorResponse(c, statusInternalServerError, err.Error())
		return
	}

	c.JSON(statusOK, users)
}

func GetUserHandler(c *gin.Context) {
	userId := c.Param("id")

	user, err := database.GetUserByID(userId)

	if err != nil {
		utils.ErrorResponse(c, statusBadRequest, err.Error())
		return
	}

	if user == nil {
		utils.ErrorResponse(c, statusNotFound, "User not found")
		return
	}

	c.JSON(statusOK, user)
}

func DeleteUserHandler(c *gin.Context) {

}

func UpdateUserHandler(c *gin.Context) {

}
