package handler

import (
	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/service"
	"github.com/Yoshikrit/fiber-test/helper"
	"github.com/Yoshikrit/fiber-test/helper/logger"
	"github.com/Yoshikrit/fiber-test/helper/errs"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authSrv service.AuthService
}

func NewAuthHandler(authSrv service.AuthService) *AuthHandler {
	return &AuthHandler{authSrv: authSrv}
}

// RegisterUser godoc
// @Summary Register User
// @Description Register user
// @Tags auths
// @Produce  json
// @param User body model.UserCreate true "User data to be register"
// @response 201 {object} model.StringResponse "Register User Successfully"
// @response 400 {object} errs.ErrorResponse "Error Bad Request"
// @response 409 {object} errs.ErrorResponse "Error Conflict Error"
// @response 500 {object} errs.ErrorResponse "Error Unexpected Error"
// @Router /auths/ [post]
func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
    ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
    ctx.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)

	userCreateReq := new(model.UserCreate)
	if err := ctx.BodyParser(userCreateReq); err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, errs.NewBadRequestError(err.Error()))
	}

	err := h.authSrv.Register(userCreateReq)
	if err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)	
	}

	logger.Info("Handler: Register User Successfully")
	webResponse := model.StringResponse{
		Code: 		201,
		Message: 	"Register User Successfully",
	}
	return ctx.Status(fiber.StatusCreated).JSON(webResponse)
}

// LoginUser godoc
// @Summary Login User
// @Description Login user
// @Tags auths
// @Produce  json
// @param User body model.LoginRequest true "User data to be login"
// @response 200 {object} model.AuthPassportResponse "Login User Successfully"
// @response 400 {object} errs.ErrorResponse "Error Bad Request"
// @response 409 {object} errs.ErrorResponse "Error Conflict Error"
// @response 500 {object} errs.ErrorResponse "Error Unexpected Error"
// @Router /auths/login [post]
func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
    ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
    ctx.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)

	loginReq := new(model.LoginRequest)
	if err := ctx.BodyParser(loginReq); err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, errs.NewBadRequestError(err.Error()))
	}

	response, err := h.authSrv.Login(loginReq)
	if err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)	
	}

	logger.Info("Handler: Login User Successfully")
	webResponse := &model.AuthPassportResponse{
		Code: 		200,
		Message: 	response,
	}
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// RefreshToken godoc
// @Summary Refresh Token
// @Description Refresh Token
// @Tags auths
// @Produce  json
// @param User body model.RefreshToken true "User data to be reflesh token"
// @response 200 {object} model.AuthPassportResponse "Reflesh Token Successfully"
// @response 400 {object} errs.ErrorResponse "Error Bad Request"
// @response 409 {object} errs.ErrorResponse "Error Conflict Error"
// @response 500 {object} errs.ErrorResponse "Error Unexpected Error"
// @Router /auths/reflesh [post]
func (h *AuthHandler) Reflesh(ctx *fiber.Ctx) error {
    ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
    ctx.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)

	refleshReq := new(model.RefreshToken)
	if err := ctx.BodyParser(refleshReq); err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, errs.NewBadRequestError(err.Error()))
	}

	response, err := h.authSrv.RefreshPassport(refleshReq)
	if err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)	
	}

	logger.Info("Handler: Reflesh Token Successfully")
	webResponse := &model.AuthPassportResponse{
		Code: 		200,
		Message: 	response,
	}
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// LogoutUser godoc
// @Summary Logout User
// @Description Logout user
// @Tags auths
// @Produce  json
// @Param        id   path      int  true  "Oauth ID"
// @response 200 {object} model.StringResponse "Logout User Successfully"
// @response 400 {object} errs.ErrorResponse "Error Bad Request"
// @response 409 {object} errs.ErrorResponse "Error Conflict Error"
// @response 500 {object} errs.ErrorResponse "Error Unexpected Error"
// @Router /auths/logout/{id} [delete]
func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
    ctx.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)

	id, err := helper.ParamsInt(ctx)
	if err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)
	}

	if err := h.authSrv.Delete(id); err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)	
	}

	logger.Info("Handler: Logout User Successfully")
	webResponse := model.StringResponse{
		Code: 		200,
		Message: 	"Logout User Successfully",
	}
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}