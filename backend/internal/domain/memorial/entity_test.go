package memorial

import (
	"testing"
	"time"
)

func TestMemorialHall_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       *MemorialHall
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid hall",
			m: &MemorialHall{
				Name:        "王氏宗祠",
				Province:    "浙江",
				TotalTablet: 100,
				UsedTablet:  50,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			m: &MemorialHall{
				Name:        "",
				Province:    "浙江",
				TotalTablet: 100,
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name: "empty province",
			m: &MemorialHall{
				Name:     "王氏宗祠",
				Province: "",
				TotalTablet: 100,
			},
			wantErr: true,
			errMsg:  "province is required",
		},
		{
			name: "negative total tablet",
			m: &MemorialHall{
				Name:        "王氏宗祠",
				Province:    "浙江",
				TotalTablet: -1,
			},
			wantErr: true,
			errMsg:  "total_tablet cannot be negative",
		},
		{
			name: "negative used tablet",
			m: &MemorialHall{
				Name:        "王氏宗祠",
				Province:    "浙江",
				TotalTablet: 100,
				UsedTablet:  -1,
			},
			wantErr: true,
			errMsg:  "used_tablet cannot be negative",
		},
		{
			name: "used exceeds total",
			m: &MemorialHall{
				Name:        "王氏宗祠",
				Province:    "浙江",
				TotalTablet: 100,
				UsedTablet:  101,
			},
			wantErr: true,
			errMsg:  "used_tablet cannot exceed total_tablet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestMemorialHall_UsageRate(t *testing.T) {
	tests := []struct {
		name       string
		m          *MemorialHall
		wantRate   float64
	}{
		{
			name: "50 percent usage",
			m: &MemorialHall{
				TotalTablet: 100,
				UsedTablet:  50,
			},
			wantRate: 50.0,
		},
		{
			name: "zero total tablet",
			m: &MemorialHall{
				TotalTablet: 0,
				UsedTablet:  0,
			},
			wantRate: 0.0,
		},
		{
			name: "100 percent usage",
			m: &MemorialHall{
				TotalTablet: 100,
				UsedTablet:  100,
			},
			wantRate: 100.0,
		},
		{
			name: "25 percent usage",
			m: &MemorialHall{
				TotalTablet: 200,
				UsedTablet:  50,
			},
			wantRate: 25.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate := tt.m.UsageRate()
			if rate != tt.wantRate {
				t.Errorf("UsageRate() = %v, want %v", rate, tt.wantRate)
			}
		})
	}
}

func TestMemorialHall_HasAvailableSpace(t *testing.T) {
	tests := []struct {
		name         string
		m            *MemorialHall
		wantAvailable bool
	}{
		{
			name: "has space available",
			m: &MemorialHall{
				TotalTablet: 100,
				UsedTablet:  50,
			},
			wantAvailable: true,
		},
		{
			name: "no space available",
			m: &MemorialHall{
				TotalTablet: 100,
				UsedTablet:  100,
			},
			wantAvailable: false,
		},
		{
			name: "empty hall",
			m: &MemorialHall{
				TotalTablet: 100,
				UsedTablet:  0,
			},
			wantAvailable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			available := tt.m.HasAvailableSpace()
			if available != tt.wantAvailable {
				t.Errorf("HasAvailableSpace() = %v, want %v", available, tt.wantAvailable)
			}
		})
	}
}

func TestMemorialTablet_Validate(t *testing.T) {
	personID := int64(1)
	tests := []struct {
		name    string
		m       *MemorialTablet
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid tablet with person_id",
			m: &MemorialTablet{
				HallID:     1,
				PersonID:   &personID,
				PersonName: "王明",
				Floor:      1,
				Row:        1,
				Number:     1,
			},
			wantErr: false,
		},
		{
			name: "valid tablet with person_name",
			m: &MemorialTablet{
				HallID:     1,
				PersonName: "王明",
				Floor:      1,
				Row:        1,
				Number:     1,
			},
			wantErr: false,
		},
		{
			name: "missing hall_id",
			m: &MemorialTablet{
				HallID:     0,
				PersonName: "王明",
				Floor:      1,
				Row:        1,
				Number:     1,
			},
			wantErr: true,
			errMsg:  "hall_id is required",
		},
		{
			name: "missing person info",
			m: &MemorialTablet{
				HallID:     1,
				PersonName: "",
				Floor:      1,
				Row:        1,
				Number:     1,
			},
			wantErr: true,
			errMsg:  "either person_name or person_id is required",
		},
		{
			name: "negative floor",
			m: &MemorialTablet{
				HallID:     1,
				PersonName: "王明",
				Floor:      -1,
				Row:        1,
				Number:     1,
			},
			wantErr: true,
			errMsg:  "floor cannot be negative",
		},
		{
			name: "negative row",
			m: &MemorialTablet{
				HallID:     1,
				PersonName: "王明",
				Floor:      1,
				Row:        -1,
				Number:     1,
			},
			wantErr: true,
			errMsg:  "row cannot be negative",
		},
		{
			name: "non-positive number",
			m: &MemorialTablet{
				HallID:     1,
				PersonName: "王明",
				Floor:      1,
				Row:        1,
				Number:     0,
			},
			wantErr: true,
			errMsg:  "number must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestMemorialTablet_GetPositionDesc(t *testing.T) {
	tests := []struct {
		name   string
		m      *MemorialTablet
		want   string
	}{
		{
			name: "normal position",
			m: &MemorialTablet{
				Floor:  1,
				Row:    2,
				Number: 3,
			},
			want: "1层 2排 3号",
		},
		{
			name: "ground floor",
			m: &MemorialTablet{
				Floor:  0,
				Row:    1,
				Number: 5,
			},
			want: "0层 1排 5号",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := tt.m.GetPositionDesc()
			if desc != tt.want {
				t.Errorf("GetPositionDesc() = %v, want %v", desc, tt.want)
			}
		})
	}
}

func TestMemorialTablet_IsOccupied(t *testing.T) {
	personID := int64(1)
	tests := []struct {
		name      string
		m         *MemorialTablet
		isOccupied bool
	}{
		{
			name: "occupied with person_id",
			m: &MemorialTablet{
				PersonID: &personID,
			},
			isOccupied: true,
		},
		{
			name: "empty hall",
			m: &MemorialTablet{
				PersonID: nil,
			},
			isOccupied: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			occupied := tt.m.IsOccupied()
			if occupied != tt.isOccupied {
				t.Errorf("IsOccupied() = %v, want %v", occupied, tt.isOccupied)
			}
		})
	}
}

func TestGender_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		g       Gender
		want    bool
	}{
		{
			name: "valid male",
			g:    GenderMale,
			want: true,
		},
		{
			name: "valid female",
			g:    GenderFemale,
			want: true,
		},
		{
			name: "invalid gender",
			g:    Gender("未知"),
			want: false,
		},
		{
			name: "empty gender",
			g:    Gender(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.g.IsValid()
			if valid != tt.want {
				t.Errorf("IsValid() = %v, want %v", valid, tt.want)
			}
		})
	}
}

func TestTabletType_Constants(t *testing.T) {
	if TabletTypeAncestor != "ancestor" {
		t.Errorf("TabletTypeAncestor = %v, want ancestor", TabletTypeAncestor)
	}
	if TabletTypeMartyr != "martyr" {
		t.Errorf("TabletTypeMartyr = %v, want martyr", TabletTypeMartyr)
	}
	if TabletTypeSage != "sage" {
		t.Errorf("TabletTypeSage = %v, want sage", TabletTypeSage)
	}
	if TabletTypeFounder != "founder" {
		t.Errorf("TabletTypeFounder = %v, want founder", TabletTypeFounder)
	}
}

func TestMemorialTablet_Enthronement(t *testing.T) {
	now := time.Now()
	m := &MemorialTablet{
		Enthronement: &now,
	}
	if m.Enthronement == nil {
		t.Error("Enthronement should not be nil")
	}
}

func TestPersonInfo(t *testing.T) {
	birthYear := 1950
	deathYear := 2020
	p := &PersonInfo{
		ID:        1,
		Name:      "王明",
		StyleName: "文正",
		Gender:    GenderMale,
		BirthYear: &birthYear,
		DeathYear: &deathYear,
	}

	if p.ID != 1 {
		t.Errorf("PersonInfo.ID = %v, want 1", p.ID)
	}
	if p.Name != "王明" {
		t.Errorf("PersonInfo.Name = %v, want 王明", p.Name)
	}
	if p.StyleName != "文正" {
		t.Errorf("PersonInfo.StyleName = %v, want 文正", p.StyleName)
	}
	if p.Gender != GenderMale {
		t.Errorf("PersonInfo.Gender = %v, want 男", p.Gender)
	}
	if *p.BirthYear != 1950 {
		t.Errorf("PersonInfo.BirthYear = %v, want 1950", *p.BirthYear)
	}
	if *p.DeathYear != 2020 {
		t.Errorf("PersonInfo.DeathYear = %v, want 2020", *p.DeathYear)
	}
}