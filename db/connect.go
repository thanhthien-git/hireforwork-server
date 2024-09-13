package dbHelper

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getURI() string {
	godotenv.Load()
	URI := os.Getenv("DATABASE_CONNECTION")
	return URI
}

func ConnectDB() (*mongo.Client, context.Context, error) {
	clientOption := options.Client().ApplyURI(getURI())

	//time > 10s cancel connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Connect to MongoDB")
	return client, ctx, nil
}
