package helper

import (
	"github.com/Yoshikrit/fiber-test/helper/errs"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	
	"strconv"
)

func ParamsInt(ctx *fiber.Ctx) (int, error) {
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, errs.NewBadRequestError("Invalid ID: " + idParam + " is not integer")
	}
	return id, nil
}

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.NewUnprocessableError(err.Error())
	}
	return hashedPassword, nil
}

func CompareHashAndPassword(passwordFromDB, password []byte) error {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordFromDB), []byte(password)); err != nil {
		return errs.NewNotFoundError("Email or Password is incorrect")
	}
	return nil
}