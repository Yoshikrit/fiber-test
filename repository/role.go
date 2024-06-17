package repository

import (
	"github.com/Yoshikrit/fiber-test/model"
)

type RoleRepository interface {
	FindByID(id int) (*model.RoleEntity, error)
}