package model

import (
	"math"
	"testing"
)

func TestNewCriteriaName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid", input: "Математика", want: "Математика"},
		{name: "trimmed", input: "  Математика  ", want: "Математика"},
		{name: "empty", input: "", wantErr: true},
		{name: "spaces", input: "   ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCriteriaName(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Value() != tt.want {
				t.Fatalf("got %q, want %q", got.Value(), tt.want)
			}
		})
	}
}

func TestNewWeight(t *testing.T) {
	tests := []struct {
		name    string
		input   float64
		wantErr bool
	}{
		{name: "zero", input: 0},
		{name: "half", input: 0.5},
		{name: "one", input: 1},
		{name: "negative", input: -0.1, wantErr: true},
		{name: "too large", input: 1.1, wantErr: true},
		{name: "nan", input: math.NaN(), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWeight(tt.input)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestNewMaxScore(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{name: "zero", input: 0, wantErr: true},
		{name: "positive", input: 10},
		{name: "negative", input: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMaxScore(tt.input)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestNewCriteria(t *testing.T) {
	criteria := mustCriteria(t, "Математика", 100, 0.5)

	if criteria.Name().Value() != "Математика" {
		t.Fatalf("got name %q", criteria.Name().Value())
	}
	if criteria.MaxScore().Value() != 100 {
		t.Fatalf("got max score %d", criteria.MaxScore().Value())
	}
	if criteria.Weight().Value() != 0.5 {
		t.Fatalf("got weight %f", criteria.Weight().Value())
	}
}

func TestNewCriteriaRejectsZeroValueName(t *testing.T) {
	if _, err := NewCriteria(CriteriaName{}, MaxScore{}, Weight{}); err == nil {
		t.Fatal("expected error for empty criteria name")
	}
}

func TestNewCriteriaRejectsZeroValueMaxScore(t *testing.T) {
	name, err := NewCriteriaName("Математика")
	if err != nil {
		t.Fatalf("NewCriteriaName: %v", err)
	}

	if _, err := NewCriteria(name, MaxScore{}, Weight{}); err == nil {
		t.Fatal("expected error for zero value max score")
	}
}

func TestCriteriaValid(t *testing.T) {
	if (&Criteria{}).Valid() {
		t.Fatal("zero value criteria must be invalid")
	}

	criteria := mustCriteria(t, "Математика", 10, 0.5)
	if !criteria.Valid() {
		t.Fatal("constructed criteria must be valid")
	}

	criteria.maxScore = MaxScore{}
	if criteria.Valid() {
		t.Fatal("criteria with zero max score must be invalid")
	}

	criteria = mustCriteria(t, "Математика", 10, 0.5)
	criteria.weight = Weight{value: 2}
	if criteria.Valid() {
		t.Fatal("criteria with invalid weight must be invalid")
	}
}
