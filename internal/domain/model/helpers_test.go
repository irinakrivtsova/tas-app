package model

import "testing"

func mustCriteria(t *testing.T, name string, max int, weight float64) *Criteria {
	t.Helper()

	criteriaName, err := NewCriteriaName(name)
	if err != nil {
		t.Fatalf("NewCriteriaName: %v", err)
	}

	maxScore, err := NewMaxScore(max)
	if err != nil {
		t.Fatalf("NewMaxScore: %v", err)
	}

	criteriaWeight, err := NewWeight(weight)
	if err != nil {
		t.Fatalf("NewWeight: %v", err)
	}

	criteria, err := NewCriteria(criteriaName, maxScore, criteriaWeight)
	if err != nil {
		t.Fatalf("NewCriteria: %v", err)
	}
	return criteria
}

func mustStudentID(t *testing.T, value int) StudentID {
	t.Helper()

	id, err := NewStudentID(value)
	if err != nil {
		t.Fatalf("NewStudentID: %v", err)
	}
	return id
}

func mustStudentName(t *testing.T, value string) StudentName {
	t.Helper()

	name, err := NewStudentName(value)
	if err != nil {
		t.Fatalf("NewStudentName: %v", err)
	}
	return name
}

func mustStudent(t *testing.T, id int, name string) *Student {
	t.Helper()

	student, err := NewStudent(mustStudentID(t, id), mustStudentName(t, name))
	if err != nil {
		t.Fatalf("NewStudent: %v", err)
	}
	return student
}
