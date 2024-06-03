package service

import (
	"github.com/Yoshikrit/fiber-test/model"
)

type ProductTypeService interface {
	Create(*model.ProductTypeCreate) error
	FindAll() ([]model.ProductType, error)
	FindByID(int) (*model.ProductType, error)
	Update(int, *model.ProductTypeUpdate) error
	Delete(int) (error)
	Count() (int64, error)
}