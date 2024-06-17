package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type OauthEntity struct {
	ID   			int    `gorm:"primaryKey; column:oauth_id;"`
	UserID   		int    `gorm:"not null;   column:oauth_user_id;"`
	AccessToken 	string `gorm:"not null;   column:access_token;"`
	RefreshToken 	string `gorm:"not null;   column:reflesh_token;"`
}

func (o OauthEntity) TableName() string {
	return "oauth"
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type UserToken struct {
	ID   			int `json:"oauth_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ServiceMapClaims struct {
	Claims *UserClaims `json:"claims"`
	jwt.RegisteredClaims
}