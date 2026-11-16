package usecase_test

import (
	"testing"

	"tas-app/internal/application/usecase"
	"tas-app/internal/infrastructure/storage"
)

func TestAddStudentUseCaseRejectsEmptyName(t *testing.T) {
	repo := storage.NewMemoryStudentRepo()
	useCase := usecase.NewAddStudentUseCase(repo)

	err := useCase.Execute(usecase.AddStudentCommand{
		ID:   1,
		Name: "   ",
	})
	if err == nil {
		t.Fatal("expected error for empty student name")
	}
}

func TestAddCriteriaUseCaseRejectsEmptyName(t *testing.T) {
	repo := storage.NewMemoryCriteriaRepo()
	useCase := usecase.NewAddCriteriaUseCase(repo)

	err := useCase.Execute(usecase.AddCriteriaCommand{
		Name:        "",
		MaxScoreInt: 10,
		WeightFloat: 0.5,
	})
	if err == nil {
		t.Fatal("expected error for empty criteria name")
	}
}
