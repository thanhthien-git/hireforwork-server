package auth

import (
	"context"
	"errors"
	"hireforwork-server/config"
	"hireforwork-server/db"
	"hireforwork-server/models"
	"hireforwork-server/utils"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// instance for authservice
func NewAuthService(dbInstance *db.DB) *AuthService {
	collections := dbInstance.GetCollections([]string{"Career", "Company"})

	jwtSecret := []byte(config.GetInstance().SecretKey)

	if len(jwtSecret) == 0 {
		log.Fatalf("Need a secret key")
	}
	return &AuthService{
		userCollection:    collections[0],
		companyCollection: collections[1],
		JwtSecret:         jwtSecret,
	}

}

// Generate token
func (a AuthService) GenerateToken(username string, id primitive.ObjectID, role string) (string, error) {
	expirationTime := time.Now().Add(1000 * time.Minute)
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

// Login
func (a *AuthService) LoginForCareer(credential Credentials) (LoginResponse, error) {
	var career models.User

	err := a.userCollection.FindOne(context.Background(), bson.D{
		{"careerEmail", credential.Username},
		{"isDeleted", false},
	}).Decode(&career)
	if err != nil {
		return LoginResponse{}, errors.New("Tên đăng nhập hoặc tài khoản sai")
	}

	if !a.CheckPasswordHash(career.Password, credential.Password) {
		return LoginResponse{}, errors.New("Tên đăng nhập hoặc tài khoản sai")
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

	err := a.companyCollection.FindOne(context.Background(), bson.D{
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

// refactor login function
func (a *AuthService) Login(credential Credentials, config LoginConfig) (LoginResponse, error) {
	var entity interface{} = config.Model

	err := config.Collection.FindOne(context.Background(), bson.D{
		{config.UsernameField, credential.Username},
		{"isDeleted", false},
	}).Decode(entity)

	if err != nil {
		return LoginResponse{}, errors.New("Invalid username or password")
	}

	var password string
	var id primitive.ObjectID
	var username string

	//check user type for token creation
	switch v := entity.(type) {
	case *models.User:
		password = v.Password
		id = v.Id
		config.Role = v.Role
	case *models.Company:
		password = v.Password
		username = v.Contact.CompanyEmail
		config.Role = "COMPANY"
	default:
		return LoginResponse{}, errors.New("Unsupported entity type")
	}

	if !a.CheckPasswordHash(password, credential.Password) {
		return LoginResponse{}, errors.New("Invalid username or password")
	}

	//create token
	token, err := a.GenerateToken(username, id, config.Role)
	if err != nil {
		return LoginResponse{}, err
	}

	response := LoginResponse{
		Token: token,
		Id:    id,
		Role:  config.Role,
	}
	return response, nil
}
