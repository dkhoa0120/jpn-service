package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Score struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title          string             `json:"title" bson:"title" validate:"required"`
	WrongWord	   []string           `json:"wrong_word" bson:"wrong_word" validate:"required"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}