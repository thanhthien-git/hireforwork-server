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

// Lấy danh sách company với phân trang
func GetCompanies(page int, pageSize int, companyName, companyEmail string) (models.PaginateDocs[models.Company], error) {

	var companies []models.Company

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize

	bsonFilter := bson.D{{"isDeleted", false}}

	if companyName != "" {
		bsonFilter = append(bsonFilter, bson.E{"companyName", bson.D{{"$regex", companyName}, {"$options", "i"}}})
	}

	if companyEmail != "" {
		bsonFilter = append(bsonFilter, bson.E{"contact.companyEmail", bson.D{{"$regex", companyEmail}, {"$options", "i"}}})
	}

	// Cấu hình phân trang
	findOptions := options.Find().SetSort(bson.D{{"companyName", 1}})
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))

	findOptions.SetProjection(bson.D{{"password", 0}})

	totalDocs, _ := companyCollection.CountDocuments(context.Background(), bsonFilter)
	totalPage := int64(math.Ceil(float64(totalDocs) / float64(pageSize)))
	cursor, err := companyCollection.Find(context.Background(), bsonFilter, findOptions)
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
			"companyName":  updatedCompany.CompanyName,
			"contact":      updatedCompany.Contact,
			"employeeSize": updatedCompany.EmployeeSize,
			"description":  updatedCompany.Description,
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
	err = jobCollection.FindOne(context.Background(), bson.M{"_id": jobObjectID, "isDeleted": false, "companyID": companyObjectID}).Decode(&job)
	if err != nil {
		log.Printf("Error finding job %v", err)
	}

	var applicants []models.UserInfo

	return applicants, nil
}

func GetJobsByCompanyID(companyID string, page int64, limit int64) (models.PaginateDocs[models.Jobs], error) {
	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	findOptions := options.Find().
		SetSort(bson.D{{"jobTitle", 1}}).
		SetLimit(limit).
		SetSkip(skip)

	totalDocs, _ := jobCollection.CountDocuments(context.Background(), bson.M{"isDeleted": false, "companyID": companyObjectID})
	totalPage := int64(math.Ceil(float64(totalDocs) / float64(limit)))
	var jobs []models.Jobs
	cursor, err := jobCollection.Find(context.Background(), bson.M{"isDeleted": false, "companyID": companyObjectID}, findOptions)
	if err != nil {
		log.Printf("Error finding jobs for company: %v", err)
	}

	if err = cursor.All(context.Background(), &jobs); err != nil {
		log.Printf("Error decoding jobs: %v", err)
	}

	result := models.PaginateDocs[models.Jobs]{
		Docs:        jobs,
		TotalDocs:   totalDocs,
		CurrentPage: int64(page),
		TotalPage:   totalPage,
	}

	return result, nil
}

func DeleteJobByID(jobID []string) error {
	ids := make([]primitive.ObjectID, len(jobID))

	for index, element := range jobID {
		objID, _ := primitive.ObjectIDFromHex(element)
		ids[index] = objID
	}

	filter := bson.M{
		"_id":       bson.M{"$in": ids},
		"isDeleted": false,
	}

	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
		},
	}

	_, err := jobCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return errors.New("Bạn không thuộc bộ phận này")
	}

	return nil
}

func GetCareersApplyJob(companyID string) ([]bson.M, error) {
	id, _ := primitive.ObjectIDFromHex(companyID)

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"companyID", id},
			{"isDeleted", false},
		}}},
		{{"$lookup", bson.D{
			{"from", "Job"},
			{"localField", "jobID"},
			{"foreignField", "_id"},
			{"as", "jobDetail"},
		}}},
		{{"$lookup", bson.D{
			{"from", "Career"},
			{"localField", "careerID"},
			{"foreignField", "_id"},
			{"as", "careerDetail"},
		}}},
		{{"$match", bson.D{
			{"jobDetail", bson.D{{"$ne", bson.A{}}}},
			{"jobDetail.isDeleted", false},
		}}},
		{{"$project", bson.D{
			{"_id", 1},
			{"jobID", 1},
			{"careerID", 1},
			{"careerEmail", bson.D{{"$arrayElemAt", bson.A{"$careerDetail.careerEmail", 0}}}},
			{"status", 1},
			{"createAt", 1},
			{"jobTitle", bson.D{{"$arrayElemAt", bson.A{"$jobDetail.jobTitle", 0}}}},
			{"jobRequirement", bson.D{{"$arrayElemAt", bson.A{"$jobDetail.jobRequirement", 0}}}},
			{"jobLevel", bson.D{{"$arrayElemAt", bson.A{"$jobDetail.jobLevel", 0}}}},
		}}},
	}
	cursor, err := careerApplyJob.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, errors.New("Chúng tôi đã cố gắng hết sức")
	}
	defer cursor.Close(context.TODO())
	var result []bson.M
	if err := cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetStatics(id primitive.ObjectID) (bson.M, error) {
	result := bson.M{}

	filter := bson.M{"isDeleted": false, "companyID": id}
	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$group", bson.D{
			{"_id", "$careerID"},
		}}},
	}

	cursor, _ := careerApplyJob.Aggregate(context.Background(), pipeline)
	defer cursor.Close(context.Background())

	totalCareer := 0
	for cursor.Next(context.Background()) {
		totalCareer++
	}

	totalResume, _ := careerApplyJob.CountDocuments(context.Background(), filter)
	totalJob, _ := jobCollection.CountDocuments(context.Background(), filter)

	result["totalCareer"] = totalCareer
	result["totalResume"] = totalResume
	result["totalJob"] = totalJob

	return result, nil
}

func UploadCompanyCover(link string, id string) error {
	obj, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": obj}
	update := bson.M{
		"$set": bson.M{
			"companyImage.coverURL": link,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := companyCollection.FindOneAndUpdate(context.Background(), filter, update, opt)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func UploadCompanyImage(link string, id string) error {
	obj, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": obj}
	update := bson.M{
		"$set": bson.M{
			"companyImage.imageURL": link,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := companyCollection.FindOneAndUpdate(context.Background(), filter, update, opt)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
