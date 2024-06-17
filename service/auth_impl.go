package service

import (
	"github.com/Yoshikrit/fiber-test/model"
	"github.com/Yoshikrit/fiber-test/repository"
	"github.com/Yoshikrit/fiber-test/helper/logger"
	"github.com/Yoshikrit/fiber-test/helper/errs"
	"github.com/Yoshikrit/fiber-test/helper"
)

const (
	UserExist = "User with this ID already exists"
	RoleExist = "Role with this ID already exists"
)

type AuthServiceImpl struct {
	UserRepo repository.UserRepository
	RoleRepo repository.RoleRepository
	OauthRepo repository.OauthRepository
}

func NewAuthServiceImpl(UserRepo repository.UserRepository, RoleRepo repository.RoleRepository, OauthRepo repository.OauthRepository) AuthService {
	return &AuthServiceImpl{
		UserRepo: UserRepo,
		RoleRepo: RoleRepo,
		OauthRepo: OauthRepo,
	}
}

func (s *AuthServiceImpl) Register(userCreateReq *model.UserCreate) error {
	if err := helper.ValidateUserCreate(userCreateReq); err != nil {
		logger.Error("User data is not valid")
		return errs.NewValidateBadRequestError(err)
	}

	//check user id
	userFromDB, _ := s.UserRepo.FindByID(userCreateReq.ID)
    if userFromDB != nil && userFromDB.ID == userCreateReq.ID {
        logger.Error(UserExist)
        return errs.NewConflictError(UserExist)
    }

	//check role id
	_, err := s.RoleRepo.FindByID(userCreateReq.RoleID)
    if err != nil {
		logger.Error(err.Error())
		return err
	}

	hashedPassword, err := helper.HashPassword(userCreateReq.Password)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	userEntity := &model.UserEntity{
		ID: 	  userCreateReq.ID,
		RoleID:   userCreateReq.RoleID,
		Name:     userCreateReq.Name,
		Email:    userCreateReq.Email,
		Password: string(hashedPassword),
	}

	if err := s.UserRepo.Create(userEntity); err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("Service: Register User Successfully")
	return nil
}

func (s *AuthServiceImpl) Login(loginReq *model.LoginRequest) (*model.UserPassport, error) {
	userEntity, err := s.UserRepo.FindByEmail(loginReq.Email)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := helper.CompareHashAndPassword([]byte(userEntity.Password), []byte(loginReq.Password)); err != nil {
		logger.Error(err)
		return nil, err
	}

	roleEntity, err := s.RoleRepo.FindByID(userEntity.RoleID)
    if err != nil {
		logger.Error(err)
		return nil, err
	}

	userClaims := &model.UserClaims{
		ID:    	userEntity.ID,
		RoleID:	userEntity.RoleID,
	}

	pairTokens, err := helper.GeneratePairTokens(userClaims, roleEntity.Title)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	oauthEntity := &model.OauthEntity{
		UserID:     	userEntity.ID,
		AccessToken: 	pairTokens.AccessToken,
		RefreshToken: 	pairTokens.RefreshToken,
	}

	if err := s.OauthRepo.Create(oauthEntity); err != nil {
		logger.Error(err)
		return nil, err
	}

	oauthFromDB, err := s.OauthRepo.FindByRefleshToken(oauthEntity.RefreshToken)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	userDTO := &model.UserDTO{
		ID: 		userEntity.ID,
		RoleID: 	userEntity.RoleID,
		Name: 		userEntity.Name,
		Email: 		userEntity.Email,
	}

	tokens := &model.UserToken{
		ID: oauthFromDB.ID,
		AccessToken: 	oauthEntity.AccessToken,
		RefreshToken: 	oauthEntity.RefreshToken,
	}

	userPassport := &model.UserPassport{
		User: 	userDTO,
		Tokens: tokens,
	}

	logger.Info("Service: Login User Successfully")
	return userPassport, nil
}

func (s *AuthServiceImpl) RefreshPassport(refreshToken *model.RefreshToken) (*model.UserPassport, error) {
	claims, err := helper.ParseToken(refreshToken.RefreshToken)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	oauthEntity, err := s.OauthRepo.FindByRefleshToken(refreshToken.RefreshToken)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	userEntity, err := s.UserRepo.FindByID(oauthEntity.UserID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	newUserClaims := &model.UserClaims{
		ID:     userEntity.ID,
		RoleID: userEntity.RoleID,
	}

	newUserDTO := &model.UserDTO{
		ID:     userEntity.ID,
		RoleID: userEntity.RoleID,
		Name: 	userEntity.Name,
		Email: 	userEntity.Email,
	}

	roleEntity, err := s.RoleRepo.FindByID(userEntity.RoleID)
    if err != nil {
		logger.Error(err)
		return nil, err
	}

	newAccessToken, err := helper.NewAccessToken(roleEntity.Title, newUserClaims) 
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	newRefreshToken, err := helper.RepeatToken(roleEntity.Title, newUserClaims, claims.ExpiresAt.Unix()) 
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	newPassport := &model.UserPassport{
		User: newUserDTO,
		Tokens: &model.UserToken{
			ID:           oauthEntity.ID,
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
		},
	}

	newOauthEntity := &model.OauthEntity{
		ID:     		oauthEntity.ID,
		UserID: 		userEntity.ID,
		AccessToken: 	newAccessToken,
		RefreshToken: 	newRefreshToken,
	}

	if err := s.OauthRepo.Update(newOauthEntity); err != nil {
		return nil, err
	}

	logger.Info("Service: Reflesh Token Successfully")
	return newPassport, nil
}

func (s *AuthServiceImpl) Delete(id int) error {
	_, err := s.OauthRepo.FindByID(id)
	if err != nil {
		logger.Error(err)
		return err
	}
	
	if err := s.OauthRepo.Delete(id); err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("Service: Logout User Successfully")
	return nil
}