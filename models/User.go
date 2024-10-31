package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	UserCV []string `bson:"userCV" json:"userCV"`
	Skills []string `bson:"skills" json:"skills"`
}

type User struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	FirstName     string             `bson:"careerFirstName" json:"careerFirstName" validate:"required"`
	LastName      string             `bson:"lastName" json:"lastName" validate:"required"`
	CareerPhone   string             `bson:"careerPhone" json:"careerPhone" validate:"required"`
	CareerEmail   string             `bson:"careerEmail" json:"careerEmail" validate:"required"`
	CareerPicture string             `bson:"careerPicture,omitempty" json:"careerPicture,omitempty"`
	CreateAt      primitive.DateTime `bson:"createAt" json:"createAt"`
	IsDeleted     bool               `bson:"isDeleted" json:"isDeleted"`
	Languages     []string           `bson:"languages,omitempty" json:"languages,omitempty"`
	Password      string             `bson:"password" json:"password"`
	Role          string             `bson:"role" json:"role"`
	Profile       Profile            `bson:"profile,omitempty" json:"profile,omitempty"`
}
