package jobs

import (
	"context"
	"fmt"
	"hireforwork-server/db"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"log"
	"math"
	"time"

	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JobRepository struct {
	jobCollection        *mongo.Collection
	careerSaveCollection *mongo.Collection
	cache                *cache.Cache
}

func NewJobRepository(dbInstance *db.DB) *JobRepository {
	jobCollection := dbInstance.GetCollection("Job")
	careerSaveCollection := dbInstance.GetCollection("CareerSaveJob")
	// Tạo cache với defaultExpiration là 5 phút và cleanupInterval là 10 phút
	jobCache := cache.New(5*time.Minute, 10*time.Minute)

	return &JobRepository{
		jobCollection:        jobCollection,
		careerSaveCollection: careerSaveCollection,
		cache:                jobCache,
	}
}

func (j *JobRepository) GetJob(page, pageSize int, filter interfaces.IJobFilter) (bson.M, error) {
	skip := (page - 1) * pageSize
	matchStage := bson.M{"isDeleted": false, "isClosed": false}

	facetStage := bson.D{
		{"$facet", bson.D{
			{"totalCount", []bson.D{{{"$count", "count"}}}},
			{"data", []bson.D{
				{{"$skip", int64(skip)}},
				{{"$limit", int64(pageSize)}},
				{{"$addFields", bson.D{
					{"companyName", "$companyDetails.companyName"},
					{"companyImage", "$companyDetails.companyImage"},
				}}},
			}},
		}},
	}

	matchOption := bson.M{}

	if filter.IsExpire {
		currentDate := time.Now()
		matchOption["expireDate"] = bson.M{"$gt": currentDate}
	}

	if filter.Query != "" {
		matchStage["$or"] = bson.A{
			bson.M{"jobTitle": bson.M{"$regex": filter.Query, "$options": "i"}},
			bson.M{"jobRequirement": bson.M{"$elemMatch": bson.M{"$regex": filter.Query, "$options": "i"}}},
		}
	}

	if filter.JobTitle != "" {
		matchStage["jobTitle"] = bson.M{"$regex": filter.JobTitle, "$options": "i"}
	}
	//filter by create date
	if filter.DateCreateFrom != "" && filter.DateCreateTo != "" {
		matchOption["createAt"] = bson.M{
			"$gte": filter.DateCreateFrom,
			"$lte": filter.DateCreateTo,
		}
	}
	//filter by expire date
	if filter.EndDateFrom != "" && filter.EndDateTo != "" {
		matchOption["createAt"] = bson.M{
			"$gte": filter.EndDateFrom,
			"$lte": filter.EndDateTo,
		}
	}
	//filter by salary
	if filter.SalaryFrom != 0 && filter.SalaryTo != 0 {
		matchOption["$and"] = []bson.M{
			{"jobSalaryMin": bson.M{"$gte": filter.SalaryFrom}},
			{"jobSalaryMax": bson.M{"$lte": filter.SalaryTo}},
		}
	}
	//filter by category
	if len(filter.JobCategory) > 0 {
		matchOption["jobCategory"] = bson.M{"$all": filter.JobCategory}
	}
	//filter by working location
	if len(filter.WorkingLocation) > 0 {
		matchOption["workingLocation"] = bson.M{"$in": filter.WorkingLocation}
	}

	//filter by job require
	if len(filter.JobRequirement) > 0 {
		matchOption["jobRequirement"] = bson.M{"$in": filter.JobRequirement}
	}
	//filter by hot
	if filter.IsHot {
		matchOption["isHot"] = filter.IsHot
	}
	//filter by job level
	if filter.JobLevel != "" {
		matchOption["jobLevel"] = filter.JobLevel
	}

	projectStage := bson.D{
		{"$project", bson.D{
			{"totalCount", 1},
			{"data", bson.D{
				{"$map", bson.D{
					{"input", "$data"},
					{"as", "doc"},
					{"in", bson.D{
						{"_id", "$$doc._id"},
						{"companyID", "$$doc.companyID"},
						{"companyName", "$$doc.companyName"},
						{"companyImage", "$$doc.companyImage"},
						{"createAt", "$$doc.createAt"},
						{"expireDate", "$$doc.expireDate"},
						{"isHot", "$$doc.isHot"},
						{"jobCategory", "$$doc.jobCategory"},
						{"jobDescription", "$$doc.jobDescription"},
						{"jobLevel", "$$doc.jobLevel"},
						{"jobRequirement", "$$doc.jobRequirement"},
						{"jobSalaryMax", "$$doc.jobSalaryMax"},
						{"jobSalaryMin", "$$doc.jobSalaryMin"},
						{"jobTitle", "$$doc.jobTitle"},
						{"quantity", "$$doc.quantity"},
						{"workingLocation", "$$doc.workingLocation"},
					}},
				}},
			}},
		}},
	}

	//default pipeline
	pipeline := mongo.Pipeline{
		{{"$match", matchStage}},
		{{"$lookup", bson.D{
			{"from", "Company"},
			{"localField", "companyID"},
			{"foreignField", "_id"},
			{"as", "companyDetails"},
		}}},
		{{"$unwind", bson.D{
			{"path", "$companyDetails"},
			{"preserveNullAndEmptyArrays", true},
		}}},
		{{"$match", bson.D{
			{"$or", bson.A{
				bson.D{{"companyDetails.isDeleted", false}},
				bson.D{{"companyDetails", bson.D{{"$exists", false}}}},
			}},
		}}},
		{{"$match", matchOption}},
	}

	if filter.CompanyName != "" {
		pipeline = append(pipeline, bson.D{{
			"$match", bson.D{{
				"companyDetails.companyName", bson.D{
					{"$regex", filter.CompanyName},
					{"$options", "i"},
				},
			}},
		}})
	}

	pipeline = append(pipeline, facetStage, projectStage)

	var result []bson.M
	cursor, err := j.jobCollection.Aggregate(context.Background(), pipeline)

	if err != nil {
		log.Printf("Error finding documents: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}

	totalDocs := int64(0)
	if len(result[0]["totalCount"].(bson.A)) > 0 {
		countVal := result[0]["totalCount"].(bson.A)[0].(bson.M)["count"].(int32)
		totalDocs = int64(countVal) // Convert int32 to int64
	}

	jobs := result[0]["data"].(bson.A)
	totalPage := int64(math.Ceil(float64(totalDocs) / float64(pageSize)))

	return bson.M{
		"docs":        jobs,
		"totalDocs":   totalDocs,
		"currentPage": page,
		"totalPage":   totalPage,
	}, nil
}

func (j *JobRepository) CreateJob(job models.Jobs) (models.Jobs, error) {
	currentTime := time.Now()
	job.Id = primitive.NewObjectID()
	job.CreateAt = primitive.NewDateTimeFromTime(currentTime)
	job.IsClosed = false
	result, err := j.jobCollection.InsertOne(context.Background(), job)
	fmt.Println(err)
	if err != nil {
		return models.Jobs{}, fmt.Errorf("Đã có lỗi xảy ra khi tạo bài đăng")
	}
	job.Id = result.InsertedID.(primitive.ObjectID)
	return job, nil
}

func (j *JobRepository) UpdateJob(job models.Jobs) (models.Jobs, error) {
	filter := bson.M{"_id": job.Id}

	update := bson.M{
		"$set": bson.M{
			"jobTitle":         job.JobTitle,
			"jobSalaryMin":     job.JobSalaryMin,
			"jobSalaryMax":     job.JobSalaryMax,
			"jobRequirement":   job.JobRequirement,
			"workingLocation":  job.WorkingLocation,
			"isHot":            job.IsHot,
			"isClosed":         job.IsClosed,
			"isDeleted":        job.IsDeleted,
			"expireDate":       job.ExpireDate,
			"jobCategory":      job.JobCategory,
			"jobDescription":   job.JobDescription,
			"jobLevel":         job.JobLevel,
			"recruitmentCount": job.RecruitmentCount,
			"workingType":      job.WorkType,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := j.jobCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&job)
	if err != nil {
		fmt.Println(err)
		return models.Jobs{}, fmt.Errorf("Có lỗi xảy ra khi cập nhập lại thông tin")
	}
	return job, nil
}

func (j *JobRepository) GetLatestJobs() ([]models.Jobs, error) {
	var jobs []models.Jobs

	filter := bson.M{
		"isDeleted": false,
		"expireDate": bson.M{
			"$gt": time.Now(),
		},
	}
	opts := options.Find().SetSort(bson.D{{"createAt", -1}}).SetLimit(10)

	cursor, err := j.jobCollection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (j *JobRepository) GetJobByID(jobID string, userId string) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check cache
	cacheKey := fmt.Sprintf("job:%s:user:%s", jobID, userId)
	if cached, found := j.cache.Get(cacheKey); found {
		fmt.Println("cache hit")
		return cached.(bson.M), nil
	}

	_id, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, fmt.Errorf("invalid job ID format: %v", err)
	}

	pipeline := j.buildJobPipeline(_id, userId)

	opts := options.Aggregate().SetBatchSize(100)
	cursor, err := j.jobCollection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, fmt.Errorf("aggregate error: %v", err)
	}
	defer cursor.Close(ctx)

	var job bson.M
	if !cursor.Next(ctx) {
		return nil, mongo.ErrNoDocuments
	}

	if err := cursor.Decode(&job); err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}

	result := bson.M{"doc": job}

	// Store in cache
	j.cache.Set(cacheKey, result, 5*time.Minute)

	return result, nil
}

func (j *JobRepository) buildJobPipeline(jobID primitive.ObjectID, userId string) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"_id", jobID}}}},
		{{"$lookup", bson.D{
			{"from", "Company"},
			{"localField", "companyID"},
			{"foreignField", "_id"},
			{"as", "company"},
		}}},
		{{"$unwind", bson.D{{"path", "$company"}, {"preserveNullAndEmptyArrays", true}}}},
		{{"$addFields", bson.D{
			{"contact", "$company.contact"},
			{"companyName", "$company.companyName"},
			{"companyImage", "$company.companyImage"},
			{"employeeSize", "$company.employeeSize"},
			{"_id", 1},
			{"jobTitle", 1},
			{"jobSalaryMin", "$company.jobSalaryMin"},
			{"jobSalaryMax", "$company.jobSalaryMax"},
			{"jobRequirement", "$company.jobRequirement"},
			{"jobCategory", "$company.jobCategory"},
			{"jobDescription", "$company.jobDescription"},
			{"jobLevel", "$company.jobLevel"},
			{"quantity", "$company.quantity"},
			{"workingLocation", "$company.workingLocation"},
		}}},
		{{"$unset", bson.A{"company"}}},
	}

	if userId != "" {
		decodeId, _ := primitive.ObjectIDFromHex(userId)
		lookUpApplyStage := bson.D{
			{"$lookup", bson.D{
				{"from", "CareerApplyJob"},
				{"let", bson.D{
					{"jobID", "$_id"},
					{"careerID", decodeId},
				}},
				{"pipeline", mongo.Pipeline{
					{{"$match", bson.D{
						{"$expr", bson.D{
							{"$and", bson.A{
								bson.D{{"$eq", bson.A{"$jobID", "$$jobID"}}},
								bson.D{{"$eq", bson.A{"$careerID", "$$careerID"}}},
							}},
						}},
					}}},
				}},
				{"as", "applications"},
			}},
		}
		lookUpSaveStage := bson.D{
			{"$lookup", bson.D{
				{"from", "CareerSaveJob"},
				{"let", bson.D{
					{"jobID", "$_id"},
					{"careerID", decodeId},
				}},
				{"pipeline", mongo.Pipeline{
					{{"$match", bson.D{
						{"$expr", bson.D{
							{"$and", bson.A{
								bson.D{{"$eq", bson.A{"$careerID", "$$careerID"}}},
								bson.D{{"$in", bson.A{"$$jobID", "$saveJob"}}},
							}},
						}},
					}}},
				}},
				{"as", "saved"},
			}},
		}

		projectStage := bson.D{
			{"$set", bson.D{
				{"isApplied", bson.D{
					{"$gt", bson.A{bson.D{{"$size", "$applications"}}, 0}},
				}},
				{"isSaved", bson.D{
					{"$gt", bson.A{bson.D{{"$size", "$saved"}}, 0}},
				}},
			}},
		}

		unsetStage := bson.D{
			{"$unset", bson.A{"applications", "saved"}},
		}

		pipeline = append(pipeline, lookUpApplyStage, lookUpSaveStage, projectStage, unsetStage)
	}

	return pipeline
}

func (j *JobRepository) SaveJob(careerID string, jobID string) (bson.M, error) {
	careerObjID, _ := primitive.ObjectIDFromHex(careerID)
	jobObjID, _ := primitive.ObjectIDFromHex(jobID)

	filter := bson.M{
		"careerID": careerObjID,
		"saveJob":  bson.M{"$in": []primitive.ObjectID{jobObjID}},
	}

	var existingDoc bson.M
	err := j.careerSaveCollection.FindOne(context.Background(), filter).Decode(&existingDoc)

	if err == mongo.ErrNoDocuments {
		update := bson.M{
			"$push": bson.M{
				"saveJob": jobObjID,
			},
		}
		opts := options.Update().SetUpsert(true)
		result, err := j.careerSaveCollection.UpdateOne(
			context.Background(),
			bson.M{"careerID": careerObjID},
			update,
			opts,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to save job: %v", err)
		}

		return bson.M{
			"success":  true,
			"message":  "Job saved successfully",
			"modified": result.ModifiedCount,
			"upserted": result.UpsertedCount,
		}, nil
	} else if err != nil {
		return nil, fmt.Errorf("error checking saved jobs: %v", err)
	}

	//if existed
	update := bson.M{
		"$pull": bson.M{
			"saveJob": jobObjID,
		},
	}

	result, err := j.careerSaveCollection.UpdateOne(
		context.Background(),
		bson.M{"careerID": careerObjID},
		update,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to unsave job: %v", err)
	}

	return bson.M{
		"success":  true,
		"message":  "Job unsaved successfully",
		"modified": result.ModifiedCount,
	}, nil

	return existingDoc, nil
}
