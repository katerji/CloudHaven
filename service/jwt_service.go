package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/katerji/UserAuthKit/envs"
	"github.com/katerji/UserAuthKit/model"
	"github.com/katerji/UserAuthKit/utils"
	"os"
	"strconv"
	"time"
)

type jwtService struct{}

type customJWTClaims struct {
	UserOutput model.UserOutput `json:"user"`
	ExpiresAt  int64            `json:"expires_at"`
}

func (jwtService jwtService) VerifyToken(token string) (model.User, error) {
	jwtSecret := envs.GetInstance().GetJWTToken()
	return jwtService.validateToken(token, jwtSecret)
}

func (jwtService jwtService) VerifyRefreshToken(token string) (model.User, error) {
	jwtSecret := envs.GetInstance().GetJWTRefreshToken()
	return jwtService.validateToken(token, jwtSecret)
}

func (jwtService jwtService) validateToken(token, jwtSecret string) (model.User, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return model.User{}, errors.New("error parsing token")
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		jsonClaims, err := json.Marshal(claims)
		if err != nil {
			return model.User{}, errors.New("error parsing token")
		}
		var customClaims customJWTClaims
		if err := json.Unmarshal(jsonClaims, &customClaims); err != nil {
			return model.User{}, errors.New("error parsing token")
		}
		expiresAt := time.Unix(customClaims.ExpiresAt, 0)
		if expiresAt.Before(time.Now()) {
			return model.User{}, errors.New("token expired")
		}
		return customClaims.UserOutput.ToUser(), nil
	}

	return model.User{}, errors.New("invalid token")
}

func (jwtService jwtService) CreateJwt(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":       user.ToOutput(),
		"expires_at": getJWTExpiry(),
	})
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (jwtService jwtService) CreateRefreshJwt(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":       user.ToOutput(),
		"expires_at": getJWTRefreshExpiry(),
	})
	jwtSecret := os.Getenv("JWT_REFRESH_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func getJWTExpiry() int64 {
	expiryString := os.Getenv("JWT_EXPIRY")
	expiry, _ := strconv.Atoi(expiryString)
	return utils.IntToUnixTime(expiry)
}

func getJWTRefreshExpiry() int64 {
	expiryString := os.Getenv("JWT_REFRESH_EXPIRY")
	expiry, _ := strconv.Atoi(expiryString)
	return utils.IntToUnixTime(expiry)
}
