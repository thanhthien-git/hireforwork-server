package service

import (
	"context"
	"errors"
	"fmt"
	"hireforwork-server/constants"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetJob(page, pageSize int, jobTitle string, jobCategory string, companyName string) (models.PaginateDocs[models.Jobs], error) {

	var jobs []models.Jobs
	skip := (page - 1) * pageSize

	findOption := options.Find()
	findOption.SetLimit(int64(pageSize))
	findOption.SetSkip(int64(skip))

	filter := bson.M{
		"isDeleted": false,
	}
	if jobTitle != "" {
		filter["jobTitle"] = bson.M{
			"$regex": primitive.Regex{Pattern: jobTitle, Options: "i"},
		}
	}
	// if workingLocation != "" {
	// 	filter["workingLocation"] = bson.M{
	// 		"$elemMatch": bson.M{
	// 			"$regex": primitive.Regex{Pattern: workingLocation, Options: "i"},
	// 		},
	// 	}
	// }

	if jobCategory != "" {
		filter["jobCategory"] = bson.M{
			"$elemMatch": bson.M{
				"$regex": primitive.Regex{Pattern: jobCategory, Options: "i"},
			},
		}
	}
	if companyName != "" {
		companyID, err := GetCompanyIDByName(companyName)
		if err != nil {
			return models.PaginateDocs[models.Jobs]{}, err
		}

		filter["companyID"] = companyID
	}

	totalDocs, _ := jobCollection.CountDocuments(context.Background(), filter)
	totalPage := int64(math.Ceil(float64(totalDocs) / float64(pageSize)))
	cursor, err := jobCollection.Find(context.Background(), filter, findOption)
	if err != nil {
		log.Printf("Error finding documents: %v", err)
		return models.PaginateDocs[models.Jobs]{}, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &jobs); err != nil {
		log.Printf("Error parsing job documents: %v", err)
		return models.PaginateDocs[models.Jobs]{}, err
	}

	result := models.PaginateDocs[models.Jobs]{
		Docs:        jobs,
		TotalDocs:   totalDocs,
		CurrentPage: int64(page),
		TotalPage:   totalPage,
	}

	return result, nil
}
func GetCompanyIDByName(companyName string) (primitive.ObjectID, error) {
	if companyCollection == nil {
		return primitive.NilObjectID, errors.New("CompanyCollection is not initialized")
	}

	var company struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	filter := bson.M{
		"companyName": bson.M{
			"$regex": primitive.Regex{Pattern: companyName, Options: "i"},
		},
	}

	err := companyCollection.FindOne(context.Background(), filter).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return primitive.NilObjectID, errors.New("company not found")
		}
		return primitive.NilObjectID, err
	}

	return company.ID, nil
}

func CreateJob(job models.Jobs) (models.Jobs, error) {
	currentTime := time.Now()
	job.ID = primitive.NewObjectID()
	job.CreateAt = primitive.NewDateTimeFromTime(currentTime)
	job.ExpireDate = primitive.NewDateTimeFromTime(currentTime.AddDate(0, 0, 14))
	job.IsClosed = false
	job.IsHot = false
	result, err := jobCollection.InsertOne(context.Background(), job)
	fmt.Println(err)
	if err != nil {
		return models.Jobs{}, fmt.Errorf("Đã có lỗi xảy ra khi tạo bài đăng")
	}
	job.ID = result.InsertedID.(primitive.ObjectID)
	return job, nil
}

func CheckIsExistedJob(request interfaces.IUserJob, collection *mongo.Collection) bool {
	filter := bson.D{
		{"isDeleted", false},
		{"careerID", request.IDCareer},
		{"jobID", request.JobID},
	}
	var result bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Printf("Error occurred while checking job existence: %v", err)
		return false
	}
	return true
}

func ApplyForJob(request interfaces.IJobApply) (models.Jobs, error) {
	newApply := models.CareerApplyJob{
		ID:        primitive.NewObjectID(),
		CareerID:  request.IDCareer,
		JobID:     request.JobID,
		CreateAt:  primitive.NewDateTimeFromTime(time.Now()),
		IsDeleted: false,
		Status:    constants.PENDING,
	}

	_, err := careerApplyJob.InsertOne(context.Background(), newApply)
	if err != nil {
		log.Printf("Error inserting apply data into CareerApplyJob: %v", err)
		return models.Jobs{}, err
	}

	filter := bson.M{"_id": request.JobID, "isDeleted": false}
	update := bson.M{"$push": bson.M{"userApply": request.IDCareer}}

	var job models.Jobs
	err = jobCollection.FindOneAndUpdate(context.Background(), filter, update).Decode(&job)
	if err != nil {
		log.Printf("Error updating job with user info: %v", err)
		return models.Jobs{}, err
	}

	return job, nil
}

func GetSavedJobsByCareerID(careerID string) ([]models.SavedJob, error) {

	CareerID, err := primitive.ObjectIDFromHex(careerID)
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"careerID", CareerID},
			}},
		},
		{
			{"$project", bson.D{
				{"saveJob", bson.D{
					{"$filter", bson.D{
						{"input", "$saveJob"},
						{"as", "job"},
						{"cond", bson.D{
							{"$eq", bson.A{"$$job.isDeleted", false}},
						}},
					}},
				}},
			}},
		},
	}

	result, err := careerSaveJob.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Printf("Error during aggregation : %v", err)
	}
	defer result.Close(context.Background())

	var results []models.CareerSaveJob
	if err := result.All(context.Background(), &results); err != nil {
		log.Printf("Error decoding results: %v", err)
		return nil, err
	}
	return results[0].SaveJob, nil
}
func GetJobApplyHistoryByCareerID(careerID string) (models.CareerApplyJob, error) {
	CareerID, err := primitive.ObjectIDFromHex(careerID)
	var applyJobs models.CareerApplyJob
	filter := bson.M{"careerID": CareerID, "isDeleted": false}
	err = careerApplyJob.FindOne(context.Background(), filter).Decode(&applyJobs)
	if err != nil {
		log.Printf("Error for job history: %v", err)
		return models.CareerApplyJob{}, err
	}
	return applyJobs, nil
}

type JobService struct {
	Collection *mongo.Collection
}

func NewJobService(db *mongo.Database) *JobService {
	return &JobService{
		Collection: db.Collection("Job"),
	}
}

func (s *JobService) GetLatestJobs() ([]models.Jobs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"isDeleted": false, "isClosed": false}
	opts := options.Find().SetSort(bson.D{{"createAt", -1}}).SetLimit(10)

	cursor, err := s.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []models.Jobs
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func GetJobByID(jobID string) (models.Jobs, error) {
	_id, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return models.Jobs{}, err
	}

	var job models.Jobs
	projection := bson.D{{"userApply", 0}}
	findOptions := options.FindOne().SetProjection(projection)

	err = jobCollection.FindOne(context.Background(), bson.D{{"_id", _id}}, findOptions).Decode(&job)
	if err != nil {
		return models.Jobs{}, err
	}
	return job, nil
}
