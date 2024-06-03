package errs

import (
	"net/http"
	"github.com/goccy/go-json"
	"fmt"
)

type ErrorResponse struct {
	Code    int		`json:"code"`
	Message string	`json:"message"`
}

type ValErrorResponse struct {
	Code    int				`json:"code"`
	Message []ErrorMessage	`json:"message"`
}

type ErrorMessage struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func(e ErrorResponse) Error() string {
	return e.Message
}

func (v ValErrorResponse) Error() string {
	messageBytes, err := json.Marshal(v.Message)
	if err != nil {
		return fmt.Sprintf("code: %d, message: %+v", v.Code, v.Message)
	}
	return fmt.Sprintf("code: %d, message: %s", v.Code, string(messageBytes))
}

func NewUnauthorizedError(message string) error {
	return ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func NewBadRequestError(message string) error {
	return ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewValidateBadRequestError(message []ErrorMessage) error {
	return ValErrorResponse{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewNotFoundError(message string) error {
	return ErrorResponse{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewConflictError(message string) error {
	return ErrorResponse{
		Code:    http.StatusConflict,
		Message: message,
	}
}

func NewInternalServerError(message string) error {
	return ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func NewUnprocessableError(message string) error {
	return ErrorResponse{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}
