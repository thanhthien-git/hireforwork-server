package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CareerApplyJob struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	CareerID  primitive.ObjectID `bson:"careerID" json:"careerID"`
	JobID     primitive.ObjectID `bson:"jobID" json:"jobID"`
	CreateAt  primitive.DateTime `bson:"createAt" json:"createAt"`
	IsDeleted bool               `bson:"isDeleted" json:"isDeleted"`
	Status    string             `bson:"status" json:"status"`
}
