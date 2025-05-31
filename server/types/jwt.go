package types

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/global"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWT struct {
	SigningKey []byte
}

type AccessToken struct {
	jwt.RegisteredClaims
	AuthUser
}

func NewJwt() *JWT {
	return &JWT{SigningKey: []byte(global.Config.Jwt.Secret)}
}

func (c *JWT) Generate(userId uint, roleId uint) (string, error) {
	claims := AccessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)),
		},
		AuthUser: AuthUser{
			UserId: userId,
			RoleId: roleId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(c.SigningKey)
}

func (c *JWT) Parse(tokenString string) (*AccessToken, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &AccessToken{}, func(token *jwt.Token) (interface{}, error) {
		return c.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsedToken.Claims.(*AccessToken)
	if ok && parsedToken.Valid {
		return claims, nil
	} else {
		// Token无效
		return nil, errors.New("无效的access_token")
	}
}
