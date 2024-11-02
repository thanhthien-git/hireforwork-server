package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CareerSaveJob struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	CareerID primitive.ObjectID   `bson:"careerID" json:"careerID"`
	SaveJob  []primitive.ObjectID `bson:"saveJob" json:"saveJob"`
}
