package service

import (
	"context"
	"errors"
	"hireforwork-server/models"
	"math"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTech(page, pageSize int, TechName string) (models.PaginateDocs[models.Tech], error) {
	var techList []models.Tech

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	bsonFilter := bson.D{{"isDeleted", false}}
	skip := (page - 1) * pageSize

	if TechName != "" {
		bsonFilter = append(bsonFilter, bson.E{"technology", bson.D{{"$regex", TechName}, {"$options", "i"}}})
	}

	findOption := options.Find().SetProjection(bson.D{{"isDeleted", 0}})
	findOption.SetLimit(int64(pageSize))
	findOption.SetSkip(int64(skip))
	findOption.SetSort(bson.D{{"technology", 1}})

	totalDocs, _ := techCollection.CountDocuments(context.Background(), bsonFilter)
	totalPage := int(math.Ceil(float64(totalDocs) / float64(pageSize)))
	cursor, err := techCollection.Find(context.Background(), bsonFilter, findOption)
	if err != nil {
		return models.PaginateDocs[models.Tech]{}, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &techList); err != nil {
		return models.PaginateDocs[models.Tech]{}, err
	}

	result := models.PaginateDocs[models.Tech]{
		Docs:        techList,
		TotalDocs:   totalDocs,
		CurrentPage: int64(page),
		TotalPage:   int64(totalPage),
	}
	return result, nil
}

func CreateTech(tech models.Tech) (models.Tech, error) {

	result, err := techCollection.InsertOne(context.Background(), tech)
	if err != nil {
		return models.Tech{}, err
	}
	tech.Id = result.InsertedID.(primitive.ObjectID)
	return tech, nil
}

func DeleteTechByID(techID string) http.Response {
	_id, _ := primitive.ObjectIDFromHex(techID)

	filter := bson.M{"_id": _id}

	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := techCollection.FindOneAndUpdate(context.Background(), filter, update, opts)

	if result.Err() != nil {
		return http.Response{
			StatusCode: http.StatusBadRequest,
		}
	}
	return http.Response{
		StatusCode: http.StatusAccepted,
	}
}

func UpdateTechByID(techID string, updatedTech models.Tech) (models.Tech, error) {
	_id, err := primitive.ObjectIDFromHex(techID)
	if err != nil {
		return models.Tech{}, errors.New("invalid tech ID format")
	}

	filter := bson.M{"_id": _id, "isDeleted": false}

	update := bson.M{
		"$set": bson.M{
			"technology": updatedTech.TechName,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedDoc models.Tech
	err = techCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Tech{}, errors.New("tech not found or already deleted")
		}
		return models.Tech{}, err
	}

	return updatedDoc, nil
}
