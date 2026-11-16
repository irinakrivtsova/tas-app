package usecase_test

import (
	"testing"

	"tas-app/internal/application/usecase"
	"tas-app/internal/infrastructure/storage"
)

func TestAddStudentUseCaseRejectsDuplicateID(t *testing.T) {
	repo := storage.NewMemoryStudentRepo()
	useCase := usecase.NewAddStudentUseCase(repo)

	first := usecase.AddStudentCommand{ID: 1, Name: "Иван Иванов"}
	if err := useCase.Execute(first); err != nil {
		t.Fatalf("first Execute: %v", err)
	}

	duplicate := usecase.AddStudentCommand{ID: 1, Name: "Пётр Петров"}
	if err := useCase.Execute(duplicate); err == nil {
		t.Fatal("expected error for duplicate student id")
	}
}

func TestAddCriteriaUseCaseRejectsDuplicateName(t *testing.T) {
	repo := storage.NewMemoryCriteriaRepo()
	useCase := usecase.NewAddCriteriaUseCase(repo)

	first := usecase.AddCriteriaCommand{
		Name:        "Математика",
		MaxScoreInt: 100,
		WeightFloat: 0.5,
	}
	if err := useCase.Execute(first); err != nil {
		t.Fatalf("first Execute: %v", err)
	}

	duplicate := usecase.AddCriteriaCommand{
		Name:        "Математика",
		MaxScoreInt: 50,
		WeightFloat: 0.5,
	}
	if err := useCase.Execute(duplicate); err == nil {
		t.Fatal("expected error for duplicate criteria name")
	}
}
