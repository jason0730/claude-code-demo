package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jason0730/claude-code-demo/internal/auth/jwt"
	"github.com/jason0730/claude-code-demo/internal/model"
	log "github.com/sirupsen/logrus"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	tokenManager *jwt.TokenManager
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(tokenManager *jwt.TokenManager) *AuthHandler {
	return &AuthHandler{
		tokenManager: tokenManager,
	}
}

// Login 用户登录
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// 验证用户名和密码（这里是示例，实际应该查询数据库）
	user := h.authenticateUser(req.Username, req.Password)
	if user == nil {
		log.WithField("username", req.Username).Warn("login failed: invalid credentials")
		respondError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	// 生成 token
	accessToken, refreshToken, err := h.tokenManager.GenerateToken(user)
	if err != nil {
		log.WithError(err).Error("failed to generate token")
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	log.WithFields(log.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("user logged in successfully")

	respondJSON(w, http.StatusOK, model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes in seconds
		TokenType:    "Bearer",
	})
}

// Refresh 刷新令牌
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// 验证刷新令牌
	userID, err := h.tokenManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		log.WithError(err).Warn("invalid refresh token")
		respondError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	// 获取用户信息（这里是示例，实际应该查询数据库）
	user := h.getUserByID(userID)
	if user == nil {
		respondError(w, http.StatusUnauthorized, "user not found")
		return
	}

	// 生成新的 token
	accessToken, refreshToken, err := h.tokenManager.GenerateToken(user)
	if err != nil {
		log.WithError(err).Error("failed to generate token")
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	log.WithField("user_id", userID).Info("token refreshed successfully")

	respondJSON(w, http.StatusOK, model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
		TokenType:    "Bearer",
	})
}

// authenticateUser 验证用户（示例实现，实际应该查询数据库）
func (h *AuthHandler) authenticateUser(username, password string) *model.User {
	// 这是一个示例用户，实际应该从数据库查询并验证密码哈希
	// 密码应该使用 bcrypt 等算法加密存储
	mockUsers := map[string]*model.User{
		"admin": {
			ID:       "1",
			Username: "admin",
			Email:    "admin@example.com",
			Password: "admin123", // 实际应该是加密后的密码
			Roles:    []string{"admin"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		"editor": {
			ID:       "2",
			Username: "editor",
			Email:    "editor@example.com",
			Password: "editor123",
			Roles:    []string{"editor"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		"viewer": {
			ID:       "3",
			Username: "viewer",
			Email:    "viewer@example.com",
			Password: "viewer123",
			Roles:    []string{"viewer"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	user, exists := mockUsers[username]
	if !exists || user.Password != password {
		return nil
	}

	return user
}

// getUserByID 根据 ID 获取用户（示例实现）
func (h *AuthHandler) getUserByID(userID string) *model.User {
	mockUsers := map[string]*model.User{
		"1": {
			ID:       "1",
			Username: "admin",
			Email:    "admin@example.com",
			Roles:    []string{"admin"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		"2": {
			ID:       "2",
			Username: "editor",
			Email:    "editor@example.com",
			Roles:    []string{"editor"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		"3": {
			ID:       "3",
			Username: "viewer",
			Email:    "viewer@example.com",
			Roles:    []string{"viewer"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return mockUsers[userID]
}
