package service

import (
	"github.com/Yoshikrit/fiber-test/model"
)

type AuthService interface {
	Register(*model.UserCreate) error
	Login(*model.LoginRequest) (*model.UserPassport, error)
	RefreshPassport(*model.RefreshToken) (*model.UserPassport, error)
	Delete(int) (error)
}