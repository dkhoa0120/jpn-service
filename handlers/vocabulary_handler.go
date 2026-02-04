package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"jpn-service/common"
	"jpn-service/models"

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

	collection := common.GetDBCollection("vocabularies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(size)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error fetching vocabularies: " + err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var vocabularies []models.Vocabulary
	if err := cursor.All(ctx, &vocabularies); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error decoding vocabularies: " + err.Error(),
		})
		return
	}

	total, _ := collection.CountDocuments(ctx, filter)

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
