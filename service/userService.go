package service

import (
	"context"
	dbHelper "hireforwork-server/db"
	"hireforwork-server/models"
	"hireforwork-server/utils"
	"log"
	"math"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var jobCollection *mongo.Collection
var companyCollection *mongo.Collection

func init() {
	client, ctx, err := dbHelper.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	collection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CAREER"), client)

	jobCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_JOB"), client)

	companyCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_COMPANY"), client)

}

func GetUser(page, pageSize int) (models.PaginateDocs[models.User], error) {
	var users []models.User

	skip := (page - 1) * pageSize

	findOption := options.Find().SetProjection(bson.D{{"password", 0}})
	findOption.SetLimit(int64(pageSize))
	findOption.SetSkip(int64(skip))
	///100 documents -> 5

	totalDocs, _ := collection.CountDocuments(context.Background(), bson.D{})
	totalPage := int64(math.Ceil(float64(totalDocs) / float64(pageSize)))
	cursor, err := collection.Find(context.Background(), bson.D{{"isDeleted", false}}, findOption)
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
		TotalPage:   totalPage,
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

func DeleteUserByID(careerID string) http.Response {
	_id, _ := primitive.ObjectIDFromHex(careerID)

	filter := bson.M{"_id": _id}

	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(context.Background(), filter, update, opts)

	if result.Err() != nil {
		return http.Response{
			StatusCode: http.StatusBadRequest,
		}
	}
	return http.Response{
		StatusCode: http.StatusAccepted,
	}
}

func CreateUser(user models.User) (models.User, error) {
	user.Password = utils.EncodeToSHA(user.Password)

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return models.User{}, err
	}
	user.Id = result.InsertedID.(primitive.ObjectID)
	return user, nil
}
