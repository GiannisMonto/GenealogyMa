package cemetery

import (
	"testing"
	"time"
)

func TestCemeteryValidate(t *testing.T) {
	tests := []struct {
		name    string
		c       *Cemetery
		wantErr bool
	}{
		{
			name: "valid cemetery",
			c: &Cemetery{
				Name:      "八里桥公墓",
				Province:  "北京市",
				TotalGrave: 1000,
				UsedGrave: 500,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			c: &Cemetery{
				Name:     "",
				Province: "北京市",
			},
			wantErr: true,
		},
		{
			name: "empty province",
			c: &Cemetery{
				Name:     "八里桥公墓",
				Province: "",
			},
			wantErr: true,
		},
		{
			name: "negative total grave",
			c: &Cemetery{
				Name:      "八里桥公墓",
				Province:  "北京市",
				TotalGrave: -1,
			},
			wantErr: true,
		},
		{
			name: "negative used grave",
			c: &Cemetery{
				Name:      "八里桥公墓",
				Province:  "北京市",
				TotalGrave: 1000,
				UsedGrave: -1,
			},
			wantErr: true,
		},
		{
			name: "used grave exceeds total",
			c: &Cemetery{
				Name:      "八里桥公墓",
				Province:  "北京市",
				TotalGrave: 1000,
				UsedGrave: 1001,
			},
			wantErr: true,
		},
		{
			name: "zero total is valid",
			c: &Cemetery{
				Name:      "八里桥公墓",
				Province:  "北京市",
				TotalGrave: 0,
				UsedGrave: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cemetery.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCemeteryUsageRate(t *testing.T) {
	tests := []struct {
		name      string
		totalGrave int
		usedGrave  int
		want       float64
	}{
		{"50% usage", 100, 50, 50.0},
		{"0% usage", 100, 0, 0.0},
		{"100% usage", 100, 100, 100.0},
		{"75% usage", 1000, 750, 75.0},
		{"zero total", 0, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cemetery{TotalGrave: tt.totalGrave, UsedGrave: tt.usedGrave}
			got := c.UsageRate()
			if got != tt.want {
				t.Errorf("Cemetery.UsageRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCemeteryHasAvailableSpace(t *testing.T) {
	tests := []struct {
		name       string
		totalGrave int
		usedGrave  int
		want       bool
	}{
		{"has space", 100, 50, true},
		{"no space", 100, 100, false},
		{"empty cemetery", 100, 0, true},
		{"zero capacity", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cemetery{TotalGrave: tt.totalGrave, UsedGrave: tt.usedGrave}
			got := c.HasAvailableSpace()
			if got != tt.want {
				t.Errorf("Cemetery.HasAvailableSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraveAvailability(t *testing.T) {
	tests := []struct {
		name   string
		status GraveStatus
		want   bool
	}{
		{"available", GraveStatusAvailable, true},
		{"occupied", GraveStatusOccupied, false},
		{"reserved", GraveStatusReserved, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Grave{Status: tt.status}
			got := g.Availability()
			if got != tt.want {
				t.Errorf("Grave.Availability() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraveIsOccupied(t *testing.T) {
	personID := int64(123)
	tests := []struct {
		name     string
		status   GraveStatus
		personID *int64
		want     bool
	}{
		{"occupied with person", GraveStatusOccupied, &personID, true},
		{"occupied without person", GraveStatusOccupied, nil, false},
		{"available", GraveStatusAvailable, nil, false},
		{"reserved", GraveStatusReserved, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Grave{Status: tt.status, PersonID: tt.personID}
			got := g.IsOccupied()
			if got != tt.want {
				t.Errorf("Grave.IsOccupied() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraveValidate(t *testing.T) {
	tests := []struct {
		name    string
		g       *Grave
		wantErr bool
	}{
		{
			name: "valid grave",
			g: &Grave{
				CemeteryID: 1,
				Section:    "A区",
				Row:        1,
				Number:     10,
			},
			wantErr: false,
		},
		{
			name: "zero cemetery id",
			g: &Grave{
				CemeteryID: 0,
				Section:    "A区",
				Row:        1,
				Number:     10,
			},
			wantErr: true,
		},
		{
			name: "negative cemetery id",
			g: &Grave{
				CemeteryID: -1,
				Section:    "A区",
				Row:        1,
				Number:     10,
			},
			wantErr: true,
		},
		{
			name: "empty section",
			g: &Grave{
				CemeteryID: 1,
				Section:    "",
				Row:        1,
				Number:     10,
			},
			wantErr: true,
		},
		{
			name: "zero row",
			g: &Grave{
				CemeteryID: 1,
				Section:    "A区",
				Row:        0,
				Number:     10,
			},
			wantErr: true,
		},
		{
			name: "negative row",
			g: &Grave{
				CemeteryID: 1,
				Section:    "A区",
				Row:        -1,
				Number:     10,
			},
			wantErr: true,
		},
		{
			name: "zero number",
			g: &Grave{
				CemeteryID: 1,
				Section:    "A区",
				Row:        1,
				Number:     0,
			},
			wantErr: true,
		},
		{
			name: "negative number",
			g: &Grave{
				CemeteryID: 1,
				Section:    "A区",
				Row:        1,
				Number:     -5,
			},
			wantErr: true,
		},
		{
			name: "valid with all fields",
			g: &Grave{
				CemeteryID:  1,
				Section:     "A区",
				Row:         5,
				Number:      20,
				Status:      GraveStatusOccupied,
				BuriedName:  "张三",
				BuriedYear:  intPtr(2020),
				Note:        "已安葬",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.g.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Grave.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}

func TestGraveStatusConstants(t *testing.T) {
	if GraveStatusAvailable != "available" {
		t.Errorf("GraveStatusAvailable = %v, want available", GraveStatusAvailable)
	}
	if GraveStatusOccupied != "occupied" {
		t.Errorf("GraveStatusOccupied = %v, want occupied", GraveStatusOccupied)
	}
	if GraveStatusReserved != "reserved" {
		t.Errorf("GraveStatusReserved = %v, want reserved", GraveStatusReserved)
	}
}

func TestCemeteryUsageRateBoundary(t *testing.T) {
	c := &Cemetery{TotalGrave: 0, UsedGrave: 0}
	if got := c.UsageRate(); got != 0.0 {
		t.Errorf("UsageRate with 0 total should be 0, got %v", got)
	}
}

func TestCemeteryHasAvailableSpaceBoundary(t *testing.T) {
	c := &Cemetery{TotalGrave: 0, UsedGrave: 0}
	if !c.HasAvailableSpace() {
		t.Error("empty cemetery should have available space")
	}
}

func TestGraveIsOccupiedBoundary(t *testing.T) {
	g := &Grave{Status: GraveStatusOccupied, PersonID: nil}
	if g.IsOccupied() {
		t.Error("occupied grave without person_id should not be occupied")
	}
}

func TestGraveWithAllFields(t *testing.T) {
	now := time.Now()
	personID := int64(999)
	buriedYear := 2023

	g := &Grave{
		ID:          1,
		CemeteryID:  10,
		PersonID:    &personID,
		Section:     "B区",
		Row:         3,
		Number:      15,
		Status:      GraveStatusOccupied,
		BuriedName:  "李四",
		BuriedDate:  &now,
		BuriedYear:  &buriedYear,
		Note:        "测试备注",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := g.Validate(); err != nil {
		t.Errorf("valid grave with all fields should pass validation: %v", err)
	}
	if !g.IsOccupied() {
		t.Error("grave with person_id should be occupied")
	}
	if g.Availability() {
		t.Error("occupied grave should not be available")
	}
}