package usecase_test

import (
	"testing"

	"tas-app/internal/application/usecase"
	"tas-app/internal/domain/model"
	"tas-app/internal/infrastructure/storage"
)

func testCriteria(t *testing.T, name string, max int, weight float64) *model.Criteria {
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

func testStudent(t *testing.T, id int, name string) *model.Student {
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

func TestRemoveEvaluationUseCase(t *testing.T) {
	criteria := testCriteria(t, "Математика", 100, 0.5)
	student := testStudent(t, 1, "Иван Иванов")

	evaluation, err := model.NewEvaluation(criteria, 80)
	if err != nil {
		t.Fatalf("NewEvaluation: %v", err)
	}
	if err := student.AddEvaluation(evaluation); err != nil {
		t.Fatalf("AddEvaluation: %v", err)
	}

	studentRepo := storage.NewMemoryStudentRepo()
	criteriaRepo := storage.NewMemoryCriteriaRepo()
	if err := studentRepo.Save(student); err != nil {
		t.Fatalf("Save student: %v", err)
	}
	if err := criteriaRepo.Save(criteria); err != nil {
		t.Fatalf("Save criteria: %v", err)
	}

	useCase := usecase.NewRemoveEvaluationUseCase(studentRepo, criteriaRepo)
	err = useCase.Execute(usecase.RemoveEvaluationCommand{
		StudentID:    1,
		CriteriaName: "Математика",
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if _, exists := student.GetScore(criteria.Name()); exists {
		t.Fatal("evaluation still exists")
	}
}

func TestRemoveEvaluationUseCaseMissingScore(t *testing.T) {
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

	useCase := usecase.NewRemoveEvaluationUseCase(studentRepo, criteriaRepo)
	err := useCase.Execute(usecase.RemoveEvaluationCommand{
		StudentID:    1,
		CriteriaName: "Математика",
	})
	if err == nil {
		t.Fatal("expected error for missing evaluation")
	}
}

func TestRemoveEvaluationUseCaseUnknownCriteria(t *testing.T) {
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

	useCase := usecase.NewRemoveEvaluationUseCase(studentRepo, criteriaRepo)
	err := useCase.Execute(usecase.RemoveEvaluationCommand{
		StudentID:    1,
		CriteriaName: "Физика",
	})
	if err == nil {
		t.Fatal("expected error for unknown criteria")
	}
}

func TestRemoveEvaluationUseCaseEmptyCriteriaName(t *testing.T) {
	studentRepo := storage.NewMemoryStudentRepo()
	criteriaRepo := storage.NewMemoryCriteriaRepo()

	useCase := usecase.NewRemoveEvaluationUseCase(studentRepo, criteriaRepo)
	err := useCase.Execute(usecase.RemoveEvaluationCommand{
		StudentID:    1,
		CriteriaName: "   ",
	})
	if err == nil {
		t.Fatal("expected error for empty criteria name")
	}
}
