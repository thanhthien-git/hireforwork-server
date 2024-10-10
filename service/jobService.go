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
var CareerSaveJobCollection *mongo.Collection

func init() {
	client, ctx, err := dbHelper.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	JobCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_JOB"), client)
	CareerApplyJobCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CAREERAPPLYJOB"), client)
	CareerSaveJobCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CAREERSAVEJOB"), client)

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

func GetSavedJobsByCareerID(careerID string) ([]struct {
	JobID     string `json:"jobID"`
	IsDeleted bool   `json:"isDeleted"`
}, error) {
	CareerID, err := primitive.ObjectIDFromHex(careerID)
	if err != nil {
		log.Printf("Invalid career ID: %v", err)
		return nil, err
	}

	var careerSave models.CareerSaveJob
	err = CareerSaveJobCollection.FindOne(context.Background(), bson.M{"careerID": CareerID}).Decode(&careerSave)
	if err != nil {
		log.Printf("Error finding CareerSaveJob: %v", err)
		return nil, err
	}

	var savedJobsResponse []struct {
		JobID     string `json:"jobID"`
		IsDeleted bool   `json:"isDeleted"`
	}

	for _, job := range careerSave.SaveJob {
		if !job.IsDeleted {
			savedJobsResponse = append(savedJobsResponse, struct {
				JobID     string `json:"jobID"`
				IsDeleted bool   `json:"isDeleted"`
			}{
				JobID:     job.JobID.Hex(),
				IsDeleted: job.IsDeleted,
			})
		}
	}

	return savedJobsResponse, nil
}

func GetJobApplyHistoryByCareerID(careerID string) (models.CareerApplyJob, error) {
	CareerID, err := primitive.ObjectIDFromHex(careerID)
	var applyJobs models.CareerApplyJob
	filter := bson.M{"careerID": CareerID, "isDeleted": false}
	err = CareerApplyJobCollection.FindOne(context.Background(), filter).Decode(&applyJobs)
	if err != nil {
		log.Printf("Error for job history: %v", err)
		return models.CareerApplyJob{}, err
	}
	return applyJobs, nil
}
