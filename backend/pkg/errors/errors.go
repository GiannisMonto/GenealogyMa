package errors

import (
	"fmt"
)

// 错误码定义
const (
	// 系统错误 (1000-1999)
	CodeInternalError     = 1000
	CodeInvalidParam      = 1001
	CodeUnauthorized      = 1002
	CodeForbidden         = 1003
	CodeNotFound          = 1004
	CodeTooManyRequests   = 1005

	// 用户相关错误 (2000-2999)
	CodeUserNotFound      = 2000
	CodeUserAlreadyExists = 2001
	CodeWrongPassword     = 2002
	CodeUserDisabled      = 2003

	// 族谱相关错误 (3000-3999)
	CodePersonNotFound    = 3000
	CodePersonExists      = 3001
	CodeRelationExists    = 3002
	CodeInvalidRelation   = 3003
	CodeCycleDetected     = 3004

	// 文件相关错误 (4000-4999)
	CodeFileTooLarge      = 4000
	CodeFileTypeInvalid   = 4001
	CodeFileUploadFailed  = 4002
)

// AppError 应用错误
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持错误链
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New 创建错误
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装错误
func Wrap(code int, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// 预定义错误
var (
	ErrInternal     = New(CodeInternalError, "服务器内部错误")
	ErrInvalidParam = New(CodeInvalidParam, "参数错误")
	ErrUnauthorized = New(CodeUnauthorized, "未授权访问")
	ErrForbidden    = New(CodeForbidden, "权限不足")
	ErrNotFound     = New(CodeNotFound, "资源不存在")

	ErrUserNotFound  = New(CodeUserNotFound, "用户不存在")
	ErrWrongPassword = New(CodeWrongPassword, "密码错误")

	ErrPersonNotFound = New(CodePersonNotFound, "人物不存在")
)

// IsAppError 检查是否为应用错误
func IsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}

	var appErr *AppError
	if As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// As 错误类型转换
func As(err error, target interface{}) bool {
	if e, ok := err.(*AppError); ok {
		if t, ok := target.(**AppError); ok {
			*t = e
			return true
		}
	}
	return false
}

// GetCode 获取错误码
func GetCode(err error) int {
	if appErr, ok := IsAppError(err); ok {
		return appErr.Code
	}
	return CodeInternalError
}

// GetMessage 获取错误消息
func GetMessage(err error) string {
	if appErr, ok := IsAppError(err); ok {
		return appErr.Message
	}
	return "服务器内部错误"
}
