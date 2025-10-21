package middleware

import (
	"encoding/json"
	"net/http"

	authmw "github.com/jason0730/claude-code-demo/internal/auth/middleware"
	"github.com/jason0730/claude-code-demo/internal/authz/rbac"
	log "github.com/sirupsen/logrus"
)

// AuthzMiddleware 授权中间件
type AuthzMiddleware struct {
	rbacManager *rbac.RBACManager
}

// NewAuthzMiddleware 创建授权中间件
func NewAuthzMiddleware(rbacManager *rbac.RBACManager) *AuthzMiddleware {
	return &AuthzMiddleware{
		rbacManager: rbacManager,
	}
}

// RequirePermission 要求特定权限
func (am *AuthzMiddleware) RequirePermission(permission rbac.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := authmw.GetClaims(r.Context())
			if !ok {
				am.respondError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			// 检查权限
			if !am.rbacManager.CheckPermission(claims.Roles, permission) {
				log.WithFields(log.Fields{
					"user_id":    claims.UserID,
					"username":   claims.Username,
					"roles":      claims.Roles,
					"permission": permission,
				}).Warn("permission denied")

				am.respondError(w, http.StatusForbidden, "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole 要求特定角色
func (am *AuthzMiddleware) RequireRole(role rbac.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := authmw.GetClaims(r.Context())
			if !ok {
				am.respondError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			// 检查角色
			if !am.rbacManager.HasRole(claims.Roles, role) {
				log.WithFields(log.Fields{
					"user_id":       claims.UserID,
					"username":      claims.Username,
					"roles":         claims.Roles,
					"required_role": role,
				}).Warn("role not found")

				am.respondError(w, http.StatusForbidden, "insufficient role")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole 要求任一角色
func (am *AuthzMiddleware) RequireAnyRole(roles ...rbac.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := authmw.GetClaims(r.Context())
			if !ok {
				am.respondError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			// 检查是否有任一角色
			if !am.rbacManager.HasAnyRole(claims.Roles, roles) {
				log.WithFields(log.Fields{
					"user_id":        claims.UserID,
					"username":       claims.Username,
					"roles":          claims.Roles,
					"required_roles": roles,
				}).Warn("no matching role found")

				am.respondError(w, http.StatusForbidden, "insufficient role")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// respondError 返回错误响应
func (am *AuthzMiddleware) respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
