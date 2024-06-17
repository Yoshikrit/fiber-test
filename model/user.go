package model

import (
)

type UserEntity struct {
	ID   			int    		`gorm:"primaryKey; column:user_id;"`
	RoleID   		int    		`gorm:"not null;   column:user_role_id;"`
	Name 			string 		`gorm:"not null;   column:user_name;     size:40;"`
	Email 			string 		`gorm:"not null;   column:user_email;    size:50;  unique;"`
	Password 		string 		`gorm:"not null;   column:user_password;"`
}

func (u UserEntity) TableName() string {
	return "user"
}

type UserDTO struct {
	ID   		int    
	RoleID 		int 
	Name 		string 
	Email 		string 
}

type UserCreate struct {
    ID     		int    	`json:"user_id"         validate:"required,gt=0"`
	RoleID      int 	`json:"role_id"         validate:"required,gt=0"`
	Name   		string  `json:"user_name"       validate:"required,max=40"`
	Email   	string  `json:"user_email"      validate:"required,email,max=50"`
	Password 	string 	`json:"user_password"   validate:"required,max=255"`
}

type LoginRequest struct {
	Email   	string    `json:"user_email"       validate:"required,email,max=40"`
	Password 	string 	  `json:"user_password"    validate:"required,max=255"`
}

type UserClaims struct {
	ID     		int
	RoleID      int 
}

type UserPassport struct {
	User   	*UserDTO    `json:"user"`
	Tokens 	*UserToken  `json:"token"`
}