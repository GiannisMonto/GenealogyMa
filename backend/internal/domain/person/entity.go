package person

import (
	"fmt"
	"time"
)

// Gender 性别枚举
type Gender string

const (
	GenderMale   Gender = "男"
	GenderFemale Gender = "女"
	GenderUnknown Gender = ""
)

// RelationType 关系类型
type RelationType string

const (
	RelationBiological RelationType = "biological" // 亲生
	RelationAdoptive   RelationType = "adoptive"   // 过继
	RelationStep       RelationType = "step"       // 继子女
)

// SpouseType 配偶类型
type SpouseType string

const (
	SpouseTypePrimary SpouseType = "配" // 原配
	SpouseTypeSecond  SpouseType = "继" // 继配
	SpouseTypeDeceased SpouseType = "妣" // 已故
	SpouseTypeThird   SpouseType = "三" // 三配
)

// Person 人物聚合根 - 领域模型
type Person struct {
	ID               int64          `json:"id"`
	LegacyID         string         `json:"legacy_id"` // 传统族谱9位编号
	Name             string         `json:"name"`
	StyleName        string         `json:"style_name"` // 字/号
	Gender           Gender         `json:"gender"`
	Generation       int            `json:"generation"` // 世代
	BirthOrder       string         `json:"birth_order"` // 排行文本（如"长子"、"次女"）
	FatherID         *int64         `json:"father_id"`
	LineagePath      string         `json:"lineage_path"` // ltree路径，如"1.2.5"
	DetailText       string         `json:"detail_text"` // 详细生平
	BirthTimeText    string         `json:"birth_time_text"` // 生卒时间文本
	DeathTimeText    string         `json:"death_time_text"`
	BirthGregorian   *time.Time     `json:"birth_gregorian,omitempty"`
	DeathGregorian   *time.Time     `json:"death_gregorian,omitempty"`
	BirthYear        *int           `json:"birth_year"`
	DeathYear        *int           `json:"death_year"`
	BirthPlace       string         `json:"birth_place"`
	BurialPlace      string         `json:"burial_place"`
	SonCount         int            `json:"son_count"`
	DaughterCount    int            `json:"daughter_count"`
	AdoptedHeirCount int            `json:"adopted_heir_count"`
	TotalChildren    int            `json:"total_children_count"`
	PageNumber       *int           `json:"page_number"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`

	// 关联数据（按需加载）
	Father         *Person          `json:"father,omitempty"`
	Mother         *Spouse          `json:"mother,omitempty"`
	Children       []*PersonChild   `json:"children,omitempty"`
	Spouses        []*Spouse        `json:"spouses,omitempty"`
}

// PersonChild 子女信息（包含关系类型）
type PersonChild struct {
	Person         *Person     `json:"person"`
	RelationType   RelationType `json:"relation_type"`
	BirthOrderNum  int          `json:"birth_order_num"`
	IsPrimary      bool         `json:"is_primary"`
}

// Spouse 配偶信息
type Spouse struct {
	ID             int64       `json:"id"`
	MemberID       int64       `json:"member_id"`
	SpouseType     SpouseType  `json:"spouse_type"`
	Name           string      `json:"name"`
	BirthTimeText  string      `json:"birth_time_text"`
	DeathTimeText  string      `json:"death_time_text"`
	BirthPlace     string      `json:"birth_place"`
	BurialPlace    string      `json:"burial_place"`
	BirthGregorian *time.Time  `json:"birth_gregorian,omitempty"`
	DeathGregorian *time.Time  `json:"death_gregorian,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// ParentChildRelation 父子关系实体
type ParentChildRelation struct {
	RelationID     int64        `json:"relation_id"`
	ParentID       int64        `json:"parent_id"`
	ChildID        int64        `json:"child_id"`
	RelationType   RelationType `json:"relation_type"`
	BirthOrderNum  int          `json:"birth_order_num"`
	BirthOrderText string       `json:"birth_order_text"`
	IsPrimary      bool         `json:"is_primary"`
	Note           string       `json:"note"`
	CreatedAt      time.Time    `json:"created_at"`
}

// 领域方法

// IsAlive 判断是否在世
func (p *Person) IsAlive() bool {
	return p.DeathYear == nil && p.DeathGregorian == nil
}

// Age 计算年龄
func (p *Person) Age() int {
	if p.BirthYear == nil {
		return 0
	}

	endYear := time.Now().Year()
	if p.DeathYear != nil {
		endYear = *p.DeathYear
	}

	return endYear - *p.BirthYear
}

// FullName 获取完整姓名（含字号）
func (p *Person) FullName() string {
	if p.StyleName != "" {
		return fmt.Sprintf("%s（%s）", p.Name, p.StyleName)
	}
	return p.Name
}

// HasChildren 判断是否有子女
func (p *Person) HasChildren() bool {
	return p.TotalChildren > 0
}

// HasSpouse 判断是否有配偶
func (p *Person) HasSpouse() bool {
	return len(p.Spouses) > 0
}

// Validate 验证领域 invariants
func (p *Person) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if p.Gender != GenderMale && p.Gender != GenderFemale {
		return fmt.Errorf("invalid gender: %s", p.Gender)
	}

	if p.Generation < 0 || p.Generation > 30 {
		return fmt.Errorf("generation must be between 0 and 30")
	}

	// 验证子女数量一致性
	expectedTotal := p.SonCount + p.DaughterCount + p.AdoptedHeirCount
	if expectedTotal != p.TotalChildren {
		return fmt.Errorf("children count mismatch: son=%d + daughter=%d + adopted=%d = %d, but total=%d",
			p.SonCount, p.DaughterCount, p.AdoptedHeirCount, expectedTotal, p.TotalChildren)
	}

	return nil
}
