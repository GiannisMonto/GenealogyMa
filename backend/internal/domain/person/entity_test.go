package person

import (
	"testing"
)

func TestPersonGenerationBoundary(t *testing.T) {
	tests := []struct {
		name       string
		generation int
		wantErr    bool
	}{
		{"generation 0", 0, false},
		{"generation 1", 1, false},
		{"generation 30", 30, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Person{
				Name:       "测试",
				Gender:     GenderMale,
				Generation: tt.generation,
			}
			err := p.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Person.Validate() for generation %d error = %v, wantErr %v", tt.generation, err, tt.wantErr)
			}
		})
	}
}

func TestGenderConstants(t *testing.T) {
	if GenderMale != "男" {
		t.Errorf("GenderMale = %v, want 男", GenderMale)
	}
	if GenderFemale != "女" {
		t.Errorf("GenderFemale = %v, want 女", GenderFemale)
	}
	if GenderUnknown != "" {
		t.Errorf("GenderUnknown = %v, want empty", GenderUnknown)
	}
}

func TestRelationTypeConstants(t *testing.T) {
	if RelationBiological != "biological" {
		t.Errorf("RelationBiological = %v, want biological", RelationBiological)
	}
	if RelationAdoptive != "adoptive" {
		t.Errorf("RelationAdoptive = %v, want adoptive", RelationAdoptive)
	}
	if RelationStep != "step" {
		t.Errorf("RelationStep = %v, want step", RelationStep)
	}
}

func TestSpouseTypeConstants(t *testing.T) {
	if SpouseTypePrimary != "配" {
		t.Errorf("SpouseTypePrimary = %v, want 配", SpouseTypePrimary)
	}
	if SpouseTypeSecond != "继" {
		t.Errorf("SpouseTypeSecond = %v, want 继", SpouseTypeSecond)
	}
	if SpouseTypeDeceased != "妣" {
		t.Errorf("SpouseTypeDeceased = %v, want 妣", SpouseTypeDeceased)
	}
	if SpouseTypeThird != "三" {
		t.Errorf("SpouseTypeThird = %v, want 三", SpouseTypeThird)
	}
}