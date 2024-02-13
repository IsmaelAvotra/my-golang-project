package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" binding:"required" unique:"true"`
	Email     string             `json:"email" binding:"required,email" unique:"true"`
	Password  string             `json:"password,omitempty" binding:"required"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
