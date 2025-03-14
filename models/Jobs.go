package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfo struct {
	UserId     primitive.ObjectID `bson:"userID" json:"userID"`
	IsAccepted string             `bson:"isAccepted" json:"isAccepted"`
	CreateAt   primitive.DateTime `bson:"createAt" json:"createAt"`
}

type Jobs struct {
	Id               primitive.ObjectID `bson:"_id" json:"_id" ,omitempty`
	JobTitle         string             `bson:"jobTitle" json:"jobTitle" validate:"required"`
	JobSalaryMin     int64              `bson:"jobSalaryMin" json:"jobSalaryMin"`
	JobSalaryMax     int64              `bson:"jobSalaryMax" json:"jobSalaryMax"`
	JobRequirement   []string           `bson:"jobRequirement" json:"jobRequirement"`
	JobCategory      []string           `bson:"jobCategory" json:"jobCategory"`
	WorkingLocation  []string           `bson:"workingLocation" json:"workingLocation"`
	IsHot            bool               `bson:"isHot" json:"isHot"`
	IsClosed         bool               `bson:"isClosed" json:"isClosed"`
	IsDeleted        bool               `bson:"isDeleted" json:"isDeleted"`
	CreateAt         primitive.DateTime `bson:"createAt" json:"createAt"`
	ExpireDate       primitive.DateTime `bson:"expireDate" json:"expireDate"`
	CompanyID        primitive.ObjectID `bson:"companyID" json:"companyID"`
	JobTech          []string           `bson:"jobTech" json:"jobTech"`
	Quantity         int64              `bson:"quantity" json:"quantity"`
	JobDescription   string             `bson:"jobDescription" json:"jobDescription"`
	JobLevel         string             `bson:"jobLevel" json:"jobLevel"`
	WorkType         []string           `bson:"workingType" json:"workingType"`
	RecruitmentCount int64              `bson:"recruitmentCount" json:"recruitmentCount"`
}
