package culture

import (
	"fmt"
	"time"
)

// DocumentCategory 文献类别
type DocumentCategory string

const (
	CategoryClassic   DocumentCategory = "classic"   // 经典文献
	CategoryGenealogy DocumentCategory = "genealogy" // 族谱序
	CategoryMemorial  DocumentCategory = "memorial"  // 纪念文章
	CategoryHistory   DocumentCategory = "history"   // 史料
)

// Document 文献聚合根
type Document struct {
	ID           int64            `json:"id"`
	Title        string           `json:"title"`
	Content      string           `json:"content"`
	Category     DocumentCategory `json:"category"`
	Author       string           `json:"author"`
	CreatedYear  *int             `json:"created_year"`
	Dynasty      string           `json:"dynasty"`
	Source       string           `json:"source"`
	ImageURLs    []string         `json:"image_urls"`
	ViewCount    int              `json:"view_count"`
	CollectCount int              `json:"collect_count"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// Validate 验证文献信息
func (d *Document) Validate() error {
	if d.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if d.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	if d.Category == "" {
		return fmt.Errorf("category is required")
	}
	return nil
}

// Story 故事实体
type Story struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Era         string    `json:"era"`         // 时代，如"清朝康熙年间"
	Category    string    `json:"category"`    // 分类
	Tags        []string  `json:"tags"`
	AudioURL    string    `json:"audio_url"`   // 音频链接
	ImageURL    string    `json:"image_url"`
	ViewCount   int       `json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 验证故事信息
func (s *Story) Validate() error {
	if s.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if s.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	return nil
}

// FamilyTeachings 家训实体
type FamilyTeachings struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Generation  int       `json:"generation"`  // 适用世代
	OriginText  string    `json:"origin_text"` // 原文
	Meaning     string    `json:"meaning"`     // 释义
	UsageCount  int       `json:"usage_count"` // 使用次数
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 验证家训信息
func (ft *FamilyTeachings) Validate() error {
	if ft.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if ft.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	return nil
}