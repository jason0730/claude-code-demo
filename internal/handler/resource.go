package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	authmw "github.com/jason0730/claude-code-demo/internal/auth/middleware"
	"github.com/jason0730/claude-code-demo/internal/model"
	log "github.com/sirupsen/logrus"
)

// ResourceHandler 资源处理器
type ResourceHandler struct{}

// NewResourceHandler 创建资源处理器
func NewResourceHandler() *ResourceHandler {
	return &ResourceHandler{}
}

// 模拟资源存储
var mockResources = []model.Resource{
	{
		ID:          "res-1",
		Name:        "Sample Resource 1",
		Description: "This is a sample resource",
		Type:        "compute",
		Owner:       "1",
		Metadata: map[string]string{
			"region": "us-west-2",
			"env":    "production",
		},
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	},
	{
		ID:          "res-2",
		Name:        "Sample Resource 2",
		Description: "Another sample resource",
		Type:        "storage",
		Owner:       "2",
		Metadata: map[string]string{
			"region": "us-east-1",
			"env":    "staging",
		},
		CreatedAt: time.Now().Add(-12 * time.Hour),
		UpdatedAt: time.Now().Add(-12 * time.Hour),
	},
}

// ListResources 列出所有资源
func (h *ResourceHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	claims, _ := authmw.GetClaims(r.Context())

	log.WithFields(log.Fields{
		"user_id":  claims.UserID,
		"username": claims.Username,
	}).Info("listing resources")

	respondJSON(w, http.StatusOK, mockResources)
}

// CreateResource 创建资源
func (h *ResourceHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	var req model.CreateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claims, _ := authmw.GetClaims(r.Context())

	// 创建新资源
	resource := model.Resource{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Owner:       claims.UserID,
		Metadata:    req.Metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 添加到模拟存储
	mockResources = append(mockResources, resource)

	log.WithFields(log.Fields{
		"user_id":     claims.UserID,
		"resource_id": resource.ID,
		"name":        resource.Name,
	}).Info("resource created")

	respondJSON(w, http.StatusCreated, resource)
}
