package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwtv5.RegisteredClaims
}

func GenerateToken(secret string, userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(secret, tokenString string) (*Claims, error) {
	token, err := jwtv5.ParseWithClaims(tokenString, &Claims{}, func(t *jwtv5.Token) (any, error) {
		//检查加密算法是否是HMAC，HS256属于此类
		_, ok := t.Method.(*jwtv5.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid signing method")
		}
		//返回密钥，用于与最后的签名对比
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
