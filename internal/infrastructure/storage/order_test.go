package storage_test

import (
	"testing"

	"tas-app/internal/domain/model"
	"tas-app/internal/infrastructure/storage"
)

func orderCriteria(t *testing.T, name string) *model.Criteria {
	t.Helper()

	criteriaName, err := model.NewCriteriaName(name)
	if err != nil {
		t.Fatalf("NewCriteriaName: %v", err)
	}
	maxScore, err := model.NewMaxScore(10)
	if err != nil {
		t.Fatalf("NewMaxScore: %v", err)
	}
	weight, err := model.NewWeight(0.5)
	if err != nil {
		t.Fatalf("NewWeight: %v", err)
	}
	criteria, err := model.NewCriteria(criteriaName, maxScore, weight)
	if err != nil {
		t.Fatalf("NewCriteria: %v", err)
	}
	return criteria
}

func orderStudent(t *testing.T, id int) *model.Student {
	t.Helper()

	studentID, err := model.NewStudentID(id)
	if err != nil {
		t.Fatalf("NewStudentID: %v", err)
	}
	name, err := model.NewStudentName("Студент")
	if err != nil {
		t.Fatalf("NewStudentName: %v", err)
	}
	student, err := model.NewStudent(studentID, name)
	if err != nil {
		t.Fatalf("NewStudent: %v", err)
	}
	return student
}

func TestMemoryStudentRepoGetAllReturnsSortedByID(t *testing.T) {
	repo := storage.NewMemoryStudentRepo()
	for _, id := range []int{3, 1, 2} {
		if err := repo.Save(orderStudent(t, id)); err != nil {
			t.Fatalf("Save student %d: %v", id, err)
		}
	}

	students, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if len(students) != 3 {
		t.Fatalf("got %d students, want 3", len(students))
	}

	for i, want := range []int{1, 2, 3} {
		if got := students[i].ID().Value(); got != want {
			t.Fatalf("students[%d] id = %d, want %d", i, got, want)
		}
	}
}

func TestMemoryCriteriaRepoGetAllReturnsSortedByName(t *testing.T) {
	repo := storage.NewMemoryCriteriaRepo()
	for _, name := range []string{"Химия", "Алгебра", "Биология"} {
		if err := repo.Save(orderCriteria(t, name)); err != nil {
			t.Fatalf("Save criteria %q: %v", name, err)
		}
	}

	criteria, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if len(criteria) != 3 {
		t.Fatalf("got %d criteria, want 3", len(criteria))
	}

	for i, want := range []string{"Алгебра", "Биология", "Химия"} {
		if got := criteria[i].Name().Value(); got != want {
			t.Fatalf("criteria[%d] name = %q, want %q", i, got, want)
		}
	}
}
