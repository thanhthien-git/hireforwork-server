package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SavedJob struct {
	JobID     primitive.ObjectID `bson:"jobID" json:"jobID"`
	IsDeleted bool               `bson:"isDeleted" json:"isDeleted"`
}

type CareerSaveJob struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CareerID primitive.ObjectID `bson:"careerID" json:"careerID"`
	SaveJob  []SavedJob         `bson:"saveJob" json:"saveJob"`
}

type SaveJobRequest struct {
	CareerID primitive.ObjectID `bson:"careerID" json:"careerID"`
	JobID    primitive.ObjectID `bson:"jobID" json:"jobID"`
}
