package user

import (
	"context"
)

// Repository 用户仓储接口
type Repository interface {
	// 基础CRUD
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error

	// 查询
	Search(ctx context.Context, query *UserSearchQuery) ([]*User, int64, error)
	FindByRole(ctx context.Context, roleCode string) ([]*User, error)

	// 角色管理
	GetUserRoles(ctx context.Context, userID int64) ([]*Role, error)
	AssignRole(ctx context.Context, userID, roleID, assignedBy int64) error
	RemoveRole(ctx context.Context, userID, roleID int64) error

	// 角色仓储方法
	FindRoleByID(ctx context.Context, roleID int64) (*Role, error)
	FindRoleByCode(ctx context.Context, code string) (*Role, error)
	FindAllRoles(ctx context.Context) ([]*Role, error)
	CreateRole(ctx context.Context, role *Role) error
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, roleID int64) error

	// 权限管理
	GetRolePermissions(ctx context.Context, roleID int64) ([]*Permission, error)
	AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID int64) error
	FindAllPermissions(ctx context.Context) ([]*Permission, error)
	FindPermissionByID(ctx context.Context, permissionID int64) (*Permission, error)
	FindPermissionByCode(ctx context.Context, code string) (*Permission, error)
	CreatePermission(ctx context.Context, permission *Permission) error
	UpdatePermission(ctx context.Context, permission *Permission) error
	DeletePermission(ctx context.Context, permissionID int64) error
}

// UserSearchQuery 用户搜索查询参数
type UserSearchQuery struct {
	Keyword string     `json:"keyword"`
	Status  UserStatus `json:"status"`
	RoleID  *int64     `json:"role_id"`

	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	SortBy   string `json:"sort_by"`
	SortDesc bool   `json:"sort_desc"`
}
