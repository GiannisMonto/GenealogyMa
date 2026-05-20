package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/user"
)

// mockRepository 模拟用户仓储
type mockRepository struct {
	users       map[int64]*user.User
	roles       map[int64]*user.Role
	permissions map[int64]*user.Permission
	nextUserID  int64
	nextRoleID  int64
	nextPermID  int64
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		users:       make(map[int64]*user.User),
		roles:       make(map[int64]*user.Role),
		permissions: make(map[int64]*user.Permission),
		nextUserID:  1,
		nextRoleID:  1,
		nextPermID:  1,
	}
}

func (m *mockRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (m *mockRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) Create(ctx context.Context, u *user.User) error {
	u.ID = m.nextUserID
	m.nextUserID++
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	m.users[u.ID] = u
	return nil
}

func (m *mockRepository) Update(ctx context.Context, u *user.User) error {
	u.UpdatedAt = time.Now()
	m.users[u.ID] = u
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id int64) error {
	delete(m.users, id)
	return nil
}

func (m *mockRepository) Search(ctx context.Context, query *user.UserSearchQuery) ([]*user.User, int64, error) {
	var result []*user.User
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, int64(len(result)), nil
}

func (m *mockRepository) FindByRole(ctx context.Context, roleCode string) ([]*user.User, error) {
	return nil, nil
}

func (m *mockRepository) GetUserRoles(ctx context.Context, userID int64) ([]*user.Role, error) {
	u, ok := m.users[userID]
	if !ok {
		return nil, nil
	}
	return u.Roles, nil
}

func (m *mockRepository) AssignRole(ctx context.Context, userID, roleID, assignedBy int64) error {
	u, ok := m.users[userID]
	if !ok {
		return errors.New("user not found")
	}
	r, ok := m.roles[roleID]
	if !ok {
		return errors.New("role not found")
	}
	u.Roles = append(u.Roles, r)
	return nil
}

func (m *mockRepository) RemoveRole(ctx context.Context, userID, roleID int64) error {
	u, ok := m.users[userID]
	if !ok {
		return errors.New("user not found")
	}
	for i, r := range u.Roles {
		if r.ID == roleID {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			break
		}
	}
	return nil
}

func (m *mockRepository) FindRoleByID(ctx context.Context, roleID int64) (*user.Role, error) {
	r, ok := m.roles[roleID]
	if !ok {
		return nil, nil
	}
	return r, nil
}

func (m *mockRepository) FindRoleByCode(ctx context.Context, code string) (*user.Role, error) {
	for _, r := range m.roles {
		if r.Code == code {
			return r, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) FindAllRoles(ctx context.Context) ([]*user.Role, error) {
	var result []*user.Role
	for _, r := range m.roles {
		result = append(result, r)
	}
	return result, nil
}

func (m *mockRepository) CreateRole(ctx context.Context, role *user.Role) error {
	role.ID = m.nextRoleID
	m.nextRoleID++
	role.CreatedAt = time.Now()
	m.roles[role.ID] = role
	return nil
}

func (m *mockRepository) UpdateRole(ctx context.Context, role *user.Role) error {
	m.roles[role.ID] = role
	return nil
}

func (m *mockRepository) DeleteRole(ctx context.Context, roleID int64) error {
	delete(m.roles, roleID)
	return nil
}

func (m *mockRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]*user.Permission, error) {
	r, ok := m.roles[roleID]
	if !ok {
		return nil, nil
	}
	return r.Permissions, nil
}

func (m *mockRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
	r, ok := m.roles[roleID]
	if !ok {
		return errors.New("role not found")
	}
	p, ok := m.permissions[permissionID]
	if !ok {
		return errors.New("permission not found")
	}
	r.Permissions = append(r.Permissions, p)
	return nil
}

func (m *mockRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID int64) error {
	r, ok := m.roles[roleID]
	if !ok {
		return errors.New("role not found")
	}
	for i, p := range r.Permissions {
		if p.ID == permissionID {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			break
		}
	}
	return nil
}

func (m *mockRepository) FindAllPermissions(ctx context.Context) ([]*user.Permission, error) {
	var result []*user.Permission
	for _, p := range m.permissions {
		result = append(result, p)
	}
	return result, nil
}

func (m *mockRepository) FindPermissionByID(ctx context.Context, permissionID int64) (*user.Permission, error) {
	p, ok := m.permissions[permissionID]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *mockRepository) FindPermissionByCode(ctx context.Context, code string) (*user.Permission, error) {
	for _, p := range m.permissions {
		if p.Code == code {
			return p, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) CreatePermission(ctx context.Context, p *user.Permission) error {
	p.ID = m.nextPermID
	m.nextPermID++
	m.permissions[p.ID] = p
	return nil
}

func (m *mockRepository) UpdatePermission(ctx context.Context, p *user.Permission) error {
	m.permissions[p.ID] = p
	return nil
}

func (m *mockRepository) DeletePermission(ctx context.Context, permissionID int64) error {
	delete(m.permissions, permissionID)
	return nil
}

// ===== 测试用例 =====

func newTestUserService() *UserService {
	repo := newMockRepository()
	return &UserService{
		repo:            repo,
		passwordService: nil,
		jwtService:      nil,
	}
}

func newTestUserServiceWithRepo() (*UserService, *mockRepository) {
	repo := newMockRepository()
	return &UserService{
		repo:            repo,
		passwordService: nil,
		jwtService:      nil,
	}, repo
}

func TestCreateRole(t *testing.T) {
	svc, repo := newTestUserServiceWithRepo()
	ctx := context.Background()

	t.Run("create role successfully", func(t *testing.T) {
		req := &CreateRoleRequest{
			Code:        "editor",
			Name:        "Editor",
			Description: "Content editor role",
		}

		role, err := svc.CreateRole(ctx, req)
		if err != nil {
			t.Fatalf("CreateRole() error = %v", err)
		}
		if role.Code != "editor" {
			t.Errorf("expected code 'editor', got '%s'", role.Code)
		}
		if role.Name != "Editor" {
			t.Errorf("expected name 'Editor', got '%s'", role.Name)
		}
		if role.IsSystem {
			t.Error("new role should not be system role")
		}
	})

	t.Run("create role with duplicate code", func(t *testing.T) {
		repo.CreateRole(ctx, &user.Role{Code: "editor", Name: "Existing Editor"})

		req := &CreateRoleRequest{
			Code: "editor",
			Name: "New Editor",
		}

		_, err := svc.CreateRole(ctx, req)
		if err == nil {
			t.Error("expected error for duplicate code")
		}
	})

	t.Run("create role with empty code", func(t *testing.T) {
		req := &CreateRoleRequest{
			Code: "",
			Name: "No Code",
		}

		_, err := svc.CreateRole(ctx, req)
		if err == nil {
			t.Error("expected error for empty code")
		}
	})
}

func TestUpdateRole(t *testing.T) {
	svc, repo := newTestUserServiceWithRepo()
	ctx := context.Background()

	repo.CreateRole(ctx, &user.Role{ID: 1, Code: "editor", Name: "Editor"})

	t.Run("update role successfully", func(t *testing.T) {
		req := &UpdateRoleRequest{
			Name:        "Senior Editor",
			Description: "Updated description",
		}

		role, err := svc.UpdateRole(ctx, 1, req)
		if err != nil {
			t.Fatalf("UpdateRole() error = %v", err)
		}
		if role.Name != "Senior Editor" {
			t.Errorf("expected name 'Senior Editor', got '%s'", role.Name)
		}
	})

	t.Run("update non-existent role", func(t *testing.T) {
		req := &UpdateRoleRequest{Name: "New Name"}
		_, err := svc.UpdateRole(ctx, 999, req)
		if err == nil {
			t.Error("expected error for non-existent role")
		}
	})
}

func TestDeleteRole(t *testing.T) {
	svc, repo := newTestUserServiceWithRepo()
	ctx := context.Background()

	repo.CreateRole(ctx, &user.Role{ID: 1, Code: "custom", Name: "Custom", IsSystem: false})
	repo.CreateRole(ctx, &user.Role{ID: 2, Code: "admin", Name: "Admin", IsSystem: true})

	t.Run("delete non-system role", func(t *testing.T) {
		err := svc.DeleteRole(ctx, 1)
		if err != nil {
			t.Fatalf("DeleteRole() error = %v", err)
		}
	})

	t.Run("delete system role should fail", func(t *testing.T) {
		err := svc.DeleteRole(ctx, 2)
		if err == nil {
			t.Error("expected error when deleting system role")
		}
	})

	t.Run("delete non-existent role", func(t *testing.T) {
		err := svc.DeleteRole(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent role")
		}
	})
}

func TestGetRoleByID(t *testing.T) {
	svc, repo := newTestUserServiceWithRepo()
	ctx := context.Background()

	repo.CreateRole(ctx, &user.Role{ID: 1, Code: "editor", Name: "Editor"})

	t.Run("get existing role", func(t *testing.T) {
		role, err := svc.GetRoleByID(ctx, 1)
		if err != nil {
			t.Fatalf("GetRoleByID() error = %v", err)
		}
		if role.Code != "editor" {
			t.Errorf("expected code 'editor', got '%s'", role.Code)
		}
	})

	t.Run("get non-existent role", func(t *testing.T) {
		_, err := svc.GetRoleByID(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent role")
		}
	})
}

func TestAssignPermissionsToRole(t *testing.T) {
	svc, repo := newTestUserServiceWithRepo()
	ctx := context.Background()

	repo.CreateRole(ctx, &user.Role{ID: 1, Code: "editor", Name: "Editor"})
	repo.CreatePermission(ctx, &user.Permission{ID: 1, Code: "person:read", Name: "Read Person", Module: "person"})
	repo.CreatePermission(ctx, &user.Permission{ID: 2, Code: "person:write", Name: "Write Person", Module: "person"})

	t.Run("assign permissions successfully", func(t *testing.T) {
		err := svc.AssignPermissionsToRole(ctx, 1, []int64{1, 2})
		if err != nil {
			t.Fatalf("AssignPermissionsToRole() error = %v", err)
		}
	})

	t.Run("assign to non-existent role", func(t *testing.T) {
		err := svc.AssignPermissionsToRole(ctx, 999, []int64{1})
		if err == nil {
			t.Error("expected error for non-existent role")
		}
	})

	t.Run("assign non-existent permission", func(t *testing.T) {
		err := svc.AssignPermissionsToRole(ctx, 1, []int64{999})
		if err == nil {
			t.Error("expected error for non-existent permission")
		}
	})
}

func TestRemovePermissionFromRole(t *testing.T) {
	svc, repo := newTestUserServiceWithRepo()
	ctx := context.Background()

	perm := &user.Permission{ID: 1, Code: "person:read", Name: "Read Person", Module: "person"}
	repo.CreatePermission(ctx, perm)
	repo.CreateRole(ctx, &user.Role{ID: 1, Code: "editor", Name: "Editor", Permissions: []*user.Permission{perm}})

	t.Run("remove permission successfully", func(t *testing.T) {
		err := svc.RemovePermissionFromRole(ctx, 1, 1)
		if err != nil {
			t.Fatalf("RemovePermissionFromRole() error = %v", err)
		}
	})

	t.Run("remove from non-existent role", func(t *testing.T) {
		err := svc.RemovePermissionFromRole(ctx, 999, 1)
		if err == nil {
			t.Error("expected error for non-existent role")
		}
	})
}

func TestCreatePermission(t *testing.T) {
	svc, repo := newTestUserServiceWithRepo()
	ctx := context.Background()

	t.Run("create permission successfully", func(t *testing.T) {
		req := &CreatePermissionRequest{
			Code:        "test:read",
			Name:        "Test Read",
			Module:      "test",
			Description: "Read test data",
		}

		perm, err := svc.CreatePermission(ctx, req)
		if err != nil {
			t.Fatalf("CreatePermission() error = %v", err)
		}
		if perm.Code != "test:read" {
			t.Errorf("expected code 'test:read', got '%s'", perm.Code)
		}
		if perm.Module != "test" {
			t.Errorf("expected module 'test', got '%s'", perm.Module)
		}
	})

	t.Run("create permission with duplicate code", func(t *testing.T) {
		repo.CreatePermission(ctx, &user.Permission{Code: "test:read", Name: "Existing", Module: "test"})

		req := &CreatePermissionRequest{
			Code:   "test:read",
			Name:   "Duplicate",
			Module: "test",
		}

		_, err := svc.CreatePermission(ctx, req)
		if err == nil {
			t.Error("expected error for duplicate code")
		}
	})

	t.Run("create permission with empty code", func(t *testing.T) {
		req := &CreatePermissionRequest{
			Code:   "",
			Name:   "No Code",
			Module: "test",
		}

		_, err := svc.CreatePermission(ctx, req)
		if err == nil {
			t.Error("expected error for empty code")
		}
	})
}

func TestUpdatePermission(t *testing.T) {
	svc, repo := newTestUserServiceWithRepo()
	ctx := context.Background()

	repo.CreatePermission(ctx, &user.Permission{ID: 1, Code: "test:read", Name: "Test Read", Module: "test"})

	t.Run("update permission successfully", func(t *testing.T) {
		req := &UpdatePermissionRequest{
			Name:        "Updated Read",
			Description: "Updated description",
		}

		perm, err := svc.UpdatePermission(ctx, 1, req)
		if err != nil {
			t.Fatalf("UpdatePermission() error = %v", err)
		}
		if perm.Name != "Updated Read" {
			t.Errorf("expected name 'Updated Read', got '%s'", perm.Name)
		}
	})

	t.Run("update non-existent permission", func(t *testing.T) {
		req := &UpdatePermissionRequest{Name: "New Name"}
		_, err := svc.UpdatePermission(ctx, 999, req)
		if err == nil {
			t.Error("expected error for non-existent permission")
		}
	})
}

func TestDeletePermission(t *testing.T) {
	svc, repo := newTestUserServiceWithRepo()
	ctx := context.Background()

	repo.CreatePermission(ctx, &user.Permission{ID: 1, Code: "test:read", Name: "Test Read", Module: "test"})

	t.Run("delete permission successfully", func(t *testing.T) {
		err := svc.DeletePermission(ctx, 1)
		if err != nil {
			t.Fatalf("DeletePermission() error = %v", err)
		}
	})

	t.Run("delete non-existent permission", func(t *testing.T) {
		err := svc.DeletePermission(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent permission")
		}
	})
}

func TestGetAllRoles(t *testing.T) {
	svc, repo := newTestUserServiceWithRepo()
	ctx := context.Background()

	repo.CreateRole(ctx, &user.Role{Code: "admin", Name: "Admin"})
	repo.CreateRole(ctx, &user.Role{Code: "editor", Name: "Editor"})

	roles, err := svc.GetAllRoles(ctx)
	if err != nil {
		t.Fatalf("GetAllRoles() error = %v", err)
	}
	if len(roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(roles))
	}
}

func TestGetAllPermissions(t *testing.T) {
	svc, repo := newTestUserServiceWithRepo()
	ctx := context.Background()

	repo.CreatePermission(ctx, &user.Permission{Code: "test:read", Name: "Test Read", Module: "test"})
	repo.CreatePermission(ctx, &user.Permission{Code: "test:write", Name: "Test Write", Module: "test"})

	perms, err := svc.GetAllPermissions(ctx)
	if err != nil {
		t.Fatalf("GetAllPermissions() error = %v", err)
	}
	if len(perms) != 2 {
		t.Errorf("expected 2 permissions, got %d", len(perms))
	}
}

func TestGetRolePermissions(t *testing.T) {
	svc, repo := newTestUserServiceWithRepo()
	ctx := context.Background()

	perm := &user.Permission{ID: 1, Code: "person:read", Name: "Read Person", Module: "person"}
	repo.CreatePermission(ctx, perm)
	repo.CreateRole(ctx, &user.Role{ID: 1, Code: "editor", Name: "Editor", Permissions: []*user.Permission{perm}})

	t.Run("get role permissions", func(t *testing.T) {
		perms, err := svc.GetRolePermissions(ctx, 1)
		if err != nil {
			t.Fatalf("GetRolePermissions() error = %v", err)
		}
		if len(perms) != 1 {
			t.Errorf("expected 1 permission, got %d", len(perms))
		}
		if perms[0].Code != "person:read" {
			t.Errorf("expected code 'person:read', got '%s'", perms[0].Code)
		}
	})

	t.Run("get permissions for non-existent role", func(t *testing.T) {
		_, err := svc.GetRolePermissions(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent role")
		}
	})
}
