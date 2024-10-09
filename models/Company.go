package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CompanyImage struct {
	ImageURL string `bson:"imageURL" json:"imageURL"`
	CoverURL string `bson:"coverURL" json:"coverURL"`
}

type Contact struct {
	CompanyPhone   string `bson:"companyPhone" json:"companyPhone"`
	CompanyEmail   string `bson:"companyEmail" json:"companyEmail"`
	CompanyWebsite string `bson:"companyWebsite" json:"companyWebsite"`
	CompanyAddress string `bson:"companyAddress" json:"companyAddress"`
}

type Company struct {
	Id             primitive.ObjectID   `bson:"_id" json:"_id"`
	CompanyImage   CompanyImage         `bson:"companyImage" json:"companyImage"`
	CompanyName    string               `bson:"companyName" json:"companyName"`
	CompanyViewed  int                  `bson:"companyViewed" json:"companyViewed"`
	Contact        Contact              `bson:"contact" json:"contact"`
	CreateAt       primitive.DateTime   `bson:"createAt" json:"createAt"`
	Description    string               `bson:"description" json:"description"`
	EmployeeSize   int                  `bson:"employeeSize" json:"employeeSize"`
	FieldOperation []primitive.ObjectID `bson:"fieldOperation" json:"fieldOperation"`
	IsDeleted      bool                 `bson:"isDeleted" json:"isDeleted"`
	Popularity     int                  `bson:"popularity" json:"popularity"`
	PostJob        []primitive.ObjectID `bson:"postJob" json:"postJob"`
	TypeOfCompany  []string             `bson:"typeOfCompany" json:"typeOfCompany"`
	Password       string               `bson:"password" json:"password"`
}
