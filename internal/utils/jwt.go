package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthUser struct {
	UserId       uint    `json:"user_id"`
	RoleId       uint    `json:"role_id"`
	LoginVersion *uint64 `json:"login_version,omitempty"`
}

type JWT struct {
	SigningKey []byte
}

type AccessToken struct {
	jwt.RegisteredClaims
	AuthUser
}

func NewJwt(secret string) *JWT {
	return &JWT{SigningKey: []byte(secret)}
}

func (c *JWT) Generate(userId uint, roleId uint, loginVersion uint64, lifetime time.Duration) (string, error) {
	now := time.Now()
	claims := AccessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(lifetime)),
		},
		AuthUser: AuthUser{
			UserId:       userId,
			RoleId:       roleId,
			LoginVersion: &loginVersion,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(c.SigningKey)
}

func (c *JWT) Parse(tokenString string) (*AccessToken, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &AccessToken{}, func(token *jwt.Token) (interface{}, error) {
		return c.SigningKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("登录已过期，请重新登录")
		}
		return nil, errors.New("登录凭证无效，请重新登录")
	}
	claims, ok := parsedToken.Claims.(*AccessToken)
	if ok && parsedToken.Valid && claims.LoginVersion != nil {
		return claims, nil
	} else {
		return nil, errors.New("登录凭证无效，请重新登录")
	}
}
