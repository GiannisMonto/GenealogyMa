package auth

import (
	"errors"
	"time"

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
	PermissionPersonRead   = "person:read"
	PermissionPersonWrite  = "person:write"
	PermissionPersonDelete = "person:delete"
	PermissionPersonImport = "person:import"
	PermissionCultureRead  = "culture:read"
	PermissionCultureWrite = "culture:write"
	PermissionAdminUser    = "admin:user"
	PermissionAdminConfig  = "admin:config"
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
		PermissionAdminUser, PermissionAdminConfig,
	},
	RoleGenealogyAdmin: {
		PermissionPersonRead, PermissionPersonWrite, PermissionPersonDelete, PermissionPersonImport,
	},
	RoleCultureAdmin: {
		PermissionCultureRead, PermissionCultureWrite,
	},
	RoleMemorialAdmin: {
		PermissionPersonRead,
	},
	RoleVerifiedUser: {
		PermissionPersonRead, PermissionCultureRead,
	},
	RoleGuest: {
		PermissionPersonRead,
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
