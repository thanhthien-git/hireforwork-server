package service

import (
	"context"
	"hireforwork-server/models"
	"log"
	"math"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// Lấy danh sách company với phân trang
func GetCompanies(page int, pageSize int) (models.PaginateDocs[models.Company], error) {
	var companies []models.Company

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize

	// Cấu hình phân trang
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))

	// Thực hiện truy vấn với phân trang
	totalDocs, _ := companyCollection.CountDocuments(context.Background(), bson.D{})
	totalPage := int64(math.Ceil(float64(totalDocs) / float64(pageSize)))
	cursor, err := companyCollection.Find(context.Background(), bson.D{{"isDeleted", false}}, findOptions)
	log.Print(cursor)
	if err != nil {
		log.Printf("Error finding documents: %v", err)
		return models.PaginateDocs[models.Company]{}, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &companies); err != nil {
		log.Printf("Error parsing documents: %v", err)
		return models.PaginateDocs[models.Company]{}, err
	}

	log.Print(companies)
	result := models.PaginateDocs[models.Company]{
		Docs:        companies,
		TotalDocs:   totalDocs,
		CurrentPage: int64(page),
		TotalPage:   totalPage,
	}

	return result, nil
}
