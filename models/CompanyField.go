package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Field struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	FieldName string             `bson:"fieldName" json:"fieldName"`
	IsDeleted bool               `bson:"isDeleted" json:"isDeleted"`
}
