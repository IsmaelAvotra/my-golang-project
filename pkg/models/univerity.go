package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Adress        string
	CoordinateGPS string
	Province      string
	Region        string
	City          string
}

type Program struct {
	Name            string
	Level           string
	Duration        int
	Requirements    []string
	CareerProspects []string
}

type Event struct {
	Title          string
	Descrioption   string
	Date           time.Time
	Location       string
	IsFree         bool
	AdmissionPrice float64
}

type Contact struct {
	PhoneNumber string
	Email       string
	Website     string
}

type University struct {
	ID              primitive.ObjectID `json:"univID,omitempty" bson:"univID,omitempty"`
	Name            string             `json:"univName" binding:"required" unique:"true" validate:"required"`
	Location        Location           `json:"univLocation" binding:"required" validate:"required"`
	Presentation    string             `json:"presentation"`
	IsPrivate       bool               `json:"isPrivate" validate:"required"`
	Tuition         float64            `json:"tuition"`
	Contact         Contact            `json:"contact"`
	ImageURL        string             `json:"imageUrl"`
	DocumentURL     string             `json:"documentUrl"`
	Programs        []Program          `json:"programs"`
	Infrastructure  []string           `json:"infrastructure"`
	Partnerships    []string           `json:"partnerships"`
	SuccessDiplomas float64            `json:"successDiplomas"`
	Events          []Event            `json:"events"`
	News            []string           `json:"news"`
	Photos          []string           `json:"Photos"`
}
