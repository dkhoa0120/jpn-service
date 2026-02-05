package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"jpn-service/common"
	"jpn-service/models"
)

func main() {
	// Load environment
	if err := common.LoadEnv(); err != nil {
		log.Fatal("Error loading .env:", err)
	}

	// Connect to MongoDB
	if err := common.InitDB(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer common.CloseDB()

	log.Println("✅ Connected to MongoDB")

	// Read JSON file
	jsonFile, err := os.Open("vocabularies_seed.json")
	if err != nil {
		log.Fatal("Error opening JSON file:", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Error reading JSON file:", err)
	}

	// Parse JSON
	var vocabularies []models.Vocabulary
	if err := json.Unmarshal(byteValue, &vocabularies); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	log.Printf("📖 Found %d vocabularies in JSON file\n", len(vocabularies))

	// Add timestamps
	now := time.Now()
	for i := range vocabularies {
		vocabularies[i].CreatedAt = now
		vocabularies[i].UpdatedAt = now
	}

	// Insert to MongoDB
	collection := common.GetDBCollection("vocabularies")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert to []interface{}
	docs := make([]interface{}, len(vocabularies))
	for i, v := range vocabularies {
		docs[i] = v
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Fatal("Error inserting vocabularies:", err)
	}

	fmt.Printf("🎉 Successfully imported %d vocabularies!\n", len(result.InsertedIDs))

	// Stats by topic
	fmt.Println("\n📊 Statistics by topic:")
	topics := map[string]int{}
	for _, v := range vocabularies {
		topics[string(v.Topic)]++
	}

	for topic, count := range topics {
		fmt.Printf("  - %s: %d từ\n", topic, count)
	}
}
