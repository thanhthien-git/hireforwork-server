package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tech struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	TechName  string             `json:"technology" bson:"technology,omitempty"`
	IsDeleted bool               `json:"isDeleted" bson:"isDeleted,omitempty"`
}
