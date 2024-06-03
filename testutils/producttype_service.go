package testutils

import (
	"github.com/Yoshikrit/fiber-test/model"

	"github.com/stretchr/testify/mock"
)

type ProdTypeServiceMock struct {
	mock.Mock
}

func NewProductTypeServiceMock() *ProdTypeServiceMock {
	return &ProdTypeServiceMock{}
}

func (m *ProdTypeServiceMock) Create(prodTypeCreateReq *model.ProductTypeCreate) error {
	args := m.Called(prodTypeCreateReq)
	return args.Error(0)
}

func (m *ProdTypeServiceMock) FindAll() ([]model.ProductType, error) {
	args := m.Called()
	return args.Get(0).([]model.ProductType), args.Error(1)
}

func (m *ProdTypeServiceMock) FindByID(id int) (*model.ProductType, error) {
	args := m.Called(id)
	return args.Get(0).(*model.ProductType), args.Error(1)
}

func (m *ProdTypeServiceMock) Update(id int, prodTypeUpdateReq *model.ProductTypeUpdate) error {
	args := m.Called(id, prodTypeUpdateReq)
	return args.Error(0)
}

func (m *ProdTypeServiceMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ProdTypeServiceMock) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}