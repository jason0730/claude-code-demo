package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jason0730/claude-code-demo/internal/auth/jwt"
	"github.com/jason0730/claude-code-demo/internal/model"
	log "github.com/sirupsen/logrus"
)

type contextKey string

const (
	ClaimsContextKey contextKey = "claims"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	tokenManager *jwt.TokenManager
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(tokenManager *jwt.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
	}
}

// Authenticate 认证中间件处理函数
func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取 token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			am.respondError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			am.respondError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// 验证 token
		claims, err := am.tokenManager.ValidateToken(tokenString)
		if err != nil {
			log.WithError(err).Warn("token validation failed")
			am.respondError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// 将 claims 存入 context
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)

		// 调用下一个处理器
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// respondError 返回错误响应
func (am *AuthMiddleware) respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// GetClaims 从 context 获取 claims
func GetClaims(ctx context.Context) (*jwt.CustomClaims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(*jwt.CustomClaims)
	return claims, ok
}

// GetUserFromContext 从 context 获取用户信息
func GetUserFromContext(ctx context.Context) (*model.Claims, error) {
	claims, ok := GetClaims(ctx)
	if !ok {
		return nil, http.ErrNoCookie
	}

	return &model.Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Roles:    claims.Roles,
	}, nil
}
