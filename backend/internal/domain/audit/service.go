package audit

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) LogAction(ctx context.Context, cmd *CreateAuditLogCommand) error {
	log := &AuditLog{
		UserID:       cmd.UserID,
		Username:     cmd.Username,
		Module:       cmd.Module,
		Action:       cmd.Action,
		ResourceType: cmd.ResourceType,
		ResourceID:   cmd.ResourceID,
		OldValue:     cmd.OldValue,
		NewValue:     cmd.NewValue,
		IPAddress:    cmd.IPAddress,
		UserAgent:    cmd.UserAgent,
		Description:  cmd.Description,
	}
	return s.repo.Create(ctx, log)
}

func (s *Service) GetLogByID(ctx context.Context, id int64) (*AuditLog, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListLogs(ctx context.Context, filter AuditLogFilter) (*AuditLogList, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}
	return s.repo.List(ctx, filter)
}

func (s *Service) CleanOldLogs(ctx context.Context, days int) (int64, error) {
	if days < 1 {
		return 0, fmt.Errorf("days must be positive")
	}
	before := time.Now().Add(-time.Duration(days) * 24 * time.Hour)
	return s.repo.DeleteOldLogs(ctx, before)
}