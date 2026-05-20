package audit

import (
	"time"
)

type AuditLog struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Username     string    `json:"username"`
	Module       string    `json:"module"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	OldValue     string    `json:"old_value,omitempty"`
	NewValue     string    `json:"new_value,omitempty"`
	IPAddress    string    `json:"ip_address,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	Description  string    `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateAuditLogCommand struct {
	UserID       int64
	Username     string
	Module       string
	Action       string
	ResourceType string
	ResourceID   string
	OldValue     string
	NewValue     string
	IPAddress    string
	UserAgent    string
	Description  string
}

type AuditLogFilter struct {
	UserID       *int64
	Module       string
	Action       string
	ResourceType string
	StartDate    *time.Time
	EndDate      *time.Time
	Page         int
	PageSize     int
}

type AuditLogList struct {
	Data       []AuditLog `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}