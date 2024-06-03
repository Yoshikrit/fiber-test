package repository_test

import (
	"testing"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"

	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/repository"
	"github.com/Yoshikrit/fiber-test/helper/errs"
	"github.com/Yoshikrit/fiber-test/testutils"
)

func TestCreate(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	prodTypeCreateMock := &model.ProductTypeEntity{
		ID:   1,
		Name: "A",
	}

	t.Run("test case : create producttype success", func(t *testing.T) {
		repo := repository.NewProductTypeRepositoryImpl(db)
		rows := sqlmock.NewRows([]string{"Id", "Name"}).AddRow(1, "A")

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "producttype"`).
			WithArgs("A", 1).
			WillReturnRows(rows)
		mock.ExpectCommit()

		err := repo.Save(prodTypeCreateMock)

		assert.NoError(t, err)
		assert.Nil(t, err)
	})

	t.Run("test case : create producttype fail", func(t *testing.T) {
		repo := repository.NewProductTypeRepositoryImpl(db)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "producttype"`).
        	WithArgs("A", 1).
        	WillReturnError(errs.NewInternalServerError(""))
    	mock.ExpectRollback()

		err := repo.Save(prodTypeCreateMock)

		expectedRes := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedRes, err)
	})
}

func TestFindAll(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	entityRes := []model.ProductTypeEntity{
		{
			ID:   1,
			Name: "A",
		},
		{
			ID:   2,
			Name: "B",
		},
	}
	t.Run("test case : find all producttype success", func(t *testing.T) {
		repo := repository.NewProductTypeRepositoryImpl(db)

		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A").AddRow(2, "B")
		mock.ExpectQuery(`SELECT \* FROM "producttype"`).
			WillReturnRows(rows)

		result, err := repo.FindAll()

		expectedRes := entityRes
		assert.NoError(t, err)
		assert.Equal(t, expectedRes, result)
	})
	t.Run("test case : find all producttype fail", func(t *testing.T) {
		repo := repository.NewProductTypeRepositoryImpl(db)

		mock.ExpectQuery(`SELECT \* FROM "producttype"`).
			WillReturnError(errs.NewInternalServerError(""))

		_, err := repo.FindAll()

		expectedRes := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedRes, err)
	})
}

func TestFindByID(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()
	t.Run("test case : find producttype by id pass", func(t *testing.T) {
		entityRes := &model.ProductTypeEntity{
			ID:   1,
			Name: "A",
		}

		repo := repository.NewProductTypeRepositoryImpl(db)
		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A")

		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnRows(rows)

		result, err := repo.FindByID(1)

		expectedRes := entityRes
		assert.NoError(t, err)
		assert.Equal(t, expectedRes, result)
	})
	t.Run("test case : get fail gorm not found", func(t *testing.T) {
		repo := repository.NewProductTypeRepositoryImpl(db)

		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.FindByID(1)

		expectedRes := errs.NewNotFoundError("record not found")
		assert.Error(t, err)
		assert.Equal(t, expectedRes, err)
	})

	t.Run("test case : get fail get id", func(t *testing.T) {
		repo := repository.NewProductTypeRepositoryImpl(db)

		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnError(errs.NewInternalServerError(""))

		_, err := repo.FindByID(1)

		expectedRes := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedRes, err)
	})
}

func TestUpdate(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	prodTypeEntityMock := &model.ProductTypeEntity{
		ID:   1,
		Name: "A",
	}

	t.Run("test case : update producttype success", func(t *testing.T) {
		repo := repository.NewProductTypeRepositoryImpl(db)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE`).
			WithArgs(1, "A", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(prodTypeEntityMock)

		assert.NoError(t, err)
		assert.Nil(t, err)
	})

	t.Run("test case : update producttype fail", func(t *testing.T) {
		repo := repository.NewProductTypeRepositoryImpl(db)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE`).
        	WithArgs(1, "A", 1).
        	WillReturnError(errs.NewInternalServerError(""))
    	mock.ExpectRollback()

		err := repo.Update(prodTypeEntityMock)

		expectedRes := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedRes, err)
	})
}

func TestDelete(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	t.Run("test case : delete producttype success", func(t *testing.T) {
		repo := repository.NewProductTypeRepositoryImpl(db)

		mock.ExpectBegin()
		mock.ExpectExec("DELETE").
    		WithArgs(1).
    		WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(1)

		assert.NoError(t, err)
		assert.Nil(t, err)
	})
	t.Run("test case : delete producttype fail", func(t *testing.T) {
		repo := repository.NewProductTypeRepositoryImpl(db)

		mock.ExpectBegin()
		mock.ExpectExec("DELETE").
    		WithArgs(1).
    		WillReturnError(errs.NewInternalServerError(""))
		mock.ExpectRollback()

		err := repo.Delete(1)

		expectedRes := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedRes, err)
	})
}

func TestCount(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	t.Run("test case : get count pass", func(t *testing.T) {
		repo := repository.NewProductTypeRepositoryImpl(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "producttype"`)).
      		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		result, err := repo.Count()

		expectedRes := int64(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedRes, result)
	})
	t.Run("test case : get count fail", func(t *testing.T) {
		repo := repository.NewProductTypeRepositoryImpl(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "producttype"`)).
			WillReturnError(errs.NewInternalServerError(""))

		_, err := repo.Count()

		expectedRes := errs.NewInternalServerError("")
		assert.Error(t, err)
		assert.Equal(t, expectedRes, err)
	})
}
