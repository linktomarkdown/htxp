package middleware

import (
	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims JWT声明结构（增强版，包含版本号）
type JWTClaims struct {
	UserID  int64   `json:"user_id"`
	Version int64   `json:"version"` // Token版本号，用于权限变更时使旧Token失效
	Roles   []string `json:"roles,omitempty"` // 角色列表（可选，减少RPC调用）
	jwt.RegisteredClaims
}

// ParseToken 解析JWT token（本地验证，无需RPC）
func ParseToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

