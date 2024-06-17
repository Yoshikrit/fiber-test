package helper

import (
	"github.com/Yoshikrit/fiber-test/helper/errs"
	"github.com/Yoshikrit/fiber-test/config"
	"github.com/Yoshikrit/fiber-test/model"

	"github.com/golang-jwt/jwt/v5"

	"time"
	"errors"
	"math"
)

func GeneratePairTokens(userClaims *model.UserClaims, title string) (*model.UserToken, error) {
	configData, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	accessToken, err := generateToken("access-token", title, userClaims, configData.JWTAccessExpires)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken("refresh-token", title, userClaims, configData.JWTRefleshExpires)
	if err != nil {
		return nil, err
	}

	return &model.UserToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateToken(subject, title string, userClaims *model.UserClaims, expireDuration int) (string, error) {
	configData, err := config.GetConfig()
	if err != nil {
		return "", err
	}

	claims := &model.ServiceMapClaims{
		Claims: userClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    configData.AppName,
			Subject:   subject,
			Audience:  []string{title},
			ExpiresAt: jwtTimeDurationCal(expireDuration),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signToken, err := token.SignedString([]byte(configData.JWTSecretKey))
	if err != nil {
		return "", errs.NewInternalServerError(err.Error())
	}
	return signToken, nil
}

func NewAccessToken(title string, userClaims *model.UserClaims) (string, error) {
	configData, err := config.GetConfig()
	if err != nil {
		return "", err
	}

	claims := &model.ServiceMapClaims{
		Claims: userClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    configData.AppName,
			Subject:   "access-token",
			Audience:  []string{title},
			ExpiresAt: jwtTimeDurationCal(configData.JWTAccessExpires),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signToken, err := token.SignedString([]byte(configData.JWTSecretKey))
	if err != nil {
		return "", errs.NewInternalServerError(err.Error())
	}
	return signToken, nil
}

func RepeatToken(title string, userClaims *model.UserClaims, expireDuration int64) (string, error) {
	configData, err := config.GetConfig()
	if err != nil {
		return "", err
	}

	claims := &model.ServiceMapClaims{
		Claims: userClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    configData.AppName,
			Subject:   "reflesh-token",
			Audience:  []string{title},
			ExpiresAt: jwtTimeRepeatAdapter(expireDuration),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signToken, err := token.SignedString([]byte(configData.JWTSecretKey))
	if err != nil {
		return "", errs.NewInternalServerError(err.Error())
	}
	return signToken, nil
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func ParseToken(tokenString string) (*model.ServiceMapClaims, error) {
	configData, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &model.ServiceMapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.NewInternalServerError("Signing method is invalid")
		}
		return []byte(configData.JWTSecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errs.NewBadRequestError("Token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errs.NewUnauthorizedError("Token had expired")
		} else {
			return nil, errs.NewUnauthorizedError("Parse token failed: " + err.Error())
		}
	}

	if claims, ok := token.Claims.(*model.ServiceMapClaims); ok {
		return claims, nil
	} else {
		return nil,  errs.NewInternalServerError("claims type is invalid")
	}
}