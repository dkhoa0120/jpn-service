package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Daily struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	Title       string    `json:"title" bson:"title" validate:"required"`
	Description string    `json:"description" bson:"description" validate:"required"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}
