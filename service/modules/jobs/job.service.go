package jobs

import (
	"hireforwork-server/db"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"

	"go.mongodb.org/mongo-driver/bson"
)

type JobService struct {
	repo *JobRepository
}

type IJobFilter struct {
	page, pageSize int
}

func NewJobService(dbInstance *db.DB) *JobService {
	return &JobService{repo: NewJobRepository(dbInstance)}
}

func (j *JobService) GetJob(page, pageSize int, filter interfaces.IJobFilter) (bson.M, error) {
	return j.repo.GetJob(page, pageSize, filter)
}

func (j *JobService) CreateJob(job models.Jobs) (models.Jobs, error) {
	return j.repo.CreateJob(job)
}

func (j *JobService) UpdateJob(job models.Jobs) (models.Jobs, error) {
	return j.repo.UpdateJob(job)
}

func (j *JobService) GetLatestJobs() ([]models.Jobs, error) {
	return j.repo.GetLatestJobs()
}

func (j *JobService) GetJobByID(jobID string, userId string) (bson.M, error) {
	return j.repo.GetJobByID(jobID, userId)
}

func (j *JobService) SaveJob(careerID string, jobID string) (bson.M, error) {
	return j.repo.SaveJob(careerID, jobID)
}

func (j *JobService) Apply(request interfaces.IJobApply) error {
	return j.repo.Apply(request)
}
