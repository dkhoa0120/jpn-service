package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"jpn-service/common"
	"jpn-service/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateGrammar - Tạo grammar mới
func CreateGrammar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var grammar models.Grammar
	if err := json.NewDecoder(r.Body).Decode(&grammar); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Set timestamps
	grammar.CreatedAt = time.Now()
	grammar.UpdatedAt = time.Now()

	collection := common.GetDBCollection("grammars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, grammar)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Grammar pattern already exists",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error creating grammar: " + err.Error(),
		})
		return
	}

	grammar.ID = result.InsertedID.(primitive.ObjectID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Grammar created successfully",
		Data:    grammar,
	})
}

// GetGrammars - Lấy danh sách grammar với pagination và filter
func GetGrammars(w http.ResponseWriter, r *http.Request) {
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

	// Build filter
	filter := bson.M{}

	// Filter by level (N1, N2, N3, N4, N5)
	if level := r.URL.Query().Get("level"); level != "" {
		filter["level"] = level
	}

	// Search by pattern
	if pattern := r.URL.Query().Get("pattern"); pattern != "" {
		filter["pattern"] = bson.M{"$regex": pattern, "$options": "i"}
	}

	collection := common.GetDBCollection("grammars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find options with pagination and sorting
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(size)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error fetching grammars: " + err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var grammars []models.Grammar
	if err := cursor.All(ctx, &grammars); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error decoding grammars: " + err.Error(),
		})
		return
	}

	// Count total documents
	total, _ := collection.CountDocuments(ctx, filter)

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: bson.M{
			"items": grammars,
			"pagination": bson.M{
				"page":  page,
				"size":  size,
				"total": total,
			},
		},
	})
}

// GetGrammarByID - Lấy grammar theo ID
func GetGrammarByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid grammar ID",
		})
		return
	}

	collection := common.GetDBCollection("grammars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var grammar models.Grammar
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&grammar)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Grammar not found",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error fetching grammar: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    grammar,
	})
}

// UpdateGrammar - Cập nhật grammar
func UpdateGrammar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid grammar ID",
		})
		return
	}

	var grammar models.Grammar
	if err := json.NewDecoder(r.Body).Decode(&grammar); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	collection := common.GetDBCollection("grammars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update document
	update := bson.M{
		"$set": bson.M{
			"title":       grammar.Title,
			"description": grammar.Description,
			"level":       grammar.Level,
			"updated_at":  time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Grammar pattern already exists",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error updating grammar: " + err.Error(),
		})
		return
	}

	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Grammar not found",
		})
		return
	}

	grammar.ID = id
	grammar.UpdatedAt = time.Now()

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Grammar updated successfully",
		Data:    grammar,
	})
}

// DeleteGrammar - Xóa grammar
func DeleteGrammar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid grammar ID",
		})
		return
	}

	collection := common.GetDBCollection("grammars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error deleting grammar: " + err.Error(),
		})
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Grammar not found",
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Grammar deleted successfully",
	})
}

// GetGrammarsByLevel - Lấy grammar theo JLPT level
func GetGrammarsByLevel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	level := params["level"]

	collection := common.GetDBCollection("grammars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"level": level})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error fetching grammars: " + err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var grammars []models.Grammar
	if err := cursor.All(ctx, &grammars); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error decoding grammars: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    grammars,
	})
}

// SearchGrammars - Tìm kiếm grammar theo pattern
func SearchGrammars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Search query is required",
		})
		return
	}

	collection := common.GetDBCollection("grammars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Search in pattern, description, and structure
	filter := bson.M{
		"$or": []bson.M{
			{"pattern": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
			{"structure": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error searching grammars: " + err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var grammars []models.Grammar
	if err := cursor.All(ctx, &grammars); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error decoding grammars: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    grammars,
	})
}
