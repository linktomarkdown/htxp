package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/linktomarkdown/htxp"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// AuthGuardMiddleware JWT本地验证中间件（无需RPC调用）
type AuthGuardMiddleware struct {
	jwtSecret string
	redis     *redis.Redis
}

// NewAuthGuardMiddleware 创建认证守卫中间件
// jwtSecret: JWT签名密钥，需要与iam-rpc保持一致
// redis: Redis客户端，用于Token版本检查
func NewAuthGuardMiddleware(jwtSecret string, redis *redis.Redis) *AuthGuardMiddleware {
	return &AuthGuardMiddleware{
		jwtSecret: jwtSecret,
		redis:     redis,
	}
}

// Handle 处理HTTP请求，验证JWT Token
func (m *AuthGuardMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 提取Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			htxp.ErrorWithCode(w, errors.New("authorization header required"), 401)
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			htxp.ErrorWithCode(w, errors.New("invalid authorization header format"), 401)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. 本地验证JWT（无需RPC调用）
		claims, err := ParseToken(token, m.jwtSecret)
		if err != nil {
			logx.Errorf("JWT解析失败: %v", err)
			htxp.ErrorWithCode(w, errors.New("invalid token"), 401)
			return
		}

		// 3. 检查Token版本号
		if !CheckTokenVersion(m.redis, claims.UserID, claims.Version) {
			logx.Errorf("Token版本过期: userID=%d, tokenVersion=%d", claims.UserID, claims.Version)
			htxp.ErrorWithCode(w, errors.New("token expired, please refresh"), 401)
			return
		}

		// 4. 将用户ID放到context
		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		next(w, r.WithContext(ctx))
	}
}

