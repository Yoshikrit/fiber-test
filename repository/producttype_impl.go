package repository

import (
	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/helper/errs"

	"gorm.io/gorm"
)

type ProductTypeRepositoryImpl struct {
	db *gorm.DB
}

func NewProductTypeRepositoryImpl(db *gorm.DB) ProductTypeRepository {
	return &ProductTypeRepositoryImpl{db: db}
}

func (r *ProductTypeRepositoryImpl) Save(prodTypeCreateReq *model.ProductTypeEntity) error{
	if err := r.db.Create(&prodTypeCreateReq).Error; err != nil {
		return errs.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *ProductTypeRepositoryImpl) FindAll() ([]model.ProductTypeEntity, error) {
	var prodTypesEntity []model.ProductTypeEntity
	err := r.db.Find(&prodTypesEntity).Error
	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}
	return prodTypesEntity, nil
}

func (r *ProductTypeRepositoryImpl) FindByID(id int) (*model.ProductTypeEntity, error) {
	var prodTypesEntity model.ProductTypeEntity
	err := r.db.First(&prodTypesEntity, id).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errs.NewNotFoundError(err.Error())
		}
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &prodTypesEntity, nil
}

func (r *ProductTypeRepositoryImpl) Update(prodTypeUpdateReq *model.ProductTypeEntity) error{
	if err := r.db.Model(&prodTypeUpdateReq).Updates(prodTypeUpdateReq).Error; err != nil {
		return errs.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *ProductTypeRepositoryImpl) Delete(id int) error{
	if err := r.db.Delete(&model.ProductTypeEntity{}, id).Error; err != nil {
		return errs.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *ProductTypeRepositoryImpl) Count() (int64, error) {
	var count int64
	err := r.db.Table("producttype").Count(&count).Error
	if err != nil {
		return 0, errs.NewInternalServerError(err.Error())
	}
	return count, nil
}