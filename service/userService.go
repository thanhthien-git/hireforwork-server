package service

import (
	"context"
	dbHelper "hireforwork-server/db"
	"hireforwork-server/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

func init() {
	client, ctx, err := dbHelper.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	// collection = client.Database(os.Getenv("DATABASE_NAME")).Collection(os.Getenv("COLLECTION_CAREER"))
	collection = client.Database("hideforwork").Collection("Career")
}

func GetUser(ctx context.Context, collection *mongo.Collection) ([]models.User, error) {
	var users []models.User
	//cancel get list after 15s
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//find doc
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
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
