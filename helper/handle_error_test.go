package helper_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/gofiber/fiber/v2"

	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/helper"
	"github.com/Yoshikrit/fiber-test/helper/errs"

	"net/http"
	"net/http/httptest"
	"testing"
	"errors"
	"io"
)

func TestHandleError(t *testing.T) {
	app := fiber.New()

	app.Get("/error", func(ctx *fiber.Ctx) error {
		err := ctx.Query("err")
		switch err {
		case "app":
			return helper.HandleError(ctx, errs.NewBadRequestError("Bad Request"))
		case "val":
			return helper.HandleError(ctx, errs.NewValidateBadRequestError([]errs.ErrorMessage{}))
		case "generic":
			return helper.HandleError(ctx, errors.New("example generic error"))
		default:
			return ctx.SendStatus(http.StatusOK)
		}
	})

	t.Run("test case : is AppError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/error?err=app", nil)
		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":"Bad Request"}`
		body, _ := io.ReadAll(resp.Body)
		assert.JSONEq(t, expectedBody, string(body))
	})

	t.Run("test case : is ValError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/error?err=val", nil)
		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		expectedBody := `{"code":400,"message":[]}`
		body, _ := io.ReadAll(resp.Body)
		assert.JSONEq(t, expectedBody, string(body))
	})

	t.Run("test case : not AppError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/error?err=generic", nil)
		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		expectedBody := `{"code":500,"message":"Internal Server Error"}`
		body, _ := io.ReadAll(resp.Body)
		assert.JSONEq(t, expectedBody, string(body))
	})
}

func TestValidateUserCreate(t *testing.T) {
	tests := []struct {
		name     string
		input    *model.UserCreate
		expected []errs.ErrorMessage
	}{
		{
			name:  "Valid user",
			input: &model.UserCreate{
				ID: 1, 
				RoleID: 1,
				Name: "ValidName",
				Email: "walter_white1@gmail.com",
				Password: "ValidPassword",
			},
			expected: []errs.ErrorMessage(nil),
		},
		{
			name:  "Invalid user - missing ID",
			input: &model.UserCreate{
				ID: 0, 
				RoleID: 1,
				Name: "ValidName",
				Email: "walter_white2@gmail.com",
				Password: "ValidPassword",
			},
			expected: []errs.ErrorMessage{
				{
					FailedField: "UserCreate.ID", 
					Tag: "required", Value: "",
				},
			},
		},
		{
			name:  "Invalid user - missing Name",
			input: &model.UserCreate{
				ID: 1, 
				RoleID: 1,
				Name: "",
				Email: "walter_white3@gmail.com",
				Password: "ValidPassword",
			},
			expected: []errs.ErrorMessage{
				{
					FailedField: "UserCreate.Name", 
					Tag: "required", Value: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.ValidateUserCreate(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateLoginRequest(t *testing.T) {
	tests := []struct {
		name     string
		input    *model.LoginRequest
		expected []errs.ErrorMessage
	}{
		{
			name:  "Valid login",
			input: &model.LoginRequest{
				Email: "walter_white1@gmail.com", 
				Password: "ValidPassword",
			},
			expected: []errs.ErrorMessage(nil),
		},
		{
			name:  "Invalid login - missing Email",
			input: &model.LoginRequest{
				Email: "", 
				Password: "ValidPassword",
			},
			expected: []errs.ErrorMessage{
				{
					FailedField: "LoginRequest.Email", 
					Tag: "required", Value: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.ValidateLoginRequest(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateProductTypeCreate(t *testing.T) {
	tests := []struct {
		name     string
		input    *model.ProductTypeCreate
		expected []errs.ErrorMessage
	}{
		{
			name:  "Valid product type",
			input: &model.ProductTypeCreate{
				ID: 1, 
				Name: "ValidName",
			},
			expected: []errs.ErrorMessage(nil),
		},
		{
			name:  "Invalid product type - missing ID",
			input: &model.ProductTypeCreate{
				Name: "ValidName",
			},
			expected: []errs.ErrorMessage{
				{
					FailedField: "ProductTypeCreate.ID", 
					Tag: "required", Value: "",
				},
			},
		},
		{
			name:  "Invalid product type - missing Name",
			input: &model.ProductTypeCreate{
				ID: 1,
			},
			expected: []errs.ErrorMessage{
				{
					FailedField: "ProductTypeCreate.Name", 
					Tag: "required", Value: "",
				},
			},
		},
		{
			name:  "Invalid product type - ID less than 1",
			input: &model.ProductTypeCreate{
				ID: 0, 
				Name: "ValidName",
			},
			expected: []errs.ErrorMessage{
				{
					FailedField: "ProductTypeCreate.ID", 
					Tag: "required", 
					Value: "",
				},
			},
		},
		{
			name:  "Invalid product type - Name too long",
			input: &model.ProductTypeCreate{
				ID: 1, 
				Name: "12345678901234567890123456789012345678901",
			},
			expected: []errs.ErrorMessage{
				{
					FailedField: "ProductTypeCreate.Name", 
					Tag: "max", 
					Value: "40",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.ValidateProductTypeCreate(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateProductTypeUpdate(t *testing.T) {
	tests := []struct {
		name     string
		input    *model.ProductTypeUpdate
		expected []errs.ErrorMessage
	}{
		{
			name:  "Valid product type",
			input: &model.ProductTypeUpdate{Name: "ValidName"},
			expected: []errs.ErrorMessage(nil),
		},
		{
			name:  "Invalid product type - missing Name",
			input: &model.ProductTypeUpdate{Name: ""},
			expected: []errs.ErrorMessage{
				{
					FailedField: "ProductTypeUpdate.Name", 
					Tag: "required", Value: "",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := helper.ValidateProductTypeUpdate(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
