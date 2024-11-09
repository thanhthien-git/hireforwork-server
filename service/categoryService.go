package service

import (
	"context"
	"errors"
	"fmt"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"log"
	"math"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCategory(page, pageSize int, filter interfaces.ICategoryFilter) (models.PaginateDocs[models.Category], error) {
	var categoryList []models.Category

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	bsonFilter := bson.D{{"isDeleted", false}}
	skip := (page - 1) * pageSize

	if filter.Category != "" {
		bsonFilter = append(bsonFilter, bson.E{"categoryName", bson.D{{"$regex", filter.Category}, {"$options", "i"}}})
	}

	findOption := options.Find()
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

func CreateCategory(category models.Category) (models.Category, error) {
	filter := bson.M{"categoryName": category.CategoryName}
	err := categoryCollection.FindOne(context.Background(), filter).Err()
	if err == nil {
		return models.Category{}, fmt.Errorf("danh mục với tên '%s' đã tồn tại", category.CategoryName)
	} else if err != mongo.ErrNoDocuments {
		return models.Category{}, fmt.Errorf("đã có lỗi khi kiểm tra danh mục: %v", err)
	}
	category.Id = primitive.NewObjectID()
	category.IsDeleted = false

	_, err = categoryCollection.InsertOne(context.Background(), category)
	if err != nil {
		log.Printf("Lỗi khi chèn danh mục: %v", err)
		return models.Category{}, fmt.Errorf("đã có lỗi xảy ra khi tạo danh mục: %v", err)
	}
	return category, nil
}

func UpdateCategoryByID(categoryID string, updatedCategory models.Category) (models.Category, error) {
	_id, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return models.Category{}, errors.New("id danh mục không đúng định dạng")
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
			return models.Category{}, errors.New("danh mục không tìm thấy hoặc đã xóa")
		}
		return models.Category{}, err
	}

	return updatedDoc, nil
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

func GetCategoryByID(categoryID string) (models.Category, error) {
	_id, _ := primitive.ObjectIDFromHex(categoryID)
	var category models.Category

	findOptions := options.FindOne()

	err := categoryCollection.FindOne(context.Background(), bson.D{{"_id", _id}}, findOptions).Decode(&category)
	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}
