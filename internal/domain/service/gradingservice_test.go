package service

import (
	"testing"

	"tas-app/internal/domain/model"
)

func gradingCriteria(t *testing.T, name string, max int, weight float64) *model.Criteria {
	t.Helper()

	criteriaName, err := model.NewCriteriaName(name)
	if err != nil {
		t.Fatalf("NewCriteriaName: %v", err)
	}
	maxScore, err := model.NewMaxScore(max)
	if err != nil {
		t.Fatalf("NewMaxScore: %v", err)
	}
	criteriaWeight, err := model.NewWeight(weight)
	if err != nil {
		t.Fatalf("NewWeight: %v", err)
	}
	criteria, err := model.NewCriteria(criteriaName, maxScore, criteriaWeight)
	if err != nil {
		t.Fatalf("NewCriteria: %v", err)
	}
	return criteria
}

func gradingStudent(t *testing.T, id int, name string) *model.Student {
	t.Helper()

	studentID, err := model.NewStudentID(id)
	if err != nil {
		t.Fatalf("NewStudentID: %v", err)
	}
	studentName, err := model.NewStudentName(name)
	if err != nil {
		t.Fatalf("NewStudentName: %v", err)
	}
	student, err := model.NewStudent(studentID, studentName)
	if err != nil {
		t.Fatalf("NewStudent: %v", err)
	}
	return student
}

func TestCalculateFinalGrade(t *testing.T) {
	first := gradingCriteria(t, "Первый", 10, 0.5)
	second := gradingCriteria(t, "Второй", 20, 0.5)
	student := gradingStudent(t, 1, "Иван Иванов")

	firstEvaluation, err := model.NewEvaluation(first, 10)
	if err != nil {
		t.Fatalf("NewEvaluation: %v", err)
	}
	secondEvaluation, err := model.NewEvaluation(second, 10)
	if err != nil {
		t.Fatalf("NewEvaluation: %v", err)
	}
	if err := student.AddEvaluation(firstEvaluation); err != nil {
		t.Fatalf("AddEvaluation: %v", err)
	}
	if err := student.AddEvaluation(secondEvaluation); err != nil {
		t.Fatalf("AddEvaluation: %v", err)
	}

	service := &GradingService{}
	got := service.CalculateFinalGrade(student, []*model.Criteria{first, second})
	if got != 10 {
		t.Fatalf("got final grade %f, want 10", got)
	}
}

func TestCalculateFinalGradeWithoutScore(t *testing.T) {
	first := gradingCriteria(t, "Первый", 10, 0.5)
	second := gradingCriteria(t, "Второй", 20, 0.5)
	student := gradingStudent(t, 1, "Иван Иванов")

	evaluation, err := model.NewEvaluation(first, 10)
	if err != nil {
		t.Fatalf("NewEvaluation: %v", err)
	}
	if err := student.AddEvaluation(evaluation); err != nil {
		t.Fatalf("AddEvaluation: %v", err)
	}

	service := &GradingService{}
	got := service.CalculateFinalGrade(student, []*model.Criteria{first, second})
	if got != 5 {
		t.Fatalf("got final grade %f, want 5", got)
	}
}

func TestCalculateMaxTotal(t *testing.T) {
	first := gradingCriteria(t, "Первый", 10, 0.5)
	second := gradingCriteria(t, "Второй", 20, 0.5)

	service := &GradingService{}
	got := service.CalculateMaxTotal([]*model.Criteria{first, second})
	if got != 15 {
		t.Fatalf("got max total %f, want 15", got)
	}
}
