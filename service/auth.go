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

