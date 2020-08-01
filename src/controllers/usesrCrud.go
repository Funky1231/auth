package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Funky1231/auth/connectDB"
	"github.com/Funky1231/auth/models"
	"github.com/google/uuid"
)

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	uuid := uuid.New()
	newUserIdString := uuid.String()
	var user = models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Print(err)
	}
	user.ID = newUserIdString

	collection := connectDB.ConnectDB()
	insertResult, err := collection.Collection("user").InsertOne(context.TODO(), &user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("Found a single document: %+vn", user)

	json.NewEncoder(w).Encode(insertResult.InsertedID)
}
