package handlers

import (
	"context"

	"github.com/IsmaelAvotra/pkg/database"
	"github.com/IsmaelAvotra/pkg/models"
	"github.com/IsmaelAvotra/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUniverity(c *gin.Context) {
	var univToCreate models.University

	if err := c.ShouldBindJSON(&univToCreate); err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}
	existingUniverity, err := database.GetUnivByName(univToCreate.Name)

	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, "Error checking univeristy name uniqueness")
	}

	if existingUniverity != nil {
		utils.ErrorResponse(c, StatusBadRequest, "Univeristy with this name already exists")
		return
	}

	newUniversity := models.University{
		Name:            univToCreate.Name,
		Location:        univToCreate.Location,
		Presentation:    univToCreate.Presentation,
		IsPrivate:       univToCreate.IsPrivate,
		Tuition:         univToCreate.Tuition,
		Contact:         univToCreate.Contact,
		ImageURL:        univToCreate.ImageURL,
		DocumentURL:     univToCreate.DocumentURL,
		Programs:        univToCreate.Programs,
		Infrastructure:  univToCreate.Infrastructure,
		Partnerships:    univToCreate.Partnerships,
		SuccessDiplomas: univToCreate.SuccessDiplomas,
		Events:          univToCreate.Events,
		News:            univToCreate.News,
		Photos:          univToCreate.Photos,
	}

	insertResult, err := database.DB.Collection("universities").InsertOne(context.Background(), newUniversity)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, "Could not save the university")
		return
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)

	if !ok {
		utils.ErrorResponse(c, StatusInternalServerError, "Invalid inserted ID")
		return
	}

	c.JSON(StatusOK, gin.H{"message": "University added successful", "univId": insertedID.Hex()})
}

func GetAllEvents(c *gin.Context) {

}

func GetEventById(c *gin.Context) {

}

func DeleteEvent(c *gin.Context) {

}

func UpdateEvent(c *gin.Context) {

}
