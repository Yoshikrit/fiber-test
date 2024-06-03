package service_test

import (
	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/service"
	"github.com/Yoshikrit/fiber-test/helper/errs"
	"github.com/Yoshikrit/fiber-test/testutils"
	
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Run("test case : create success", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{}, nil)
		mockRepository.On("Save", &model.ProductTypeEntity{ID:1,Name:"A"}).Return(nil)

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Create(&model.ProductTypeCreate{ID:1,Name:"A"})

		assert.NoError(t, err)
		mockRepository.AssertExpectations(t)
	})

	t.Run("test case : create fail validate no id", func(t *testing.T) {
		valError := errs.ValErrorResponse{
			Code: 400,
			Message: []errs.ErrorMessage{
			  	{
					FailedField: "ProductTypeCreate.ID",
					Tag:        "required",
					Value:      "",
			  	},
			},
		}
		mockRepository := testutils.NewProductTypeRepositoryMock()

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Create(&model.ProductTypeCreate{ID:0,Name:"A"})

		expectedBody := valError
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		mockRepository.AssertExpectations(t)
	})

	t.Run("test case : create fail validate no name", func(t *testing.T) {
		valError := errs.ValErrorResponse{
			Code: 400,
			Message: []errs.ErrorMessage{
			  	{
					FailedField: "ProductTypeCreate.Name",
					Tag:        "required",
					Value:      "",
			  	},
			},
		}
		mockRepository := testutils.NewProductTypeRepositoryMock()

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Create(&model.ProductTypeCreate{ID:1,Name:""})

		expectedBody := valError
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
	})

	

	t.Run("test case : create fail validate no id and name", func(t *testing.T) {
		valError := errs.ValErrorResponse{
			Code: 400,
			Message: []errs.ErrorMessage{
				{
				  FailedField: "ProductTypeCreate.ID",
				  Tag:        "required",
				  Value:      "",
				},
			  	{
					FailedField: "ProductTypeCreate.Name",
					Tag:        "required",
					Value:      "",
			  	},
			},
		}

		mockRepository := testutils.NewProductTypeRepositoryMock()

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Create(&model.ProductTypeCreate{ID:0,Name:""})

		expectedBody := valError
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
	})

	t.Run("test case : create fail conflict", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{ID:1,Name:"A"}, nil)

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Create(&model.ProductTypeCreate{ID:1,Name:"A"})

		expectedBody := errs.NewConflictError("ProductType with this ID already exists")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		mockRepository.AssertExpectations(t)
	})

	t.Run("test case : create fail from repository", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{}, nil)
		mockRepository.On("Save", &model.ProductTypeEntity{ID:1,Name:"A"}).Return(errs.NewInternalServerError(""))

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Create(&model.ProductTypeCreate{ID:1,Name:"A"})

		expectedBody := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		mockRepository.AssertExpectations(t)
	})
}

func TestFindAll(t *testing.T) {
	t.Run("test case : find all success", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindAll").Return([]model.ProductTypeEntity{{ID:1,Name:"A",},{ID:2,Name:"B",}}, nil)

		service := service.NewProductTypeServiceImpl(mockRepository)
		prodTypeRes, err := service.FindAll()

		expectedBody := []model.ProductType{{ID:1,Name:"A",},{ID:2,Name:"B",}}
		assert.NoError(t, err)
		assert.Equal(t, expectedBody, prodTypeRes)
		mockRepository.AssertExpectations(t)
	})

	
	t.Run("test case : find all fail from repository", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindAll").Return([]model.ProductTypeEntity{}, errs.NewInternalServerError(""))

		service := service.NewProductTypeServiceImpl(mockRepository)
		prodTypesRes, err := service.FindAll()

		expectedBody := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		assert.Nil(t, prodTypesRes)
		mockRepository.AssertExpectations(t)
	})
}

func TestFindByID(t *testing.T) {
	t.Run("test case : find by ID success", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{ID:1,Name:"A"}, nil)

		service := service.NewProductTypeServiceImpl(mockRepository)
		prodTypeRes, err := service.FindByID(1)

		expectedBody := &model.ProductType{ID:1,Name:"A"}
		assert.NoError(t, err)
		assert.Equal(t, expectedBody, prodTypeRes)
		mockRepository.AssertExpectations(t)
	})

	
	t.Run("test case : find by ID fail from repository", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{}, errs.NewInternalServerError(""))

		service := service.NewProductTypeServiceImpl(mockRepository)
		prodTypesRes, err := service.FindByID(1)

		expectedBody := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		assert.Nil(t, prodTypesRes)
		mockRepository.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("test case : update success", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{ID:1,Name:"A"}, nil)
		mockRepository.On("Update", &model.ProductTypeEntity{ID:1,Name:"B"}).Return(nil)

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Update(1, &model.ProductTypeUpdate{Name:"B"})

		assert.NoError(t, err)
		mockRepository.AssertExpectations(t)
	})

	t.Run("test case : update fail validate no name", func(t *testing.T) {
		valError := errs.ValErrorResponse{
			Code: 400,
			Message: []errs.ErrorMessage{
			  	{
					FailedField: "ProductTypeUpdate.Name",
					Tag:        "required",
					Value:      "",
			  	},
			},
		}

		mockRepository := testutils.NewProductTypeRepositoryMock()

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Update(1, &model.ProductTypeUpdate{Name:""})

		expectedBody := errs.ValErrorResponse(valError)
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
	})

	t.Run("test case : update fail not found from repository", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{}, errs.NewNotFoundError(""))

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Update(1, &model.ProductTypeUpdate{Name:"B"})

		expectedBody := errs.NewNotFoundError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		mockRepository.AssertExpectations(t)
	})
	
	t.Run("test case : update fail from repository", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{ID:1,Name:"A"}, nil)
		mockRepository.On("Update", &model.ProductTypeEntity{ID:1,Name:"B"}).Return(errs.NewNotFoundError(""))

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Update(1, &model.ProductTypeUpdate{Name:"B"})

		expectedBody := errs.NewNotFoundError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		mockRepository.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	t.Run("test case : delete success", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{ID:1,Name:"A"}, nil)
		mockRepository.On("Delete", 1).Return(nil)

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Delete(1)

		assert.NoError(t, err)
		mockRepository.AssertExpectations(t)
	})

	t.Run("test case : delete fail not found from repository", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{}, errs.NewNotFoundError(""))

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Delete(1)

		expectedBody := errs.NewNotFoundError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		mockRepository.AssertExpectations(t)
	})
	
	t.Run("test case : delete fail from repository", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{ID:1,Name:"A"}, nil)
		mockRepository.On("Delete", 1).Return(errs.NewInternalServerError(""))

		service := service.NewProductTypeServiceImpl(mockRepository)
		err := service.Delete(1)

		expectedBody := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		mockRepository.AssertExpectations(t)
	})
}

func TestCount(t *testing.T) {
	t.Run("test case : getcount success", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("Count").Return(int64(1), nil)

		service := service.NewProductTypeServiceImpl(mockRepository)
		count, err := service.Count()

		expectedBody := int64(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedBody, count)
		mockRepository.AssertExpectations(t)
	})

	
	t.Run("test case : delete fail from repository", func(t *testing.T) {
		mockRepository := testutils.NewProductTypeRepositoryMock()
		mockRepository.On("Count").Return(int64(0), errs.NewInternalServerError(""))

		service := service.NewProductTypeServiceImpl(mockRepository)
		count, err := service.Count()

		expectedBody := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		assert.Equal(t, int64(0), count)
		mockRepository.AssertExpectations(t)
	})
}

