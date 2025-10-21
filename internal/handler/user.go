package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	authmw "github.com/jason0730/claude-code-demo/internal/auth/middleware"
	"github.com/jason0730/claude-code-demo/internal/model"
	log "github.com/sirupsen/logrus"
)

// UserHandler 用户处理器
type UserHandler struct{}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// ListUsers 列出所有用户
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	claims, _ := authmw.GetClaims(r.Context())

	log.WithFields(log.Fields{
		"user_id":  claims.UserID,
		"username": claims.Username,
	}).Info("listing users")

	// 模拟用户列表
	users := []model.User{
		{
			ID:       "1",
			Username: "admin",
			Email:    "admin@example.com",
			Roles:    []string{"admin"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:       "2",
			Username: "editor",
			Email:    "editor@example.com",
			Roles:    []string{"editor"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:       "3",
			Username: "viewer",
			Email:    "viewer@example.com",
			Roles:    []string{"viewer"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	respondJSON(w, http.StatusOK, users)
}

// GetUser 获取用户详情
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	claims, _ := authmw.GetClaims(r.Context())

	log.WithFields(log.Fields{
		"requester_id": claims.UserID,
		"target_id":    userID,
	}).Info("getting user details")

	// 模拟用户数据
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

	user, exists := mockUsers[userID]
	if !exists {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}
