package repository

import (
	"github.com/Yoshikrit/fiber-test/model"
)

type OauthRepository interface {
	Create(oauthEntity *model.OauthEntity) error
	FindByID(int) (*model.OauthEntity, error)
	FindByUserID(id int) (*model.OauthEntity, error)
	FindByAccessToken(id int, accessToken string) (*model.OauthEntity, error)
	FindByRefleshToken(refleshToken string) (*model.OauthEntity, error)
	Update(*model.OauthEntity) error
	Delete(id int) error
}