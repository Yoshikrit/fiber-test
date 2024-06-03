package integration_test

import (
	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/repository"
	"github.com/Yoshikrit/fiber-test/service"
	"github.com/Yoshikrit/fiber-test/helper/errs"
	"github.com/Yoshikrit/fiber-test/testutils"

	"testing"
	"regexp"
	
	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateServiceRepository(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := repository.NewProductTypeRepositoryImpl(db)
	service := service.NewProductTypeServiceImpl(repo)

	t.Run("test case : create success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"Id", "Name"}).AddRow(1, "A")

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "producttype"`).
			WithArgs("A", 1).
			WillReturnRows(rows)
		mock.ExpectCommit()

		err := service.Create(&model.ProductTypeCreate{ID:1,Name:"A"})

		assert.NoError(t, err)
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
		err := service.Create(&model.ProductTypeCreate{ID:0,Name:"A"})

		expectedBody := valError
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
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
		err := service.Create(&model.ProductTypeCreate{ID:0,Name:""})

		expectedBody := valError
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
	})

	t.Run("test case : create fail conflict from repository", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A")
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnRows(rows)

		err := service.Create(&model.ProductTypeCreate{ID:1,Name:"A"})

		expectedBody := errs.NewConflictError("ProductType with this ID already exists")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
	})

	t.Run("test case : create fail from repository", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "producttype"`).
        	WithArgs("A", 1).
        	WillReturnError(errs.NewInternalServerError(""))
    	mock.ExpectRollback()

		err := service.Create(&model.ProductTypeCreate{ID:1,Name:"A"})

		expectedBody := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
	})
}

func TestFindAllServiceRepository(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := repository.NewProductTypeRepositoryImpl(db)
	service := service.NewProductTypeServiceImpl(repo)

	t.Run("test case : find all success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A").AddRow(2, "B")
		mock.ExpectQuery(`SELECT \* FROM "producttype"`).
			WillReturnRows(rows)

		prodTypeRes, err := service.FindAll()

		expectedBody := []model.ProductType{{ID:1,Name:"A",},{ID:2,Name:"B",}}
		assert.NoError(t, err)
		assert.Equal(t, expectedBody, prodTypeRes)
	})

	
	t.Run("test case : find all fail from repository", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "producttype"`).
			WillReturnError(errs.NewInternalServerError(""))

		prodTypesRes, err := service.FindAll()

		expectedBody := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		assert.Nil(t, prodTypesRes)
	})
}

func TestFindByIDServiceRepository(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := repository.NewProductTypeRepositoryImpl(db)
	service := service.NewProductTypeServiceImpl(repo)

	t.Run("test case : find by ID success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A")
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnRows(rows)

		prodTypeRes, err := service.FindByID(1)

		expectedBody := &model.ProductType{ID:1,Name:"A"}
		assert.NoError(t, err)
		assert.Equal(t, expectedBody, prodTypeRes)
	})

	t.Run("test case : find by ID fail not found from repository", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnError(gorm.ErrRecordNotFound)

		prodTypesRes, err := service.FindByID(1)

		expectedBody := errs.NewNotFoundError(recordNotFound)
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		assert.Nil(t, prodTypesRes)
	})
	
	
	t.Run("test case : find by ID fail from repository", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnError(errs.NewInternalServerError(""))

		prodTypesRes, err := service.FindByID(1)

		expectedBody := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		assert.Nil(t, prodTypesRes)
	})
}

func TestUpdateServiceRepository(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := repository.NewProductTypeRepositoryImpl(db)
	service := service.NewProductTypeServiceImpl(repo)

	t.Run("test case : update success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A")
		mock.ExpectQuery(`SELECT \* FROM "producttype"`).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE`).
			WithArgs(1, "B", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.Update(1, &model.ProductTypeUpdate{Name:"B"})

		assert.NoError(t, err)
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

		err := service.Update(1, &model.ProductTypeUpdate{Name:""})

		expectedBody := errs.ValErrorResponse(valError)
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
	})

	t.Run("test case : update fail not found from repository", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnError(gorm.ErrRecordNotFound)

		err := service.Update(1, &model.ProductTypeUpdate{Name:"B"})

		expectedBody := errs.NewNotFoundError(recordNotFound)
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
	})
	
	t.Run("test case : update fail from repository", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A")
		mock.ExpectQuery(`SELECT \* FROM "producttype"`).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE`).
        	WithArgs(1, "B", 1).
        	WillReturnError(errs.NewInternalServerError(""))
    	mock.ExpectRollback()
		
		err := service.Update(1, &model.ProductTypeUpdate{Name:"B"})

		expectedBody := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
	})
}

func TestDeleteServiceRepository(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := repository.NewProductTypeRepositoryImpl(db)
	service := service.NewProductTypeServiceImpl(repo)

	t.Run("test case : delete success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A")
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec("DELETE").
    		WithArgs(1).
    		WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.Delete(1)

		assert.NoError(t, err)
	})

	t.Run("test case : delete fail not found from repository", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnError(gorm.ErrRecordNotFound)
			
		err := service.Delete(1)

		expectedBody := errs.NewNotFoundError(recordNotFound)
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
	})

	t.Run("test case : delete fail from repository", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnError(errs.NewInternalServerError(""))

		mock.ExpectBegin()
		mock.ExpectExec("DELETE").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.Delete(1)

		expectedBody := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
	})
}

func TestCountServiceRepository(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := repository.NewProductTypeRepositoryImpl(db)
	service := service.NewProductTypeServiceImpl(repo)

	t.Run("test case : getcount success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "producttype"`)).
      		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		count, err := service.Count()

		expectedBody := int64(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedBody, count)
	})

	
	t.Run("test case : delete fail from repository", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "producttype"`)).
			WillReturnError(errs.NewInternalServerError(""))

		count, err := service.Count()

		expectedBody := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedBody, err)
		assert.Equal(t, int64(0), count)
	})
}

