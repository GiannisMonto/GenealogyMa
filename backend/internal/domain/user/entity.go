package user

import (
	"fmt"
	"time"
)

// UserStatus 用户状态
type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusBanned   UserStatus = "banned"
)

// User 用户聚合根
type User struct {
	ID           int64       `json:"id"`
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"-"` // 不序列化到JSON
	DisplayName  string      `json:"display_name"`
	AvatarURL    string      `json:"avatar_url"`
	Status       UserStatus  `json:"status"`
	LastLoginAt  *time.Time  `json:"last_login_at"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`

	// 关联数据
	Roles []*Role `json:"roles,omitempty"`
}

// Role 角色实体
type Role struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsSystem    bool   `json:"is_system"` // 系统角色不可删除
	CreatedAt   time.Time `json:"created_at"`

	// 关联数据
	Permissions []*Permission `json:"permissions,omitempty"`
}

// Permission 权限实体
type Permission struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Module      string `json:"module"`
	Description string `json:"description"`
}

// UserRole 用户角色关联
type UserRole struct {
	UserID     int64     `json:"user_id"`
	RoleID     int64     `json:"role_id"`
	AssignedAt time.Time `json:"assigned_at"`
	AssignedBy int64     `json:"assigned_by"`

	// 关联数据
	Role *Role `json:"role,omitempty"`
}

// 领域方法

// IsActive 检查用户是否活跃
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// CanLogin 检查用户是否可以登录
func (u *User) CanLogin() error {
	if u.Status == StatusBanned {
		return fmt.Errorf("user is banned")
	}
	if u.Status == StatusInactive {
		return fmt.Errorf("user is inactive")
	}
	return nil
}

// UpdateLoginTime 更新登录时间
func (u *User) UpdateLoginTime() {
	now := time.Now()
	u.LastLoginAt = &now
}

// HasRole 检查用户是否拥有指定角色
func (u *User) HasRole(roleCode string) bool {
	for _, role := range u.Roles {
		if role.Code == roleCode {
			return true
		}
	}
	return false
}

// GetRoleCodes 获取用户所有角色代码
func (u *User) GetRoleCodes() []string {
	codes := make([]string, 0, len(u.Roles))
	for _, role := range u.Roles {
		codes = append(codes, role.Code)
	}
	return codes
}

// GetPermissions 获取用户所有权限（基于角色）
func (u *User) GetPermissions() []string {
	permMap := make(map[string]bool)
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			permMap[perm.Code] = true
		}
	}

	perms := make([]string, 0, len(permMap))
	for p := range permMap {
		perms = append(perms, p)
	}
	return perms
}

// HasPermission 检查用户是否拥有指定权限
func (u *User) HasPermission(permissionCode string) bool {
	// 超级管理员拥有所有权限
	if u.HasRole("super_admin") {
		return true
	}

	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if perm.Code == permissionCode {
				return true
			}
		}
	}
	return false
}

// Validate 验证角色领域 invariants
func (r *Role) Validate() error {
	if r.Code == "" {
		return fmt.Errorf("role code cannot be empty")
	}
	if len(r.Code) > 50 {
		return fmt.Errorf("role code must be at most 50 characters")
	}
	if r.Name == "" {
		return fmt.Errorf("role name cannot be empty")
	}
	return nil
}

// Validate 验证权限领域 invariants
func (p *Permission) Validate() error {
	if p.Code == "" {
		return fmt.Errorf("permission code cannot be empty")
	}
	if len(p.Code) > 100 {
		return fmt.Errorf("permission code must be at most 100 characters")
	}
	if p.Name == "" {
		return fmt.Errorf("permission name cannot be empty")
	}
	if p.Module == "" {
		return fmt.Errorf("permission module cannot be empty")
	}
	return nil
}

// Validate 验证用户领域 invariants
func (u *User) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if len(u.Username) < 3 || len(u.Username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}
	if u.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if u.Status != StatusActive && u.Status != StatusInactive && u.Status != StatusBanned {
		return fmt.Errorf("invalid user status: %s", u.Status)
	}
	return nil
}
