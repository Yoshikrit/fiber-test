package integration_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"regexp"
	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/goccy/go-json"

	"github.com/Yoshikrit/fiber-test/handler"
	"github.com/Yoshikrit/fiber-test/service"
	"github.com/Yoshikrit/fiber-test/repository"
	"github.com/Yoshikrit/fiber-test/testutils"
	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/helper/errs"
)

const (
	endpointPath = "/producttype"
	recordNotFound = "record not found"
)

func TestCreateHandlerServiceRepository(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	prodTypeRepository := repository.NewProductTypeRepositoryImpl(db)
	prodTypeService := service.NewProductTypeServiceImpl(prodTypeRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)
	
	app := fiber.New()
	app.Post(endpointPath, prodTypeHandler.Create)

	prodTypeReqMock := &model.ProductTypeCreate{
		ID:   1,
		Name: "A",
	}
	
	prodTypeReqError1Mock := &model.ProductTypeCreate{
		ID:   0,
		Name: "A",
	}
	prodTypeReqError2Mock := &model.ProductTypeCreate{
		ID:   1,
		Name: "",
	}
	prodTypeReqError3Mock := &model.ProductTypeCreate{
		ID:   0,
		Name: "",
	}
	
	prodTypeReqJSON, _ := json.Marshal(prodTypeReqMock)
	prodTypeReqError1JSON, _ := json.Marshal(prodTypeReqError1Mock)
	prodTypeReqError2JSON, _ := json.Marshal(prodTypeReqError2Mock)
	prodTypeReqError3JSON, _ := json.Marshal(prodTypeReqError3Mock)

	t.Run("test case : create success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"Id", "Name"}).AddRow(1, "A")

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "producttype"`).
			WithArgs("A", 1).
			WillReturnRows(rows)
		mock.ExpectCommit()

		req := httptest.NewRequest(fiber.MethodPost, endpointPath, strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusCreated, resp.StatusCode)

		expectedBody := `{"code":201,"message":"Create ProductType Successfully"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})

	t.Run("test case : create fail body parser", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, endpointPath, strings.NewReader(`invalid json`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":"invalid character 'i' looking for beginning of value"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})

	t.Run("test case : create fail validate from service no id", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, endpointPath, strings.NewReader(string(prodTypeReqError1JSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":[{"failed_field":"ProductTypeCreate.ID","tag":"required","value":""}]}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
	
	t.Run("test case : create fail validate from service no name", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, endpointPath, strings.NewReader(string(prodTypeReqError2JSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":[{"failed_field":"ProductTypeCreate.Name","tag":"required","value":""}]}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
	
	t.Run("test case : create fail validate from service no id and no name", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, endpointPath, strings.NewReader(string(prodTypeReqError3JSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":[{"failed_field":"ProductTypeCreate.ID","tag":"required","value":""},{"failed_field":"ProductTypeCreate.Name","tag":"required","value":""}]}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})

	t.Run("test case : create fail conflict from repository", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A")
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnRows(rows)

		req := httptest.NewRequest(fiber.MethodPost, endpointPath, strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusConflict, resp.StatusCode)

		expectedBody := `{"code":409,"message":"ProductType with this ID already exists"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})

	t.Run("test case : create fail from repository", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "producttype"`).
			WithArgs("A", 1).
			WillReturnError(errs.NewInternalServerError(""))
		mock.ExpectRollback()

		req := httptest.NewRequest(fiber.MethodPost, endpointPath, strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
}

func TestFindAllHandlerServiceRepository(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	prodTypeRepository := repository.NewProductTypeRepositoryImpl(db)
	prodTypeService := service.NewProductTypeServiceImpl(prodTypeRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)
	
	app := fiber.New()
	app.Get(endpointPath, prodTypeHandler.FindAll)

	prodTypesResMock := []model.ProductType {
		{
			ID:   1,
			Name: "A",
		},
		{
			ID:   2,
			Name: "B",
		},
	}

	prodTypesResJSON, _ := json.Marshal(prodTypesResMock)

	t.Run("test case : find all success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A").AddRow(2, "B")
		mock.ExpectQuery(`SELECT \* FROM "producttype"`).
			WillReturnRows(rows)

		req := httptest.NewRequest(fiber.MethodGet, endpointPath, nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":` + string(prodTypesResJSON) + `}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
	
	t.Run("test case : find all fail from repository", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "producttype"`).
			WillReturnError(errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodGet, endpointPath, nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
}

func TestFindByIDHandlerServiceRepository(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	prodTypeRepository := repository.NewProductTypeRepositoryImpl(db)
	prodTypeService := service.NewProductTypeServiceImpl(prodTypeRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)
	
	app := fiber.New()
	app.Get(endpointPath  + "/:id", prodTypeHandler.FindByID)

	prodTypeResMock := model.ProductType {
		ID:   1,
		Name: "A",
	}

	prodTypeResJSON, _ := json.Marshal(prodTypeResMock)

	t.Run("test case : find by id success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A")
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnRows(rows)

		req := httptest.NewRequest(fiber.MethodGet, endpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":` + string(prodTypeResJSON) + `}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})

	t.Run("test case : find by id fail param", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, endpointPath + "/a", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":"Invalid ID: a is not integer"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
	
	t.Run("test case : find by id fail record not found from repository", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest(fiber.MethodGet, endpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusNotFound, resp.StatusCode)

		expectedBody := `{"code":404,"message":"record not found"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
	
	t.Run("test case : find by id fail from repository", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnError(errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodGet, endpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
}

func TestUpdateHandlerServiceRepository(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	prodTypeRepository := repository.NewProductTypeRepositoryImpl(db)
	prodTypeService := service.NewProductTypeServiceImpl(prodTypeRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)
	
	app := fiber.New()
	app.Put(endpointPath  + "/:id", prodTypeHandler.Update)

	prodTypeReqMock := &model.ProductTypeUpdate {
		Name: "B",
	}
	prodTypeReqErrorMock := &model.ProductTypeUpdate {
		Name: "",
	}

	prodTypeReqJSON, _ := json.Marshal(prodTypeReqMock)
	prodTypeReqErrorJSON, _ := json.Marshal(prodTypeReqErrorMock)

	t.Run("test case : update success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A")
		mock.ExpectQuery(`SELECT \* FROM "producttype"`).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE`).
			WithArgs(1, "B", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest(fiber.MethodPut, endpointPath + "/1", strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":"Update ProductType Successfully"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})

	t.Run("test case : update fail param", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPut, endpointPath + "/a", strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":"Invalid ID: a is not integer"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})

	t.Run("test case : update fail body parser", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPut, endpointPath + "/1", strings.NewReader(`invalid json`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":"invalid character 'i' looking for beginning of value"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
	
	t.Run("test case : update fail validate from service no name", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPut, endpointPath + "/1", strings.NewReader(string(prodTypeReqErrorJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":[{"failed_field":"ProductTypeUpdate.Name","tag":"required","value":""}]}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
	
	t.Run("test case : update fail not found from repository", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest(fiber.MethodPut, endpointPath + "/1", strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusNotFound, resp.StatusCode)

		expectedBody := `{"code":404,"message":"record not found"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
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

		req := httptest.NewRequest(fiber.MethodPut, endpointPath + "/1", strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
}

func TestDeleteHandlerServiceRepository(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	prodTypeRepository := repository.NewProductTypeRepositoryImpl(db)
	prodTypeService := service.NewProductTypeServiceImpl(prodTypeRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)
	
	app := fiber.New()
	app.Delete(endpointPath  + "/:id", prodTypeHandler.Delete)

	t.Run("test case : delete success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "A")
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec("DELETE").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest(fiber.MethodDelete, endpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":"Delete ProductType Successfully"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})

	t.Run("test case : delete fail param", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodDelete, endpointPath + "/a", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":"Invalid ID: a is not integer"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
	
	t.Run("test case : delete fail not found from repository", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest(fiber.MethodDelete, endpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusNotFound, resp.StatusCode)

		expectedBody := `{"code":404,"message":"record not found"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
	
	t.Run("test case : delete fail from repository", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "producttype" WHERE`).
			WillReturnError(errs.NewInternalServerError(""))

		mock.ExpectBegin()
		mock.ExpectExec("DELETE").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest(fiber.MethodDelete, endpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
}

func TestCountHandlerServiceRepository(t *testing.T) {
	db, mock := testutils.SetupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	prodTypeRepository := repository.NewProductTypeRepositoryImpl(db)
	prodTypeService := service.NewProductTypeServiceImpl(prodTypeRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)
	
	app := fiber.New()
	app.Get(endpointPath, prodTypeHandler.Count)

	t.Run("test case : get count success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "producttype"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		req := httptest.NewRequest(fiber.MethodGet, endpointPath, nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":1}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
	
	t.Run("test case : get count fail from repository", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "producttype"`)).
			WillReturnError(errs.NewInternalServerError(""))
		
		req := httptest.NewRequest(fiber.MethodGet, endpointPath, nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
}