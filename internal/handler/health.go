package handler

import (
	"net/http"
	"time"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// Health 健康检查端点（存活探针）
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "healthy",
		"uptime": time.Since(h.startTime).String(),
	})
}

// Ready 就绪检查端点（就绪探针）
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// 这里可以检查数据库连接、依赖服务等
	// 如果服务未就绪，返回 503 状态码
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ready",
	})
}

// Metrics Prometheus 指标端点（简化版）
func (h *HealthHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	// 这里可以集成 Prometheus 客户端库
	// 返回基本的指标信息
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("# HELP api_server_uptime_seconds API server uptime in seconds\n"))
	w.Write([]byte("# TYPE api_server_uptime_seconds gauge\n"))
	w.Write([]byte("api_server_uptime_seconds " + time.Since(h.startTime).String() + "\n"))
}
