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

	err := companyCollection.FindOne(context.Background(), bson.D{{"_id", _id}}).Decode(&company)
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

func UpdateCompanyByID(companyID string, updatedData models.Company) (models.Company, error) {
	_id, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return models.Company{}, err
	}

	filter := bson.M{"_id": _id, "isDeleted": false}

	updateFields := bson.M{}

	if updatedData.CompanyImage.ImageURL != "" {
		updateFields["companyImage.imageURL"] = updatedData.CompanyImage.ImageURL
	}
	if updatedData.CompanyImage.CoverURL != "" {
		updateFields["companyImage.coverURL"] = updatedData.CompanyImage.CoverURL
	}

	if updatedData.Contact.CompanyPhone != "" {
		updateFields["contact.companyPhone"] = updatedData.Contact.CompanyPhone
	}
	if updatedData.Contact.CompanyEmail != "" {
		updateFields["contact.companyEmail"] = updatedData.Contact.CompanyEmail
	}
	if updatedData.Contact.CompanyWebsite != "" {
		updateFields["contact.companyWebsite"] = updatedData.Contact.CompanyWebsite
	}
	if updatedData.Contact.CompanyAddress != "" {
		updateFields["contact.companyAddress"] = updatedData.Contact.CompanyAddress
	}

	if updatedData.CompanyName != "" {
		updateFields["companyName"] = updatedData.CompanyName
	}
	if updatedData.CompanyViewed != 0 {
		updateFields["companyViewed"] = updatedData.CompanyViewed
	}
	if updatedData.Description != "" {
		updateFields["description"] = updatedData.Description
	}
	if updatedData.EmployeeSize != 0 {
		updateFields["employeeSize"] = updatedData.EmployeeSize
	}
	if len(updatedData.FieldOperation) > 0 {
		updateFields["fieldOperation"] = updatedData.FieldOperation
	}
	if updatedData.Popularity != 0 {
		updateFields["popularity"] = updatedData.Popularity
	}
	if len(updatedData.PostJob) > 0 {
		updateFields["postJob"] = updatedData.PostJob
	}
	if len(updatedData.TypeOfCompany) > 0 {
		updateFields["typeOfCompany"] = updatedData.TypeOfCompany
	}

	if len(updateFields) == 0 {
		return models.Company{}, errors.New("no fields to update")
	}

	update := bson.M{
		"$set": updateFields,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedCompany models.Company
	err = companyCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedCompany)
	if err != nil {
		return models.Company{}, err
	}

	return updatedCompany, nil
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
