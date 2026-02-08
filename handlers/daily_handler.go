package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"jpn-service/common"
	"jpn-service/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllDailies - Lấy tất cả dailies
func GetAllDailies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := common.GetDBCollection("dailies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error fetching dailies: " + err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var dailies []models.Daily
	if err = cursor.All(ctx, &dailies); err != nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error decoding dailies: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    dailies,
	})
}

// GetDaily - Lấy daily theo ID
func GetDaily(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid daily ID",
		})
		return
	}

	collection := common.GetDBCollection("dailies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var daily models.Daily
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&daily)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Daily not found",
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    daily,
	})
}

// CreateDaily - Tạo daily mới
func CreateDaily(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var daily models.Daily
	if err := json.NewDecoder(r.Body).Decode(&daily); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Set timestamps
	daily.CreatedAt = time.Now()
	daily.UpdatedAt = time.Now()

	collection := common.GetDBCollection("dailies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, daily)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error creating daily: " + err.Error(),
		})
		return
	}

	daily.ID = result.InsertedID.(primitive.ObjectID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Daily created successfully",
		Data:    daily,
	})
}

// UpdateDaily - Cập nhật daily
func UpdateDaily(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid daily ID",
		})
		return
	}

	var daily models.Daily
	if err := json.NewDecoder(r.Body).Decode(&daily); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	collection := common.GetDBCollection("dailies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":            daily.Title,
			"content":          daily.Content,
			"fix_content":      daily.FixContent,
			"note":             daily.Note,
			"vocabulary_notes": daily.VocabularyNotes,
			"updated_at":       time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error updating daily: " + err.Error(),
		})
		return
	}

	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Daily not found",
		})
		return
	}

	daily.ID = id
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Daily updated successfully",
		Data:    daily,
	})
}

// DeleteDaily - Xóa daily
func DeleteDaily(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid daily ID",
		})
		return
	}

	collection := common.GetDBCollection("dailies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error deleting daily: " + err.Error(),
		})
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Daily not found",
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Daily deleted successfully",
	})
}
