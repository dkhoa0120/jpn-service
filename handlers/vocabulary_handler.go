package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"jpn-service/common"
	"jpn-service/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateVocabulary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var vocabulary models.Vocabulary
	if err := json.NewDecoder(r.Body).Decode(&vocabulary); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Set timestamps
	vocabulary.CreatedAt = time.Now()
	vocabulary.UpdatedAt = time.Now()

	collection := common.GetDBCollection("vocabularies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, vocabulary)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			w.WriteHeader(http.StatusConflict) // 409
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Already Exist",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error creating vocabulary: " + err.Error(),
		})
		return
	}

	vocabulary.ID = result.InsertedID.(primitive.ObjectID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Vocabulary created successfully",
		Data:    vocabulary,
	})
}

func GetVocabularies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query params
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	skip := (page - 1) * size

	// Filter
	filter := bson.M{}
	if category := r.URL.Query().Get("category"); category != "" {
		filter["category"] = category
	}

	if topic := r.URL.Query().Get("topic"); topic != "" {
		filter["topic"] = topic
	}

	collection := common.GetDBCollection("vocabularies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ✅ Sử dụng goroutine để chạy song song
	var vocabularies []models.Vocabulary
	var total int64
	var wg sync.WaitGroup
	var findErr, countErr error

	wg.Add(2)

	// Query 1: Get documents
	go func() {
		defer wg.Done()
		opts := options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(size)).
			SetSort(bson.M{"created_at": -1})

		cursor, err := collection.Find(ctx, filter, opts)
		if err != nil {
			findErr = err
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &vocabularies); err != nil {
			findErr = err
		}
	}()

	// Query 2: Count documents (parallel)
	go func() {
		defer wg.Done()
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			countErr = err
			return
		}
		total = count
	}()

	wg.Wait()

	// Check for errors
	if findErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error fetching vocabularies: " + findErr.Error(),
		})
		return
	}

	if countErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error counting vocabularies: " + countErr.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: bson.M{
			"items": vocabularies,
			"pagination": bson.M{
				"page":  page,
				"size":  size,
				"total": total,
			},
		},
	})
}

// UpdateVocabulary - Cập nhật vocabulary theo ID
func UpdateVocabulary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Lấy ID từ URL params
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid vocabulary ID",
		})
		return
	}

	// Parse request body
	var vocabulary models.Vocabulary
	if err := json.NewDecoder(r.Body).Decode(&vocabulary); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	collection := common.GetDBCollection("vocabularies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update document
	update := bson.M{
		"$set": bson.M{
			"name_vi":    vocabulary.NameVI,
			"name_jpn":   vocabulary.NameJPN,
			"phonetic":   vocabulary.Phonetic,
			"category":   vocabulary.Category,
			"topic":      vocabulary.Topic,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Vocabulary with same data already exists",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error updating vocabulary: " + err.Error(),
		})
		return
	}

	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Vocabulary not found",
		})
		return
	}

	// Return updated data
	vocabulary.ID = id
	vocabulary.UpdatedAt = time.Now()

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Vocabulary updated successfully",
		Data:    vocabulary,
	})
}
