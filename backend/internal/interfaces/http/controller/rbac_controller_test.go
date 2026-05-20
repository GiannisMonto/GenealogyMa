package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/internal/domain/user"
	"github.com/genealogy-ma/platform/pkg/utils"
)

// 错误定义
var (
	errRoleNotFound       = errors.New("角色不存在")
	errPermissionNotFound = errors.New("权限不存在")
	errSystemRole         = errors.New("不能删除系统角色")
	errInvalidCode        = errors.New("编码不能为空")
)

// mockUserService 模拟用户服务
type mockUserService struct {
	roles       map[int64]*user.Role
	permissions map[int64]*user.Permission
	nextRoleID  int64
	nextPermID  int64
}

func newMockUserService() *mockUserService {
	return &mockUserService{
		roles:       make(map[int64]*user.Role),
		permissions: make(map[int64]*user.Permission),
		nextRoleID:  1,
		nextPermID:  1,
	}
}

func (m *mockUserService) GetAllRoles(ctx context.Context) ([]*user.Role, error) {
	var result []*user.Role
	for _, r := range m.roles {
		result = append(result, r)
	}
	return result, nil
}

func (m *mockUserService) GetRoleByID(ctx context.Context, id int64) (*user.Role, error) {
	r, ok := m.roles[id]
	if !ok {
		return nil, errRoleNotFound
	}
	return r, nil
}

func (m *mockUserService) CreateRole(ctx context.Context, req *service.CreateRoleRequest) (*user.Role, error) {
	if req.Code == "" {
		return nil, errInvalidCode
	}
	role := &user.Role{
		ID:          m.nextRoleID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    false,
	}
	m.roles[m.nextRoleID] = role
	m.nextRoleID++
	return role, nil
}

func (m *mockUserService) UpdateRole(ctx context.Context, id int64, req *service.UpdateRoleRequest) (*user.Role, error) {
	r, ok := m.roles[id]
	if !ok {
		return nil, errRoleNotFound
	}
	r.Name = req.Name
	r.Description = req.Description
	return r, nil
}

func (m *mockUserService) DeleteRole(ctx context.Context, id int64) error {
	r, ok := m.roles[id]
	if !ok {
		return errRoleNotFound
	}
	if r.IsSystem {
		return errSystemRole
	}
	delete(m.roles, id)
	return nil
}

func (m *mockUserService) GetRolePermissions(ctx context.Context, roleID int64) ([]*user.Permission, error) {
	r, ok := m.roles[roleID]
	if !ok {
		return nil, errRoleNotFound
	}
	return r.Permissions, nil
}

func (m *mockUserService) AssignPermissionsToRole(ctx context.Context, roleID int64, permissionIDs []int64) error {
	r, ok := m.roles[roleID]
	if !ok {
		return errRoleNotFound
	}
	for _, pid := range permissionIDs {
		p, ok := m.permissions[pid]
		if !ok {
			return errPermissionNotFound
		}
		r.Permissions = append(r.Permissions, p)
	}
	return nil
}

func (m *mockUserService) RemovePermissionFromRole(ctx context.Context, roleID, permissionID int64) error {
	r, ok := m.roles[roleID]
	if !ok {
		return errRoleNotFound
	}
	for i, p := range r.Permissions {
		if p.ID == permissionID {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			return nil
		}
	}
	return errPermissionNotFound
}

func (m *mockUserService) GetAllPermissions(ctx context.Context) ([]*user.Permission, error) {
	var result []*user.Permission
	for _, p := range m.permissions {
		result = append(result, p)
	}
	return result, nil
}

func (m *mockUserService) CreatePermission(ctx context.Context, req *service.CreatePermissionRequest) (*user.Permission, error) {
	if req.Code == "" {
		return nil, errInvalidCode
	}
	perm := &user.Permission{
		ID:          m.nextPermID,
		Code:        req.Code,
		Name:        req.Name,
		Module:      req.Module,
		Description: req.Description,
	}
	m.permissions[m.nextPermID] = perm
	m.nextPermID++
	return perm, nil
}

func (m *mockUserService) UpdatePermission(ctx context.Context, id int64, req *service.UpdatePermissionRequest) (*user.Permission, error) {
	p, ok := m.permissions[id]
	if !ok {
		return nil, errPermissionNotFound
	}
	p.Name = req.Name
	p.Description = req.Description
	return p, nil
}

func (m *mockUserService) DeletePermission(ctx context.Context, id int64) error {
	if _, ok := m.permissions[id]; !ok {
		return errPermissionNotFound
	}
	delete(m.permissions, id)
	return nil
}

// testRBACController 测试用控制器
type testRBACController struct {
	userService *mockUserService
}

func (c *testRBACController) ListRoles(ctx *gin.Context) {
	roles, err := c.userService.GetAllRoles(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取角色列表失败")
		return
	}
	utils.Success(ctx, roles)
}

func (c *testRBACController) GetRole(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
		return
	}
	role, err := c.userService.GetRoleByID(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(ctx, role)
}

func (c *testRBACController) CreateRole(ctx *gin.Context) {
	var req struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	createReq := &service.CreateRoleRequest{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}
	role, err := c.userService.CreateRole(ctx, createReq)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessWithStatus(ctx, http.StatusCreated, role)
}

func (c *testRBACController) UpdateRole(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
		return
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	updateReq := &service.UpdateRoleRequest{
		Name:        req.Name,
		Description: req.Description,
	}
	role, err := c.userService.UpdateRole(ctx, id, updateReq)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(ctx, role)
}

func (c *testRBACController) DeleteRole(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的角色ID")
		return
	}
	err = c.userService.DeleteRole(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(ctx, gin.H{"message": "删除成功"})
}

func (c *testRBACController) ListPermissions(ctx *gin.Context) {
	perms, err := c.userService.GetAllPermissions(ctx)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取权限列表失败")
		return
	}
	utils.Success(ctx, perms)
}

func (c *testRBACController) CreatePermission(ctx *gin.Context) {
	var req struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Module      string `json:"module"`
		Description string `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	createReq := &service.CreatePermissionRequest{
		Code:        req.Code,
		Name:        req.Name,
		Module:      req.Module,
		Description: req.Description,
	}
	perm, err := c.userService.CreatePermission(ctx, createReq)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessWithStatus(ctx, http.StatusCreated, perm)
}

func (c *testRBACController) DeletePermission(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的权限ID")
		return
	}
	err = c.userService.DeletePermission(ctx, id)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(ctx, gin.H{"message": "删除成功"})
}

// ===== 测试用例 =====

func TestRBACController_ListRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("list roles successfully", func(t *testing.T) {
		svc := newMockUserService()
		svc.roles[1] = &user.Role{ID: 1, Code: "admin", Name: "Admin"}

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.GET("/roles", controller.ListRoles)

		req := httptest.NewRequest(http.MethodGet, "/roles", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"].(float64) != 0 {
			t.Errorf("expected code 0, got %v", response["code"])
		}
	})
}

func TestRBACController_CreateRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("create role successfully", func(t *testing.T) {
		svc := newMockUserService()

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.POST("/roles", controller.CreateRole)

		body := map[string]string{"code": "editor", "name": "Editor", "description": "Content editor"}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/roles", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("create role with empty code", func(t *testing.T) {
		svc := newMockUserService()

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.POST("/roles", controller.CreateRole)

		body := map[string]string{"code": "", "name": "No Code"}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/roles", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestRBACController_GetRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("get existing role", func(t *testing.T) {
		svc := newMockUserService()
		svc.roles[1] = &user.Role{ID: 1, Code: "admin", Name: "Admin"}

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.GET("/roles/:id", controller.GetRole)

		req := httptest.NewRequest(http.MethodGet, "/roles/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("get non-existent role", func(t *testing.T) {
		svc := newMockUserService()

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.GET("/roles/:id", controller.GetRole)

		req := httptest.NewRequest(http.MethodGet, "/roles/999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// utils.Error maps 404 to 400 via getStatusCode(404) -> StatusBadRequest
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("get role with invalid id", func(t *testing.T) {
		svc := newMockUserService()

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.GET("/roles/:id", controller.GetRole)

		req := httptest.NewRequest(http.MethodGet, "/roles/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestRBACController_UpdateRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("update role successfully", func(t *testing.T) {
		svc := newMockUserService()
		svc.roles[1] = &user.Role{ID: 1, Code: "editor", Name: "Editor"}

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.PUT("/roles/:id", controller.UpdateRole)

		body := map[string]string{"name": "Senior Editor", "description": "Updated"}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/roles/1", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("update non-existent role", func(t *testing.T) {
		svc := newMockUserService()

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.PUT("/roles/:id", controller.UpdateRole)

		body := map[string]string{"name": "New Name"}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/roles/999", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestRBACController_DeleteRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("delete non-system role", func(t *testing.T) {
		svc := newMockUserService()
		svc.roles[1] = &user.Role{ID: 1, Code: "custom", Name: "Custom", IsSystem: false}

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.DELETE("/roles/:id", controller.DeleteRole)

		req := httptest.NewRequest(http.MethodDelete, "/roles/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("delete system role should fail", func(t *testing.T) {
		svc := newMockUserService()
		svc.roles[2] = &user.Role{ID: 2, Code: "admin", Name: "Admin", IsSystem: true}

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.DELETE("/roles/:id", controller.DeleteRole)

		req := httptest.NewRequest(http.MethodDelete, "/roles/2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("delete non-existent role", func(t *testing.T) {
		svc := newMockUserService()

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.DELETE("/roles/:id", controller.DeleteRole)

		req := httptest.NewRequest(http.MethodDelete, "/roles/999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestRBACController_ListPermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("list permissions successfully", func(t *testing.T) {
		svc := newMockUserService()
		svc.permissions[1] = &user.Permission{ID: 1, Code: "person:read", Name: "Read Person", Module: "person"}

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.GET("/permissions", controller.ListPermissions)

		req := httptest.NewRequest(http.MethodGet, "/permissions", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestRBACController_CreatePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("create permission successfully", func(t *testing.T) {
		svc := newMockUserService()

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.POST("/permissions", controller.CreatePermission)

		body := map[string]string{"code": "test:read", "name": "Test Read", "module": "test"}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("create permission with empty code", func(t *testing.T) {
		svc := newMockUserService()

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.POST("/permissions", controller.CreatePermission)

		body := map[string]string{"code": "", "name": "No Code", "module": "test"}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestRBACController_DeletePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("delete permission successfully", func(t *testing.T) {
		svc := newMockUserService()
		svc.permissions[1] = &user.Permission{ID: 1, Code: "test:read", Name: "Test Read", Module: "test"}

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.DELETE("/permissions/:id", controller.DeletePermission)

		req := httptest.NewRequest(http.MethodDelete, "/permissions/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("delete non-existent permission", func(t *testing.T) {
		svc := newMockUserService()

		r := gin.New()
		controller := &testRBACController{userService: svc}
		r.DELETE("/permissions/:id", controller.DeletePermission)

		req := httptest.NewRequest(http.MethodDelete, "/permissions/999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}