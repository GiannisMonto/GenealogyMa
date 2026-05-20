package auth

import (
	"errors"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 角色常量
const (
	RoleSuperAdmin     = "super_admin"
	RoleGenealogyAdmin = "genealogy_admin"
	RoleCultureAdmin   = "culture_admin"
	RoleMemorialAdmin  = "memorial_admin"
	RoleVerifiedUser   = "verified_user"
	RoleGuest          = "guest"
)

// 权限常量
const (
	PermissionPersonRead    = "person:read"
	PermissionPersonWrite   = "person:write"
	PermissionPersonDelete  = "person:delete"
	PermissionPersonImport  = "person:import"
	PermissionCultureRead   = "culture:read"
	PermissionCultureWrite  = "culture:write"
	PermissionGenealogyRead = "genealogy:read"
	PermissionGenealogyWrite = "genealogy:write"
	PermissionRoleRead      = "role:read"
	PermissionRoleWrite    = "role:write"
	PermissionRoleDelete   = "role:delete"
	PermissionPermRead     = "permission:read"
	PermissionPermWrite    = "permission:write"
	PermissionPermDelete   = "permission:delete"
	PermissionAdminUser    = "admin:user"
	PermissionAdminConfig  = "admin:config"
	PermissionCommunityRead  = "community:read"
	PermissionCommunityWrite = "community:write"
)

// Claims JWT载荷
type Claims struct {
	UserID   int64    `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// JWTService JWT服务
type JWTService struct {
	secret          []byte
	expireHours     int
	refreshExpireHours int
}

// NewJWTService 创建JWT服务
func NewJWTService(secret string, expireHours, refreshExpireHours int) *JWTService {
	return &JWTService{
		secret:              []byte(secret),
		expireHours:         expireHours,
		refreshExpireHours:  refreshExpireHours,
	}
}

// GenerateToken 生成访问令牌
func (j *JWTService) GenerateToken(userID int64, username string, roles []string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.expireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "genealogy-ma",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// GenerateRefreshToken 生成刷新令牌
func (j *JWTService) GenerateRefreshToken(userID int64) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.refreshExpireHours))),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "genealogy-ma-refresh",
		Subject:   string(rune(userID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ParseToken 解析令牌
func (j *JWTService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// 角色权限映射
var rolePermissions = map[string][]string{
	RoleSuperAdmin: {
		PermissionPersonRead, PermissionPersonWrite, PermissionPersonDelete, PermissionPersonImport,
		PermissionCultureRead, PermissionCultureWrite,
		PermissionGenealogyRead, PermissionGenealogyWrite,
		PermissionCommunityRead, PermissionCommunityWrite,
		PermissionRoleRead, PermissionRoleWrite, PermissionRoleDelete,
		PermissionPermRead, PermissionPermWrite, PermissionPermDelete,
		PermissionAdminUser, PermissionAdminConfig,
	},
	RoleGenealogyAdmin: {
		PermissionPersonRead, PermissionPersonWrite, PermissionPersonDelete, PermissionPersonImport,
		PermissionGenealogyRead, PermissionGenealogyWrite,
	},
	RoleCultureAdmin: {
		PermissionCultureRead, PermissionCultureWrite,
	},
	RoleMemorialAdmin: {
		PermissionPersonRead,
	},
	RoleVerifiedUser: {
		PermissionPersonRead, PermissionCultureRead,
		PermissionCommunityRead, PermissionCommunityWrite,
	},
	RoleGuest: {
		PermissionPersonRead,
		PermissionCommunityRead,
	},
}

// HasRole 检查是否具有指定角色
func HasRole(userRoles []string, requiredRole string) bool {
	for _, role := range userRoles {
		if role == requiredRole {
			return true
		}
	}
	return false
}

// HasPermission 检查是否具有指定权限
func HasPermission(userRoles []string, requiredPermission string) bool {
	for _, role := range userRoles {
		if permissions, ok := rolePermissions[role]; ok {
			for _, perm := range permissions {
				if perm == requiredPermission {
					return true
				}
			}
		}
	}
	return false
}

// IsSuperAdmin 检查是否超级管理员
func IsSuperAdmin(userRoles []string) bool {
	return HasRole(userRoles, RoleSuperAdmin)
}

// ===== 全局实例和辅助函数 =====

var globalJWTService *JWTService

// SetGlobalJWTService 设置全局JWT服务实例
func SetGlobalJWTService(service *JWTService) {
	globalJWTService = service
}

// GenerateToken 生成访问令牌（全局函数）
func GenerateToken(userID int64, username string, roles []string) (string, error) {
	if globalJWTService == nil {
		return "", nil
	}
	return globalJWTService.GenerateToken(userID, username, roles)
}

// GenerateRefreshToken 生成刷新令牌（全局函数）
func GenerateRefreshToken(userID int64) (string, error) {
	if globalJWTService == nil {
		return "", nil
	}
	return globalJWTService.GenerateRefreshToken(userID)
}

// ParseToken 解析令牌（全局函数）
func ParseToken(tokenString string) (*Claims, error) {
	if globalJWTService == nil {
		return nil, nil
	}
	return globalJWTService.ParseToken(tokenString)
}

// RefreshClaims 刷新令牌声明
type RefreshClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// ParseRefreshToken 解析刷新令牌
func ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
	if globalJWTService == nil {
		return nil, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return globalJWTService.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, nil
}

// ===== Gin 上下文辅助函数 =====

// userContextKey 用户信息在 Gin 上下文中的键
type userContextKey string

const (
	userIDKey     userContextKey = "user_id"
	usernameKey   userContextKey = "username"
	userRolesKey  userContextKey = "roles"
)

// SetUserToContext 将用户信息设置到 Gin 上下文
func SetUserToContext(ctx *gin.Context, userID int64, username string, roles []string) {
	ctx.Set(string(userIDKey), userID)
	ctx.Set(string(usernameKey), username)
	ctx.Set(string(userRolesKey), roles)
}

// GetUserID 从上下文中获取用户ID
func GetUserID(ctx *gin.Context) int64 {
	value, exists := ctx.Get(string(userIDKey))
	if !exists {
		return 0
	}
	if userID, ok := value.(int64); ok {
		return userID
	}
	return 0
}

// GetUsername 从上下文中获取用户名
func GetUsername(ctx *gin.Context) string {
	value, exists := ctx.Get(string(usernameKey))
	if !exists {
		return ""
	}
	if username, ok := value.(string); ok {
		return username
	}
	return ""
}

// GetUserRoles 从上下文中获取用户角色
func GetUserRoles(ctx *gin.Context) []string {
	value, exists := ctx.Get(string(userRolesKey))
	if !exists {
		return nil
	}
	if roles, ok := value.([]string); ok {
		return roles
	}
	return nil
}

// ExtractToken 从请求中提取令牌
func ExtractToken(ctx *gin.Context) string {
	// 首先从 Authorization header 提取
	authHeader := ctx.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	// 其次从 query 参数提取
	if token := ctx.Query("token"); token != "" {
		return token
	}

	// 最后从 cookie 提取
	if token, err := ctx.Cookie("access_token"); err == nil {
		return token
	}

	return ""
}


// RequirePermission 返回要求指定权限的中间件
func RequirePermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无法获取用户权限",
			})
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok || !HasPermission(userRoles, requiredPermission) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
