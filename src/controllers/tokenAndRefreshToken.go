package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Funky1231/auth/connectDB"
	"github.com/Funky1231/auth/models"
	"github.com/Funky1231/auth/utils"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var connect = connectDB.ConnectDB()

func Signin(w http.ResponseWriter, r *http.Request) {
	var err error

	paramId, ok := r.URL.Query()["id"]

	if !ok || len(paramId[0]) < 1 {
		log.Println("Url Param '_id' is missing")
		return
	}

	_id, err := primitive.ObjectIDFromHex(paramId[0])
	var dbUser = models.User{}
	filter := bson.D{primitive.E{Key: "_id", Value: _id}}
	err = connect.Collection("user").FindOne(context.TODO(), filter).Decode(&dbUser)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(51231 * time.Minute)

	claims := utils.Claims{
		UserId: dbUser.ID,
		StandardClaims: jwt.StandardClaims{
			// время истечения токена в милисекундах
			ExpiresAt: expirationTime.Unix(),
		},
	}

	tokenString, refreshToken, err := utils.CreateTokenAndRefreshToken(&claims, jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	hash, err := utils.HashRefreshToken(refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = connect.Client().UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			fmt.Println(err)
			return err
		}

		var ref = models.RefreshToken{RefreshToken: hash, AccessToken: tokenString, UserID: dbUser.ID, IsUsed: false}
		_, err = connect.Collection("token").InsertOne(sessionContext, &ref)
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			fmt.Println(err)
			return err
		} else {
			fmt.Println(sessionContext)
			sessionContext.CommitTransaction(sessionContext)
		}
		return nil
	})

	w.Header().Set("Content-Type", "application/json")

	type Body struct {
		AccessToken  string
		Refreshtoken string
	}
	body := Body{tokenString, hash}
	json.NewEncoder(w).Encode(body)
}
