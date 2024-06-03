package helper

import (
	"github.com/Yoshikrit/fiber-test/helper/errs"
	"github.com/Yoshikrit/fiber-test/model"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	
	"strconv"
)

func HandleError(ctx *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case errs.ErrorResponse:
		return ctx.Status(e.Code).JSON(fiber.Map{
			"code":    e.Code,
			"message": e.Message,
		})
	case errs.ValErrorResponse:
		return ctx.Status(e.Code).JSON(fiber.Map{
			"code":    e.Code,
			"message": e.Message,
		})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Internal Server Error",
		})
	}
}

func ParamsInt(ctx *fiber.Ctx) (int, error) {
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, errs.NewBadRequestError("Invalid ID: " + idParam + " is not integer")
	}
	return id, nil
}

func ValidateProductTypeCreate(prod *model.ProductTypeCreate) []errs.ErrorMessage {
    var errors []errs.ErrorMessage
    validate := validator.New()
    err := validate.Struct(prod)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            var element errs.ErrorMessage
            element.FailedField = err.StructNamespace()
            element.Tag = err.Tag()
            element.Value = err.Param()
            errors = append(errors, element)
        }
    }
    return errors
}

func ValidateProductTypeUpdate(prod *model.ProductTypeUpdate) []errs.ErrorMessage {
    var errors []errs.ErrorMessage
    validate := validator.New()
    err := validate.Struct(prod)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            var element errs.ErrorMessage
            element.FailedField = err.StructNamespace()
            element.Tag = err.Tag()
            element.Value = err.Param()
            errors = append(errors, element)
        }
    }
    return errors
}
