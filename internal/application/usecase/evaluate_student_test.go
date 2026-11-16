package usecase_test

import (
	"testing"

	"tas-app/internal/application/usecase"
	"tas-app/internal/domain/model"
	"tas-app/internal/infrastructure/storage"
)

func newEvaluateFixture(t *testing.T) (
	*storage.MemoryStudentRepo,
	*storage.MemoryCriteriaRepo,
	*usecase.EvaluateStudentUseCase,
	*model.Student,
	*model.Criteria,
) {
	t.Helper()

	criteria := testCriteria(t, "Математика", 100, 0.5)
	student := testStudent(t, 1, "Иван Иванов")

	studentRepo := storage.NewMemoryStudentRepo()
	criteriaRepo := storage.NewMemoryCriteriaRepo()
	if err := studentRepo.Save(student); err != nil {
		t.Fatalf("Save student: %v", err)
	}
	if err := criteriaRepo.Save(criteria); err != nil {
		t.Fatalf("Save criteria: %v", err)
	}

	useCase := usecase.NewEvaluateStudentUseCase(studentRepo, criteriaRepo)
	return studentRepo, criteriaRepo, useCase, student, criteria
}

func TestEvaluateStudentUseCaseSetsScore(t *testing.T) {
	_, _, useCase, student, criteria := newEvaluateFixture(t)

	err := useCase.Execute(usecase.StudentEvaluationCommand{
		StudentID:    1,
		CriteriaName: "Математика",
		Score:        80,
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	score, exists := student.GetScore(criteria.Name())
	if !exists {
		t.Fatal("score not found")
	}
	if score.Value() != 80 {
		t.Fatalf("got score %d, want 80", score.Value())
	}
}

func TestEvaluateStudentUseCaseRejectsInvalidScore(t *testing.T) {
	_, _, useCase, student, criteria := newEvaluateFixture(t)

	for _, score := range []int{-1, 101} {
		err := useCase.Execute(usecase.StudentEvaluationCommand{
			StudentID:    1,
			CriteriaName: "Математика",
			Score:        score,
		})
		if err == nil {
			t.Fatalf("expected error for score %d", score)
		}
	}

	if _, exists := student.GetScore(criteria.Name()); exists {
		t.Fatal("invalid score was added")
	}
}

func TestEvaluateStudentUseCaseRejectsUnknownCriteria(t *testing.T) {
	_, _, useCase, _, _ := newEvaluateFixture(t)

	err := useCase.Execute(usecase.StudentEvaluationCommand{
		StudentID:    1,
		CriteriaName: "Физика",
		Score:        10,
	})
	if err == nil {
		t.Fatal("expected error for unknown criteria")
	}
}

func TestEvaluateStudentUseCaseRejectsUnknownStudent(t *testing.T) {
	criteriaRepo := storage.NewMemoryCriteriaRepo()
	if err := criteriaRepo.Save(testCriteria(t, "Математика", 100, 0.5)); err != nil {
		t.Fatalf("Save criteria: %v", err)
	}

	studentRepo := storage.NewMemoryStudentRepo()
	useCase := usecase.NewEvaluateStudentUseCase(studentRepo, criteriaRepo)

	err := useCase.Execute(usecase.StudentEvaluationCommand{
		StudentID:    42,
		CriteriaName: "Математика",
		Score:        10,
	})
	if err == nil {
		t.Fatal("expected error for unknown student")
	}
}

func TestEvaluateStudentUseCaseOverwritesScore(t *testing.T) {
	_, _, useCase, student, criteria := newEvaluateFixture(t)

	for _, score := range []int{80, 90} {
		err := useCase.Execute(usecase.StudentEvaluationCommand{
			StudentID:    1,
			CriteriaName: "Математика",
			Score:        score,
		})
		if err != nil {
			t.Fatalf("Execute score %d: %v", score, err)
		}
	}

	got, exists := student.GetScore(criteria.Name())
	if !exists {
		t.Fatal("score not found")
	}
	if got.Value() != 90 {
		t.Fatalf("got score %d, want 90", got.Value())
	}
}
