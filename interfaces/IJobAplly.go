package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type IJobApply struct {
	JobID    primitive.ObjectID `json:"jobID"`
	IDCareer primitive.ObjectID `json:"idCareer"`
	CreateAt primitive.DateTime `json:"createAt"`
}
