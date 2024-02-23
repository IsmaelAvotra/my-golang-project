package handlers

import (
	"context"

	"github.com/IsmaelAvotra/pkg/database"
	"github.com/IsmaelAvotra/pkg/models"
	"github.com/IsmaelAvotra/pkg/utils"
	"github.com/gin-gonic/gin"
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
		ProgramIDs:      univToCreate.ProgramIDs,
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

func GetFilteredUniversitiesHandler(c *gin.Context) {
	programName := c.Query("programName")
	careerProspect := c.Query("careerProspect")

	filter := bson.M{}
	if programName != "" {
		filter["programs.programName"] = programName
	}
	if careerProspect != "" {
		filter["careerProspects"] = careerProspect
	}

	universities, err := database.GetFilteredUniversities(filter)
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
	university := models.University{}
	univID := c.Param("univId")

	err := c.BindJSON(&university)
	if err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}

	update := bson.M{
		"$set": bson.M{
			"univName":     university.Name,
			"univLocation": university.Location,
			"presentation": university.Presentation,
			"isPrivate":    university.IsPrivate,
			"tuition":      university.Tuition,
			"contact":      university.Contact,
			"imageUrl":     university.ImageURL,
			"documentUrl":  university.DocumentURL,
			// "programsId":        university.Programs,
			"infrastructure":  university.Infrastructure,
			"partnerships":    university.Partnerships,
			"successDiplomas": university.SuccessDiplomas,
			"events":          university.Events,
			"news":            university.News,
			"photos":          university.Photos,
		},
	}

	err = database.UpdateUniversity(univID, update)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}
	c.JSON(StatusOK, gin.H{"message": "university updated successfully"})
}

// for program's university
func CreateProgramHandler(c *gin.Context) {
	var programToCreate models.Program

	if err := c.ShouldBindJSON(&programToCreate); err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}

	existingProgram, err := database.GetProgramByName(programToCreate.ProgramName)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, "error checking program name uniqueness")
		return
	}
	if existingProgram != nil {
		utils.ErrorResponse(c, StatusBadRequest, "program with this name already exists")
		return
	}

	insertResult, err := database.DB.Collection("programs").InsertOne(context.Background(), programToCreate)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, "could not save the program")
		return
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		utils.ErrorResponse(c, StatusInternalServerError, "invalid inserted ID")
		return
	}

	c.JSON(StatusOK, gin.H{"message": "program added successfully", "programId": insertedID.Hex()})
}

func GetProgramsHandler(c *gin.Context) {
	programs, err := database.GetAllPrograms()
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}
	c.JSON(StatusOK, programs)
}

func GetProgramHandler(c *gin.Context) {
	programID := c.Param("programId")

	program, err := database.GetProgramById(programID)
	if err != nil {
		utils.ErrorResponse(c, StatusNotFound, "program not found.")
		return
	}
	c.JSON(StatusOK, program)
}

func DeleteProgramHandler(c *gin.Context) {
	programID := c.Param("programId")

	err := database.DeleteProgram(programID)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}

	c.JSON(StatusOK, gin.H{"message": "program deleted successfully"})
}

func UpdateProgramHandler(c *gin.Context) {
	program := models.Program{}
	programID := c.Param("programId")

	err := c.BindJSON(&program)
	if err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}

	update := bson.M{
		"$set": bson.M{
			"programName":     program.ProgramName,
			"level":           program.Level,
			"duration":        program.Duration,
			"requirements":    program.Requirements,
			"careerProspects": program.CareerProspects,
		},
	}

	err = database.UpdateProgram(programID, update)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}
	c.JSON(StatusOK, gin.H{"message": "program updated successfully"})
}
