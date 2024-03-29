package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Adress        string `json:"adress"`
	CoordinateGPS string `json:"coordinateGPS"`
	Province      string `json:"province"`
	Region        string `json:"region"`
	City          string `json:"city"`
}

type Event struct {
	Title          string    `json:"eventTitle"`
	Descrioption   string    `json:"description"`
	Date           time.Time `json:"eventDate"`
	Location       string    `json:"eventLocation"`
	IsFree         bool      `json:"isFree"`
	AdmissionPrice float64   `json:"admissionPrice"`
}

type Contact struct {
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Website     string `json:"website"`
}

type University struct {
	ID              primitive.ObjectID   `json:"univID,omitempty" bson:"_id,omitempty"`
	Name            string               `json:"univName" bson:"univName,omitempty" binding:"required" unique:"true" validate:"required"`
	Location        Location             `json:"univLocation" binding:"required" validate:"required"`
	Presentation    string               `json:"presentation"`
	IsPrivate       bool                 `json:"isPrivate" validate:"required"`
	Tuition         float64              `json:"tuition"`
	Contact         Contact              `json:"contact"`
	ImageURL        string               `json:"imageUrl"`
	DocumentURL     string               `json:"documentUrl"`
	ProgramIDs      []primitive.ObjectID `json:"programIDs,omitempty" bson:"programIDs,omitempty"`
	Infrastructure  []string             `json:"infrastructure"`
	Partnerships    []string             `json:"partnerships"`
	SuccessDiplomas float64              `json:"successDiplomas"`
	Events          []Event              `json:"events"`
	News            []string             `json:"news"`
	Photos          []string             `json:"Photos"`
	Ratings         []Rating             `json:"ratings"`
}

type Rating struct {
	UserID  primitive.ObjectID `json:"userID,omitempty" bson:"userID,omitempty"`
	Rating  int                `json:"rating"`
	Comment string             `json:"comment"`
}

type Program struct {
	ID              primitive.ObjectID `json:"programID,omitempty" bson:"_id,omitempty"`
	ProgramName     string             `json:"programName"`
	Level           string             `json:"level"`
	Duration        int                `json:"duration"`
	Requirements    []string           `json:"requirements"`
	CareerProspects []string           `json:"careerProspects"`
}
