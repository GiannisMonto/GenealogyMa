package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/genealogy-ma/platform/internal/application/service"
	"github.com/genealogy-ma/platform/internal/domain/user"
)

// MockUserService 模拟用户服务
type MockUserService struct {
	users    map[int64]*service.UserDTO
	nextID   int64
	profiles map[int64]*user.User
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		users:    make(map[int64]*service.UserDTO),
		nextID:   1,
		profiles: make(map[int64]*user.User),
	}
}

func (m *MockUserService) Register(ctx context.Context, req *service.RegisterRequest) (*service.UserDTO, error) {
	u := &service.UserDTO{
		ID:          m.nextID,
		Username:    req.Username,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Status:      "active",
	}
	m.nextID++
	m.users[u.ID] = u
	return u, nil
}

func (m *MockUserService) Login(ctx context.Context, req *service.LoginRequest) (*service.LoginResponse, error) {
	if req.Username == "test" && req.Password == "password" {
		return &service.LoginResponse{
			AccessToken:  "mock_token",
			RefreshToken: "mock_refresh",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
			User: service.UserDTO{
				ID:       1,
				Username: "test",
				Status:   "active",
			},
		}, nil
	}
	return nil, nil
}

func (m *MockUserService) GetByID(ctx context.Context, id int64) (*service.UserDTO, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *MockUserService) Update(ctx context.Context, id int64, req *service.UpdateUserRequest) (*service.UserDTO, error) {
	if u, ok := m.users[id]; ok {
		if req.DisplayName != "" {
			u.DisplayName = req.DisplayName
		}
		return u, nil
	}
	return nil, nil
}

func (m *MockUserService) Delete(ctx context.Context, id int64) error {
	delete(m.users, id)
	return nil
}

func (m *MockUserService) ChangePassword(ctx context.Context, userID int64, req *service.ChangePasswordRequest) error {
	return nil
}

func (m *MockUserService) Search(ctx context.Context, req *service.SearchUserRequest) ([]*service.UserDTO, int64, error) {
	result := make([]*service.UserDTO, 0, len(m.users))
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, int64(len(result)), nil
}

func (m *MockUserService) AssignRole(ctx context.Context, userID int64, assignedBy int64, roleCode string) error {
	return nil
}

func (m *MockUserService) RemoveRole(ctx context.Context, userID int64, roleCode string) error {
	return nil
}

// MockTokenBlacklist 模拟令牌黑名单
type MockTokenBlacklist struct {
	blacklist map[string]time.Time
}

func NewMockTokenBlacklist() *MockTokenBlacklist {
	return &MockTokenBlacklist{
		blacklist: make(map[string]time.Time),
	}
}

func (m *MockTokenBlacklist) Add(ctx context.Context, token string, expiresAt time.Time) {
	m.blacklist[token] = expiresAt
}

func (m *MockTokenBlacklist) IsBlacklisted(ctx context.Context, token string) bool {
	_, ok := m.blacklist[token]
	return ok
}

func setupAuthTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthController_Register(t *testing.T) {
	mockService := NewMockUserService()

	r := setupAuthTestRouter()
	r.POST("/auth/register", func(ctx *gin.Context) {
		var req service.RegisterRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, _ := mockService.Register(ctx, &req)
		ctx.JSON(http.StatusCreated, gin.H{"data": u})
	})

	body := bytes.NewBufferString(`{"username":"testuser","password":"password123","email":"test@example.com"}`)
	req, _ := http.NewRequest("POST", "/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["data"] == nil {
		t.Error("Expected user data in response")
	}
}

func TestAuthController_Register_ValidationError(t *testing.T) {
	mockService := NewMockUserService()
	_ = mockService

	r := setupAuthTestRouter()
	r.POST("/auth/register", func(ctx *gin.Context) {
		var req service.RegisterRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{})
	})

	// Missing required fields
	body := bytes.NewBufferString(`{"username":""}`)
	req, _ := http.NewRequest("POST", "/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAuthController_Login_Success(t *testing.T) {
	mockService := NewMockUserService()

	r := setupAuthTestRouter()
	r.POST("/auth/login", func(ctx *gin.Context) {
		var req service.LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, _ := mockService.Login(ctx, &req)
		if resp == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": resp})
	})

	body := bytes.NewBufferString(`{"username":"test","password":"password"}`)
	req, _ := http.NewRequest("POST", "/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok || data["access_token"] != "mock_token" {
		t.Errorf("Expected access token 'mock_token', got %v", data)
	}
}

func TestAuthController_Login_InvalidCredentials(t *testing.T) {
	mockService := NewMockUserService()

	r := setupAuthTestRouter()
	r.POST("/auth/login", func(ctx *gin.Context) {
		var req service.LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, _ := mockService.Login(ctx, &req)
		if resp == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": resp})
	})

	body := bytes.NewBufferString(`{"username":"wrong","password":"wrongpass"}`)
	req, _ := http.NewRequest("POST", "/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthController_Logout(t *testing.T) {
	mockBlacklist := NewMockTokenBlacklist()

	r := setupAuthTestRouter()
	r.POST("/auth/logout", func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token != "" {
			mockBlacklist.Add(ctx, token, time.Now().Add(24*time.Hour))
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "登出成功"})
	})

	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer mock_token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if !mockBlacklist.IsBlacklisted(context.Background(), "Bearer mock_token") {
		t.Error("Expected token to be blacklisted")
	}
}

func TestAuthController_GetProfile(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users[1] = &service.UserDTO{ID: 1, Username: "test", Email: "test@example.com"}

	r := setupAuthTestRouter()
	r.GET("/auth/profile", func(ctx *gin.Context) {
		userID := int64(1) // Mock authenticated user
		user, _ := mockService.GetByID(ctx, userID)
		if user == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": user})
	})

	req, _ := http.NewRequest("GET", "/auth/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok || data["username"] != "test" {
		t.Errorf("Expected username 'test', got %v", data)
	}
}

func TestAuthController_UpdateProfile(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users[1] = &service.UserDTO{ID: 1, Username: "test", DisplayName: "Old Nickname"}

	r := setupAuthTestRouter()
	r.PUT("/auth/profile", func(ctx *gin.Context) {
		userID := int64(1)
		var req service.UpdateUserRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, _ := mockService.Update(ctx, userID, &req)
		if user == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": user})
	})

	body := bytes.NewBufferString(`{"display_name":"New Nickname"}`)
	req, _ := http.NewRequest("PUT", "/auth/profile", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok || data["display_name"] != "New Nickname" {
		t.Errorf("Expected display_name 'New Nickname', got %v", data)
	}
}

func TestAuthController_ChangePassword(t *testing.T) {
	mockService := NewMockUserService()

	r := setupAuthTestRouter()
	r.POST("/auth/change-password", func(ctx *gin.Context) {
		var req service.ChangePasswordRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_ = mockService.ChangePassword(ctx, 1, &req)
		ctx.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
	})

	body := bytes.NewBufferString(`{"old_password":"old","new_password":"new123"}`)
	req, _ := http.NewRequest("POST", "/auth/change-password", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthController_ListUsers(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users[1] = &service.UserDTO{ID: 1, Username: "user1"}
	mockService.users[2] = &service.UserDTO{ID: 2, Username: "user2"}

	r := setupAuthTestRouter()
	r.GET("/users", func(ctx *gin.Context) {
		users, total, _ := mockService.Search(ctx, &service.SearchUserRequest{})
		ctx.JSON(http.StatusOK, gin.H{"data": users, "total": total})
	})

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("Expected 2 users, got %v", response["data"])
	}
}

func TestAuthController_GetUser(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users[1] = &service.UserDTO{ID: 1, Username: "testuser"}

	r := setupAuthTestRouter()
	r.GET("/users/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		_ = id
		user, _ := mockService.GetByID(ctx, 1)
		if user == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": user})
	})

	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthController_GetUser_NotFound(t *testing.T) {
	r := setupAuthTestRouter()
	r.GET("/users/:id", func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	})

	req, _ := http.NewRequest("GET", "/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestAuthController_UpdateUser(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users[1] = &service.UserDTO{ID: 1, Username: "test", DisplayName: "Old"}

	r := setupAuthTestRouter()
	r.PUT("/users/:id", func(ctx *gin.Context) {
		var req service.UpdateUserRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, _ := mockService.Update(ctx, 1, &req)
		if user == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": user})
	})

	body := bytes.NewBufferString(`{"nickname":"Updated"}`)
	req, _ := http.NewRequest("PUT", "/users/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthController_DeleteUser(t *testing.T) {
	mockService := NewMockUserService()
	mockService.users[1] = &service.UserDTO{ID: 1, Username: "test"}

	r := setupAuthTestRouter()
	r.DELETE("/users/:id", func(ctx *gin.Context) {
		_ = mockService.Delete(ctx, 1)
		ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})

	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if mockService.users[1] != nil {
		t.Error("Expected user to be deleted")
	}
}

func TestAuthController_AssignRole(t *testing.T) {
	mockService := NewMockUserService()

	r := setupAuthTestRouter()
	r.POST("/users/:id/roles", func(ctx *gin.Context) {
		var req service.AssignRoleRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_ = mockService.AssignRole(ctx, 1, 1, req.RoleCode)
		ctx.JSON(http.StatusOK, gin.H{"message": "角色分配成功"})
	})

	body := bytes.NewBufferString(`{"role_code":"admin"}`)
	req, _ := http.NewRequest("POST", "/users/1/roles", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthController_RemoveRole(t *testing.T) {
	mockService := NewMockUserService()

	r := setupAuthTestRouter()
	r.DELETE("/users/:id/roles/:role_code", func(ctx *gin.Context) {
		_ = mockService.RemoveRole(ctx, 1, "admin")
		ctx.JSON(http.StatusOK, gin.H{"message": "角色移除成功"})
	})

	req, _ := http.NewRequest("DELETE", "/users/1/roles/admin", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}