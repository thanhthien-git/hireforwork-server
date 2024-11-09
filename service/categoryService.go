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

func GetCategory(page, pageSize int, CategoryName string) (models.PaginateDocs[models.Category], error) {
	var categoryList []models.Category

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	skip := (page - 1) * pageSize

	bsonFilter := bson.D{{"isDeleted", false}}

	if CategoryName != "" {
		bsonFilter = append(bsonFilter, bson.E{"categoryName", bson.D{{"$regex", CategoryName}, {"$options", "i"}}})
	}

	findOption := options.Find().SetProjection(bson.D{{"isDeleted", 0}})
	findOption.SetLimit(int64(pageSize))
	findOption.SetSkip(int64(skip))
	findOption.SetSort(bson.D{{"categoryName", 1}})

	totalDocs, _ := categoryCollection.CountDocuments(context.Background(), bsonFilter)
	totalPage := int(math.Ceil(float64(totalDocs) / float64(pageSize)))
	cursor, err := categoryCollection.Find(context.Background(), bsonFilter, findOption)
	if err != nil {
		return models.PaginateDocs[models.Category]{}, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &categoryList); err != nil {
		return models.PaginateDocs[models.Category]{}, err
	}

	result := models.PaginateDocs[models.Category]{
		Docs:        categoryList,
		TotalDocs:   totalDocs,
		CurrentPage: int64(page),
		TotalPage:   int64(totalPage),
	}
	return result, nil
}

func CreateCategory(category models.Category) (models.Category, error) {

	result, err := categoryCollection.InsertOne(context.Background(), category)
	if err != nil {
		return models.Category{}, err
	}
	category.Id = result.InsertedID.(primitive.ObjectID)
	return category, nil
}

func DeleteCategoryByID(categoryID string) http.Response {
	_id, _ := primitive.ObjectIDFromHex(categoryID)

	filter := bson.M{"_id": _id}

	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := categoryCollection.FindOneAndUpdate(context.Background(), filter, update, opts)

	if result.Err() != nil {
		return http.Response{
			StatusCode: http.StatusBadRequest,
		}
	}
	return http.Response{
		StatusCode: http.StatusAccepted,
	}
}

func UpdateCategoryByID(categoryID string, updatedCategory models.Category) (models.Category, error) {
	_id, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return models.Category{}, errors.New("invalid category ID format")
	}

	filter := bson.M{"_id": _id, "isDeleted": false}

	update := bson.M{
		"$set": bson.M{
			"categoryName": updatedCategory.CategoryName,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedDoc models.Category
	err = categoryCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Category{}, errors.New("category not found or already deleted")
		}
		return models.Category{}, err
	}

	return updatedDoc, nil
}
