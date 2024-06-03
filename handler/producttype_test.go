package handler_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/goccy/go-json"

	"github.com/Yoshikrit/fiber-test/handler"
	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/testutils"
	"github.com/Yoshikrit/fiber-test/helper/errs"
)

const (
	EndpointPath = "/producttype"
)

func TestCreate(t *testing.T) {
	mockService := testutils.NewProductTypeServiceMock()
	prodTypeHandler := handler.NewProductTypeHandler(mockService)
	
	app := fiber.New()
	app.Post(EndpointPath, prodTypeHandler.Create)

	prodTypeReqMock := &model.ProductTypeCreate{
		ID:   1,
		Name: "A",
	}
	
	prodTypeReqJSON, _ := json.Marshal(prodTypeReqMock)

	t.Run("test case : create success", func(t *testing.T) {
		mockService.On("Create", prodTypeReqMock).Return(nil)

		req := httptest.NewRequest(fiber.MethodPost, EndpointPath, strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusCreated, resp.StatusCode)

		expectedBody := `{"code":201,"message":"Create ProductType Successfully"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})

	t.Run("test case : create fail body parser", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, EndpointPath, strings.NewReader(`invalid json`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":"invalid character 'i' looking for beginning of value"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})

	t.Run("test case : create fail from service", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Create", prodTypeReqMock).Return(errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodPost, EndpointPath, strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})
}

func TestFindAll(t *testing.T) {
	mockService := testutils.NewProductTypeServiceMock()
	prodTypeHandler := handler.NewProductTypeHandler(mockService)
	
	app := fiber.New()
	app.Get(EndpointPath, prodTypeHandler.FindAll)

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
		mockService.On("FindAll").Return(prodTypesResMock, nil)

		req := httptest.NewRequest(fiber.MethodGet, EndpointPath, nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":` + string(prodTypesResJSON) + `}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})
	
	t.Run("test case : find all fail from service", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("FindAll").Return([]model.ProductType{}, errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodGet, EndpointPath, nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})
}

func TestFindByID(t *testing.T) {
	mockService := testutils.NewProductTypeServiceMock()
	prodTypeHandler := handler.NewProductTypeHandler(mockService)
	
	app := fiber.New()
	app.Get(EndpointPath  + "/:id", prodTypeHandler.FindByID)

	prodTypeResMock := model.ProductType {
		ID:   1,
		Name: "A",
	}

	prodTypeResJSON, _ := json.Marshal(prodTypeResMock)

	t.Run("test case : find by id success", func(t *testing.T) {
		mockService.On("FindByID", 1).Return(&prodTypeResMock, nil)

		req := httptest.NewRequest(fiber.MethodGet, EndpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":` + string(prodTypeResJSON) + `}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})

	t.Run("test case : find by id fail param", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, EndpointPath + "/a", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":"Invalid ID: a is not integer"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
	
	t.Run("test case : find by id fail from service", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("FindByID", 1).Return(&model.ProductType{}, errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodGet, EndpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockService := testutils.NewProductTypeServiceMock()
	prodTypeHandler := handler.NewProductTypeHandler(mockService)
	
	app := fiber.New()
	app.Put(EndpointPath  + "/:id", prodTypeHandler.Update)

	prodTypeReqMock := &model.ProductTypeUpdate {
		Name: "B",
	}

	prodTypeReqJSON, _ := json.Marshal(prodTypeReqMock)

	t.Run("test case : update success", func(t *testing.T) {
		mockService.On("Update", 1, prodTypeReqMock).Return(nil)

		req := httptest.NewRequest(fiber.MethodPut, EndpointPath + "/1", strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":"Update ProductType Successfully"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})

	t.Run("test case : update fail param", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPut, EndpointPath + "/a", strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":"Invalid ID: a is not integer"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})

	t.Run("test case : update fail body parser", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPut, EndpointPath + "/1", strings.NewReader(`invalid json`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":"invalid character 'i' looking for beginning of value"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})
	
	t.Run("test case : update fail from service", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Update", 1, prodTypeReqMock).Return(errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodPut, EndpointPath + "/1", strings.NewReader(string(prodTypeReqJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockService := testutils.NewProductTypeServiceMock()
	prodTypeHandler := handler.NewProductTypeHandler(mockService)
	
	app := fiber.New()
	app.Delete(EndpointPath  + "/:id", prodTypeHandler.Delete)

	t.Run("test case : delete success", func(t *testing.T) {
		mockService.On("Delete", 1).Return(nil)

		req := httptest.NewRequest(fiber.MethodDelete, EndpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":"Delete ProductType Successfully"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})

	t.Run("test case : delete fail param", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodDelete, EndpointPath + "/a", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":"Invalid ID: a is not integer"}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
	})
	
	t.Run("test case : delete fail from service", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Delete", 1).Return(errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodDelete, EndpointPath + "/1", nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})
}

func TestCount(t *testing.T) {
	mockService := testutils.NewProductTypeServiceMock()
	prodTypeHandler := handler.NewProductTypeHandler(mockService)
	
	app := fiber.New()
	app.Get(EndpointPath, prodTypeHandler.Count)

	t.Run("test case : get count success", func(t *testing.T) {
		mockService.On("Count").Return(int64(5), nil)

		req := httptest.NewRequest(fiber.MethodGet, EndpointPath, nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusOK, resp.StatusCode)

		expectedBody := `{"code":200,"message":5}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})
	
	t.Run("test case : get count fail from service", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Count").Return(int64(0), errs.NewInternalServerError(""))

		req := httptest.NewRequest(fiber.MethodGet, EndpointPath, nil)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		utils.AssertEqual(t, fiber.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":""}`
		body, _ := io.ReadAll(resp.Body)
		utils.AssertEqual(t, expectedBody, string(body))
		mockService.AssertExpectations(t)
	})
}