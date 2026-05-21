package memorial

import (
	"fmt"
	"time"
)

// MemorialHall 宗祠/纪念馆聚合根
type MemorialHall struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`           // 宗祠名称
	Description  string     `json:"description"`    // 描述
	Province     string     `json:"province"`       // 省份
	City         string     `json:"city"`           // 城市
	District     string     `json:"district"`       // 区/县
	Address      string     `json:"address"`        // 详细地址
	Latitude     float64    `json:"latitude"`       // 纬度
	Longitude    float64    `json:"longitude"`      // 经度
	BuildYear    *int       `json:"build_year"`     // 建造年份
	Style        string     `json:"style"`          // 建筑风格
	ImageURL     string     `json:"image_url"`      // 图片
	TotalTablet  int        `json:"total_tablet"`   // 总牌位数
	UsedTablet   int        `json:"used_tablet"`    // 已用牌位数
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// 关联数据
	Tables       []*MemorialTablet `json:"tablets,omitempty"`
}

// MemorialTablet 牌位实体
type MemorialTablet struct {
	ID           int64       `json:"id"`
	HallID       int64       `json:"hall_id"`         // 所属宗祠
	PersonID     *int64      `json:"person_id"`       // 关联人物ID
	PersonName   string      `json:"person_name"`     // 人物姓名（冗余）
	Generation   int         `json:"generation"`      // 世代
	TabletType   TabletType  `json:"tablet_type"`     // 牌位类型
	Position     string      `json:"position"`       // 位置描述
	Floor        int         `json:"floor"`           // 层
	Row          int         `json:"row"`            // 排
	Number       int         `json:"number"`         // 号
	Enthronement *time.Time  `json:"entronement"`     // 安放日期
	Note         string      `json:"note"`           // 备注
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`

	// 关联数据
	Person       *PersonInfo `json:"person,omitempty"`
}

// PersonInfo 人物简要信息
type PersonInfo struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	StyleName  string  `json:"style_name"`  // 字/号
	Gender     Gender  `json:"gender"`
	BirthYear  *int    `json:"birth_year"`
	DeathYear  *int    `json:"death_year"`
}

// TabletType 牌位类型
type TabletType string

const (
	TabletTypeAncestor  TabletType = "ancestor"  // 祖先牌位
	TabletTypeMartyr    TabletType = "martyr"    // 烈士牌位
	TabletTypeSage      TabletType = "sage"      // 圣贤牌位
	TabletTypeFounder   TabletType = "founder"   // 创始人牌位
)

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

// Validate 验证宗祠信息
func (m *MemorialHall) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if m.Province == "" {
		return fmt.Errorf("province is required")
	}
	if m.TotalTablet < 0 {
		return fmt.Errorf("total_tablet cannot be negative")
	}
	if m.UsedTablet < 0 {
		return fmt.Errorf("used_tablet cannot be negative")
	}
	if m.UsedTablet > m.TotalTablet {
		return fmt.Errorf("used_tablet cannot exceed total_tablet")
	}
	return nil
}

// UsageRate 计算使用率
func (m *MemorialHall) UsageRate() float64 {
	if m.TotalTablet == 0 {
		return 0
	}
	return float64(m.UsedTablet) / float64(m.TotalTablet) * 100
}

// HasAvailableSpace 判断是否有可用空间
func (m *MemorialHall) HasAvailableSpace() bool {
	return m.UsedTablet < m.TotalTablet
}

// Validate 验证牌位信息
func (m *MemorialTablet) Validate() error {
	if m.HallID <= 0 {
		return fmt.Errorf("hall_id is required")
	}
	if m.PersonName == "" && m.PersonID == nil {
		return fmt.Errorf("either person_name or person_id is required")
	}
	if m.Floor < 0 {
		return fmt.Errorf("floor cannot be negative")
	}
	if m.Row < 0 {
		return fmt.Errorf("row cannot be negative")
	}
	if m.Number <= 0 {
		return fmt.Errorf("number must be positive")
	}
	return nil
}

// GetPositionDesc 获取位置描述
func (m *MemorialTablet) GetPositionDesc() string {
	return fmt.Sprintf("%d层 %d排 %d号", m.Floor, m.Row, m.Number)
}

// IsOccupied 判断是否已被占用
func (m *MemorialTablet) IsOccupied() bool {
	return m.PersonID != nil
}