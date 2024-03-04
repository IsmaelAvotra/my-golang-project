package handlers

import (
	"net/http"
	"net/url"
	"strings"

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

	insertResult, err := database.DB.Collection("universities").InsertOne(c, newUniversity)
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
	encodedProgramName := c.Query("programName")
	encodedUnivName := c.Query("univName")
	encodedProvince := c.Query("province")
	encodedRegion := c.Query("region")
	encodedCity := c.Query("city")

	programName, err := url.QueryUnescape(encodedProgramName)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}

	univName, err := url.QueryUnescape(encodedUnivName)

	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}

	province, err := url.QueryUnescape(encodedProvince)

	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}

	region, err := url.QueryUnescape(encodedRegion)

	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}

	city, err := url.QueryUnescape(encodedCity)

	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}

	filter := bson.M{}
	if programName != "" {
		program, err := database.GetProgramByName(programName)
		if err != nil {
			utils.ErrorResponse(c, StatusInternalServerError, err.Error())
			return
		}
		if program != nil {
			filter["programIDs"] = program.ID
		} else {
			utils.ErrorResponse(c, StatusNotFound, "Program not found")
			return
		}
	}

	if univName != "" {
		normalizedUnivName := utils.RemoveAccents(univName)
		filter["univName"] = bson.M{"$regex": primitive.Regex{Pattern: normalizedUnivName, Options: "i"}}
	}

	if province != "" {
		filter["location.province"] = bson.M{"$regex": primitive.Regex{Pattern: province, Options: "i"}}
	}

	if region != "" {
		filter["location.region"] = bson.M{"$regex": primitive.Regex{Pattern: region, Options: "i"}}

	}

	if city != "" {
		filter["location.city"] = bson.M{"$regex": primitive.Regex{Pattern: city, Options: "i"}}

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
			"univName":        university.Name,
			"univLocation":    university.Location,
			"presentation":    university.Presentation,
			"isPrivate":       university.IsPrivate,
			"tuition":         university.Tuition,
			"contact":         university.Contact,
			"imageUrl":        university.ImageURL,
			"documentUrl":     university.DocumentURL,
			"programsId":      university.ProgramIDs,
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

	insertResult, err := database.DB.Collection("programs").InsertOne(c, programToCreate)
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

func GetProgramsFilteredHandler(c *gin.Context) {
	careerProspect, err := url.QueryUnescape(c.Query("careerProspect"))

	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}

	filter := bson.M{}
	if careerProspect != "" {
		regexPattern := primitive.Regex{Pattern: strings.ToLower(careerProspect), Options: "i"}
		filter["careerprospects"] = regexPattern
	}

	programs, err := database.GetAllPrograms(filter)
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
	update := bson.M{}

	if program.ProgramName != "" {
		update["$set"] = bson.M{"programname": program.ProgramName}
	}
	if program.Level != "" {
		if update["$set"] == nil {
			update["$set"] = bson.M{}
		}
		update["$set"].(bson.M)["level"] = program.Level
	}
	if program.Duration != 0 {
		if update["$set"] == nil {
			update["$set"] = bson.M{}
		}
		update["$set"].(bson.M)["duration"] = program.Duration
	}
	if program.Requirements != nil {
		if update["$set"] == nil {
			update["$set"] = bson.M{}
		}
		update["$set"].(bson.M)["requirements"] = program.Requirements
	}
	if program.CareerProspects != nil {
		if update["$set"] == nil {
			update["$set"] = bson.M{}
		}
		update["$set"].(bson.M)["careerprospects"] = program.CareerProspects
	}

	err = database.UpdateProgram(programID, update)
	if err != nil {
		utils.ErrorResponse(c, StatusInternalServerError, err.Error())
		return
	}
	c.JSON(StatusOK, gin.H{"message": "program updated successfully"})
}

// Add to favorites
func AddUniversityToFavoritesHandler(c *gin.Context) {
	userID := c.Param("userId")
	univID := c.Param("univId")

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	univIDObj, err := primitive.ObjectIDFromHex(univID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	err = database.AddFavoriteUniversity(userIDObj, univIDObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add university to favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "University added to favorites successfully"})
}

func RemoveUniversityToFavoritesHandler(c *gin.Context) {
	userID := c.Param("userId")
	univID := c.Param("univId")

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	univIDObj, err := primitive.ObjectIDFromHex(univID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	err = database.RemoveFavoriteUniversity(userIDObj, univIDObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove university to favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "University removed to favorites successfully"})
}
