package models

import (
	"jpn-service/enum"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vocabulary struct {
	ID        primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	NameVI    string              `json:"name_vi" bson:"name_vi" validate:"required"`
	NameJPN   string              `json:"name_jpn" bson:"name_jpn" validate:"required"`
	Phonetic  string              `json:"phonetic" bson:"phonetic" validate:"required"`
	Category  enum.VocabularyType `json:"category" bson:"category" validate:"required"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
}
