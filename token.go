package htxp

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v4"
)

type JWTToken struct {
}

func NewJWTTokenLogic() *JWTToken {
	return &JWTToken{}
}

// GenToken 生成jwt（包级函数，无需声明实例）
func GenToken(secretKey string, iat, seconds int64, claims jwt.MapClaims) (string, error) {
	if claims == nil {
		claims = make(jwt.MapClaims)
	}
	claims["exp"] = iat + seconds
	claims["iat"] = iat

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

// GenTokenWithPayload 生成jwt（兼容旧版本，payload作为字符串）
func GenTokenWithPayload(secretKey string, iat, seconds int64, payload string) (string, error) {
	claims := jwt.MapClaims{}
	claims["payload"] = payload
	return GenToken(secretKey, iat, seconds, claims)
}

// GenTokenWithUser 生成包含用户信息的jwt（便捷方法）
func GenTokenWithUser(secretKey string, iat, seconds int64, uid, cid interface{}) (string, error) {
	claims := jwt.MapClaims{}
	if uid != nil {
		claims["uid"] = uid
	}
	if cid != nil {
		claims["cid"] = cid
	}
	return GenToken(secretKey, iat, seconds, claims)
}

// GenToken 生成jwt（实例方法，保持向后兼容）
func (j *JWTToken) GenToken(secretKey string, iat, seconds int64, payload string) (string, error) {
	return GenTokenWithPayload(secretKey, iat, seconds, payload)
}

// GenRefreshToken 生成刷新token
func (j *JWTToken) GenRefreshToken() (string, error) {
	b := make([]byte, 16) // 16字节（128位）随机数据
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// GenRefreshToken 生成刷新token（包级函数）
func GenRefreshToken() (string, error) {
	b := make([]byte, 16) // 16字节（128位）随机数据
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
