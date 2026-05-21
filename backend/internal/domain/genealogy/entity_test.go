package genealogy

import (
	"testing"
	"time"
)

func TestGenealogy_Validate(t *testing.T) {
	tests := []struct {
		name    string
		g       *Genealogy
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid genealogy",
			g: &Genealogy{
				Name:             "王氏族谱",
				Surname:          "王",
				TotalGenerations: 20,
				TotalMembers:     1000,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			g: &Genealogy{
				Name:     "",
				Surname:  "王",
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name: "empty surname",
			g: &Genealogy{
				Name:    "王氏族谱",
				Surname: "",
			},
			wantErr: true,
			errMsg:  "surname cannot be empty",
		},
		{
			name: "negative generations",
			g: &Genealogy{
				Name:             "王氏族谱",
				Surname:          "王",
				TotalGenerations: -1,
			},
			wantErr: true,
			errMsg:  "total_generations cannot be negative",
		},
		{
			name: "negative members",
			g: &Genealogy{
				Name:          "王氏族谱",
				Surname:       "王",
				TotalMembers:  -1,
			},
			wantErr: true,
			errMsg:  "total_members cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.g.Validate()
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

func TestGenealogy_HasBranches(t *testing.T) {
	tests := []struct {
		name       string
		g          *Genealogy
		hasBranches bool
	}{
		{
			name: "has branches",
			g: &Genealogy{
				Branches: []*Branch{{ID: 1}, {ID: 2}},
			},
			hasBranches: true,
		},
		{
			name: "no branches",
			g: &Genealogy{
				Branches: nil,
			},
			hasBranches: false,
		},
		{
			name: "empty branches",
			g: &Genealogy{
				Branches: []*Branch{},
			},
			hasBranches: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.g.HasBranches() != tt.hasBranches {
				t.Errorf("HasBranches() = %v, want %v", tt.g.HasBranches(), tt.hasBranches)
			}
		})
	}
}

func TestBranch_Validate(t *testing.T) {
	tests := []struct {
		name    string
		b       *Branch
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid branch",
			b: &Branch{
				GenealogyID:      1,
				Name:             "长房",
				GenerationStart:  1,
				GenerationEnd:    10,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			b: &Branch{
				GenealogyID: 1,
				Name:        "",
			},
			wantErr: true,
			errMsg:  "branch name cannot be empty",
		},
		{
			name: "missing genealogy_id",
			b: &Branch{
				GenealogyID: 0,
				Name:        "长房",
			},
			wantErr: true,
			errMsg:  "genealogy_id is required",
		},
		{
			name: "negative generation start",
			b: &Branch{
				GenealogyID:     1,
				Name:           "长房",
				GenerationStart: -1,
			},
			wantErr: true,
			errMsg:  "generation cannot be negative",
		},
		{
			name: "generation start greater than end",
			b: &Branch{
				GenealogyID:     1,
				Name:           "长房",
				GenerationStart: 10,
				GenerationEnd:   5,
			},
			wantErr: true,
			errMsg:  "generation_start cannot be greater than generation_end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.b.Validate()
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

func TestGeneration_Validate(t *testing.T) {
	tests := []struct {
		name    string
		g       *Generation
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid generation",
			g: &Generation{
				GenealogyID: 1,
				Generation:  1,
				Name:        "道",
				Sequence:    1,
			},
			wantErr: false,
		},
		{
			name: "missing genealogy_id",
			g: &Generation{
				GenealogyID: 0,
				Generation:  1,
				Name:        "道",
			},
			wantErr: true,
			errMsg:  "genealogy_id is required",
		},
		{
			name: "negative generation",
			g: &Generation{
				GenealogyID: 1,
				Generation:  -1,
				Name:        "道",
			},
			wantErr: true,
			errMsg:  "generation must be non-negative",
		},
		{
			name: "empty name",
			g: &Generation{
				GenealogyID: 1,
				Generation:  1,
				Name:        "",
			},
			wantErr: true,
			errMsg:  "generation name cannot be empty",
		},
		{
			name: "negative sequence",
			g: &Generation{
				GenealogyID: 1,
				Generation: 1,
				Name:       "道",
				Sequence:   -1,
			},
			wantErr: true,
			errMsg:  "sequence must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.g.Validate()
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

func TestGeneration_IsActive(t *testing.T) {
	startYear := 2000
	endYear := 2020

	tests := []struct {
		name     string
		g        *Generation
		isActive bool
	}{
		{
			name: "active - no end year",
			g: &Generation{
				EndYear: nil,
			},
			isActive: true,
		},
		{
			name: "inactive - has end year",
			g: &Generation{
				EndYear: &endYear,
			},
			isActive: false,
		},
		{
			name: "active - start year only",
			g: &Generation{
				StartYear: &startYear,
				EndYear:   nil,
			},
			isActive: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.g.IsActive() != tt.isActive {
				t.Errorf("IsActive() = %v, want %v", tt.g.IsActive(), tt.isActive)
			}
		})
	}
}

func TestGeneration_GetGenerationRange(t *testing.T) {
	startYear := 2000
	endYear := 2020

	tests := []struct {
		name     string
		g        *Generation
		expected string
	}{
		{
			name: "both nil - returns 至今",
			g: &Generation{
				StartYear: nil,
				EndYear:   nil,
			},
			expected: "至今",
		},
		{
			name: "start only - returns start年至至今",
			g: &Generation{
				StartYear: &startYear,
				EndYear:    nil,
			},
			expected: "2000年至至今",
		},
		{
			name: "both years",
			g: &Generation{
				StartYear: &startYear,
				EndYear:   &endYear,
			},
			expected: "2000年至2020年",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.g.GetGenerationRange() != tt.expected {
				t.Errorf("GetGenerationRange() = %v, want %v", tt.g.GetGenerationRange(), tt.expected)
			}
		})
	}
}

func TestGender_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		g     Gender
		valid bool
	}{
		{"valid male", GenderMale, true},
		{"valid female", GenderFemale, true},
		{"invalid gender", Gender("未知"), false},
		{"empty gender", Gender(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.g.IsValid() != tt.valid {
				t.Errorf("IsValid() = %v, want %v", tt.g.IsValid(), tt.valid)
			}
		})
	}
}

func TestGenealogy_Fields(t *testing.T) {
	now := time.Now()
	g := &Genealogy{
		ID:               1,
		Name:             "王氏族谱",
		Surname:          "王",
		Description:      "王氏家族族谱",
		OriginPlace:      "山西",
		TotalGenerations: 20,
		TotalMembers:     1000,
		Version:          "1.0.0",
		ImageURL:         "https://example.com/cover.jpg",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if g.ID != 1 {
		t.Errorf("Genealogy.ID = %v, want 1", g.ID)
	}
	if g.Name != "王氏族谱" {
		t.Errorf("Genealogy.Name = %v, want 王氏族谱", g.Name)
	}
	if g.Surname != "王" {
		t.Errorf("Genealogy.Surname = %v, want 王", g.Surname)
	}
	if g.TotalGenerations != 20 {
		t.Errorf("Genealogy.TotalGenerations = %v, want 20", g.TotalGenerations)
	}
	if g.TotalMembers != 1000 {
		t.Errorf("Genealogy.TotalMembers = %v, want 1000", g.TotalMembers)
	}
}

func TestBranch_Fields(t *testing.T) {
	now := time.Now()
	b := &Branch{
		ID:              1,
		GenealogyID:     1,
		Name:            "长房",
		Code:            "A",
		Description:     "长子后裔",
		GenerationStart: 1,
		GenerationEnd:   10,
		MemberCount:     100,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if b.ID != 1 {
		t.Errorf("Branch.ID = %v, want 1", b.ID)
	}
	if b.Name != "长房" {
		t.Errorf("Branch.Name = %v, want 长房", b.Name)
	}
	if b.Code != "A" {
		t.Errorf("Branch.Code = %v, want A", b.Code)
	}
	if b.MemberCount != 100 {
		t.Errorf("Branch.MemberCount = %v, want 100", b.MemberCount)
	}
}

func TestGeneration_Fields(t *testing.T) {
	startYear := 2000
	endYear := 2020
	now := time.Now()
	g := &Generation{
		ID:           1,
		GenealogyID:  1,
		Generation:   1,
		Name:         "道",
		Sequence:     1,
		Description:  "道德文章",
		StartYear:    &startYear,
		EndYear:      &endYear,
		MemberCount:  50,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if g.ID != 1 {
		t.Errorf("Generation.ID = %v, want 1", g.ID)
	}
	if g.Name != "道" {
		t.Errorf("Generation.Name = %v, want 道", g.Name)
	}
	if g.Sequence != 1 {
		t.Errorf("Generation.Sequence = %v, want 1", g.Sequence)
	}
	if g.MemberCount != 50 {
		t.Errorf("Generation.MemberCount = %v, want 50", g.MemberCount)
	}
}

func TestPersonSummary_Fields(t *testing.T) {
	birthYear := 1950
	deathYear := 2020
	ps := &PersonSummary{
		ID:         1,
		Name:       "王明",
		Gender:     GenderMale,
		Generation: 10,
		BirthYear:  &birthYear,
		DeathYear:  &deathYear,
		StyleName:  "文正",
	}

	if ps.ID != 1 {
		t.Errorf("PersonSummary.ID = %v, want 1", ps.ID)
	}
	if ps.Name != "王明" {
		t.Errorf("PersonSummary.Name = %v, want 王明", ps.Name)
	}
	if ps.Gender != GenderMale {
		t.Errorf("PersonSummary.Gender = %v, want 男", ps.Gender)
	}
	if ps.Generation != 10 {
		t.Errorf("PersonSummary.Generation = %v, want 10", ps.Generation)
	}
}