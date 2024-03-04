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

	if err := c.BindJSON(&university); err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}

	update := bson.M{}
	set := bson.M{}

	if university.Name != "" {
		set["univName"] = university.Name
	}
	if university.Location != (models.Location{}) {
		set["univLocation"] = university.Location
	}
	if university.Presentation != "" {
		set["presentation"] = university.Presentation
	}
	if university.IsPrivate {
		set["isPrivate"] = university.IsPrivate
	}
	if university.Tuition != 0 {
		set["tuition"] = university.Tuition
	}
	if university.Contact != (models.Contact{}) {
		set["contact"] = university.Contact
	}
	if university.ImageURL != "" {
		set["imageUrl"] = university.ImageURL
	}
	if university.DocumentURL != "" {
		set["documentUrl"] = university.DocumentURL
	}
	if len(university.ProgramIDs) > 0 {
		set["programIDs"] = university.ProgramIDs
	}
	if len(university.Infrastructure) > 0 {
		set["infrastructure"] = university.Infrastructure
	}
	if len(university.Partnerships) > 0 {
		set["partnerships"] = university.Partnerships
	}
	if university.SuccessDiplomas != 0 {
		set["successDiplomas"] = university.SuccessDiplomas
	}
	if len(university.Events) > 0 {
		set["events"] = university.Events
	}
	if len(university.News) > 0 {
		set["news"] = university.News
	}
	if len(university.Photos) > 0 {
		set["photos"] = university.Photos
	}

	if len(set) > 0 {
		update["$set"] = set
	}

	if err := database.UpdateUniversity(univID, update); err != nil {
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

	if err := c.BindJSON(&program); err != nil {
		utils.ErrorResponse(c, StatusBadRequest, err.Error())
		return
	}

	update := bson.M{}
	set := bson.M{}

	if program.ProgramName != "" {
		set["programname"] = program.ProgramName
	}
	if program.Level != "" {
		set["level"] = program.Level
	}
	if program.Duration != 0 {
		set["duration"] = program.Duration
	}
	if program.Requirements != nil {
		set["requirements"] = program.Requirements
	}
	if program.CareerProspects != nil {
		set["careerprospects"] = program.CareerProspects
	}

	if len(set) > 0 {
		update["$set"] = set
	}

	if err := database.UpdateProgram(programID, update); err != nil {
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
