package middleware

import (
	"github.com/Yoshikrit/fiber-test/repository"
	"github.com/Yoshikrit/fiber-test/helper"
	"github.com/Yoshikrit/fiber-test/helper/errs"
	"github.com/Yoshikrit/fiber-test/helper/logger"
	"github.com/gofiber/fiber/v2"

	"strings"
)

func NewJWTMiddleware(userRepo repository.UserRepository, oauthRepo repository.OauthRepository, roleRepo repository.RoleRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenString := strings.TrimPrefix(ctx.Get("Authorization"), "Bearer ")
		logger.Info("idk " + tokenString)
		claims, err := helper.ParseToken(tokenString)
		if err != nil {
			logger.Error(err.Error())
			return helper.HandleError(ctx, err)
		}

		_, err = oauthRepo.FindByAccessToken(claims.Claims.ID, tokenString)
        if err != nil {
			logger.Error(err.Error())
			return helper.HandleError(ctx, err)
		}

		roleEntity, err := roleRepo.FindByID(claims.Claims.RoleID)
        if err != nil {
			logger.Error(err.Error())
			return helper.HandleError(ctx, err)
		}

		if (roleEntity.Title != "Manager") {
			logger.Error("Unauthorized")
			return helper.HandleError(ctx, errs.NewUnauthorizedError("Unauthorized"))
		}

		return ctx.Next()
	}
}

