package models

import (
	"github.com/google/uuid"
)

var uud = uuid.New()

//user Struct
type User struct {
	ID       string `json:"id" bson:"id,omitempty"`
	UserName string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

//token Struct
type RefreshToken struct {
	UserID       string `json:"userid" bson:"userid"`
	AccessToken  string `json:"accessToken" bson:"accessToken"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
	IsUsed       bool   `json:"isUsed" bson:"isUsed"`
}
