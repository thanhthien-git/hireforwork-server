package utils

import (
	"encoding/hex"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/sha3"
)

func GetQueryID(r *http.Request) primitive.ObjectID {
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}

	return objectID
}

func EncodeToSHA(password string) string {
	data := []byte(password)

	hash := sha3.New256()
	hash.Write(data)
	hashedPassword := hex.EncodeToString(hash.Sum(nil))

	return hashedPassword
}
