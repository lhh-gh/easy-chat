package ctxdata

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

const Identity = "imooc.com"

func GetJwtToken1(secretKey string, iat, seconds int64, uid string) (string, error) {
	// 输入验证
	if len(secretKey) < 32 {
		return "", fmt.Errorf("secretKey must be at least 32 bytes long")
	}
	if iat <= 0 {
		return "", fmt.Errorf("iat must be positive")
	}
	if seconds <= 0 {
		return "", fmt.Errorf("seconds must be positive")
	}
	if uid == "" {
		return "", fmt.Errorf("uid cannot be empty")
	}

	// 使用自定义Claims结构体
	claims := jwt.MapClaims{
		"exp":    iat + seconds, // 失效时间
		"iat":    iat,           // 生成时间
		Identity: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds // 失效时间
	claims["iat"] = iat           // 生成时间
	claims[Identity] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
