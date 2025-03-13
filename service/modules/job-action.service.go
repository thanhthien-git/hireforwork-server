package service

import (
	"context"
	"hireforwork-server/constants"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JobActionService struct {
	careerSaveJob  *mongo.Collection
	careerApplyJob *mongo.Collection
}

func NewJobActionService(careerSaveJob, careerApplyJob *mongo.Collection) *JobActionService {
	if careerSaveJob == nil || careerApplyJob == nil {
		log.Fatalf("Invailid connection: Job Action")
	}
	return &JobActionService{careerApplyJob: careerApplyJob, careerSaveJob: careerSaveJob}
}

func (a *JobActionService) Apply(request interfaces.IJobApply) error {
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
		IsChange:  false,
		CompanyID: companyID,
	}

	_, err := a.careerApplyJob.InsertOne(context.Background(), newApply)
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
            <p>Chào bạn,</p>
            <p>Cảm ơn bạn đã ứng tuyển công việc trên hệ thống của chúng tôi, nhà tuyển dụng sẽ liên lạc với bạn qua gmail sớm nhất.</p>
            <p>Nhà tuyển dụng sẽ phản hồi bạn sớm nhất thông qua email. Nếu bạn có bất kỳ câu hỏi nào, đừng ngần ngại liên hệ với chúng tôi qua email này.</p>
            <p>Chúc bạn một ngày tuyệt vời!</p>
            <p>Cheers, </p>
            <p>The NHIEUViec Team</p>
        </div>
    </body>
    </html>
    `
	if err = SendEmail(request.CareerEmail, subject, body); err != nil {
		return err
	}
	return nil
}

func (a *JobActionService) SaveJob(careerID string, jobID string) (models.CareerSaveJob, error) {
	careerObjID, _ := primitive.ObjectIDFromHex(careerID)
	jobObjID, _ := primitive.ObjectIDFromHex(jobID)

	filter := bson.M{"careerID": careerObjID}

	update := bson.M{
		"$setOnInsert": bson.M{
			"careerID": careerObjID,
		},
		"$addToSet": bson.M{
			"saveJob": jobObjID,
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var careerListSave models.CareerSaveJob
	err := a.careerSaveJob.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&careerListSave)
	if err != nil {
		return models.CareerSaveJob{}, err
	}
	return careerListSave, nil
}

func (a *JobActionService) RemoveSaveJob(careerID string, jobID string) (models.CareerSaveJob, error) {
	careerObjID, _ := primitive.ObjectIDFromHex(careerID)
	jobObjID, err := primitive.ObjectIDFromHex(jobID)

	filter := bson.M{"careerID": careerObjID}

	update := bson.M{
		"$pull": bson.M{
			"saveJob": jobObjID,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var careerListSave models.CareerSaveJob

	err = a.careerSaveJob.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&careerListSave)
	if err != nil {
		return models.CareerSaveJob{}, err
	}
	return careerListSave, nil
}
