package service

import (
	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/repository"
	"github.com/Yoshikrit/fiber-test/helper/errs"
	"github.com/Yoshikrit/fiber-test/helper/logger"
	"github.com/Yoshikrit/fiber-test/helper"
)

type ProductTypeServiceImpl struct {
	ProdTypeRepo 	repository.ProductTypeRepository
}

func NewProductTypeServiceImpl(prodTypeRepo repository.ProductTypeRepository) ProductTypeService {
	return &ProductTypeServiceImpl{
		ProdTypeRepo: 	prodTypeRepo,
	}
}

func (s *ProductTypeServiceImpl) Create(prodTypeCreateReq *model.ProductTypeCreate) error {
	if err := helper.ValidateProductTypeCreate(prodTypeCreateReq); err != nil {
		logger.Error("ProductType data is not valid")
		return errs.NewValidateBadRequestError(err)
	}

	prodTypeFromDB, _ := s.ProdTypeRepo.FindByID(prodTypeCreateReq.ID)
    if prodTypeFromDB != nil && prodTypeFromDB.ID == prodTypeCreateReq.ID {
        logger.Error("ProductType with this ID already exists")
        return errs.NewConflictError("ProductType with this ID already exists")
    }

	prodTypeEntity := &model.ProductTypeEntity{
		ID:       prodTypeCreateReq.ID,
		Name:     prodTypeCreateReq.Name,
	}
	
	if err := s.ProdTypeRepo.Save(prodTypeEntity); err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("Service: Create ProductType Successfully")
	return nil
}

func (s *ProductTypeServiceImpl) FindAll() ([]model.ProductType, error) {
	prodTypeEntities, err := s.ProdTypeRepo.FindAll()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var prodTypesRes []model.ProductType
	for _, prodTypeEntity := range prodTypeEntities {
		prodTypeRes := &model.ProductType{
			ID:       prodTypeEntity.ID,
			Name:     prodTypeEntity.Name,
		}
		prodTypesRes = append(prodTypesRes, *prodTypeRes)
	}

	logger.Info("Service: Find All ProductTypes Successfully")
	return prodTypesRes, nil
}

func (s *ProductTypeServiceImpl) FindByID(id int) (*model.ProductType, error) {
	prodTypeEntity, err := s.ProdTypeRepo.FindByID(id);
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	prodTypeRes := &model.ProductType{
		ID:       prodTypeEntity.ID,
		Name:     prodTypeEntity.Name,
	}

	logger.Info("Service: Find ProductType By ID Successfully")
	return prodTypeRes, nil
}

func (s *ProductTypeServiceImpl) Update(id int, prodTypeUpdateReq *model.ProductTypeUpdate) error {
	if err := helper.ValidateProductTypeUpdate(prodTypeUpdateReq); err != nil {
		logger.Error("ProductType Update data is not valid")
		return errs.NewValidateBadRequestError(err)
	}

	_, err := s.ProdTypeRepo.FindByID(id)
	if err != nil {
		logger.Error(err)
		return err
	}

	prodTypeEntity := &model.ProductTypeEntity{
		ID:       id,
		Name:     prodTypeUpdateReq.Name,
	}
	if err := s.ProdTypeRepo.Update(prodTypeEntity); err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("Service: Update ProductType Successfully")
	return nil
}

func (s *ProductTypeServiceImpl) Delete(id int) error {
	_, err := s.ProdTypeRepo.FindByID(id)
	if err != nil {
		logger.Error(err)
		return err
	}
	
	if err := s.ProdTypeRepo.Delete(id); err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("Service: Delete ProductType Successfully")
	return nil
}

func (s *ProductTypeServiceImpl) Count() (int64, error) {
	count, err := s.ProdTypeRepo.Count();
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	logger.Info("Service: Get ProductType'Count Successfully")
	return count, nil
}