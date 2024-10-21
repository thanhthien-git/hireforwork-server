package service

import (
	"context"
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

func GetJob(page, pageSize int) (models.PaginateDocs[models.Jobs], error) {
	var jobs []models.Jobs

	skip := (page - 1) * pageSize

	findOption := options.Find()
	findOption.SetLimit(int64(pageSize))
	findOption.SetSkip(int64(skip))

	totalDocs, _ := jobCollection.CountDocuments(context.Background(), bson.D{})
	totalPage := int64(math.Ceil(float64(totalDocs) / float64(pageSize)))
	cursor, err := jobCollection.Find(context.Background(), bson.D{{"isDeleted", false}}, findOption)
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

func CreateJob(job models.Jobs) (models.Jobs, error) {
	currentTime := time.Now()
	job.Id = primitive.NewObjectID()
	job.CreateAt = primitive.NewDateTimeFromTime(currentTime)
	job.ExpireDate = primitive.NewDateTimeFromTime(currentTime.AddDate(0, 0, 14))
	job.IsClosed = false
	job.IsHot = false
	result, err := jobCollection.InsertOne(context.Background(), job)
	fmt.Println(err)
	if err != nil {
		return models.Jobs{}, fmt.Errorf("Đã có lỗi xảy ra khi tạo bài đăng")
	}
	job.Id = result.InsertedID.(primitive.ObjectID)
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

func (s *JobService) GetFilteredJobs(ctx context.Context, createDate, expireDate time.Time) ([]models.Jobs, error) {
	filter := bson.M{
		"createAt":   bson.M{"$gte": createDate},
		"expireDate": bson.M{"$lte": expireDate},
	}
	opts := options.Find().SetSort(bson.D{{"createAt", -1}})

	cursor, err := s.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []models.Jobs
	if err = cursor.All(ctx, &jobs); err != nil {
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
