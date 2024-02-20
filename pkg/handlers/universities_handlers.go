package handlers

import (
	"context"

	"github.com/IsmaelAvotra/pkg/database"
	"github.com/IsmaelAvotra/pkg/models"
	"github.com/IsmaelAvotra/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
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
		utils.ErrorResponse(c, StatusInternalServerError, "error checking univeristy name uniqueness")
	}

	if existingUniverity != nil {
		utils.ErrorResponse(c, StatusBadRequest, "univeristy with this name already exists")
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
		utils.ErrorResponse(c, StatusInternalServerError, "could not save the university")
		return
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)

	if !ok {
		utils.ErrorResponse(c, StatusInternalServerError, "invalid inserted ID")
		return
	}

	c.JSON(StatusOK, gin.H{"message": "university added successful", "univId": insertedID.Hex()})
}

func GetUniversitiesHandler(c *gin.Context) {
	universities, err := database.GetAllUniversities()
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}
	c.JSON(StatusOK, universities)
}

func GetUniversityHandler(c *gin.Context) {
	univId := c.Param("univId")

	university, err := database.GetUnivById(univId)
	if err != nil {
		utils.ErrorResponse(c, StatusNotFound, "university not found.")
		return
	}
	c.JSON(StatusOK, university)
}

func DeleteUniversityHandler(c *gin.Context) {
	univID := c.Param("univId")

	err := database.DeleteUniversity(univID)

	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}

	c.JSON(StatusOK, gin.H{"message": "university deleted with success"})
}

func UpdateUniversityHandler(c *gin.Context) {
	univID := c.Param("univId")
	var update bson.M

	err := c.ShouldBindBodyWith(&update, binding.JSON)
	if err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}

	err = database.UpdateUniversity(univID, update)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}
	c.JSON(StatusOK, gin.H{"message": "university updated successfully"})
}
