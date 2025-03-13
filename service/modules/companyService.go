package service

import (
	"context"
	"errors"
	"fmt"
	"hireforwork-server/db"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/utils"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CompanyService struct {
	companyCollection, jobCollection, careerApplyJob *mongo.Collection
}

func NewCompanyService(dbInstance *db.DB) *CompanyService {
	c := dbInstance.GetCollections([]string{"Company", "Job", "CareerApplyJob"})
	return &CompanyService{
		companyCollection: c[0],
		jobCollection:     c[1],
		careerApplyJob:    c[2],
	}
}

// Lấy danh sách company với phân trang
func (c *CompanyService) GetCompanies(page int, pageSize int, filter interfaces.ICompanyFilter) (models.PaginateDocs[models.Company], error) {

	var companies []models.Company

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize

	bsonFilter := bson.D{{"isDeleted", false}}

	if filter.CompanyName != "" {
		bsonFilter = append(bsonFilter, bson.E{"companyName", bson.D{{"$regex", filter.CompanyName}, {"$options", "i"}}})
	}

	if filter.CompanyEmail != "" {
		bsonFilter = append(bsonFilter, bson.E{"contact.companyEmail", bson.D{{"$regex", filter.CompanyEmail}, {"$options", "i"}}})
	}

	if filter.StartDate != nil || filter.EndDate != nil {
		dateFilter := bson.D{}

		if filter.StartDate != nil {
			start, _ := time.Parse("2006-01-02", *filter.StartDate)
			startOfDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
			dateFilter = append(dateFilter, bson.E{"$gte", startOfDay})
		}

		if filter.EndDate != nil {
			end, _ := time.Parse("2006-01-02", *filter.EndDate)
			endOfDay := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999, end.Location())
			dateFilter = append(dateFilter, bson.E{"$lte", endOfDay})
		} else {
			now := time.Now()
			endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999, now.Location())
			dateFilter = append(dateFilter, bson.E{"$lte", endOfDay})
		}

		bsonFilter = append(bsonFilter, bson.E{"createAt", dateFilter})
	}

	// Cấu hình phân trang
	findOptions := options.Find().SetSort(bson.D{{"companyName", 1}})
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))

	findOptions.SetProjection(bson.D{{"password", 0}})

	totalDocs, _ := c.companyCollection.CountDocuments(context.Background(), bsonFilter)
	totalPage := int64(math.Ceil(float64(totalDocs) / float64(pageSize)))
	cursor, err := c.companyCollection.Find(context.Background(), bsonFilter, findOptions)
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

func (c *CompanyService) GetCompanyByID(companyID string) (models.Company, error) {
	_id, _ := primitive.ObjectIDFromHex(companyID)
	var company models.Company

	findOptions := options.FindOne().SetProjection(bson.D{{"password", 0}})

	err := c.companyCollection.FindOne(context.Background(), bson.D{{"_id", _id}}, findOptions).Decode(&company)
	if err != nil {
		return models.Company{}, err
	}
	return company, nil
}

func (c *CompanyService) GetRandomCompany() (models.Company, error) {
	var company models.Company
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{{"isDeleted", false}}},
		},
		{
			{"$sample", bson.D{{"size", 1}}},
		},
	}

	cursor, err := c.companyCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return company, err
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&company); err != nil {
			return company, err
		}
	} else {
		return company, mongo.ErrNoDocuments
	}

	return company, nil
}

func (c *CompanyService) CreateCompany(company models.Company) (models.Company, error) {

	result, err := c.companyCollection.InsertOne(context.Background(), company)
	if err != nil {
		return models.Company{}, err
	}
	company.Id = result.InsertedID.(primitive.ObjectID)
	return company, nil
}

func (c *CompanyService) DeleteCompanyByID(companyID string) http.Response {
	_id, _ := primitive.ObjectIDFromHex(companyID)

	filter := bson.M{"_id": _id}

	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := c.companyCollection.FindOneAndUpdate(context.Background(), filter, update, opts)

	if result.Err() != nil {
		return http.Response{
			StatusCode: http.StatusBadRequest,
		}
	}
	return http.Response{
		StatusCode: http.StatusAccepted,
	}
}

func (c *CompanyService) UpdateCompanyByID(companyID string, updatedCompany models.Company) (models.Company, error) {
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
			"companyField": updatedCompany.CompanyField,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedDoc models.Company
	err = c.companyCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Company{}, errors.New("company not found or already deleted")
		}
		return models.Company{}, err
	}

	return updatedDoc, nil
}

func (c *CompanyService) GetCareersByJobID(jobID string, companyID string) ([]models.UserInfo, error) {
	jobObjectID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		log.Printf("Invalid job ID: %v", err)
	}

	companyObjectID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		log.Printf("Invalid company ID: %v", err)
	}

	var job models.Jobs
	err = c.jobCollection.FindOne(context.Background(), bson.M{"_id": jobObjectID, "isDeleted": false, "companyID": companyObjectID}).Decode(&job)
	if err != nil {
		log.Printf("Error finding job %v", err)
	}

	var applicants []models.UserInfo

	return applicants, nil
}

func (c *CompanyService) GetJobsByCompanyID(companyID string, page int, pageSize int, filter interfaces.IJobFilter) (models.PaginateDocs[models.Jobs], error) {

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize

	filterOption := bson.M{"isDeleted": false, "companyID": companyObjectID}

	if filter.JobTitle != "" {
		filterOption["jobTitle"] = bson.M{"$regex": filter.JobTitle, "$options": "i"}
	}

	if filter.DateCreateFrom != "" && filter.DateCreateTo != "" {
		createFromTime, _ := time.Parse("2006-01-02", filter.DateCreateFrom)
		createToTime, _ := time.Parse("2006-01-02", filter.DateCreateTo)

		createFromTime = time.Date(createFromTime.Year(), createFromTime.Month(), createFromTime.Day(), 0, 0, 0, 0, time.UTC)
		createToTime = time.Date(createToTime.Year(), createToTime.Month(), createToTime.Day(), 23, 59, 59, 0, time.UTC)

		filterOption["createAt"] = bson.M{
			"$gte": primitive.NewDateTimeFromTime(createFromTime),
			"$lte": primitive.NewDateTimeFromTime(createToTime),
		}
	}

	if filter.EndDateFrom != "" && filter.EndDateTo != "" {
		endFromTime, _ := time.Parse("2006-01-02", filter.EndDateFrom)
		endToTime, _ := time.Parse("2006-01-02", filter.EndDateTo)

		endFromTime = time.Date(endFromTime.Year(), endFromTime.Month(), endFromTime.Day(), 0, 0, 0, 0, time.UTC)
		endToTime = time.Date(endToTime.Year(), endToTime.Month(), endToTime.Day(), 23, 59, 59, 0, time.UTC)

		filterOption["expireDate"] = bson.M{
			"$gte": primitive.NewDateTimeFromTime(endFromTime),
			"$lte": primitive.NewDateTimeFromTime(endToTime),
		}
	}

	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.D{{"jobTitle", 1}})

	totalDocs, _ := c.jobCollection.CountDocuments(context.Background(), filterOption)

	totalPages := int64(math.Ceil(float64(totalDocs) / float64(pageSize)))

	var jobs []models.Jobs

	cursor, err := c.jobCollection.Find(context.Background(), filterOption, findOptions)

	if err != nil {
		return models.PaginateDocs[models.Jobs]{}, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &jobs); err != nil {
		return models.PaginateDocs[models.Jobs]{}, err
	}

	result := models.PaginateDocs[models.Jobs]{
		Docs:        jobs,
		TotalDocs:   totalDocs,
		CurrentPage: int64(page),
		TotalPage:   totalPages,
	}
	return result, nil
}

func (c *CompanyService) DeleteJobByID(jobID []string) error {
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

	_, err := c.jobCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return errors.New("Bạn không thuộc bộ phận này")
	}
	return nil
}

func (c *CompanyService) GetCareersApplyJob(companyID string, filter interfaces.IJobApplicationFilter) (map[string]interface{}, error) {
	id, _ := primitive.ObjectIDFromHex(companyID)

	skip := (filter.Page - 1) * filter.PageSize
	limit := filter.PageSize

	filterStage := bson.M{}

	matchStage := bson.M{
		"companyID": id,
		"isDeleted": false,
	}
	//filter by date
	if filter.CreateFrom != "" && filter.CreateTo != "" {
		createFromTime, _ := time.Parse("2006-01-02", filter.CreateFrom)
		createToTime, _ := time.Parse("2006-01-02", filter.CreateTo)

		createFromTime = time.Date(createFromTime.Year(), createFromTime.Month(), createFromTime.Day(), 0, 0, 0, 0, time.UTC)
		createToTime = time.Date(createToTime.Year(), createToTime.Month(), createToTime.Day(), 23, 59, 59, 0, time.UTC)

		matchStage["createAt"] = bson.M{
			"$gte": primitive.NewDateTimeFromTime(createFromTime),
			"$lte": primitive.NewDateTimeFromTime(createToTime),
		}
	}
	//filter by status
	if filter.Status != "" {
		matchStage["status"] = filter.Status
	}
	//filter by mail
	if filter.CareerEmail != " " {
		filterStage["careerDetail.careerEmail"] = bson.M{
			"$regex":   filter.CareerEmail,
			"$options": "i",
		}
	}
	//filter by level
	if filter.JobLevel != "" {
		filterStage["jobDetail.jobLevel"] = filter.JobLevel
	}
	//filter by jobTitle
	if filter.JobTitle != " " {
		filterStage["jobDetail.jobTitle"] = bson.M{
			"$regex":   filter.JobTitle,
			"$options": "i",
		}
	}

	pipeline := mongo.Pipeline{
		{{"$match", matchStage}},
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
		///filter stage here
		{{"$match", filterStage}},
		{{"$project", bson.D{
			{"_id", 1},
			{"jobID", 1},
			{"careerID", 1},
			{"careerEmail", bson.D{{"$arrayElemAt", bson.A{"$careerDetail.careerEmail", 0}}}},
			{"status", 1},
			{"createAt", 1},
			{"careerCV", 1},
			{"isChange", 1},
			{"jobTitle", bson.D{{"$arrayElemAt", bson.A{"$jobDetail.jobTitle", 0}}}},
			{"jobRequirement", bson.D{{"$arrayElemAt", bson.A{"$jobDetail.jobRequirement", 0}}}},
			{"jobLevel", bson.D{{"$arrayElemAt", bson.A{"$jobDetail.jobLevel", 0}}}},
		}}},
		{{"$skip", skip}},
		{{"$limit", limit}},
	}
	cursor, err := c.careerApplyJob.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var result []bson.M
	if err := cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"docs": result,
		"page": filter.Page,
	}
	return response, nil
}

func (c *CompanyService) GetStatics(id primitive.ObjectID) (bson.M, error) {
	result := bson.M{}

	filter := bson.M{"isDeleted": false, "companyID": id}
	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$group", bson.D{
			{"_id", "$careerID"},
		}}},
	}

	cursor, _ := c.careerApplyJob.Aggregate(context.Background(), pipeline)
	defer cursor.Close(context.Background())

	totalCareer := 0
	for cursor.Next(context.Background()) {
		totalCareer++
	}

	totalResume, _ := c.careerApplyJob.CountDocuments(context.Background(), filter)
	totalJob, _ := c.jobCollection.CountDocuments(context.Background(), filter)

	result["totalCareer"] = totalCareer
	result["totalResume"] = totalResume
	result["totalJob"] = totalJob

	return result, nil
}

func (c *CompanyService) UploadCompanyCover(link string, id string) error {
	obj, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": obj}
	update := bson.M{
		"$set": bson.M{
			"companyImage.coverURL": link,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := c.companyCollection.FindOneAndUpdate(context.Background(), filter, update, opt)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (c *CompanyService) UploadCompanyImage(link string, id string) error {
	obj, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": obj}
	update := bson.M{
		"$set": bson.M{
			"companyImage.imageURL": link,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := c.companyCollection.FindOneAndUpdate(context.Background(), filter, update, opt)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (c *CompanyService) ChangeResumeStatus(resumeID string, status string) error {
	_id, _ := primitive.ObjectIDFromHex(resumeID)

	update := bson.M{
		"$set": bson.M{
			"status":   status,
			"isChange": true,
		},
	}

	result, err := c.careerApplyJob.UpdateOne(
		context.Background(),
		bson.M{
			"_id":      _id,
			"isChange": bson.M{"$ne": true},
		},
		update,
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("Không thể thay đổi trạng thái hoặc trạng thái đã thay đổi")
	}

	return nil
}
func (c *CompanyService) RequestPasswordResetCompany(email string) (string, error) {
	var company models.Company

	err := c.companyCollection.FindOne(context.Background(), bson.M{"contact.companyEmail": email}).Decode(&company)
	if err != nil {
		return "", fmt.Errorf("Không tìm thấy công ty với email %s: %v", email, err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := rng.Intn(9000) + 1000
	verificationCode := fmt.Sprintf("%d", code)

	subject := "Mã xác nhận khôi phục mật khẩu"
	body := fmt.Sprintf("Mã xác nhận của bạn là: %s", verificationCode)
	if err := SendEmail(email, subject, body); err != nil {
		return "", err
	}

	_, err = c.companyCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": company.Id},
		bson.M{"$set": bson.M{"verificationCode": verificationCode}},
	)
	if err != nil {
		return "", fmt.Errorf("Lỗi khi cập nhật mã xác nhận: %v", err)
	}

	return verificationCode, nil
}

func (c *CompanyService) ResetPasswordCompany(email string, code string, newPassword string) error {
	var company models.Company

	err := c.companyCollection.FindOne(context.Background(), bson.M{"contact.companyEmail": email, "verificationCode": code}).Decode(&company)
	if err != nil {
		return fmt.Errorf("Sai mã xác nhận: %v", err)
	}

	hashedPassword := utils.EncodeToSHA(newPassword)

	_, err = c.companyCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": company.Id},
		bson.M{"$set": bson.M{"password": hashedPassword, "verificationCode": ""}},
	)
	if err != nil {
		return fmt.Errorf("Lỗi khi khôi phục mật khẩu: %v", err)
	}

	return nil
}
