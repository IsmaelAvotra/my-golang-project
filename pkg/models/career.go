package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Job struct {
	JobId              primitive.ObjectID `json:"jobId,omitempty" bson:"_id,omitempty"`
	Name               string             `json:"jobName"`
	About              About              `json:"about"`
	WorkingEnvironment WorkingEnvironment `json:"workingEnvironment"`
	Formation          string             `json:"formation"`
	SectorID           primitive.ObjectID `json:"sectorID,omitempty" bson:"sectorID,omitempty"`
}

type About struct {
	Description           string        `json:"description"`
	Missions              []string      `json:"missions"`
	Skills                QualitySkills `json:"skills"`
	ProfessionalEvolution string        `json:"professionalEvolution"`
}

type WorkingEnvironment struct {
	Presentation  string `json:"presentation"`
	ExercicePlace string `json:"exercicePlace"`
}

type QualitySkills struct {
	Knowledges []string `json:"knowledges"`
	KnowHow    []string `json:"knowHow"`
}

type Sector struct {
	SectorId primitive.ObjectID `json:"sectorId,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"sectorName"`
}
