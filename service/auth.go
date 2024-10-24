package service

import (
	"context"
	"errors"
	dbHelper "hireforwork-server/db"
	"hireforwork-server/models"
	"hireforwork-server/utils"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AuthService struct {
	userCollection    *mongo.Collection
	companyCollection *mongo.Collection
	JwtSecret         []byte
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type LoginResponse struct {
	Token string             `json:"token"`
	Id    primitive.ObjectID `json:"_id"`
	Role  string             `json:"role"`
}

var userCollection *mongo.Collection

func init() {
	client, ctx, err := dbHelper.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	userCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_CAREER"), client)
	companyCollection = dbHelper.GetCollection(ctx, os.Getenv("COLLECTION_COMPANY"), client)
}

// Generate token
func (a AuthService) GenerateToken(username string, id primitive.ObjectID, role string) (string, error) {
	expirationTime := time.Now().Add(100 * time.Minute)
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			Subject:   id.Hex(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.JwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validate the token
func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return a.JwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("Invalid token: " + err.Error())
	}
	return claims, nil
}

// check hashpassword
func (a *AuthService) CheckPasswordHash(hashedPassword, password string) bool {
	return hashedPassword == utils.EncodeToSHA(password)
}

// user authentication

func (a *AuthService) LoginForCareer(credential Credentials) (LoginResponse, error) {
	var career models.User

	err := userCollection.FindOne(context.Background(), bson.D{
		{"careerEmail", credential.Username},
		{"isDeleted", false},
	}).Decode(&career)

	if err != nil {
		return LoginResponse{}, errors.New("Invalid username or password")
	}

	if !a.CheckPasswordHash(career.Password, credential.Password) {
		return LoginResponse{}, errors.New("Invalid username or password")
	}
	token, _ := a.GenerateToken(career.CareerEmail, career.Id, career.Role)

	response := LoginResponse{
		Token: token,
		Id:    career.Id,
		Role:  career.Role,
	}

	return response, nil
}

func (a *AuthService) LoginForCompany(credential Credentials) (LoginResponse, error) {
	var company models.Company

	err := companyCollection.FindOne(context.Background(), bson.D{
		{"contact.companyEmail", credential.Username},
		{"isDeleted", false},
	}).Decode(&company)

	if err != nil {
		return LoginResponse{}, errors.New("Invalid username or password")
	}
	if !a.CheckPasswordHash(company.Password, credential.Password) {
		return LoginResponse{}, errors.New("Wrong password")
	}
	token, _ := a.GenerateToken(company.Contact.CompanyEmail, company.Id, "COMPANY")

	response := LoginResponse{
		Token: token,
		Id:    company.Id,
		Role:  "COMPANY",
	}
	return response, nil
}

func (a *AuthService) RegisterForCareer(user models.User) error {
	var existingUser models.User
	err := userCollection.FindOne(context.Background(), bson.D{
		{"careerEmail", user.CareerEmail},
	}).Decode(&existingUser)

	if err == nil {
		return errors.New("User with this email already exists")
	}

	hashedPassword := utils.EncodeToSHA(user.Password)
	user.Password = hashedPassword

	if user.Role == "" {
		user.Role = "CAREER"
	}

	user.Id = primitive.NewObjectID()
	user.CreateAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Error inserting user into MongoDB: %v", err)
		return errors.New("Error creating user")
	}

	return nil
}

func (a *AuthService) RegisterForCompany(company models.Company) error {
	var existingCompany models.Company

	err := companyCollection.FindOne(context.Background(), bson.D{
		{"companyEmail", company.Contact.CompanyEmail},
	}).Decode(&existingCompany)

	if err == nil {
		return errors.New("Company with this email already exists")
	}

	hashedPassword := utils.EncodeToSHA(company.Password)
	company.Password = hashedPassword

	company.Id = primitive.NewObjectID()
	company.CreateAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = companyCollection.InsertOne(context.Background(), company)
	if err != nil {
		log.Printf("Error inserting company into MongoDB: %v", err)
		return errors.New("Error creating company")
	}

	return nil
}
