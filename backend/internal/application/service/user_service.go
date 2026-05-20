package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/user"
	"github.com/genealogy-ma/platform/pkg/auth"
)

// UserService 用户应用服务
type UserService struct {
	repo            user.Repository
	passwordService *auth.PasswordService
	jwtService      *auth.JWTService
}

// NewUserService 创建用户应用服务
func NewUserService(repo user.Repository, jwtService *auth.JWTService) *UserService {
	return &UserService{
		repo:            repo,
		passwordService: auth.NewPasswordService(),
		jwtService:      jwtService,
	}
}

// ===== DTO 定义 =====

// UserDTO 用户详情DTO
type UserDTO struct {
	ID          int64       `json:"id"`
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	DisplayName string      `json:"display_name"`
	AvatarURL   string      `json:"avatar_url"`
	Status      string      `json:"status"`
	LastLoginAt string      `json:"last_login_at,omitempty"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Roles       []*RoleDTO  `json:"roles,omitempty"`
}

// RoleDTO 角色DTO
type RoleDTO struct {
	ID          int64           `json:"id"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	IsSystem    bool            `json:"is_system"`
	CreatedAt   string          `json:"created_at"`
	Permissions []*PermissionDTO `json:"permissions,omitempty"`
}

// PermissionDTO 权限DTO
type PermissionDTO struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Module      string `json:"module"`
	Description string `json:"description"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"display_name"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    int64   `json:"expires_in"`
	User         UserDTO `json:"user"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	DisplayName string      `json:"display_name"`
	AvatarURL   string      `json:"avatar_url"`
	Status      user.UserStatus `json:"status"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// SearchUserRequest 搜索用户请求
type SearchUserRequest struct {
	Keyword  string          `form:"keyword"`
	Status   user.UserStatus `form:"status"`
	RoleID   *int64          `form:"role_id"`
	Page     int             `form:"page,default=1"`
	PageSize int             `form:"page_size,default=20"`
	SortBy   string          `form:"sort_by"`
	SortDesc bool            `form:"sort_desc,default=false"`
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	RoleCode string `json:"role_code" binding:"required"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Code        string `json:"code" binding:"required,max=50"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Code        string `json:"code" binding:"required,max=100"`
	Name        string `json:"name" binding:"required"`
	Module      string `json:"module" binding:"required"`
	Description string `json:"description"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	Module      string `json:"module"`
	Description string `json:"description"`
}

// AssignPermissionsRequest 分配权限请求
type AssignPermissionsRequest struct {
	PermissionIDs []int64 `json:"permission_ids" binding:"required"`
}

// ===== 模型转换 =====

func toUserDTO(u *user.User) UserDTO {
	var lastLoginAt string
	if u.LastLoginAt != nil {
		lastLoginAt = u.LastLoginAt.Format(time.RFC3339)
	}

	roles := make([]*RoleDTO, len(u.Roles))
	for i, r := range u.Roles {
		roles[i] = toRoleDTO(r)
	}

	return UserDTO{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		AvatarURL:   u.AvatarURL,
		Status:      string(u.Status),
		LastLoginAt: lastLoginAt,
		CreatedAt:   u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   u.UpdatedAt.Format(time.RFC3339),
		Roles:       roles,
	}
}

func toRoleDTO(r *user.Role) *RoleDTO {
	perms := make([]*PermissionDTO, len(r.Permissions))
	for i, p := range r.Permissions {
		perms[i] = toPermissionDTO(p)
	}

	return &RoleDTO{
		ID:          r.ID,
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsSystem:    r.IsSystem,
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		Permissions: perms,
	}
}

func toPermissionDTO(p *user.Permission) *PermissionDTO {
	return &PermissionDTO{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Module:      p.Module,
		Description: p.Description,
	}
}

// ===== 服务方法 =====

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*UserDTO, error) {
	// 检查用户名是否已存在
	existingUser, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	existingUser, err = s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// 验证密码强度
	if err := s.passwordService.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// 加密密码
	passwordHash, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	u := &user.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		DisplayName:  req.DisplayName,
		Status:       user.StatusActive,
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	// 分配默认角色（verified_user）
	defaultRole, err := s.repo.FindRoleByCode(ctx, auth.RoleVerifiedUser)
	if err != nil {
		return nil, err
	}
	if defaultRole != nil {
		_ = s.repo.AssignRole(ctx, u.ID, defaultRole.ID, u.ID)
	}

	// 重新加载用户（包含角色）
	u, err = s.repo.FindByID(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	dto := toUserDTO(u)
	return &dto, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 查找用户（用户名或邮箱）
	u, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if u == nil {
		u, err = s.repo.FindByEmail(ctx, req.Username)
		if err != nil {
			return nil, err
		}
	}
	if u == nil {
		return nil, errors.New("invalid credentials")
	}

	// 检查用户状态
	if err := u.CanLogin(); err != nil {
		return nil, err
	}

	// 验证密码
	if !s.passwordService.CheckPassword(req.Password, u.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// 更新登录时间
	u.UpdateLoginTime()
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	// 生成令牌
	accessToken, err := s.jwtService.GenerateToken(u.ID, u.Username, u.GetRoleCodes())
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	userDTO := toUserDTO(u)
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    24 * 60 * 60, // 24小时
		User:         userDTO,
	}, nil
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(ctx context.Context, id int64) (*UserDTO, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}
	dto := toUserDTO(u)
	return &dto, nil
}

// Update 更新用户信息
func (s *UserService) Update(ctx context.Context, userID int64, req *UpdateUserRequest) (*UserDTO, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	if req.DisplayName != "" {
		u.DisplayName = req.DisplayName
	}
	if req.AvatarURL != "" {
		u.AvatarURL = req.AvatarURL
	}
	if req.Status != "" {
		u.Status = req.Status
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	dto := toUserDTO(u)
	return &dto, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, userID int64, req *ChangePasswordRequest) error {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.New("user not found")
	}

	// 验证旧密码
	if !s.passwordService.CheckPassword(req.OldPassword, u.PasswordHash) {
		return errors.New("old password is incorrect")
	}

	// 验证新密码强度
	if err := s.passwordService.ValidatePasswordStrength(req.NewPassword); err != nil {
		return err
	}

	// 加密新密码
	newHash, err := s.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	u.PasswordHash = newHash
	return s.repo.Update(ctx, u)
}

// Search 搜索用户
func (s *UserService) Search(ctx context.Context, req *SearchUserRequest) ([]*UserDTO, int64, error) {
	query := &user.UserSearchQuery{
		Keyword:  req.Keyword,
		Status:   req.Status,
		RoleID:   req.RoleID,
		Page:     req.Page,
		PageSize: req.PageSize,
		SortBy:   req.SortBy,
		SortDesc: req.SortDesc,
	}

	users, total, err := s.repo.Search(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]*UserDTO, len(users))
	for i, u := range users {
		dto := toUserDTO(u)
		dtos[i] = &dto
	}

	return dtos, total, nil
}

// Delete 删除用户
func (s *UserService) Delete(ctx context.Context, userID int64) error {
	return s.repo.Delete(ctx, userID)
}

// AssignRole 分配角色给用户
func (s *UserService) AssignRole(ctx context.Context, userID, assignedBy int64, roleCode string) error {
	role, err := s.repo.FindRoleByCode(ctx, roleCode)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	return s.repo.AssignRole(ctx, userID, role.ID, assignedBy)
}

// RemoveRole 移除用户角色
func (s *UserService) RemoveRole(ctx context.Context, userID int64, roleCode string) error {
	role, err := s.repo.FindRoleByCode(ctx, roleCode)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	return s.repo.RemoveRole(ctx, userID, role.ID)
}

// GetAllRoles 获取所有角色
func (s *UserService) GetAllRoles(ctx context.Context) ([]*RoleDTO, error) {
	roles, err := s.repo.FindAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]*RoleDTO, len(roles))
	for i, r := range roles {
		dtos[i] = toRoleDTO(r)
	}

	return dtos, nil
}

// GetAllPermissions 获取所有权限
func (s *UserService) GetAllPermissions(ctx context.Context) ([]*PermissionDTO, error) {
	perms, err := s.repo.FindAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]*PermissionDTO, len(perms))
	for i, p := range perms {
		dtos[i] = toPermissionDTO(p)
	}

	return dtos, nil
}

// GetRoleByID 根据ID获取角色
func (s *UserService) GetRoleByID(ctx context.Context, roleID int64) (*RoleDTO, error) {
	role, err := s.repo.FindRoleByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}
	return toRoleDTO(role), nil
}

// CreateRole 创建角色
func (s *UserService) CreateRole(ctx context.Context, req *CreateRoleRequest) (*RoleDTO, error) {
	existing, err := s.repo.FindRoleByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("role code already exists")
	}

	role := &user.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    false,
	}

	if err := role.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}

	return toRoleDTO(role), nil
}

// UpdateRole 更新角色
func (s *UserService) UpdateRole(ctx context.Context, roleID int64, req *UpdateRoleRequest) (*RoleDTO, error) {
	role, err := s.repo.FindRoleByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	if err := s.repo.UpdateRole(ctx, role); err != nil {
		return nil, err
	}

	return toRoleDTO(role), nil
}

// DeleteRole 删除角色
func (s *UserService) DeleteRole(ctx context.Context, roleID int64) error {
	role, err := s.repo.FindRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}
	if role.IsSystem {
		return errors.New("cannot delete system role")
	}

	return s.repo.DeleteRole(ctx, roleID)
}

// GetRolePermissions 获取角色权限
func (s *UserService) GetRolePermissions(ctx context.Context, roleID int64) ([]*PermissionDTO, error) {
	role, err := s.repo.FindRoleByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	perms, err := s.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*PermissionDTO, len(perms))
	for i, p := range perms {
		dtos[i] = toPermissionDTO(p)
	}

	return dtos, nil
}

// AssignPermissionsToRole 批量分配权限给角色
func (s *UserService) AssignPermissionsToRole(ctx context.Context, roleID int64, permissionIDs []int64) error {
	role, err := s.repo.FindRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	for _, pid := range permissionIDs {
		perm, err := s.repo.FindPermissionByID(ctx, pid)
		if err != nil {
			return err
		}
		if perm == nil {
			return fmt.Errorf("permission with id %d not found", pid)
		}
		if err := s.repo.AssignPermissionToRole(ctx, roleID, pid); err != nil {
			return err
		}
	}

	return nil
}

// RemovePermissionFromRole 移除角色权限
func (s *UserService) RemovePermissionFromRole(ctx context.Context, roleID, permissionID int64) error {
	role, err := s.repo.FindRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	return s.repo.RemovePermissionFromRole(ctx, roleID, permissionID)
}

// CreatePermission 创建权限
func (s *UserService) CreatePermission(ctx context.Context, req *CreatePermissionRequest) (*PermissionDTO, error) {
	existing, err := s.repo.FindPermissionByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("permission code already exists")
	}

	perm := &user.Permission{
		Code:        req.Code,
		Name:        req.Name,
		Module:      req.Module,
		Description: req.Description,
	}

	if err := perm.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreatePermission(ctx, perm); err != nil {
		return nil, err
	}

	return toPermissionDTO(perm), nil
}

// UpdatePermission 更新权限
func (s *UserService) UpdatePermission(ctx context.Context, permissionID int64, req *UpdatePermissionRequest) (*PermissionDTO, error) {
	perm, err := s.repo.FindPermissionByID(ctx, permissionID)
	if err != nil {
		return nil, err
	}
	if perm == nil {
		return nil, errors.New("permission not found")
	}

	if req.Name != "" {
		perm.Name = req.Name
	}
	if req.Module != "" {
		perm.Module = req.Module
	}
	if req.Description != "" {
		perm.Description = req.Description
	}

	if err := s.repo.UpdatePermission(ctx, perm); err != nil {
		return nil, err
	}

	return toPermissionDTO(perm), nil
}

// DeletePermission 删除权限
func (s *UserService) DeletePermission(ctx context.Context, permissionID int64) error {
	perm, err := s.repo.FindPermissionByID(ctx, permissionID)
	if err != nil {
		return err
	}
	if perm == nil {
		return errors.New("permission not found")
	}

	return s.repo.DeletePermission(ctx, permissionID)
}
