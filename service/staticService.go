package service

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewUser struct using primitive.ObjectID
type NewUser struct {
	Id            primitive.ObjectID `json:"_id"` // Use ObjectID directly
	CareerPicture string             `json:"careerPicture"`
	CareerEmail   string             `json:"careerEmail"`
}

// NewPost struct using primitive.ObjectID
type NewPost struct {
	Id         primitive.ObjectID `json:"_id"` // Use ObjectID directly
	JobTitle   string             `json:"jobTitle"`
	ExpireDate primitive.DateTime `json:"expireDate"`
}

// Static struct containing the counts and new users and posts
type Static struct {
	TotalUser    int64     `json:"totalUser"`
	TotalCompany int64     `json:"totalCompany"`
	TotalPost    int64     `json:"totalPost"`
	TotalResume  int64     `json:"totalResume"`
	NewUser      []NewUser `json:"newUser"`
	NewPost      []NewPost `json:"newPost"`
}

// GetStatic function to gather statistics
func GetStatic() Static {
	var static Static

	// Filter to exclude deleted entries
	filter := bson.M{
		"isDeleted": false,
	}

	// Count documents for users, companies, and job posts
	static.TotalUser, _ = userCollection.CountDocuments(context.Background(), filter)
	static.TotalCompany, _ = companyCollection.CountDocuments(context.Background(), filter)
	static.TotalPost, _ = jobCollection.CountDocuments(context.Background(), filter)

	// Filter for resumes
	resumeFilter := bson.M{
		"isDeleted": false,
		"careerCV":  bson.M{"$ne": nil},
	}
	static.TotalResume, _ = userCollection.CountDocuments(context.Background(), resumeFilter)

	// Options for limiting to 5 documents and sorting by createdAt in descending order
	opt := options.Find().SetLimit(5).SetSort(bson.M{"createAt": -1})

	// Fetch the newest users
	cur1, err := userCollection.Find(context.Background(), filter, opt)
	if err != nil {
		log.Fatal(err)
	}
	defer cur1.Close(context.Background())

	var newUsers []NewUser
	for cur1.Next(context.Background()) {
		var user struct {
			ID            primitive.ObjectID `bson:"_id"`
			CareerPicture string             `bson:"careerPicture"`
			CareerEmail   string             `bson:"careerEmail"`
		}
		if err := cur1.Decode(&user); err != nil {
			log.Fatal(err)
		}

		// Add the decoded user to the newUsers slice
		newUser := NewUser{
			Id:            user.ID, // Directly use the primitive.ObjectID
			CareerPicture: user.CareerPicture,
			CareerEmail:   user.CareerEmail,
		}
		newUsers = append(newUsers, newUser)
	}
	static.NewUser = newUsers

	// Fetch the newest job posts
	cur2, err := jobCollection.Find(context.Background(), filter, opt)
	if err != nil {
		log.Fatal(err)
	}
	defer cur2.Close(context.Background())

	var newPosts []NewPost
	for cur2.Next(context.Background()) {
		var post struct {
			ID         primitive.ObjectID `bson:"_id"`
			JobTitle   string             `bson:"jobTitle"`
			ExpireDate primitive.DateTime `bson:"expireDate"`
		}
		if err := cur2.Decode(&post); err != nil {
			log.Fatal(err)
		}

		// Add the decoded job post to the newPosts slice
		newPost := NewPost{
			Id:         post.ID, // Directly use the primitive.ObjectID
			JobTitle:   post.JobTitle,
			ExpireDate: post.ExpireDate,
		}
		newPosts = append(newPosts, newPost)
	}
	static.NewPost = newPosts

	return static
}
