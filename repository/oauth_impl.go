package repository

import (
	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/helper/errs"

	"gorm.io/gorm"
)

type OauthRepositoryImpl struct {
	db *gorm.DB
}

func NewOauthRepositoryImpl(db *gorm.DB) OauthRepository {
	return &OauthRepositoryImpl{db: db}
}

func (r *OauthRepositoryImpl) Create(oauthReq *model.OauthEntity) error {
	if err := r.db.Create(&oauthReq).Error; err != nil {
		return errs.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *OauthRepositoryImpl) FindByID(id int) (*model.OauthEntity, error) {
	var oauthEntity model.OauthEntity
	err := r.db.First(&oauthEntity, id).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errs.NewNotFoundError(err.Error())
		}
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &oauthEntity, nil
}

func (r *OauthRepositoryImpl) FindByUserID(id int) (*model.OauthEntity, error) {
	var oauthEntity model.OauthEntity
	err := r.db.Where("user_id = ?", id).First(&oauthEntity).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errs.NewNotFoundError("Email or Password is incorrect")
		}
		return nil, errs.NewInternalServerError(err.Error())
	}
	return &oauthEntity, nil
}

func (r *OauthRepositoryImpl) FindByAccessToken(id int, accessToken string) (*model.OauthEntity, error) {
	var oauthEntity model.OauthEntity
	err := r.db.Where("oauth_id = ? AND access_token = ?", id, accessToken).First(&oauthEntity).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errs.NewUnauthorizedError(err.Error())
		}
		return nil, errs.NewInternalServerError(err.Error())
	}
	return &oauthEntity, nil
}

func (r *OauthRepositoryImpl) FindByRefleshToken(refleshToken string) (*model.OauthEntity, error) {
	var oauthEntity model.OauthEntity
	err := r.db.Where("reflesh_token = ?", refleshToken).First(&oauthEntity).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errs.NewUnauthorizedError("Reflesh Token is incorrect")
		}
		return nil, errs.NewInternalServerError(err.Error())
	}
	return &oauthEntity, nil
}

func (r *OauthRepositoryImpl) Update(oauthUpdateReq *model.OauthEntity) error{
	if err := r.db.Model(&oauthUpdateReq).Updates(oauthUpdateReq).Error; err != nil {
		return errs.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *OauthRepositoryImpl) Delete(id int) error{
	if err := r.db.Delete(&model.OauthEntity{}, id).Error; err != nil {
		return errs.NewInternalServerError(err.Error())
	}
	return nil
}