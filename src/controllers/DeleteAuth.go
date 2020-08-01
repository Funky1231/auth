package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"../utils"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteTokenUser(w http.ResponseWriter, r *http.Request) {
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
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	var body struct {
		RefreshToken string
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	filter := bson.M{"userid": userId, "isUsed": false}

	deleteResult, err := connect.Collection("token").DeleteOne(context.TODO(), filter)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("ok")
}
