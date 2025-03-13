package auth

import (
	"github.com/dgrijalva/jwt-go"
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

type LoginConfig struct {
	Collection    *mongo.Collection
	UsernameField string
	Role          string
	Model         interface{}
}
