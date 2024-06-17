package repository

import (
	"github.com/Yoshikrit/fiber-test/model"
)

type UserRepository interface {
	Create(userCreateReq *model.UserEntity) error
	FindByID(id int) (*model.UserEntity, error)
	FindByEmail(email string) (*model.UserEntity, error)
}