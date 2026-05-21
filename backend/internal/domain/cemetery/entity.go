package cemetery

import (
	"fmt"
	"time"
)

// Cemetery 墓园聚合根
type Cemetery struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Province    string    `json:"province"`
	City        string    `json:"city"`
	District    string    `json:"district"`
	Address     string    `json:"address"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	TotalGrave  int       `json:"total_grave"`
	UsedGrave   int       `json:"used_grave"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Grave 墓位实体
type Grave struct {
	ID         int64      `json:"id"`
	CemeteryID  int64      `json:"cemetery_id"`
	PersonID   *int64     `json:"person_id"`
	Section    string     `json:"section"`    // 区
	Row         int        `json:"row"`       // 排
	Number     int        `json:"number"`    // 号
	Status     GraveStatus `json:"status"`   // available, occupied, reserved
	BuriedName  string     `json:"buried_name"`
	BuriedDate  *time.Time `json:"buried_date,omitempty"`
	BuriedYear  *int       `json:"buried_year"`
	Note        string     `json:"note"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// GraveStatus 墓位状态
type GraveStatus string

const (
	GraveStatusAvailable GraveStatus = "available" // 可用
	GraveStatusOccupied  GraveStatus = "occupied"  // 已占用
	GraveStatusReserved  GraveStatus = "reserved"  // 已预约
)

// Availability 计算墓位可用状态
func (g *Grave) Availability() bool {
	return g.Status == GraveStatusAvailable
}

// IsOccupied 判断是否已安葬
func (g *Grave) IsOccupied() bool {
	return g.Status == GraveStatusOccupied && g.PersonID != nil
}

// Validate 验证墓位信息
func (g *Grave) Validate() error {
	if g.CemeteryID <= 0 {
		return fmt.Errorf("cemetery_id is required")
	}
	if g.Section == "" {
		return fmt.Errorf("section is required")
	}
	if g.Row <= 0 {
		return fmt.Errorf("row must be positive")
	}
	if g.Number <= 0 {
		return fmt.Errorf("number must be positive")
	}
	return nil
}

// CemeteryUsageRate 计算墓园使用率
func (c *Cemetery) UsageRate() float64 {
	if c.TotalGrave == 0 {
		return 0
	}
	return float64(c.UsedGrave) / float64(c.TotalGrave) * 100
}

// HasAvailableSpace 判断是否有可用墓位
func (c *Cemetery) HasAvailableSpace() bool {
	return c.TotalGrave == 0 || c.UsedGrave < c.TotalGrave
}

// Validate 验证墓园信息
func (c *Cemetery) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if c.Province == "" {
		return fmt.Errorf("province is required")
	}
	if c.TotalGrave < 0 {
		return fmt.Errorf("total_grave cannot be negative")
	}
	if c.UsedGrave < 0 {
		return fmt.Errorf("used_grave cannot be negative")
	}
	if c.UsedGrave > c.TotalGrave {
		return fmt.Errorf("used_grave cannot exceed total_grave")
	}
	return nil
}