package handler

import (
	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/service"
	"github.com/Yoshikrit/fiber-test/helper"
	"github.com/Yoshikrit/fiber-test/helper/logger"
	"github.com/Yoshikrit/fiber-test/helper/errs"

	"github.com/gofiber/fiber/v2"
)

type ProductTypeHandler struct {
	productTypeSrv service.ProductTypeService
}

func NewProductTypeHandler(productTypeSrv service.ProductTypeService) *ProductTypeHandler {
	return &ProductTypeHandler{productTypeSrv: productTypeSrv}
}

// CreateProductType godoc
// @Summary Create ProductType
// @Description Create producttype
// @Tags producttype
// @Produce  json
// @param ProductType body model.ProductTypeCreate true "ProductType data to be create"
// @response 201 {object} model.StringResponse "Create ProductType Successfully"
// @response 400 {object} errs.ErrorResponse "Error Bad Request"
// @response 409 {object} errs.ErrorResponse "Error Conflict Error"
// @response 500 {object} errs.ErrorResponse "Error Unexpected Error"
// @Router /producttype/ [post]
func (h *ProductTypeHandler) Create(ctx *fiber.Ctx) error {
    ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
    ctx.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)

	prodTypeReq := new(model.ProductTypeCreate)
	if err := ctx.BodyParser(prodTypeReq); err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, errs.NewBadRequestError(err.Error()))
	}

	err := h.productTypeSrv.Create(prodTypeReq)
	if err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)	
	}

	logger.Info("Handler: Create ProductType Successfully")
	webResponse := model.StringResponse{
		Code: 		201,
		Message: 	"Create ProductType Successfully",
	}
	return ctx.Status(fiber.StatusCreated).JSON(webResponse)
}

// GetAllProductTypes godoc
// @Summary Get All ProductType
// @Description Get all producttype
// @Tags producttype
// @Produce  json
// @response 200 {array} model.ProductTypesResponse "Get ProductTypes Successfully"
// @response 404 {object} errs.ErrorResponse "Error Not Found"
// @response 500 {object} errs.ErrorResponse "Error Unexpected Error"
// @Router /producttype/ [get]
func (h *ProductTypeHandler) FindAll(ctx *fiber.Ctx) error {
    ctx.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)
	
	prodTypesRes, err := h.productTypeSrv.FindAll()
	if err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)	
	}

	logger.Info("Handler: Find All ProductTypes Successfully")
	webResponse := model.ProductTypesResponse{
		Code: 		200,
		Message: 	prodTypesRes,
	}
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// GetProductTypeByID godoc
// @Summary Get ProductType
// @Description Get producttype by id
// @Tags producttype
// @Produce  json
// @Param        id   path      int  true  "ProductType ID"
// @response 200 {object} model.ProductTypeResponse "Get ProductType Successfully"
// @response 400 {object} errs.ErrorResponse "Error Bad Request"
// @response 404 {object} errs.ErrorResponse "Error Not Found"
// @response 500 {object} errs.ErrorResponse "Error Unexpected Error"
// @Router /producttype/{id} [get]
func (h *ProductTypeHandler) FindByID(ctx *fiber.Ctx) error {
    ctx.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)
	
	id, err := helper.ParamsInt(ctx)
	if err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)
	}

	prodTypeRes, err := h.productTypeSrv.FindByID(id)
	if err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)	
	}

	logger.Info("Handler: Find ProductType By ID Successfully")
	webResponse := model.ProductTypeResponse{
		Code: 		200,
		Message: 	*prodTypeRes,
	}
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// UpdateProductTypeByID godoc
// @Summary Update ProductType
// @Description Update producttype by id
// @Tags producttype
// @Produce  json
// @Param        id   path      int  true  "ProductType ID"
// @param ProductType body model.ProductTypeUpdate true "ProductType data to be update"
// @response 200 {object} model.StringResponse "Update ProductType Successfully"
// @response 400 {object} errs.ErrorResponse "Error Bad Request"
// @response 404 {object} errs.ErrorResponse "Error Not Found"
// @response 500 {object} errs.ErrorResponse "Error Unexpected Error"
// @Router /producttype/{id} [put]
func (h *ProductTypeHandler) Update(ctx *fiber.Ctx) error {
    ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
    ctx.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)

	id, err := helper.ParamsInt(ctx)
	if err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)
	}

	prodTypeReq := new(model.ProductTypeUpdate)
	if err := ctx.BodyParser(prodTypeReq); err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, errs.NewBadRequestError(err.Error()))
	}

	if err := h.productTypeSrv.Update(id, prodTypeReq); err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)	
	}

	logger.Info("Handler: Update ProductType Successfully")
	webResponse := model.StringResponse{
		Code: 		200,
		Message: 	"Update ProductType Successfully",
	}
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// DeleteProductTypeByID godoc
// @Summary Delete ProductType
// @Description Delete producttype by id
// @Tags producttype
// @Produce  json
// @Param        id   path      int  true  "ProductType ID"
// @response 200 {object} model.StringResponse "Delete ProductType Successfully"
// @response 400 {object} errs.ErrorResponse "Error Bad Request"
// @response 404 {object} errs.ErrorResponse "Error Not Found"
// @response 500 {object} errs.ErrorResponse "Error Unexpected Error"
// @Router /producttype/{id} [delete]
func (h *ProductTypeHandler) Delete(ctx *fiber.Ctx) error {
    ctx.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)

	id, err := helper.ParamsInt(ctx)
	if err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)
	}

	if err := h.productTypeSrv.Delete(id); err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)	
	}

	logger.Info("Handler: Delete ProductType Successfully")
	webResponse := model.StringResponse{
		Code: 		200,
		Message: 	"Delete ProductType Successfully",
	}
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// GetProductTypeCount godoc
// @Summary Get ProductType Count
// @Description Get producttype's count from database
// @Tags producttype
// @Produce  json
// @response 200 {object} model.CountResponse "Get ProductType'Count Successfully"
// @response 500 {object} errs.ErrorResponse "Error Unexpected Error"
// @Router /producttype/count [get]
func (h *ProductTypeHandler) Count(ctx *fiber.Ctx) error {
    ctx.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)

	count, err := h.productTypeSrv.Count()
	if err != nil {
		logger.Error(err.Error())
		return helper.HandleError(ctx, err)	
	}

	logger.Info("Handler: Get ProductType'Count Successfully")
	webResponse := model.CountResponse{
		Code: 		200,
		Message: 	int(count),
	}
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}