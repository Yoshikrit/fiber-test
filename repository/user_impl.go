package repository

import (
	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/helper/errs"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(userCreateReq *model.UserEntity) error {
	if err := r.db.Create(&userCreateReq).Error; err != nil {
		return errs.NewInternalServerError(err.Error())
	}
	return nil
}


func (r *UserRepositoryImpl) FindByID(id int) (*model.UserEntity, error) {
	var userEntity model.UserEntity
	err := r.db.First(&userEntity, id).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errs.NewNotFoundError(err.Error())
		}
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &userEntity, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*model.UserEntity, error) {
	var user model.UserEntity
	err := r.db.Where("user_email = ?", email).First(&user).Error
	if err != nil {
		return nil, errs.NewNotFoundError("Email or Password is incorrect")
	}
	return &user, nil
}