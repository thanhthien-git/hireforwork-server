package service

import (
	"context"
	"fmt"
	"hireforwork-server/models"
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
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func GetUser(page, pageSize int, careerFirstName, lastName, careerEmail, careerPhone string) (models.PaginateDocs[models.User], error) {
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

	totalDocs, _ := userCollection.CountDocuments(context.Background(), bsonFilter)
	totalPage := int64(math.Ceil(float64(totalDocs) / float64(pageSize)))
	cursor, err := userCollection.Find(context.Background(), bsonFilter, findOption)
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

func GetUserByID(careerID string) (models.User, error) {
	_id, _ := primitive.ObjectIDFromHex(careerID)
	var user models.User

	err := userCollection.FindOne(context.Background(), bson.D{{"_id", _id}}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetUserByEmail(careerEmail string) (models.User, error) {
	var user models.User

	err := userCollection.FindOne(context.Background(), bson.D{{"careerEmail", careerEmail}, {"isDeleted", false}}).Decode(&user)
	if err != nil {
		return models.User{}, nil
	}
	return user, nil
}

func DeleteUserByID(careerID string) http.Response {
	_id, _ := primitive.ObjectIDFromHex(careerID)

	filter := bson.M{"_id": _id}

	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := userCollection.FindOneAndUpdate(context.Background(), filter, update, opts)

	if result.Err() != nil {
		return http.Response{
			StatusCode: http.StatusBadRequest,
		}
	}
	return http.Response{
		StatusCode: http.StatusAccepted,
	}
}

func CreateUser(user models.User) (models.User, error) {

	if user.Profile.UserCV == nil {
		user.Profile.UserCV = []string{} // Default to an empty array
	}
	if user.Profile.Skills == nil {
		user.Profile.Skills = []string{} // Default to an empty array
	}

	// Sử dụng session để đảm bảo tính toàn vẹn của dữ liệu
	wc := writeconcern.New(writeconcern.WMajority())
	opts := options.Transaction().SetWriteConcern(wc)

	// Khởi tạo session
	session, err := userCollection.Database().Client().StartSession()
	if err != nil {
		return models.User{}, fmt.Errorf("Error starting session: %v", err)
	}
	defer session.EndSession(context.Background())

	// Thực hiện transaction
	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		err := session.StartTransaction(opts)
		if err != nil {
			return fmt.Errorf("Error starting transaction: %v", err)
		}

		// Kiểm tra xem email đã tồn tại chưa
		filter := bson.M{"careerEmail": user.CareerEmail}
		count, err := userCollection.CountDocuments(sc, filter)
		if err != nil {
			_ = session.AbortTransaction(sc)
			return fmt.Errorf("Error checking email existence: %v", err)
		}

		if count > 0 {
			_ = session.AbortTransaction(sc)
			return fmt.Errorf("Account has already been registered")
		}

		// Nếu email chưa tồn tại, thêm user
		_, err = userCollection.InsertOne(sc, user)
		if err != nil {
			_ = session.AbortTransaction(sc)
			return fmt.Errorf("Error inserting user: %v", err)
		}

		// Commit transaction
		err = session.CommitTransaction(sc)
		if err != nil {
			_ = session.AbortTransaction(sc)
			return fmt.Errorf("Error committing transaction: %v", err)
		}

		return nil
	})

	if err != nil {
		// Kiểm tra lỗi trùng lặp email
		if err.Error() == "Account has already been registered" {
			return models.User{}, fmt.Errorf("duplicate_email") // Trả về lỗi trùng lặp email
		}
		return models.User{}, err
	}

	return user, nil
}
func UpdateUserByID(userID string, updatedUser models.User) (models.User, error) {
	_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.User{}, fmt.Errorf("invalid user ID format: %v", err)
	}

	filter := bson.M{"_id": _id, "isDeleted": false}

	var existingUser models.User
	err = userCollection.FindOne(context.Background(), filter).Decode(&existingUser)
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

	result, err := userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return models.User{}, fmt.Errorf("error updating user: %v", err)
	}

	if result.ModifiedCount == 0 {
		return models.User{}, fmt.Errorf("no changes were made to the user with ID %s", userID)
	}

	return updatedUser, nil
}
func SaveJob(careerID string, jobID string) (models.CareerSaveJob, error) {
	careerObjID, _ := primitive.ObjectIDFromHex(careerID)
	jobObjID, _ := primitive.ObjectIDFromHex(jobID)

	filter := bson.M{"careerID": careerObjID}

	update := bson.M{
		"$setOnInsert": bson.M{
			"careerID": careerObjID,
		},
		"$addToSet": bson.M{
			"saveJob": jobObjID,
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var careerListSave models.CareerSaveJob
	err := careerSaveJob.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&careerListSave)
	if err != nil {
		return models.CareerSaveJob{}, err
	}
	return careerListSave, nil
}

func RemoveSaveJob(careerID string, jobID string) (models.CareerSaveJob, error) {
	careerObjID, _ := primitive.ObjectIDFromHex(careerID)
	jobObjID, err := primitive.ObjectIDFromHex(jobID)

	filter := bson.M{"careerID": careerObjID}

	update := bson.M{
		"$pull": bson.M{
			"saveJob": jobObjID,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var careerListSave models.CareerSaveJob

	err = careerSaveJob.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&careerListSave)
	if err != nil {
		return models.CareerSaveJob{}, err
	}
	return careerListSave, nil
}

func CareerViewedJob(careerID string, jobID string) (models.CareerViewedJob, error) {
	careerObjID, err := primitive.ObjectIDFromHex(careerID)
	if err != nil {
		return models.CareerViewedJob{}, fmt.Errorf("invalid career ID: %v", err)
	}

	jobObjID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return models.CareerViewedJob{}, fmt.Errorf("invalid job ID: %v", err)
	}

	viewedJob := models.ViewedJob{
		JobID: jobObjID,
	}

	// Kiểm tra cập nhật document
	filter := bson.M{"careerID": careerObjID}
	update := bson.M{
		"$addToSet": bson.M{"viewedJob": viewedJob}, // Add to set to avoid duplicates
	}

	opts := options.Update().SetUpsert(true)
	result, err := careerViewedJob.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return models.CareerViewedJob{}, fmt.Errorf("failed to update viewed jobs: %v", err)
	}

	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return models.CareerViewedJob{}, fmt.Errorf("no document matched or inserted")
	}

	// Lấy lại document đã cập nhật
	var careerViewed models.CareerViewedJob
	err = careerViewedJob.FindOne(context.Background(), filter).Decode(&careerViewed)
	if err != nil {
		return models.CareerViewedJob{}, fmt.Errorf("failed to retrieve updated document: %v", err)
	}

	return careerViewed, nil

}

func GetSavedJobByCareerID(careerID string) (models.PaginateDocs[models.Jobs], error) {
	var saveJob []models.Jobs
	careerObjID, _ := primitive.ObjectIDFromHex(careerID)

	filter := bson.M{"careerID": careerObjID}

	var career models.CareerSaveJob
	err := careerSaveJob.FindOne(context.Background(), filter).Decode(&career)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.PaginateDocs[models.Jobs]{}, fmt.Errorf("Bạn chưa lưu công việc nào", careerID)
		}
	}

	jobFilter := bson.M{
		"_id": bson.M{
			"$in": career.SaveJob,
		},
		"isDeleted": false,
	}

	totalDocs, _ := jobCollection.CountDocuments(context.Background(), jobFilter)
	totalPage := int64(math.Ceil(float64(totalDocs)) / float64(10))

	cursor, err := jobCollection.Find(context.Background(), jobFilter)
	if err != nil {
		return models.PaginateDocs[models.Jobs]{}, fmt.Errorf("Đã có lỗi xảy ra")
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &saveJob); err != nil {
		log.Printf("Error parsing documents: %v", err)
		return models.PaginateDocs[models.Jobs]{}, err
	}

	result := models.PaginateDocs[models.Jobs]{
		Docs:        saveJob,
		TotalDocs:   totalDocs,
		CurrentPage: int64(1),
		TotalPage:   totalPage,
	}
	return result, nil
}

func GetViewedJobByCareerID(careerID string) ([]models.ViewedJob, error) {
	careerObjID, err := primitive.ObjectIDFromHex(careerID)
	if err != nil {
		return nil, fmt.Errorf("invalid career ID: %v", err)
	}

	filter := bson.M{"careerID": careerObjID}

	var careerViewed models.CareerViewedJob
	err = careerViewedJob.FindOne(context.Background(), filter).Decode(&careerViewed)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no viewed jobs found for career ID: %s", careerID)
		}
		return nil, fmt.Errorf("error retrieving viewed jobs: %v", err)
	}

	return careerViewed.ViewedJob, nil
}

func UpdateCareerImage(link string, id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"careerPicture": link,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

	res := userCollection.FindOneAndUpdate(context.Background(), filter, update, opt)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func UpdateCareerResume(link string, id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$push": bson.M{
			"profile.userCV": link,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

	res := userCollection.FindOneAndUpdate(context.Background(), filter, update, opt)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func RemoveResume(id string, link string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$pull": bson.M{
			"profile.userCV": link,
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := userCollection.FindOneAndUpdate(context.Background(), filter, update, opt)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
func RequestPasswordReset(email string) (string, error) {
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"careerEmail": email}).Decode(&user)
	if err != nil {
		return "", fmt.Errorf("no user found with email %s: %v", email, err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := rng.Intn(9000) + 1000
	user.VerificationCode = fmt.Sprintf("%d", code)
	subject := "Mã xác nhận khôi phục mật khẩu"
	body := fmt.Sprintf("Mã xác nhận của bạn là: %s", user.VerificationCode)
	if err := SendEmail(email, subject, body); err != nil {
		return "", err
	}
	_, err = userCollection.UpdateOne(context.Background(), bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"verificationCode": user.VerificationCode}})
	if err != nil {
		return "", fmt.Errorf("Lỗi khi update mã xác nhận: %v", err)
	}

	return user.VerificationCode, nil
}

func ResetPassword(email string, code string, newPassword string) error {
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"careerEmail": email, "verificationCode": code}).Decode(&user)
	if err != nil {
		return fmt.Errorf("Sai mã xác nhận: %v", err)
	}
	hashedPassword := utils.EncodeToSHA(newPassword)

	_, err = userCollection.UpdateOne(context.Background(), bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"password": hashedPassword, "verificationCode": ""}})
	if err != nil {
		return fmt.Errorf("Lỗi khi cố lấy lại mật khẩu: %v", err)
	}

	return nil
}
