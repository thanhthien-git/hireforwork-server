package service

import (
	"context"
	"fmt"
	"hireforwork-server/db"
	"hireforwork-server/models"
	"hireforwork-server/service/modules/unit_of_work"
	"hireforwork-server/utils"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct {
	userCollection        *mongo.Collection
	userSaveJobCollection *mongo.Collection
	jobCollection         *mongo.Collection
	userApplyCollection   *mongo.Collection
	uow                   *unit_of_work.UnitOfWork
}

// Dependency Injection (DI)
func NewUserService(dbInstance *db.DB) *UserService {
	collections := dbInstance.GetCollections([]string{"Career", "Job", "CareerSaveJob", "CareerApplyJob"})
	return &UserService{userCollection: collections[0], userSaveJobCollection: collections[1], jobCollection: collections[2], userApplyCollection: collections[3], uow: unit_of_work.NewUnitOfWork(dbInstance)}
}

func (u *UserService) GetUser(page, pageSize int, careerFirstName, lastName, careerEmail, careerPhone string) (models.PaginateDocs[models.User], error) {
	var users []models.User
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	bsonFilter := bson.D{{"isDeleted", false}}

	skip := (page - 1) * pageSize

	findOption := options.Find().SetProjection(bson.D{{"password", 0}})
	findOption.SetLimit(int64(pageSize))
	findOption.SetSkip(int64(skip))

	if careerFirstName != "" {
		bsonFilter = append(bsonFilter, bson.E{"careerFirstName", bson.D{{"$regex", careerFirstName}, {"$options", "i"}}})
	}

	if lastName != "" {
		bsonFilter = append(bsonFilter, bson.E{"lastName", bson.D{{"$regex", lastName}, {"$options", "i"}}})
	}

	if careerEmail != "" {
		bsonFilter = append(bsonFilter, bson.E{"careerEmail", bson.D{{"$regex", careerEmail}, {"$options", "i"}}})
	}

	if careerPhone != "" {
		bsonFilter = append(bsonFilter, bson.E{"careerPhone", bson.D{{"$regex", careerPhone}, {"$options", "i"}}})
	}

	totalDocs, _ := u.userCollection.CountDocuments(context.Background(), bsonFilter)
	totalPage := int64(math.Ceil(float64(totalDocs) / float64(pageSize)))
	cursor, err := u.userCollection.Find(context.Background(), bsonFilter, findOption)
	if err != nil {
		log.Printf("Error finding documents: %v", err)
		return models.PaginateDocs[models.User]{}, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &users); err != nil {
		log.Printf("Error parsing documents: %v", err)
		return models.PaginateDocs[models.User]{}, err
	}

	result := models.PaginateDocs[models.User]{
		Docs:        users,
		TotalDocs:   totalDocs,
		CurrentPage: int64(page),
		TotalPage:   totalPage,
	}

	return result, nil
}

func (u *UserService) GetUserByID(careerID string) (models.User, error) {
	_id, _ := primitive.ObjectIDFromHex(careerID)
	var user models.User

	err := u.userCollection.FindOne(context.Background(), bson.D{{"_id", _id}}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *UserService) GetUserByEmail(careerEmail string) (models.User, error) {
	var user models.User

	err := u.userCollection.FindOne(context.Background(), bson.D{{"careerEmail", careerEmail}, {"isDeleted", false}}).Decode(&user)
	if err != nil {
		return models.User{}, nil
	}
	return user, nil
}

func (u *UserService) DeleteUserByID(careerID string) http.Response {
	_id, _ := primitive.ObjectIDFromHex(careerID)

	filter := bson.M{"_id": _id}

	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := u.userCollection.FindOneAndUpdate(context.Background(), filter, update, opts)

	if result.Err() != nil {
		return http.Response{
			StatusCode: http.StatusBadRequest,
		}
	}
	return http.Response{
		StatusCode: http.StatusAccepted,
	}
}

func (u *UserService) CreateUser(user models.User) error {
	u.uow.RegisterChange(func(ctx mongo.SessionContext) error {
		filter := bson.M{"careerEmail": user.CareerEmail}
		count, err := u.userCollection.CountDocuments(ctx, filter)
		if err != nil {
			return fmt.Errorf("Error checking email existence: %v", err)
		}
		if count > 0 {
			return fmt.Errorf("Account has already been registered")
		}

		_, err = u.userCollection.InsertOne(ctx, user)
		if err != nil {
			return fmt.Errorf("Error inserting user: %v", err)
		}

		return nil
	})

	// Execute the transaction
	if err := u.uow.Commit(); err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateUserByID(userID string, updatedUser models.User) (models.User, error) {
	_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.User{}, fmt.Errorf("invalid user ID format: %v", err)
	}

	filter := bson.M{"_id": _id, "isDeleted": false}

	var existingUser models.User
	err = u.userCollection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err != nil {
		return models.User{}, fmt.Errorf("no user found with ID %s: %v", userID, err)
	}

	update := bson.M{
		"$set": bson.M{
			"careerFirstName": updatedUser.FirstName,
			"lastName":        updatedUser.LastName,
			"careerEmail":     updatedUser.CareerEmail,
			"careerPhone":     updatedUser.CareerPhone,
			"careerPicture":   updatedUser.CareerPicture,
			"profile":         updatedUser.Profile,
			"languages":       updatedUser.Languages,
		},
	}

	result, err := u.userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return models.User{}, fmt.Errorf("error updating user: %v", err)
	}

	if result.ModifiedCount == 0 {
		return models.User{}, fmt.Errorf("no changes were made to the user with ID %s", userID)
	}

	return updatedUser, nil
}

func (u *UserService) GetSavedJobByCareerID(careerID string) models.PaginateDocs[models.Jobs] {
	var saveJob []models.Jobs
	careerObjID, _ := primitive.ObjectIDFromHex(careerID)
	fmt.Println(careerObjID)
	filter := bson.M{"careerID": careerObjID}

	var career models.CareerSaveJob
	err := u.userSaveJobCollection.FindOne(context.Background(), filter).Decode(&career)

	jobFilter := bson.M{
		"_id": bson.M{
			"$in": career.SaveJob,
		},
		"isDeleted": false,
	}

	totalDocs, _ := u.jobCollection.CountDocuments(context.Background(), jobFilter)
	totalPage := int64(math.Ceil(float64(totalDocs)) / float64(10))

	cursor, err := u.jobCollection.Find(context.Background(), jobFilter)
	if err != nil {
		return models.PaginateDocs[models.Jobs]{}
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &saveJob); err != nil {
		log.Printf("Error parsing documents: %v", err)
		return models.PaginateDocs[models.Jobs]{}
	}

	result := models.PaginateDocs[models.Jobs]{
		Docs:        saveJob,
		TotalDocs:   totalDocs,
		CurrentPage: int64(1),
		TotalPage:   totalPage,
	}
	return result
}

func (u *UserService) UpdateCareerImage(link string, id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"careerPicture": link,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

	res := u.userCollection.FindOneAndUpdate(context.Background(), filter, update, opt)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (u *UserService) UpdateCareerResume(link string, id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$push": bson.M{
			"profile.userCV": link,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

	res := u.userCollection.FindOneAndUpdate(context.Background(), filter, update, opt)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (u *UserService) RemoveResume(id string, link string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$pull": bson.M{
			"profile.userCV": link,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := u.userCollection.FindOneAndUpdate(context.Background(), filter, update, opt)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (u *UserService) RequestPasswordReset(email string) (string, error) {
	var user models.User
	err := u.userCollection.FindOne(context.Background(), bson.M{"careerEmail": email}).Decode(&user)
	if err != nil {
		return "", fmt.Errorf("Không tìm thấy người dùng với email %s: %v", email, err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := rng.Intn(9000) + 1000
	verificationCode := fmt.Sprintf("%d", code)

	subject := "Mã xác nhận khôi phục mật khẩu"
	body := fmt.Sprintf("Mã xác nhận của bạn là: %s", verificationCode)
	if err := SendEmail(email, subject, body); err != nil {
		return "", fmt.Errorf("Lỗi khi gửi email: %v", err)
	}

	_, err = u.userCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.Id},
		bson.M{"$set": bson.M{"verificationCode": verificationCode}},
	)
	if err != nil {
		return "", fmt.Errorf("Lỗi khi tạo mã xác nhận: %v", err)
	}

	return verificationCode, nil
}

func (u *UserService) ResetPassword(email string, code string, newPassword string) error {
	var user models.User

	// Kiểm tra mã xác nhận
	err := u.userCollection.FindOne(context.Background(), bson.M{"careerEmail": email, "verificationCode": code}).Decode(&user)
	if err != nil {
		return fmt.Errorf("Sai mã xác nhận hoặc email không hợp lệ: %v", err)
	}
	hashedPassword := utils.EncodeToSHA(newPassword)
	_, err = u.userCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.Id},
		bson.M{
			"$set":   bson.M{"password": hashedPassword},
			"$unset": bson.M{"verificationCode": ""},
		},
	)
	if err != nil {
		return fmt.Errorf("Lỗi khi cập nhật mật khẩu mới: %v", err)
	}

	return nil
}

func (u *UserService) GetAppliedJob(id string, page int, pageSize int) ([]bson.M, error) {
	_id, _ := primitive.ObjectIDFromHex(id)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"careerID", _id},
			{"isDeleted", false},
		}}},
		{{"$lookup", bson.D{
			{"from", "Job"},
			{"localField", "jobID"},
			{"foreignField", "_id"},
			{"as", "jobDetails"},
		}}},
		{{"$unwind", bson.D{
			{"path", "$jobDetails"},
			{"preserveNullAndEmptyArrays", false},
		}}},
		{{
			"$lookup", bson.D{
				{"from", "Company"},
				{"localField", "jobDetails.companyID"},
				{"foreignField", "_id"},
				{"as", "companyDetails"},
			},
		}},
		{{"$unwind", bson.D{
			{"path", "$companyDetails"},
			{"preserveNullAndEmptyArrays", false},
		}}},
		{{"$project", bson.D{
			{"_id", 1},
			{"careerID", 1},
			{"jobTitle", "$jobDetails.jobTitle"},
			{"jobID", "$jobDetails._id"},
			{"jobRequirement", "$jobDetails.jobRequirement"},
			{"jobSalaryMin", "$jobDetails.jobSalaryMin"},
			{"jobSalaryMax", "$jobDetails.jobSalaryMax"},
			{"companyImage", "$companyDetails.companyImage.imageURL"},
			{"companyName", "$companyDetails.companyName"},
			{"isDeleted", 1},
			{"status", 1},
		}}},
		//facet stage
		bson.D{
			{"$facet", bson.D{
				{"totalDocs", bson.A{
					bson.D{{"$count", "totalDocs"}},
				}},
				{"docs", bson.A{
					bson.D{{"$skip", skip}},
					bson.D{{"$limit", pageSize}},
				}},
			}},
		},
		bson.D{
			{"$addFields", bson.D{
				{"page", page},
			}},
		},
		bson.D{
			{"$project", bson.D{
				{"totalDocs", bson.D{{"$arrayElemAt", bson.A{"$totalDocs.totalDocs", 0}}}},
				{"docs", 1},
				{"page", page},
			}},
		},
	}

	cursor, err := u.userApplyCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
