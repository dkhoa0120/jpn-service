package models

import (
	"jpn-service/enum"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrammarExample struct {
	Japanese   string `json:"japanese" bson:"japanese"`
	Vietnamese string `json:"vietnamese" bson:"vietnamese"`
}

type Grammar struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	Title       string         `json:"title" bson:"title" validate:"required"`
	Description string         `json:"description" bson:"description" validate:"required"`
	Level       enum.JLPTLevel `json:"level" bson:"level" validate:"required"`
	CreatedAt   time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" bson:"updated_at"`
}
