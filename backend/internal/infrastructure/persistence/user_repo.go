package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/user"
	"gorm.io/gorm"
)

// UserModel GORM模型对应users表
type UserModel struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement"`
	Username     string     `gorm:"column:username;type:varchar(50);unique;not null"`
	Email        string     `gorm:"column:email;type:varchar(100);unique;not null"`
	PasswordHash string     `gorm:"column:password_hash;type:varchar(255);not null"`
	DisplayName  string     `gorm:"column:display_name;type:varchar(100)"`
	AvatarURL    string     `gorm:"column:avatar_url;type:varchar(255)"`
	Status       string     `gorm:"column:status;type:varchar(20);default:active"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserModel) TableName() string {
	return "users"
}

// RoleModel GORM模型对应roles表
type RoleModel struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Code        string    `gorm:"column:code;type:varchar(50);unique;not null"`
	Name        string    `gorm:"column:name;type:varchar(100);not null"`
	Description string    `gorm:"column:description;type:text"`
	IsSystem    bool      `gorm:"column:is_system;default:false"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (RoleModel) TableName() string {
	return "roles"
}

// UserRoleModel GORM模型对应user_roles表
type UserRoleModel struct {
	UserID     int64     `gorm:"column:user_id;primaryKey"`
	RoleID     int64     `gorm:"column:role_id;primaryKey"`
	AssignedAt time.Time `gorm:"column:assigned_at;autoCreateTime"`
	AssignedBy int64     `gorm:"column:assigned_by"`
}

func (UserRoleModel) TableName() string {
	return "user_roles"
}

// PermissionModel GORM模型对应permissions表
type PermissionModel struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Code        string `gorm:"column:code;type:varchar(100);unique;not null"`
	Name        string `gorm:"column:name;type:varchar(100);not null"`
	Module      string `gorm:"column:module;type:varchar(50)"`
	Description string `gorm:"column:description;type:text"`
}

func (PermissionModel) TableName() string {
	return "permissions"
}

// RolePermissionModel GORM模型对应role_permissions表
type RolePermissionModel struct {
	RoleID       int64 `gorm:"column:role_id;primaryKey"`
	PermissionID int64 `gorm:"column:permission_id;primaryKey"`
}

func (RolePermissionModel) TableName() string {
	return "role_permissions"
}

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepositoryImpl{db: db}
}

// 模型转换函数
func modelToUser(m *UserModel) *user.User {
	return &user.User{
		ID:           m.ID,
		Username:     m.Username,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		DisplayName:  m.DisplayName,
		AvatarURL:    m.AvatarURL,
		Status:       user.UserStatus(m.Status),
		LastLoginAt:  m.LastLoginAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func userToModel(u *user.User) *UserModel {
	return &UserModel{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		DisplayName:  u.DisplayName,
		AvatarURL:    u.AvatarURL,
		Status:       string(u.Status),
		LastLoginAt:  u.LastLoginAt,
	}
}

func modelToRole(m *RoleModel) *user.Role {
	return &user.Role{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		IsSystem:    m.IsSystem,
		CreatedAt:   m.CreatedAt,
	}
}

func modelToPermission(m *PermissionModel) *user.Permission {
	return &user.Permission{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Module:      m.Module,
		Description: m.Description,
	}
}

// FindByID 根据ID查找用户
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int64) (*user.User, error) {
	var m UserModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	u := modelToUser(&m)
	// 加载用户角色
	roles, err := r.GetUserRoles(ctx, id)
	if err != nil {
		return nil, err
	}
	u.Roles = roles
	return u, nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var m UserModel
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	u := modelToUser(&m)
	// 加载用户角色
	roles, err := r.GetUserRoles(ctx, m.ID)
	if err != nil {
		return nil, err
	}
	u.Roles = roles
	return u, nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var m UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	u := modelToUser(&m)
	// 加载用户角色
	roles, err := r.GetUserRoles(ctx, m.ID)
	if err != nil {
		return nil, err
	}
	u.Roles = roles
	return u, nil
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(ctx context.Context, u *user.User) error {
	m := userToModel(u)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	u.ID = m.ID
	u.CreatedAt = m.CreatedAt
	u.UpdatedAt = m.UpdatedAt
	return nil
}

// Update 更新用户
func (r *UserRepositoryImpl) Update(ctx context.Context, u *user.User) error {
	m := userToModel(u)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return err
	}
	u.UpdatedAt = m.UpdatedAt
	return nil
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除用户角色关联
		if err := tx.Where("user_id = ?", id).Delete(&UserRoleModel{}).Error; err != nil {
			return err
		}
		// 再删除用户
		return tx.Delete(&UserModel{}, id).Error
	})
}

// Search 搜索用户
func (r *UserRepositoryImpl) Search(ctx context.Context, query *user.UserSearchQuery) ([]*user.User, int64, error) {
	db := r.db.WithContext(ctx).Model(&UserModel{})

	if query.Keyword != "" {
		db = db.Where("username LIKE ? OR email LIKE ? OR display_name LIKE ?",
			"%"+query.Keyword+"%",
			"%"+query.Keyword+"%",
			"%"+query.Keyword+"%")
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.RoleID != nil {
		subQuery := r.db.Model(&UserRoleModel{}).
			Select("user_id").
			Where("role_id = ?", *query.RoleID)
		db = db.Where("id IN (?)", subQuery)
	}

	// 统计总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if query.SortBy != "" {
		order := query.SortBy
		if query.SortDesc {
			order += " DESC"
		} else {
			order += " ASC"
		}
		db = db.Order(order)
	} else {
		db = db.Order("id DESC")
	}

	// 分页
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	var models []*UserModel
	if err := db.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*user.User, len(models))
	for i, m := range models {
		users[i] = modelToUser(m)
		// 加载角色
		roles, err := r.GetUserRoles(ctx, m.ID)
		if err == nil {
			users[i].Roles = roles
		}
	}

	return users, total, nil
}

// FindByRole 根据角色查找用户
func (r *UserRepositoryImpl) FindByRole(ctx context.Context, roleCode string) ([]*user.User, error) {
	var models []*UserModel

	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ur ON users.id = ur.user_id").
		Joins("JOIN roles r ON ur.role_id = r.id").
		Where("r.code = ?", roleCode).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	users := make([]*user.User, len(models))
	for i, m := range models {
		users[i] = modelToUser(m)
	}
	return users, nil
}

// GetUserRoles 获取用户角色
func (r *UserRepositoryImpl) GetUserRoles(ctx context.Context, userID int64) ([]*user.Role, error) {
	var roleModels []*RoleModel

	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ur ON roles.id = ur.role_id").
		Where("ur.user_id = ?", userID).
		Find(&roleModels).Error

	if err != nil {
		return nil, err
	}

	roles := make([]*user.Role, len(roleModels))
	for i, m := range roleModels {
		roles[i] = modelToRole(m)
		// 加载角色权限
		perms, err := r.GetRolePermissions(ctx, m.ID)
		if err == nil {
			roles[i].Permissions = perms
		}
	}

	return roles, nil
}

// AssignRole 分配角色给用户
func (r *UserRepositoryImpl) AssignRole(ctx context.Context, userID, roleID, assignedBy int64) error {
	ur := &UserRoleModel{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: assignedBy,
	}
	return r.db.WithContext(ctx).FirstOrCreate(ur).Error
}

// RemoveRole 移除用户角色
func (r *UserRepositoryImpl) RemoveRole(ctx context.Context, userID, roleID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&UserRoleModel{}).Error
}

// FindRoleByID 根据ID查找角色
func (r *UserRepositoryImpl) FindRoleByID(ctx context.Context, roleID int64) (*user.Role, error) {
	var m RoleModel
	if err := r.db.WithContext(ctx).First(&m, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return modelToRole(&m), nil
}

// FindRoleByCode 根据代码查找角色
func (r *UserRepositoryImpl) FindRoleByCode(ctx context.Context, code string) (*user.Role, error) {
	var m RoleModel
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return modelToRole(&m), nil
}

// FindAllRoles 获取所有角色
func (r *UserRepositoryImpl) FindAllRoles(ctx context.Context) ([]*user.Role, error) {
	var models []*RoleModel
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	roles := make([]*user.Role, len(models))
	for i, m := range models {
		roles[i] = modelToRole(m)
		perms, err := r.GetRolePermissions(ctx, m.ID)
		if err == nil {
			roles[i].Permissions = perms
		}
	}
	return roles, nil
}

// CreateRole 创建角色
func (r *UserRepositoryImpl) CreateRole(ctx context.Context, role *user.Role) error {
	m := &RoleModel{
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	role.ID = m.ID
	role.CreatedAt = m.CreatedAt
	return nil
}

// UpdateRole 更新角色
func (r *UserRepositoryImpl) UpdateRole(ctx context.Context, role *user.Role) error {
	m := &RoleModel{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
	}
	return r.db.WithContext(ctx).Save(m).Error
}

// DeleteRole 删除角色
func (r *UserRepositoryImpl) DeleteRole(ctx context.Context, roleID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查是否是系统角色
		var role RoleModel
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}
		if role.IsSystem {
			return errors.New("cannot delete system role")
		}

		// 删除角色权限关联
		if err := tx.Where("role_id = ?", roleID).Delete(&RolePermissionModel{}).Error; err != nil {
			return err
		}
		// 删除用户角色关联
		if err := tx.Where("role_id = ?", roleID).Delete(&UserRoleModel{}).Error; err != nil {
			return err
		}
		// 删除角色
		return tx.Delete(&RoleModel{}, roleID).Error
	})
}

// GetRolePermissions 获取角色权限
func (r *UserRepositoryImpl) GetRolePermissions(ctx context.Context, roleID int64) ([]*user.Permission, error) {
	var models []*PermissionModel

	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions rp ON permissions.id = rp.permission_id").
		Where("rp.role_id = ?", roleID).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	perms := make([]*user.Permission, len(models))
	for i, m := range models {
		perms[i] = modelToPermission(m)
	}
	return perms, nil
}

// AssignPermissionToRole 分配权限给角色
func (r *UserRepositoryImpl) AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
	rp := &RolePermissionModel{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return r.db.WithContext(ctx).FirstOrCreate(rp).Error
}

// RemovePermissionFromRole 移除角色权限
func (r *UserRepositoryImpl) RemovePermissionFromRole(ctx context.Context, roleID, permissionID int64) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&RolePermissionModel{}).Error
}

// FindAllPermissions 获取所有权限
func (r *UserRepositoryImpl) FindAllPermissions(ctx context.Context) ([]*user.Permission, error) {
	var models []*PermissionModel
	if err := r.db.WithContext(ctx).Order("module, code ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	perms := make([]*user.Permission, len(models))
	for i, m := range models {
		perms[i] = modelToPermission(m)
	}
	return perms, nil
}

// FindPermissionByID 根据ID查找权限
func (r *UserRepositoryImpl) FindPermissionByID(ctx context.Context, permissionID int64) (*user.Permission, error) {
	var m PermissionModel
	if err := r.db.WithContext(ctx).First(&m, permissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return modelToPermission(&m), nil
}

// FindPermissionByCode 根据代码查找权限
func (r *UserRepositoryImpl) FindPermissionByCode(ctx context.Context, code string) (*user.Permission, error) {
	var m PermissionModel
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return modelToPermission(&m), nil
}

// CreatePermission 创建权限
func (r *UserRepositoryImpl) CreatePermission(ctx context.Context, p *user.Permission) error {
	m := &PermissionModel{
		Code:        p.Code,
		Name:        p.Name,
		Module:      p.Module,
		Description: p.Description,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	p.ID = m.ID
	return nil
}

// UpdatePermission 更新权限
func (r *UserRepositoryImpl) UpdatePermission(ctx context.Context, p *user.Permission) error {
	m := &PermissionModel{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Module:      p.Module,
		Description: p.Description,
	}
	return r.db.WithContext(ctx).Save(m).Error
}

// DeletePermission 删除权限
func (r *UserRepositoryImpl) DeletePermission(ctx context.Context, permissionID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除角色权限关联
		if err := tx.Where("permission_id = ?", permissionID).Delete(&RolePermissionModel{}).Error; err != nil {
			return err
		}
		// 删除权限
		return tx.Delete(&PermissionModel{}, permissionID).Error
	})
}
