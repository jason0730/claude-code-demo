package rbac

import (
	"errors"
)

var (
	ErrPermissionDenied = errors.New("permission denied")
)

// Role 角色定义
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleEditor Role = "editor"
	RoleViewer Role = "viewer"
	RoleUser   Role = "user"
)

// Permission 权限定义
type Permission string

const (
	PermissionUserRead   Permission = "user:read"
	PermissionUserWrite  Permission = "user:write"
	PermissionUserDelete Permission = "user:delete"
	PermissionUserList   Permission = "user:list"

	PermissionResourceRead   Permission = "resource:read"
	PermissionResourceWrite  Permission = "resource:write"
	PermissionResourceDelete Permission = "resource:delete"
	PermissionResourceList   Permission = "resource:list"
)

// RBACManager RBAC 管理器
type RBACManager struct {
	rolePermissions map[Role][]Permission
}

// NewRBACManager 创建 RBAC 管理器
func NewRBACManager() *RBACManager {
	return &RBACManager{
		rolePermissions: map[Role][]Permission{
			RoleAdmin: {
				// 管理员拥有所有权限
				PermissionUserRead,
				PermissionUserWrite,
				PermissionUserDelete,
				PermissionUserList,
				PermissionResourceRead,
				PermissionResourceWrite,
				PermissionResourceDelete,
				PermissionResourceList,
			},
			RoleEditor: {
				// 编辑者可以读写资源
				PermissionResourceRead,
				PermissionResourceWrite,
				PermissionResourceList,
				PermissionUserRead, // 可以读取用户信息
			},
			RoleViewer: {
				// 查看者只能读取
				PermissionResourceRead,
				PermissionResourceList,
			},
			RoleUser: {
				// 普通用户可以读取自己的信息
				PermissionUserRead,
				PermissionResourceRead,
			},
		},
	}
}

// CheckPermission 检查用户是否有指定权限
func (rm *RBACManager) CheckPermission(userRoles []string, permission Permission) bool {
	for _, roleStr := range userRoles {
		role := Role(roleStr)
		permissions, exists := rm.rolePermissions[role]
		if !exists {
			continue
		}

		for _, perm := range permissions {
			if perm == permission {
				return true
			}
		}
	}

	return false
}

// HasRole 检查用户是否有指定角色
func (rm *RBACManager) HasRole(userRoles []string, requiredRole Role) bool {
	for _, roleStr := range userRoles {
		if Role(roleStr) == requiredRole {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否有任一指定角色
func (rm *RBACManager) HasAnyRole(userRoles []string, requiredRoles []Role) bool {
	for _, required := range requiredRoles {
		if rm.HasRole(userRoles, required) {
			return true
		}
	}
	return false
}
