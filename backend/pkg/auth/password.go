package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordService 密码加密服务
type PasswordService struct {
	cost int
}

// NewPasswordService 创建密码服务
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: bcrypt.DefaultCost,
	}
}

// HashPassword 加密密码
func (p *PasswordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	return string(bytes), err
}

// CheckPassword 验证密码
func (p *PasswordService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength 验证密码强度
func (p *PasswordService) ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return ErrPasswordTooShort
	}
	if len(password) > 72 {
		return ErrPasswordTooLong
	}
	return nil
}

// 密码错误类型
var (
	ErrPasswordTooShort = &PasswordError{Message: "password must be at least 6 characters"}
	ErrPasswordTooLong  = &PasswordError{Message: "password must be less than 72 characters"}
	ErrPasswordMismatch = &PasswordError{Message: "password mismatch"}
)

// PasswordError 密码错误类型
type PasswordError struct {
	Message string
}

func (e *PasswordError) Error() string {
	return e.Message
}
