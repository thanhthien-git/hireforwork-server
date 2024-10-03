package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ViewedJob struct {
    JobID primitive.ObjectID `bson:"jobID" json:"jobID"`
}

type CareerViewedJob struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
    CareerID  primitive.ObjectID `bson:"careerID" json:"careerID"`
    ViewedJob []ViewedJob        `bson:"viewedJob" json:"viewedJob"`
}
