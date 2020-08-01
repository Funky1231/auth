package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Funky1231/auth/models"
	"github.com/Funky1231/auth/utils"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	var err error

	token := utils.GetHeaderToken(w, r)
	if len(token) == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	isValid := utils.VerifyToken(token, jwtKey)
	if isValid != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId, _ := utils.DecodeTokenAndGetUserId(token)
	if err != nil {
		log.Println(userId)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var body struct {
		RefreshToken string
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	var dbRef = models.RefreshToken{}
	filter := bson.D{primitive.E{Key: "userid", Value: userId},
		{Key: "refreshToken", Value: body.RefreshToken}, {Key: "accessToken", Value: token}, {Key: "isUsed", Value: false}}

	err = connect.Collection("token").FindOne(context.TODO(), filter).Decode(&dbRef)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	match := utils.CheckRefreshTokenHash(dbRef.RefreshToken, body.RefreshToken)
	if match != false {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	_, err = connect.Collection("token").UpdateMany(context.TODO(), bson.M{"userid": userId}, bson.M{"$set": bson.M{"isUsed": true}})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	// время жизни токена (5 мин)
	expirationTime := time.Now().Add(5123132 * time.Minute)

	// создаем притензию которая включает в себя имя пользователя и время истечения jwt
	claims := utils.Claims{
		UserId: userId,
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
	hash, _ := utils.HashRefreshToken(refreshToken)

	var ref = models.RefreshToken{RefreshToken: hash, AccessToken: tokenString, UserID: userId}
	_, err = connect.Collection("token").InsertOne(context.TODO(), &ref)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	fmt.Println(tokenString, refreshToken)
	type Resopnse struct {
		AccessToken  string
		Refreshtoken string
	}

	res := Resopnse{tokenString, hash}
	json.NewEncoder(w).Encode(res)
}
