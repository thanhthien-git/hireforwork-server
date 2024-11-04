package service

import (
	dbHelper "hireforwork-server/db"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

var collection, jobCollection, companyCollection, careerSaveJob, careerViewedJob, careerApplyJob, techCollection, companyFieldCollection, categoryCollection *mongo.Collection

func init() {
	client, ctx, err := dbHelper.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	collection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CAREER"), client)

	jobCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_JOB"), client)

	companyCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_COMPANY"), client)

	careerSaveJob = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CAREER_SAVEDJOB"), client)

	careerViewedJob = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CAREER_VIEWEDJOB"), client)

	careerApplyJob = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CAREER_APPLYJOB"), client)

	categoryCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CATEGORY"), client)

	techCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_TECHNOLOGIES"), client)
	
	categoryCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CATEGORY"), client)

	companyFieldCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_COMPANY_FIELD"), client)
}
