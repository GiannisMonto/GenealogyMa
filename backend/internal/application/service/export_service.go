package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/person"
)

// ExportService 导出服务
type ExportService struct {
	personRepo person.Repository
}

// NewExportService 创建导出服务
func NewExportService(personRepo person.Repository) *ExportService {
	return &ExportService{
		personRepo: personRepo,
	}
}

// ExportFormat 导出格式
type ExportFormat string

const (
	ExportFormatCSV   ExportFormat = "csv"
	ExportFormatJSON  ExportFormat = "json"
)

// ExportEntity 导出实体类型
type ExportEntity string

const (
	ExportEntityPerson    ExportEntity = "person"
	ExportEntityGenealogy ExportEntity = "genealogy"
)

// ExportFilter 导出筛选条件
type ExportFilter struct {
	Entity      ExportEntity `json:"entity"`
	Format      ExportFormat `json:"format"`
	Name        string       `json:"name"`
	Gender      string       `json:"gender"`
	Generation  *int         `json:"generation"`
	StartDate   string       `json:"start_date"`
	EndDate     string       `json:"end_date"`
	LineagePath string       `json:"lineage_path"`
	IsAlive     *bool        `json:"is_alive"`
}

// ExportResult 导出结果
type ExportResult struct {
	Data   []byte `json:"data"`
	Format string `json:"format"`
	Count  int    `json:"count"`
}

// ExportPersonsToCSV 导出人物数据到CSV
func (s *ExportService) ExportPersonsToCSV(ctx context.Context, filter *ExportFilter) (*ExportResult, error) {
	query := &person.SearchQuery{
		Name:       filter.Name,
		PageSize:   10000, // 限制单次导出数量
	}

	if filter.Generation != nil {
		query.Generation = filter.Generation
	}

	if filter.Gender != "" {
		g := person.Gender(filter.Gender)
		query.Gender = &g
	}

	persons, _, err := s.personRepo.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("查询人物数据失败: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 写入表头
	headers := []string{"ID", "姓名", "字辈", "性别", "代际", "出生地", "埋葬地", "生卒时间", "是否在世", "创建时间"}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("写入CSV表头失败: %w", err)
	}

	// 写入数据
	for _, p := range persons {
		isAlive := "否"
		if p.IsAlive() {
			isAlive = "是"
		}

		birthDeath := ""
		if p.BirthTimeText != "" || p.DeathTimeText != "" {
			birthDeath = p.BirthTimeText
			if p.DeathTimeText != "" {
				birthDeath += " - " + p.DeathTimeText
			}
		}

		row := []string{
			fmt.Sprintf("%d", p.ID),
			p.Name,
			p.StyleName,
			string(p.Gender),
			fmt.Sprintf("%d", p.Generation),
			p.BirthPlace,
			p.BurialPlace,
			birthDeath,
			isAlive,
			p.CreatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("写入CSV数据行失败: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("刷新CSV写入器失败: %w", err)
	}

	return &ExportResult{
		Data:   buf.Bytes(),
		Format: string(ExportFormatCSV),
		Count:  len(persons),
	}, nil
}

// ExportPersonsToJSON 导出人物数据到JSON
func (s *ExportService) ExportPersonsToJSON(ctx context.Context, filter *ExportFilter) (*ExportResult, error) {
	query := &person.SearchQuery{
		Name:       filter.Name,
		PageSize:   10000,
	}

	if filter.Generation != nil {
		query.Generation = filter.Generation
	}

	if filter.Gender != "" {
		g := person.Gender(filter.Gender)
		query.Gender = &g
	}

	persons, _, err := s.personRepo.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("查询人物数据失败: %w", err)
	}

	// 转换为DTO
	dtos := make([]PersonExportDTO, 0, len(persons))
	for _, p := range persons {
		dto := PersonExportDTO{
			ID:            p.ID,
			Name:          p.Name,
			StyleName:     p.StyleName,
			Gender:        string(p.Gender),
			Generation:    p.Generation,
			BirthPlace:    p.BirthPlace,
			BurialPlace:   p.BurialPlace,
			BirthTimeText: p.BirthTimeText,
			DeathTimeText: p.DeathTimeText,
			IsAlive:       p.IsAlive(),
			LineagePath:   p.LineagePath,
			CreatedAt:     p.CreatedAt.Format(time.RFC3339),
		}
		dtos = append(dtos, dto)
	}

	// 手动构建JSON避免外部依赖
	jsonBytes, err := buildJSON(dtos)
	if err != nil {
		return nil, fmt.Errorf("构建JSON失败: %w", err)
	}

	return &ExportResult{
		Data:   jsonBytes,
		Format: string(ExportFormatJSON),
		Count:  len(persons),
	}, nil
}

// PersonExportDTO 人物导出DTO
type PersonExportDTO struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	StyleName     string `json:"style_name"`
	Gender        string `json:"gender"`
	Generation    int    `json:"generation"`
	BirthPlace    string `json:"birth_place"`
	BurialPlace   string `json:"burial_place"`
	BirthTimeText string `json:"birth_time_text"`
	DeathTimeText string `json:"death_time_text"`
	IsAlive       bool   `json:"is_alive"`
	LineagePath   string `json:"lineage_path"`
	CreatedAt     string `json:"created_at"`
}

// buildJSON 简单JSON序列化
func buildJSON(data []PersonExportDTO) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString("[")
	for i, d := range data {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString("{")
		sb.WriteString(fmt.Sprintf(`"id":%d,`, d.ID))
		sb.WriteString(fmt.Sprintf(`"name":"%s",`, escapeJSON(d.Name)))
		sb.WriteString(fmt.Sprintf(`"style_name":"%s",`, escapeJSON(d.StyleName)))
		sb.WriteString(fmt.Sprintf(`"gender":"%s",`, escapeJSON(d.Gender)))
		sb.WriteString(fmt.Sprintf(`"generation":%d,`, d.Generation))
		sb.WriteString(fmt.Sprintf(`"birth_place":"%s",`, escapeJSON(d.BirthPlace)))
		sb.WriteString(fmt.Sprintf(`"burial_place":"%s",`, escapeJSON(d.BurialPlace)))
		sb.WriteString(fmt.Sprintf(`"birth_time_text":"%s",`, escapeJSON(d.BirthTimeText)))
		sb.WriteString(fmt.Sprintf(`"death_time_text":"%s",`, escapeJSON(d.DeathTimeText)))
		sb.WriteString(fmt.Sprintf(`"is_alive":%t,`, d.IsAlive))
		sb.WriteString(fmt.Sprintf(`"lineage_path":"%s",`, escapeJSON(d.LineagePath)))
		sb.WriteString(fmt.Sprintf(`"created_at":"%s"`, d.CreatedAt))
		sb.WriteString("}")
	}
	sb.WriteString("]")
	return []byte(sb.String()), nil
}

// escapeJSON 转义JSON字符串
func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// Export 执行导出
func (s *ExportService) Export(ctx context.Context, filter *ExportFilter) (*ExportResult, error) {
	if filter.Entity == "" {
		filter.Entity = ExportEntityPerson
	}
	if filter.Format == "" {
		filter.Format = ExportFormatCSV
	}

	switch filter.Entity {
	case ExportEntityPerson:
		if filter.Format == ExportFormatJSON {
			return s.ExportPersonsToJSON(ctx, filter)
		}
		return s.ExportPersonsToCSV(ctx, filter)
	default:
		return nil, fmt.Errorf("不支持的实体类型: %s", filter.Entity)
	}
}

// GetExportableEntities 获取可导出的实体类型
func (s *ExportService) GetExportableEntities() []ExportEntity {
	return []ExportEntity{ExportEntityPerson, ExportEntityGenealogy}
}

// GetExportFormats 获取支持的导出格式
func (s *ExportService) GetExportFormats() []ExportFormat {
	return []ExportFormat{ExportFormatCSV, ExportFormatJSON}
}
