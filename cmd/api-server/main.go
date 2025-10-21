package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	authmw "github.com/jason0730/claude-code-demo/internal/auth/middleware"
	authjwt "github.com/jason0730/claude-code-demo/internal/auth/jwt"
	authzmw "github.com/jason0730/claude-code-demo/internal/authz/middleware"
	"github.com/jason0730/claude-code-demo/internal/authz/rbac"
	"github.com/jason0730/claude-code-demo/internal/config"
	"github.com/jason0730/claude-code-demo/internal/handler"
	log "github.com/sirupsen/logrus"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 配置日志
	setupLogger(cfg.Log)

	log.Info("Starting API Server...")

	// 初始化组件
	tokenManager := authjwt.NewTokenManager(&cfg.Auth)
	rbacManager := rbac.NewRBACManager()

	// 初始化中间件
	authMiddleware := authmw.NewAuthMiddleware(tokenManager)
	authzMiddleware := authzmw.NewAuthzMiddleware(rbacManager)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(tokenManager)
	userHandler := handler.NewUserHandler()
	resourceHandler := handler.NewResourceHandler()
	healthHandler := handler.NewHealthHandler()

	// 创建路由
	router := setupRouter(
		authMiddleware,
		authzMiddleware,
		authHandler,
		userHandler,
		resourceHandler,
		healthHandler,
	)

	// 创建 HTTP 服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器
	go func() {
		log.WithField("address", addr).Info("Server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start server")
		}
	}()

	// 优雅关闭
	gracefulShutdown(srv, cfg.Server.ShutdownTimeout)
}

// setupRouter 设置路由
func setupRouter(
	authMw *authmw.AuthMiddleware,
	authzMw *authzmw.AuthzMiddleware,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	resourceHandler *handler.ResourceHandler,
	healthHandler *handler.HealthHandler,
) *mux.Router {
	router := mux.NewRouter()

	// 添加日志中间件
	router.Use(loggingMiddleware)

	// 健康检查端点（无需认证）
	router.HandleFunc("/health", healthHandler.Health).Methods("GET")
	router.HandleFunc("/ready", healthHandler.Ready).Methods("GET")
	router.HandleFunc("/metrics", healthHandler.Metrics).Methods("GET")

	// API 路由
	api := router.PathPrefix("/api/v1").Subrouter()

	// 认证端点（无需认证）
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/refresh", authHandler.Refresh).Methods("POST")

	// 需要认证的端点
	authenticated := api.PathPrefix("").Subrouter()
	authenticated.Use(authMw.Authenticate)

	// 用户端点
	authenticated.Handle("/users",
		authzMw.RequirePermission(rbac.PermissionUserList)(
			http.HandlerFunc(userHandler.ListUsers),
		),
	).Methods("GET")

	authenticated.Handle("/users/{id}",
		authzMw.RequireAnyRole(rbac.RoleAdmin, rbac.RoleUser)(
			http.HandlerFunc(userHandler.GetUser),
		),
	).Methods("GET")

	// 资源端点
	authenticated.Handle("/resources",
		authzMw.RequirePermission(rbac.PermissionResourceList)(
			http.HandlerFunc(resourceHandler.ListResources),
		),
	).Methods("GET")

	authenticated.Handle("/resources",
		authzMw.RequirePermission(rbac.PermissionResourceWrite)(
			http.HandlerFunc(resourceHandler.CreateResource),
		),
	).Methods("POST")

	return router
}

// setupLogger 配置日志
func setupLogger(cfg config.LogConfig) {
	// 设置日志级别
	level, err := log.ParseLevel(cfg.Level)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)

	// 设置日志格式
	if cfg.Format == "json" {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	log.SetOutput(os.Stdout)
}

// loggingMiddleware 日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 创建响应写入器包装器以捕获状态码
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		log.WithFields(log.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      wrapped.statusCode,
			"duration_ms": duration.Milliseconds(),
			"remote_addr": r.RemoteAddr,
		}).Info("HTTP request")
	})
}

// responseWriter 包装 ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// gracefulShutdown 优雅关闭
func gracefulShutdown(srv *http.Server, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server forced to shutdown")
	}

	log.Info("Server exited")
}
