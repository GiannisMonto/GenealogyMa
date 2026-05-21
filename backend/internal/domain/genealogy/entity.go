package genealogy

import (
	"fmt"
	"time"
)

// Genealogy 族谱聚合根
type Genealogy struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`          // 族谱名称，如"王氏族谱"
	Surname     string     `json:"surname"`       // 姓氏
	Description string     `json:"description"`   // 族谱描述
	OriginPlace string     `json:"origin_place"`   // 起源地
	TotalGenerations int   `json:"total_generations"` // 总世代数
	TotalMembers int       `json:"total_members"`  // 总人数
	Version     string     `json:"version"`       // 版本号
	ImageURL    string     `json:"image_url"`     // 封面图
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联数据
	Branches    []*Branch  `json:"branches,omitempty"`
	Generations []*Generation `json:"generations,omitempty"`
}

// Branch 分支实体
type Branch struct {
	ID           int64     `json:"id"`
	GenealogyID  int64     `json:"genealogy_id"`
	Name         string    `json:"name"`          // 分支名称，如"长房"
	Code         string    `json:"code"`          // 分支代号，如"A"
	Description  string    `json:"description"`   // 分支描述
	GenerationStart int    `json:"generation_start"` // 起始世代
	GenerationEnd   int    `json:"generation_end"`   // 结束世代
	MemberCount  int       `json:"member_count"`  // 分支人数
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联数据
	Members     []*PersonSummary `json:"members,omitempty"`
}

// PersonSummary 人物摘要（用于分支成员列表）
type PersonSummary struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Gender       Gender  `json:"gender"`
	Generation   int     `json:"generation"`
	BirthYear    *int    `json:"birth_year,omitempty"`
	DeathYear    *int    `json:"death_year,omitempty"`
	StyleName    string  `json:"style_name"` // 字/号
}

// Generation 世代（字辈）实体
type Generation struct {
	ID           int64     `json:"id"`
	GenealogyID  int64     `json:"genealogy_id"`
	Generation   int       `json:"generation"`    // 世代序号
	Name         string    `json:"name"`         // 世代名称，如"道"
	Sequence     int       `json:"sequence"`      // 在字辈中的顺序
	Description  string    `json:"description"`   // 说明，如"道德文章"
	StartYear    *int      `json:"start_year"`    // 开始使用年份
	EndYear      *int      `json:"end_year"`      // 结束使用年份
	MemberCount  int       `json:"member_count"`  // 该世代人数
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Gender 性别枚举
type Gender string

const (
	GenderMale   Gender = "男"
	GenderFemale Gender = "女"
)

// IsValid 验证性别是否有效
func (g Gender) IsValid() bool {
	return g == GenderMale || g == GenderFemale
}

// Validate 验证族谱信息
func (g *Genealogy) Validate() error {
	if g.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if g.Surname == "" {
		return fmt.Errorf("surname cannot be empty")
	}
	if g.TotalGenerations < 0 {
		return fmt.Errorf("total_generations cannot be negative")
	}
	if g.TotalMembers < 0 {
		return fmt.Errorf("total_members cannot be negative")
	}
	return nil
}

// HasBranches 判断是否有分支
func (g *Genealogy) HasBranches() bool {
	return len(g.Branches) > 0
}

// Validate 验证分支信息
func (b *Branch) Validate() error {
	if b.Name == "" {
		return fmt.Errorf("branch name cannot be empty")
	}
	if b.GenealogyID <= 0 {
		return fmt.Errorf("genealogy_id is required")
	}
	if b.GenerationStart < 0 || b.GenerationEnd < 0 {
		return fmt.Errorf("generation cannot be negative")
	}
	if b.GenerationStart > b.GenerationEnd && b.GenerationEnd > 0 {
		return fmt.Errorf("generation_start cannot be greater than generation_end")
	}
	return nil
}

// Validate 验证世代信息
func (g *Generation) Validate() error {
	if g.GenealogyID <= 0 {
		return fmt.Errorf("genealogy_id is required")
	}
	if g.Generation < 0 {
		return fmt.Errorf("generation must be non-negative")
	}
	if g.Name == "" {
		return fmt.Errorf("generation name cannot be empty")
	}
	if g.Sequence < 0 {
		return fmt.Errorf("sequence must be non-negative")
	}
	return nil
}

// IsActive 判断世代是否在使用中
func (g *Generation) IsActive() bool {
	if g.EndYear != nil {
		return false
	}
	return true
}

// GetGenerationRange 获取世代年份范围描述
func (g *Generation) GetGenerationRange() string {
	if g.StartYear == nil && g.EndYear == nil {
		return "至今"
	}
	if g.StartYear != nil && g.EndYear == nil {
		return fmt.Sprintf("%d年至至今", *g.StartYear)
	}
	if g.StartYear != nil && g.EndYear != nil {
		return fmt.Sprintf("%d年至%d年", *g.StartYear, *g.EndYear)
	}
	return "未知"
}