package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   *int64      `json:"total,omitempty"` // 分页总数
}

// PageInfo 分页信息
type PageInfo struct {
	Page     int   `json:"page" form:"page"`
	PageSize int   `json:"page_size" form:"page_size"`
	Total    int64 `json:"total"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应带自定义消息
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// SuccessWithPage 分页成功响应
func SuccessWithPage(c *gin.Context, data interface{}, total int64) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
		Total:   &total,
	})
}

// PageSuccess 分页成功响应（带页码信息）
func PageSuccess(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
		Total:   &total,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(getStatusCode(code), Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 错误响应带数据
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(getStatusCode(code), Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// getStatusCode 根据业务错误码获取HTTP状态码
func getStatusCode(code int) int {
	switch {
	case code == 0:
		return http.StatusOK
	case code >= 400 && code < 500:
		return http.StatusBadRequest
	case code == 401:
		return http.StatusUnauthorized
	case code == 403:
		return http.StatusForbidden
	case code == 404:
		return http.StatusNotFound
	case code >= 500:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}

// GetPageInfo 从请求中获取分页信息
func GetPageInfo(c *gin.Context) *PageInfo {
	page := c.GetInt("page")
	if page <= 0 {
		page = 1
	}

	pageSize := c.GetInt("page_size")
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return &PageInfo{
		Page:     page,
		PageSize: pageSize,
	}
}

// Offset 计算偏移量
func (p *PageInfo) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit 获取限制数
func (p *PageInfo) Limit() int {
	return p.PageSize
}
