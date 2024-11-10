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

func GetTech(page, pageSize int, filter interfaces.ITechnologyFilter) (models.PaginateDocs[models.Tech], error) {

	var techList []models.Tech

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	bsonFilter := bson.D{{"isDeleted", false}}
	skip := (page - 1) * pageSize

	if filter.Technology != "" {
		bsonFilter = append(bsonFilter, bson.E{"technology", bson.D{{"$regex", filter.Technology}, {"$options", "i"}}})
	}

	findOption := options.Find()
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
func CreateTech(technology models.Tech) (models.Tech, error) {
	filter := bson.M{"technology": technology.TechName}
	err := techCollection.FindOne(context.Background(), filter).Err()
	if err == nil {
		return models.Tech{}, fmt.Errorf("kĩ năng với tên '%s' đã tồn tại", technology.TechName)
	} else if err != mongo.ErrNoDocuments {
		return models.Tech{}, fmt.Errorf("đã có lỗi khi kiểm tra kĩ năng: %v", err)
	}
	technology.Id = primitive.NewObjectID()
	technology.IsDeleted = false

	_, err = techCollection.InsertOne(context.Background(), technology)
	if err != nil {
		log.Printf("Lỗi khi chèn công nghệ: %v", err)
		return models.Tech{}, fmt.Errorf("đã có lỗi xảy ra khi tạo kĩ năng: %v", err)
	}
	return technology, nil
}

func UpdateTechByID(techID string, updatedTech models.Tech) (models.Tech, error) {
	_id, err := primitive.ObjectIDFromHex(techID)
	if err != nil {
		return models.Tech{}, errors.New("id kỹ năng không đúng định dạng")
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
			return models.Tech{}, errors.New("kỹ năng không tìm thấy hoặc đã xóa")
		}
		return models.Tech{}, err
	}

	return updatedDoc, nil
}
func DeleteTechByID(technologyID string) http.Response {
	_id, _ := primitive.ObjectIDFromHex(technologyID)
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

func GetTechByID(technologyID string) (models.Tech, error) {
	_id, _ := primitive.ObjectIDFromHex(technologyID)
	var technology models.Tech

	findOptions := options.FindOne()

	err := techCollection.FindOne(context.Background(), bson.D{{"_id", _id}}, findOptions).Decode(&technology)
	if err != nil {
		return models.Tech{}, err
	}
	return technology, nil
}
