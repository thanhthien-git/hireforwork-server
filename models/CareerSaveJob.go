package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CareerSaveJob struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CareerID  primitive.ObjectID `bson:"careerID" json:"careerID"`
	JobID     primitive.ObjectID `json:"job_id,omitempty"`
	IsDeleted bool               `bson:"isDeleted" json:"isDeleted"`
}
type SaveJobRequest struct {
	CareerID primitive.ObjectID `bson:"careerID" json:"careerID"`
	JobID    primitive.ObjectID `bson:"jobID" json:"jobID"`
}
