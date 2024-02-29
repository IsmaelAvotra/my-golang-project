package handlers

import (
	"github.com/IsmaelAvotra/pkg/database"
	"github.com/IsmaelAvotra/pkg/models"
	"github.com/IsmaelAvotra/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateJob(c *gin.Context) {
	var jobToCreate models.Job

	if err := c.ShouldBindJSON(&jobToCreate); err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}
	existingUniverity, err := database.GetUnivByName(jobToCreate.Name)

	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, "error checking job name uniqueness")
	}

	if existingUniverity != nil {
		utils.ErrorResponse(c, StatusBadRequest, "univeristy with this name already exists")
		return
	}

	newJob := models.Job{
		Name:               jobToCreate.Name,
		About:              jobToCreate.About,
		WorkingEnvironment: jobToCreate.WorkingEnvironment,
		Formation:          jobToCreate.Formation,
		SectorID:           jobToCreate.SectorID,
	}

	insertResult, err := database.DB.Collection("jobs").InsertOne(c, newJob)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, "could not save the new job")
		return
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)

	if !ok {
		utils.ErrorResponse(c, StatusInternalServerError, "invalid inserted ID")
		return
	}

	c.JSON(StatusOK, gin.H{"message": "job added successful", "jobId": insertedID.Hex()})
}

func CreateSector(c *gin.Context) {
	sectorToCreate := models.Sector{}

	err := c.ShouldBindJSON(&sectorToCreate)

	if err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}

	newSector := models.Sector{
		Name: sectorToCreate.Name,
	}

	insertResult, err := database.DB.Collection("sectors").InsertOne(c, newSector)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, "could not save the new sector")
		return
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)

	if !ok {
		utils.ErrorResponse(c, StatusInternalServerError, "invalid inserted ID")
		return
	}

	c.JSON(StatusOK, gin.H{"message": "sector added successful", "sectorId": insertedID.Hex()})
}
