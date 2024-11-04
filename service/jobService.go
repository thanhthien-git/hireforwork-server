package service

import (
	"context"
	"fmt"
	"hireforwork-server/constants"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/utils"
	"log"
	"math"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetJob(page, pageSize int, filter interfaces.IJobFilter) (bson.M, error) {

	skip := (page - 1) * pageSize
	matchStage := bson.M{"isDeleted": false}

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
	if filter.SalaryFrom != "" && filter.SalaryTo != "" {
		matchOption["jobSalaryMin"] = bson.M{"$gte": filter.SalaryFrom}
		matchOption["jobSalaryMax"] = bson.M{"$lte": filter.SalaryTo}
	}
	//filter by working location
	if len(filter.WorkingLocation) > 0 {
		matchOption["workingLocation"] = bson.M{"$in": filter.WorkingLocation}
	}
	//filter by job require
	if len(filter.JobRequirement) > 0 {
		matchOption["jobRequirement"] = bson.M{"$in": filter.JobRequirement}
	}
	if filter.IsHot {
		matchOption["isHot"] = true
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
						{"isClosed", "$$doc.isClosed"},
						{"isDeleted", "$$doc.isDeleted"},
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
	cursor, err := jobCollection.Aggregate(context.Background(), pipeline)

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

func UpdateJob(job models.Jobs) (models.Jobs, error) {
	filter := bson.M{"_id": job.Id}

	update := bson.M{
		"$set": bson.M{
			"jobTitle":        job.JobTitle,
			"jobSalaryMin":    job.JobSalaryMin,
			"jobSalaryMax":    job.JobSalaryMax,
			"jobRequirement":  job.JobRequirement,
			"workingLocation": job.WorkingLocation,
			"isHot":           job.IsHot,
			"isClosed":        job.IsClosed,
			"isDeleted":       job.IsDeleted,
			"expireDate":      job.ExpireDate,
			"jobCategory":     job.JobCategory,
			"quantity":        job.Quantity,
			"jobDescription":  job.JobDescription,
			"jobLevel":        job.JobLevel,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := jobCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&job)
	if err != nil {
		fmt.Println(err)
		return models.Jobs{}, fmt.Errorf("Có lỗi xảy ra khi cập nhập lại thông tin")
	}
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

func ApplyForJob(request interfaces.IJobApply) error {
	id, _ := utils.ConvertoObjectID(request.IDCareer)
	companyID, _ := utils.ConvertoObjectID(request.CompanyID)
	jobID, _ := utils.ConvertoObjectID(request.JobID)

	newApply := models.CareerApplyJob{
		ID:        primitive.NewObjectID(),
		CareerID:  id,
		JobID:     jobID,
		CreateAt:  primitive.NewDateTimeFromTime(time.Now()),
		CareerCV:  request.CareerCV,
		IsDeleted: false,
		Status:    constants.PENDING,
		CompanyID: companyID,
	}

	_, err := careerApplyJob.InsertOne(context.Background(), newApply)
	if err != nil {
		log.Printf("Loi o day")
		return err
	}

	subject := "Cảm ơn bạn đã ứng tuyển"
	body := `
<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: auto;
            background: #ffffff;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            padding: 10px 0;
        }
        .header h1 {
            color: #4a4a4a;
        }
        .footer {
            margin-top: 20px;
            text-align: center;
            font-size: 0.8em;
            color: #666666;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Cảm ơn bạn đã ứng tuyển!</h1>
        </div>
        <p>Xin chào [Tên Ứng Viên],</p>
        <p>Cảm ơn bạn đã ứng tuyển vào vị trí <strong>[Tên Vị Trí]</strong> tại [Tên Công Ty]. Chúng tôi rất vui mừng khi nhận được hồ sơ của bạn.</p>
        <p>Đội ngũ tuyển dụng của chúng tôi sẽ xem xét hồ sơ của bạn và sẽ liên hệ trong thời gian sớm nhất. Nếu bạn có bất kỳ câu hỏi nào, đừng ngần ngại liên hệ với chúng tôi qua email này.</p>
        <p>Chúc bạn một ngày tuyệt vời!</p>
        <p>Trân trọng,</p>
        <p><em>Đội ngũ tuyển dụng tại [Tên Công Ty]</em></p>
        <div class="footer">
            <p>[Tên Công Ty] | [Địa chỉ Công Ty] | [Số điện thoại]</p>
        </div>
    </div>
</body>
</html>
`
	if err = SendEmail(request.CareerEmail, subject, body); err != nil {
		return err
	}

	return nil
}

func GetLatestJobs() ([]models.Jobs, error) {
	var jobs []models.Jobs

	filter := bson.M{"isDeleted": false}
	opts := options.Find().SetSort(bson.D{{"createAt", -1}}).SetLimit(10)

	cursor, err := jobCollection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func GetJobByID(jobID string, tokenString string) (bson.M, error) {
	_id, err := primitive.ObjectIDFromHex(jobID)

	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{{"_id", _id}}},
		},
		{
			{"$lookup", bson.D{
				{"from", "Company"},
				{"localField", "companyID"},
				{"foreignField", "_id"},
				{"as", "company"},
			}},
		},
		{
			{"$unwind", bson.D{{"path", "$company"}, {"preserveNullAndEmptyArrays", true}}},
		},
		{
			{"$addFields", bson.D{
				{"contact", "$company.contact"},
				{"companyName", "$company.companyName"},
				{"companyImage", "$company.companyImage"},
				{"employeeSize", "$company.employeeSize"},
			}},
		},
		{
			{"$unset", "company"},
		},
	}

	var userID primitive.ObjectID
	if tokenString != "" {
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			return nil, err
		}
		userID, _ = primitive.ObjectIDFromHex(claims["sub"].(string))
	}

	if userID != primitive.NilObjectID {

		lookUpApplyStage := bson.D{
			{"$lookup", bson.D{
				{"from", "CareerApplyJob"},
				{"let", bson.D{
					{"jobID", _id},
					{"careerID", userID},
				}},
				{"pipeline", mongo.Pipeline{
					{
						{"$match", bson.D{
							{"$expr", bson.D{
								{"$and", bson.A{
									bson.D{{"$eq", bson.A{"$jobID", "$$jobID"}}},
									bson.D{{"$eq", bson.A{"$careerID", "$$careerID"}}},
								}},
							}},
						}},
					},
				}},
				{"as", "applications"},
			}},
		}
		lookUpSaveStage := bson.D{
			{"$lookup", bson.D{
				{"from", "CareerSaveJob"},
				{"let", bson.D{
					{"jobID", _id},       // The job ID to check
					{"careerID", userID}, // The career ID to match
				}},
				{"pipeline", mongo.Pipeline{
					{
						{"$match", bson.D{
							{"$expr", bson.D{
								{"$and", bson.A{
									bson.D{{"$eq", bson.A{"$careerID", "$$careerID"}}},
									bson.D{{"$in", bson.A{"$$jobID", "$saveJob"}}},
								}},
							}},
						}},
					},
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

	cursor, err := jobCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var job bson.M
	if cursor.Next(context.Background()) {
		err = cursor.Decode(&job)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, mongo.ErrNoDocuments
	}

	result := bson.M{
		"doc": job,
	}

	return result, nil
}
