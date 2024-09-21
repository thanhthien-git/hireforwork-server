package service

import (
	"context"
	dbHelper "hireforwork-server/db"
	"hireforwork-server/models"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	client, ctx, err := dbHelper.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	collection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CAREER"), client)
}

func GetUser(page, pageSize int) (models.PaginateDocs[models.User], error) {
	var users []models.User

	skip := (page - 1) * pageSize

	findOption := options.Find()
	findOption.SetLimit(int64(pageSize))
	findOption.SetSkip(int64(skip))

	totalDocs, _ := collection.CountDocuments(context.Background(), bson.D{})
	cursor, err := collection.Find(context.Background(), bson.D{}, findOption)
	if err != nil {
		log.Printf("Error finding documents: %v", err)
		return models.PaginateDocs[models.User]{}, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &users); err != nil {
		log.Printf("Error parsing documents: %v", err)
		return models.PaginateDocs[models.User]{}, err
	}

	result := models.PaginateDocs[models.User]{
		Docs:        users,
		TotalDocs:   totalDocs,
		CurrentPage: int64(page),
	}

	return result, nil
}

func GetUserByID(careerID string) (models.User, error) {
	_id, _ := primitive.ObjectIDFromHex(careerID)
	var user models.User

	err := collection.FindOne(context.Background(), bson.D{{"_id", _id}}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
