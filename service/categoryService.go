package service

import (
	"context"
	"fmt"
	"hireforwork-server/models"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCategory(page, pageSize int) (models.PaginateDocs[models.Category], error) {
	var categoryList []models.Category

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
	findOption.SetSort(bson.D{{"categoryName", 1}})

	totalDocs, _ := categoryCollection.CountDocuments(context.Background(), bsonFilter)
	totalPage := int(math.Ceil(float64(totalDocs) / float64(pageSize)))
	cursor, err := categoryCollection.Find(context.Background(), bsonFilter, findOption)
	if err != nil {
		fmt.Println("Lỗi khi tìm kiếm công nghệ: %v", err)
		return models.PaginateDocs[models.Category]{}, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &categoryList); err != nil {
		fmt.Println("Lỗi khi tìm kiếm công nghệ: %v", err)
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
