package service

import (
	"context"

	"github.com/genealogy-ma/platform/internal/domain/audit"
)

// AuditService 审计应用服务
type AuditService struct {
	domainService *audit.Service
}

// NewAuditService 创建审计应用服务
func NewAuditService(repo audit.Repository) *AuditService {
	return &AuditService{
		domainService: audit.NewService(repo),
	}
}

// ===== DTO 定义 =====

// AuditLogDTO 审计日志DTO
type AuditLogDTO struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	Module       string `json:"module"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	OldValue     string `json:"old_value,omitempty"`
	NewValue     string `json:"new_value,omitempty"`
	IPAddress    string `json:"ip_address,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	Description  string `json:"description,omitempty"`
	CreatedAt    string `json:"created_at"`
}

// AuditLogListDTO 审计日志列表DTO
type AuditLogListDTO struct {
	Data       []*AuditLogDTO `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// LogActionRequest 记录审计日志请求
type LogActionRequest struct {
	UserID       int64  `json:"user_id" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Module       string `json:"module" binding:"required"`
	Action       string `json:"action" binding:"required"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	OldValue     string `json:"old_value"`
	NewValue     string `json:"new_value"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	Description  string `json:"description"`
}

// AuditFilterRequest 审计日志筛选请求
type AuditFilterRequest struct {
	UserID       *int64 `form:"user_id"`
	Module       string `form:"module"`
	Action       string `form:"action"`
	ResourceType string `form:"resource_type"`
	StartDate    string `form:"start_date"`
	EndDate      string `form:"end_date"`
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
}

// CleanOldLogsRequest 清理旧日志请求
type CleanOldLogsRequest struct {
	Days int `json:"days" binding:"required,min=1"`
}

// ===== 服务方法 =====

func (s *AuditService) LogAction(ctx context.Context, req *LogActionRequest) error {
	cmd := &audit.CreateAuditLogCommand{
		UserID:       req.UserID,
		Username:     req.Username,
		Module:       req.Module,
		Action:       req.Action,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		OldValue:     req.OldValue,
		NewValue:     req.NewValue,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		Description:  req.Description,
	}
	return s.domainService.LogAction(ctx, cmd)
}

func (s *AuditService) GetLogByID(ctx context.Context, id int64) (*AuditLogDTO, error) {
	log, err := s.domainService.GetLogByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if log == nil {
		return nil, nil
	}
	return auditLogToDTO(log), nil
}

func (s *AuditService) ListLogs(ctx context.Context, req *AuditFilterRequest) (*AuditLogListDTO, error) {
	filter := audit.AuditLogFilter{
		UserID:       req.UserID,
		Module:       req.Module,
		Action:       req.Action,
		ResourceType: req.ResourceType,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}

	result, err := s.domainService.ListLogs(ctx, filter)
	if err != nil {
		return nil, err
	}

	dtos := make([]*AuditLogDTO, len(result.Data))
	for i, log := range result.Data {
		dtos[i] = auditLogToDTO(&log)
	}

	return &AuditLogListDTO{
		Data:       dtos,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *AuditService) CleanOldLogs(ctx context.Context, days int) (int64, error) {
	return s.domainService.CleanOldLogs(ctx, days)
}

// ===== 辅助函数 =====

func auditLogToDTO(log *audit.AuditLog) *AuditLogDTO {
	return &AuditLogDTO{
		ID:           log.ID,
		UserID:       log.UserID,
		Username:     log.Username,
		Module:       log.Module,
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceID:   log.ResourceID,
		OldValue:     log.OldValue,
		NewValue:     log.NewValue,
		IPAddress:    log.IPAddress,
		UserAgent:    log.UserAgent,
		Description:  log.Description,
		CreatedAt:    log.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
