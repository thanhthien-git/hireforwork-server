package jobs

import (
	"hireforwork-server/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobAPI struct {
	Id           primitive.ObjectID `json:"_id" `
	JobTitle     string             `json:"jobTitle" `
	JobSalaryMin int64              `json:"jobSalaryMin"`
	JobSalaryMax int64              `json:"jobSalaryMax"`
}

type JobAdapter struct {
	repo *JobRepository
}

func (a *JobAdapter) ToAPI(job models.Jobs) JobAPI {
	return JobAPI{
		Id:           job.Id,
		JobTitle:     job.JobTitle,
		JobSalaryMin: job.JobSalaryMin,
		JobSalaryMax: job.JobSalaryMax,
	}
}
