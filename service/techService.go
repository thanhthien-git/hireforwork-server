package service

import (
	"context"
	"hireforwork-server/models"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTech(page, pageSize int, techName string) (models.PaginateDocs[models.Tech], error) {
	var techList []models.Tech

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	bsonFilter := bson.D{{"isDeleted", false}}
	if techName != "" {
		bsonFilter = append(bsonFilter, bson.E{"technology", bson.D{{"$regex", techName}, {"$options", "i"}}})
	}

	skip := (page - 1) * pageSize
	findOptions := options.Find().
		SetProjection(bson.D{{"isDeleted", 0}}).
		SetLimit(int64(pageSize)).
		SetSkip(int64(skip)).
		SetSort(bson.D{{"technology", 1}})

	totalDocs, err := techCollection.CountDocuments(context.Background(), bsonFilter)
	if err != nil {
		return models.PaginateDocs[models.Tech]{}, err
	}
	totalPage := int(math.Ceil(float64(totalDocs) / float64(pageSize)))

	cursor, err := techCollection.Find(context.Background(), bsonFilter, findOptions)
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
