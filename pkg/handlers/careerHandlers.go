package handlers

import (
	"github.com/IsmaelAvotra/pkg/database"
	"github.com/IsmaelAvotra/pkg/models"
	"github.com/IsmaelAvotra/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

func GetJobsHandler(c *gin.Context) {
	jobs, err := database.GetAllJobs()
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}
	c.JSON(StatusOK, jobs)
}

func GetJobHandler(c *gin.Context) {
	jobId := c.Param("jobId")

	job, err := database.GetJobById(jobId)

	if err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}

	if job == nil {
		utils.ErrorResponse(c, StatusNotFound, "user not found")
		return
	}

	c.JSON(StatusOK, job)
}

func UpdateJobHandler(c *gin.Context) {
	job := models.Job{}
	JobId := c.Param("jobId")

	err := c.BindJSON(&job)
	if err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}

	update := bson.M{}
	set := bson.M{}
	var emptyObjectID primitive.ObjectID

	if job.Name != "" {
		set["name"] = job.Name
	}
	if job.About.Description != "" {
		set["about.description"] = job.About.Description
	}
	if len(job.About.Missions) > 0 {
		set["about.missions"] = job.About.Missions
	}
	if len(job.About.Skills.Knowledges) > 0 {
		set["about.skills.knowledges"] = job.About.Skills.Knowledges
	}
	if len(job.About.Skills.KnowHow) > 0 {
		set["about.skills.knowhow"] = job.About.Skills.KnowHow
	}
	if job.About.ProfessionalEvolution != "" {
		set["about.professionalevolution"] = job.About.ProfessionalEvolution
	}
	if job.WorkingEnvironment.Presentation != "" {
		set["workingEnvironment.presentation"] = job.WorkingEnvironment.Presentation
	}
	if job.WorkingEnvironment.ExercicePlace != "" {
		set["workingEnvironment.exerciceplace"] = job.WorkingEnvironment.ExercicePlace
	}
	if job.Formation != "" {
		set["formation"] = job.Formation
	}

	if job.SectorID != emptyObjectID {
		set["sectorID"] = job.SectorID
	}

	if len(set) > 0 {
		update["$set"] = set
	}

	err = database.UpdateJobById(JobId, update)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}
	c.JSON(StatusOK, gin.H{"message": "job updated successfully"})
}

func DeleteJobHandler(c *gin.Context) {
	jobId := c.Param("jobId")

	err := database.DeleteJob(jobId)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}

	c.JSON(StatusOK, gin.H{"message": "job deleted successfully"})
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
