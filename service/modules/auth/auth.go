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

type LoginStrategy interface {
	Login(credential Credentials) (LoginResponse, error)
}

type CareerLoginStrategy struct {
	authService *AuthService
}

type CompanyLoginStrategy struct {
	authService *AuthService
}

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
		Id:       id.Hex(),
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
func (c *CareerLoginStrategy) Login(credential Credentials) (LoginResponse, error) {
	var career models.User

	err := c.authService.userCollection.FindOne(context.Background(), bson.D{
		{"careerEmail", credential.Username},
		{"isDeleted", false},
	}).Decode(&career)
	if err != nil {
		return LoginResponse{}, errors.New("Tên đăng nhập hoặc tài khoản sai")
	}

	if !c.authService.CheckPasswordHash(career.Password, credential.Password) {
		return LoginResponse{}, errors.New("Tên đăng nhập hoặc tài khoản sai")
	}
	token, _ := c.authService.GenerateToken(career.CareerEmail, career.Id, career.Role)

	response := LoginResponse{
		Token: token,
	}

	return response, nil
}

func (co *CompanyLoginStrategy) Login(credential Credentials) (LoginResponse, error) {
	var company models.Company

	err := co.authService.companyCollection.FindOne(context.Background(), bson.D{
		{"contact.companyEmail", credential.Username},
		{"isDeleted", false},
	}).Decode(&company)

	if err != nil {
		return LoginResponse{}, errors.New("Invalid username or password")
	}
	if !co.authService.CheckPasswordHash(company.Password, credential.Password) {
		return LoginResponse{}, errors.New("Wrong password")
	}
	token, _ := co.authService.GenerateToken(company.Contact.CompanyEmail, company.Id, "COMPANY")

	response := LoginResponse{
		Token: token,
	}
	return response, nil
}

// Login sử dụng Strategy Pattern
func (a *AuthService) Login(credential Credentials, strategy LoginStrategy) (LoginResponse, error) {
	return strategy.Login(credential)
}

// Hàm tạo Strategy
func NewCareerLoginStrategy(authService *AuthService) *CareerLoginStrategy {
	return &CareerLoginStrategy{authService: authService}
}

func NewCompanyLoginStrategy(authService *AuthService) *CompanyLoginStrategy {
	return &CompanyLoginStrategy{authService: authService}
}
