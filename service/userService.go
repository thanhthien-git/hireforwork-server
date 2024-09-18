package service

import (
	"context"
	dbHelper "hireforwork-server/db"
	"hireforwork-server/models"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

func init() {
	client, ctx, err := dbHelper.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	collection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CAREER"), client)
}

func GetUser() ([]models.User, error) {
	var users []models.User

	// Log collection name and context for debugging
	log.Printf("Using collection: %s", os.Getenv("COLLECTION_CAREER"))

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Printf("Error finding documents: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &users); err != nil {
		log.Printf("Error parsing documents: %v", err)
		return nil, err
	}

	// Log number of users fetched
	log.Printf("Number of users found: %d", len(users))

	return users, nil
}
