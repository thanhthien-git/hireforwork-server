package service

import (
	"context"
	"fmt"
	"hireforwork-server/models"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTech(page, pageSize int) (models.PaginateDocs[models.Tech], error) {
	var techList []models.Tech

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	bsonFilter := bson.D{{"isDeleted", false}}
	skip := (page - 1) * pageSize

	findOption := options.Find().SetProjection(bson.D{{"isDeleted", 0}})
	findOption.SetLimit(int64(pageSize))
	findOption.SetSkip(int64(skip))
	findOption.SetSort(bson.D{{"technology", 1}})

	totalDocs, _ := techCollection.CountDocuments(context.Background(), bsonFilter)
	totalPage := int(math.Ceil(float64(totalDocs) / float64(pageSize)))
	cursor, err := techCollection.Find(context.Background(), bsonFilter, findOption)
	if err != nil {
		fmt.Println("Lỗi khi tìm kiếm công nghệ: %v", err)
		return models.PaginateDocs[models.Tech]{}, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &techList); err != nil {
		fmt.Println("Lỗi khi tìm kiếm công nghệ: %v", err)
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
