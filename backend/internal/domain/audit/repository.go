package audit

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, log *AuditLog) error
	GetByID(ctx context.Context, id int64) (*AuditLog, error)
	List(ctx context.Context, filter AuditLogFilter) (*AuditLogList, error)
	DeleteOldLogs(ctx context.Context, before time.Time) (int64, error)
}