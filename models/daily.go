package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VocabularyNote struct {
	Word       string `json:"word" bson:"word"`             // Từ vựng
	Reading    string `json:"reading" bson:"reading"`       // Cách đọc
	Meaning    string `json:"meaning" bson:"meaning"`       // Nghĩa
	Position   int    `json:"position" bson:"position"`     // Vị trí trong text
	Length     int    `json:"length" bson:"length"`         // Độ dài của từ
}

type Daily struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title          string             `json:"title" bson:"title" validate:"required"`
	Content        string             `json:"content" bson:"content" validate:"required"`
	FixContent     string             `json:"fix_content" bson:"fix_content"`
	Note           string             `json:"note" bson:"note"`
	VocabularyNotes []VocabularyNote  `json:"vocabulary_notes" bson:"vocabulary_notes"` // Danh sách từ vựng được note
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}