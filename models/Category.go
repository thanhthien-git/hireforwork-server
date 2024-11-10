package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	CategoryName string             `json:"categoryName" bson:"categoryName,omitempty"`
	IsDeleted    bool               `json:"isDeleted" bson:"isDeleted"`
}
