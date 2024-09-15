package service

import (
	"context"
	"fmt"
	dbHelper "hireforwork-server/db"
	"hireforwork-server/models"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	client, ctx, err := dbHelper.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	collection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CAREER"), client)
}

func GetUser(ctx context.Context) ([]models.User, error) {
	var users []models.User
	//cancel get list after 15s
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//find doc
	cursor, err := collection.Find(ctx, bson.D{{"isDeleted", false}})
	if err != nil {
		return []models.User{}, fmt.Errorf(err.Error())
	}
	defer cursor.Close(ctx)

	//decode to user structor
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(ctx context.Context, id primitive.ObjectID) (models.User, error) {
	var user models.User
	filter := bson.D{
		{"_id", id},
		{"isDeleted", false},
	}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return models.User{}, fmt.Errorf(err.Error())
	}
	return user, nil
}

func DeleteUserByID(ctx context.Context, id primitive.ObjectID) (http.Response, error) {
	filter := bson.D{{"_id", id}}
	query := bson.D{{"$set", bson.D{{"isDeleted", true}}}}

	result := collection.FindOneAndUpdate(ctx, filter, query, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		return http.Response{}, result.Err()
	}

	return http.Response{StatusCode: http.StatusOK}, nil
}
