package service

import (
	"context"
	"errors"
	"hireforwork-server/models"
	"log"
	"math"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCompanies(page int, pageSize int) (models.PaginateDocs[models.Company], error) {
	var companies []models.Company

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize

	findOptions := options.Find().SetProjection(bson.D{{"password", 0}})
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))

	totalDocs, _ := companyCollection.CountDocuments(context.Background(), bson.D{{"isDeleted", false}})
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

	result := models.PaginateDocs[models.Company]{
		Docs:        companies,
		TotalDocs:   totalDocs,
		CurrentPage: int64(page),
		TotalPage:   totalPage,
	}

	return result, nil
}

func GetCompanyByID(companyID string) (models.Company, error) {
	_id, _ := primitive.ObjectIDFromHex(companyID)
	var company models.Company

	findOptions := options.FindOne().SetProjection(bson.D{{"password", 0}})

	err := companyCollection.FindOne(context.Background(), bson.D{{"_id", _id}}, findOptions).Decode(&company)
	if err != nil {
		return models.Company{}, err
	}
	return company, nil
}

func CreateCompany(company models.Company) (models.Company, error) {

	result, err := companyCollection.InsertOne(context.Background(), company)
	if err != nil {
		return models.Company{}, err
	}
	company.Id = result.InsertedID.(primitive.ObjectID)
	return company, nil
}

func DeleteCompanyByID(companyID string) http.Response {
	_id, _ := primitive.ObjectIDFromHex(companyID)

	filter := bson.M{"_id": _id}

	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := companyCollection.FindOneAndUpdate(context.Background(), filter, update, opts)

	if result.Err() != nil {
		return http.Response{
			StatusCode: http.StatusBadRequest,
		}
	}
	return http.Response{
		StatusCode: http.StatusAccepted,
	}
}

func UpdateCompanyByID(companyID string, updatedCompany models.Company) (models.Company, error) {
	_id, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return models.Company{}, errors.New("invalid company ID format")
	}

	filter := bson.M{"_id": _id, "isDeleted": false}

	update := bson.M{
		"$set": bson.M{
			"companyName":   updatedCompany.CompanyName,
			"companyImage":  updatedCompany.CompanyImage,
			"contact":       updatedCompany.Contact,
			"description":   updatedCompany.Description,
			"typeOfCompany": updatedCompany.TypeOfCompany,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedDoc models.Company
	err = companyCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Company{}, errors.New("company not found or already deleted")
		}
		return models.Company{}, err
	}

	return updatedDoc, nil
}

func GetCareersByJobID(jobID string, companyID string) ([]models.UserInfo, error) {
	jobObjectID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		log.Printf("Invalid job ID: %v", err)
	}

	companyObjectID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		log.Printf("Invalid company ID: %v", err)
	}

	var job models.Jobs
	err = JobCollection.FindOne(context.Background(), bson.M{"_id": jobObjectID, "isDeleted": false, "companyID": companyObjectID}).Decode(&job)
	if err != nil {
		log.Printf("Error finding job %v", err)
	}

	var applicants []models.UserInfo
	for _, application := range job.UserApply {
		var user models.UserInfo
		user.UserId = application.UserId
		user.IsAccepted = application.IsAccepted
		user.CreateAt = application.CreateAt

		applicants = append(applicants, user)
	}

	return applicants, nil
}

func GetJobsByCompanyID(companyID string) ([]models.Jobs, error) {
	companyObjectID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		log.Printf("Invalid company ID: %v", err)
	}

	var jobs []models.Jobs
	cursor, err := JobCollection.Find(context.Background(), bson.M{"isDeleted": false, "companyID": companyObjectID})
	if err != nil {
		log.Printf("Error finding jobs for company: %v", err)
	}

	if err = cursor.All(context.Background(), &jobs); err != nil {
		log.Printf("Error decoding jobs: %v", err)
	}

	return jobs, nil
}
