package model

import "testing"

func TestNewEvaluation(t *testing.T) {
	criteria := mustCriteria(t, "Математика", 10, 0.5)

	evaluation, err := NewEvaluation(criteria, 5)
	if err != nil {
		t.Fatalf("NewEvaluation: %v", err)
	}
	if evaluation.Criteria().Value() != "Математика" {
		t.Fatalf("got criteria %q", evaluation.Criteria().Value())
	}
	if evaluation.Score().Value() != 5 {
		t.Fatalf("got score %d", evaluation.Score().Value())
	}
}

func TestNewEvaluationRejectsInvalidScore(t *testing.T) {
	criteria := mustCriteria(t, "Математика", 10, 0.5)

	for _, value := range []int{-1, 11} {
		if _, err := NewEvaluation(criteria, value); err == nil {
			t.Fatalf("expected error for score %d", value)
		}
	}
}

func TestNewEvaluationRejectsNilCriteria(t *testing.T) {
	if _, err := NewEvaluation(nil, 5); err == nil {
		t.Fatal("expected error for nil criteria")
	}
}

func TestNewEvaluationRejectsZeroValueCriteria(t *testing.T) {
	if _, err := NewEvaluation(&Criteria{}, 5); err == nil {
		t.Fatal("expected error for zero value criteria")
	}
}
