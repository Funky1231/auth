package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"../utils"
	"go.mongodb.org/mongo-driver/bson"
)

var jwtKey = []byte("my_secret_key")

func DeleteAllTokenUser(w http.ResponseWriter, r *http.Request) {
	var err error
	token := utils.GetHeaderToken(w, r)
	if len(token) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	isValid := utils.VerifyToken(token, jwtKey)
	if isValid != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId, err := utils.DecodeTokenAndGetUserId(token)
	var body struct {
		RefreshToken string
	}
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	filter := bson.M{"userid": userId}
	deleteResult, err := connect.Collection("token").DeleteMany(context.TODO(), filter)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("ok")
}
