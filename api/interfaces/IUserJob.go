package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type IUserJob struct {
	JobID    primitive.ObjectID `json:"jobID"`
	IDCareer primitive.ObjectID `json:"idCareer"`
}
