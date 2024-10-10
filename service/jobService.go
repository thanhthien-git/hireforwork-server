package service

import (
	"context"
	"errors"
	dbHelper "hireforwork-server/db"
	"hireforwork-server/models"
	"log"
	"math"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var JobCollection *mongo.Collection
var CareerApplyJobCollection *mongo.Collection

func init() {
	client, ctx, err := dbHelper.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	JobCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_JOB"), client)
	CareerApplyJobCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CAREERAPPLYJOB"), client)

}
func GetJob(page, pageSize int) (models.PaginateDocs[models.Jobs], error) {
	var jobs []models.Jobs

	skip := (page - 1) * pageSize

	findOption := options.Find()
	findOption.SetLimit(int64(pageSize))
	findOption.SetSkip(int64(skip))
	///100 documents -> 5

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

func ApplyForJob(jobID string, userInfo models.UserInfo) (models.Jobs, error) {
	objectID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		log.Printf("Invalid job ID: %v", err)
		return models.Jobs{}, errors.New("invalid job ID")
	}

	var job models.Jobs
	filter := bson.M{"_id": objectID, "isDeleted": false}
	update := bson.M{"$push": bson.M{"userApply": userInfo}}

	err = JobCollection.FindOneAndUpdate(context.Background(), filter, update).Decode(&job)
	if err != nil {
		log.Printf("Error updating job with user info: %v", err)
		return models.Jobs{}, err
	}

	newApply := models.CareerApplyJob{
		ID:        primitive.NewObjectID(),
		CareerID:  userInfo.UserId,
		JobID:     objectID,
		CreateAt:  primitive.NewDateTimeFromTime(time.Now()),
		IsDeleted: false,
		Status:    "applied",
	}

	_, err = CareerApplyJobCollection.InsertOne(context.Background(), newApply)
	if err != nil {
		log.Printf("Error inserting apply data into CareerApplyJob: %v", err)
		return models.Jobs{}, err
	}

	// Return updated job
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

	// Lọc và sắp xếp công việc
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
