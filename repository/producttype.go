package repository

import (
	"github.com/Yoshikrit/fiber-test/model"
)

type ProductTypeRepository interface {
	Save(*model.ProductTypeEntity) error
	FindAll() ([]model.ProductTypeEntity, error)
	FindByID(int) (*model.ProductTypeEntity, error)
	Update(*model.ProductTypeEntity) error
	Delete(int) error
	Count() (int64, error)
}

