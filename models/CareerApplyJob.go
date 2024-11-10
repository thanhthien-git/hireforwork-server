package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CareerApplyJob struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	CareerID  primitive.ObjectID `bson:"careerID" json:"careerID"`
	CareerCV  string             `bson:"careerCV" json:"careerCV"`
	JobID     primitive.ObjectID `bson:"jobID" json:"jobID"`
	CompanyID primitive.ObjectID `bson:"companyID" json:"companyID"`
	CreateAt  primitive.DateTime `bson:"createAt" json:"createAt"`
	IsDeleted bool               `bson:"isDeleted" json:"isDeleted"`
	IsChange  bool               `bson:"isChange" json:"isChange"`
	Status    string             `bson:"status" json:"status"`
}
