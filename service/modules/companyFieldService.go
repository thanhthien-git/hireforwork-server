package service

import (
	"context"
	"hireforwork-server/db"
	"hireforwork-server/models"
	"log"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FieldService struct {
	collection *mongo.Collection
}

func NewFieldService(dbInstance *db.DB) *FieldService {
	return &FieldService{
		collection: db.GetInstance().GetCollection("CompanyField"),
	}
}

func (f *FieldService) GetField(page, pageSize int) (models.PaginateDocs[models.Field], error) {
	var field []models.Field

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	bsonFilter := bson.D{{"isDeleted", false}}

	skip := (page - 1) * pageSize
	findOption := options.Find().SetLimit(int64(pageSize)).SetSkip(int64(skip))
	totalDocs, _ := f.collection.CountDocuments(context.Background(), bsonFilter)
	totalPage := int64(math.Ceil(float64(totalDocs) / float64(pageSize)))

	cursor, err := f.collection.Find(context.Background(), bsonFilter, findOption)
	if err != nil {
		log.Printf("Có lỗi xảy ra")
		return models.PaginateDocs[models.Field]{}, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &field); err != nil {
		log.Printf("Có lỗi xảy ra")
		return models.PaginateDocs[models.Field]{}, err
	}

	result := models.PaginateDocs[models.Field]{
		Docs:        field,
		TotalDocs:   totalDocs,
		CurrentPage: int64(page),
		TotalPage:   totalPage,
	}

	return result, nil
}
