package testutils

import (
	"github.com/Yoshikrit/fiber-test/model"

	"github.com/stretchr/testify/mock"
)

type ProdTypeRepositoryMock struct {
	mock.Mock
}

func NewProductTypeRepositoryMock() *ProdTypeRepositoryMock {
	return &ProdTypeRepositoryMock{}
}

func (m *ProdTypeRepositoryMock) Save(prodTypeCreateReq *model.ProductTypeEntity) error {
	args := m.Called(prodTypeCreateReq)
	return args.Error(0)
}

func (m *ProdTypeRepositoryMock) FindAll() ([]model.ProductTypeEntity, error) {
	args := m.Called()
	return args.Get(0).([]model.ProductTypeEntity), args.Error(1)
}

func (m *ProdTypeRepositoryMock) FindByID(id int) (*model.ProductTypeEntity, error) {
	args := m.Called(id)
	return args.Get(0).(*model.ProductTypeEntity), args.Error(1)
}

func (m *ProdTypeRepositoryMock) Update(prodTypeUpdateReq *model.ProductTypeEntity) error {
	args := m.Called(prodTypeUpdateReq)
	return args.Error(0)
}

func (m *ProdTypeRepositoryMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ProdTypeRepositoryMock) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}