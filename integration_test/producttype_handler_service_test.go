package integration_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/goccy/go-json"

	"github.com/Yoshikrit/fiber-test/handler"
	"github.com/Yoshikrit/fiber-test/service"
	"github.com/Yoshikrit/fiber-test/testutils"
	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/helper/errs"
)

func TestCreateHandlerService(t *testing.T) {
	mockRepository := testutils.NewProductTypeRepositoryMock()
	prodTypeService := service.NewProductTypeServiceImpl(mockRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)
	
	app := fiber.New()
	app.Post(endpointPath, prodTypeHandler.Create)

	prodTypeReqMock := &model.ProductTypeCreate{
		ID:   1,
		Name: "A",
	}
	
	prodTypeReqJSON, _ := json.Marshal(prodTypeReqMock)

	t.Run("test case : create success", func(t *testing.T) {
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{}, nil)
		mockRepository.On("Save", &model.ProductTypeEntity{ID:1,Name:"A"}).Return(nil)

		req := httptest.NewRequest(fiber.MethodPost, endpointPath, strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusCreated, resp.StatusCode)

		expectedBody := `{"code":201,"message":"Create ProductType Successfully"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockRepository.AssertExpectations(t)
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
		mockRepository.AssertExpectations(t)
	})

	t.Run("test case : create fail from repository", func(t *testing.T) {
		mockRepository.ExpectedCalls = nil
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{}, nil)
		mockRepository.On("Save", &model.ProductTypeEntity{ID:1,Name:"A"}).Return(errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodPost, endpointPath, strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockRepository.AssertExpectations(t)
	})
}

func TestFindAllHandlerService(t *testing.T) {
	mockRepository := testutils.NewProductTypeRepositoryMock()
	prodTypeService := service.NewProductTypeServiceImpl(mockRepository)
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
		mockRepository.On("FindAll").Return([]model.ProductTypeEntity{{ID:1,Name:"A",},{ID:2,Name:"B",}}, nil)

		req := httptest.NewRequest(fiber.MethodGet, endpointPath, nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":` + string(prodTypesResJSON) + `}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockRepository.AssertExpectations(t)
	})
	
	t.Run("test case : find all fail from repository", func(t *testing.T) {
		mockRepository.ExpectedCalls = nil
		mockRepository.On("FindAll").Return([]model.ProductTypeEntity{}, errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodGet, endpointPath, nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockRepository.AssertExpectations(t)
	})
}

func TestFindByIDHandlerService(t *testing.T) {
	mockRepository := testutils.NewProductTypeRepositoryMock()
	prodTypeService := service.NewProductTypeServiceImpl(mockRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)
	
	app := fiber.New()
	app.Get(endpointPath  + "/:id", prodTypeHandler.FindByID)

	prodTypeResMock := model.ProductType {
		ID:   1,
		Name: "A",
	}

	prodTypeResJSON, _ := json.Marshal(prodTypeResMock)

	t.Run("test case : find by id success", func(t *testing.T) {
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{ID:1,Name:"A"}, nil)

		req := httptest.NewRequest(fiber.MethodGet, endpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":` + string(prodTypeResJSON) + `}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockRepository.AssertExpectations(t)
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
	
	t.Run("test case : find by id fail from repository", func(t *testing.T) {
		mockRepository.ExpectedCalls = nil
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{}, errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodGet, endpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockRepository.AssertExpectations(t)
	})
}

func TestUpdateHandlerService(t *testing.T) {
	mockRepository := testutils.NewProductTypeRepositoryMock()
	prodTypeService := service.NewProductTypeServiceImpl(mockRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)
	
	app := fiber.New()
	app.Put(endpointPath  + "/:id", prodTypeHandler.Update)

	prodTypeReqMock := &model.ProductTypeUpdate {
		Name: "B",
	}

	prodTypeReqJSON, _ := json.Marshal(prodTypeReqMock)

	t.Run("test case : update success", func(t *testing.T) {
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{ID:1,Name:"A"}, nil)
		mockRepository.On("Update", &model.ProductTypeEntity{ID:1,Name:"B"}).Return(nil)

		req := httptest.NewRequest(fiber.MethodPut, endpointPath + "/1", strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":"Update ProductType Successfully"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockRepository.AssertExpectations(t)
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
		mockRepository.AssertExpectations(t)
	})
	
	t.Run("test case : update fail from repository", func(t *testing.T) {
		mockRepository.ExpectedCalls = nil
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{ID:1,Name:"A"}, nil)
		mockRepository.On("Update", &model.ProductTypeEntity{ID:1,Name:"B"}).Return(errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodPut, endpointPath + "/1", strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockRepository.AssertExpectations(t)
	})
}

func TestDeleteHandlerService(t *testing.T) {
	mockRepository := testutils.NewProductTypeRepositoryMock()
	prodTypeService := service.NewProductTypeServiceImpl(mockRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)
	
	app := fiber.New()
	app.Delete(endpointPath  + "/:id", prodTypeHandler.Delete)

	t.Run("test case : delete success", func(t *testing.T) {
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{ID:1,Name:"A"}, nil)
		mockRepository.On("Delete", 1).Return(nil)

		req := httptest.NewRequest(fiber.MethodDelete, endpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":"Delete ProductType Successfully"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockRepository.AssertExpectations(t)
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
	
	t.Run("test case : delete fail from repository", func(t *testing.T) {
		mockRepository.ExpectedCalls = nil
		mockRepository.On("FindByID", 1).Return(&model.ProductTypeEntity{ID:1,Name:"A"}, nil)
		mockRepository.On("Delete", 1).Return(errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodDelete, endpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockRepository.AssertExpectations(t)
	})
}

func TestCountHandlerService(t *testing.T) {
	mockRepository := testutils.NewProductTypeRepositoryMock()
	prodTypeService := service.NewProductTypeServiceImpl(mockRepository)
	prodTypeHandler := handler.NewProductTypeHandler(prodTypeService)
	
	app := fiber.New()
	app.Get(endpointPath, prodTypeHandler.Count)

	t.Run("test case : get count success", func(t *testing.T) {
		mockRepository.On("Count").Return(int64(5), nil)

		req := httptest.NewRequest(fiber.MethodGet, endpointPath, nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":5}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockRepository.AssertExpectations(t)
	})
	
	t.Run("test case : get count fail from repository", func(t *testing.T) {
		mockRepository.ExpectedCalls = nil
		mockRepository.On("Count").Return(int64(0), errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodGet, endpointPath, nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockRepository.AssertExpectations(t)
	})
}