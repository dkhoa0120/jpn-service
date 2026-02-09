package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"jpn-service/common"
	"jpn-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllUsers - Lấy tất cả users
func GetAllScores(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := common.GetDBCollection("scores")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error fetching scores: " + err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var scores []models.Score
	if err = cursor.All(ctx, &scores); err != nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error decoding scores: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    scores,
	})
}


// CreateUser - Tạo user mới
func CreateScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var score models.Score
	if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Set timestamps
	score.CreatedAt = time.Now()
	score.UpdatedAt = time.Now()

	collection := common.GetDBCollection("scores")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, score)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error creating score: " + err.Error(),
		})
		return
	}

	score.ID = result.InsertedID.(primitive.ObjectID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Score created successfully",
		Data:    score,
	})
}


