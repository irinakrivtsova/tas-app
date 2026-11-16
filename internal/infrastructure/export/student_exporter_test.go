package export_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"tas-app/internal/application/usecase"
	"tas-app/internal/domain/model"
	"tas-app/internal/domain/service"
	"tas-app/internal/infrastructure/export"
	"tas-app/internal/infrastructure/storage"
)

func exportCriteria(t *testing.T, name string, max int, weight float64) *model.Criteria {
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

func exportStudent(t *testing.T, id int, name string) *model.Student {
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

func newExportReport(t *testing.T) *usecase.GetReportUseCase {
	t.Helper()

	criteria := exportCriteria(t, "Математика", 100, 0.5)
	student := exportStudent(t, 1, "Иван Иванов")

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

	return usecase.NewGetReportUseCase(studentRepo, criteriaRepo, &service.GradingService{})
}

func TestFileExporterWritesCSV(t *testing.T) {
	exporter := export.NewFileExporter(newExportReport(t))
	exporter.RegisterFormatter(".csv", &export.CSVFormatter{})

	path := filepath.Join(t.TempDir(), "students.csv")
	if err := exporter.ExportStudents(path); err != nil {
		t.Fatalf("ExportStudents: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "ID,Имя,Математика,Итог,Макс. Итог") {
		t.Fatalf("unexpected CSV headers: %q", content)
	}
	if !strings.Contains(content, "1,Иван Иванов,80,40.00,50.00") {
		t.Fatalf("unexpected CSV row: %q", content)
	}
}

func TestFileExporterWritesJSON(t *testing.T) {
	exporter := export.NewFileExporter(newExportReport(t))
	exporter.RegisterFormatter(".json", &export.JSONFormatter{})

	path := filepath.Join(t.TempDir(), "students.json")
	if err := exporter.ExportStudents(path); err != nil {
		t.Fatalf("ExportStudents: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var students []usecase.StudentDTO
	if err := json.Unmarshal(data, &students); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(students) != 1 {
		t.Fatalf("got %d students, want 1", len(students))
	}

	student := students[0]
	if student.Total != 40 {
		t.Fatalf("got total %f, want 40", student.Total)
	}
	if student.MaxTotal != 50 {
		t.Fatalf("got max total %f, want 50", student.MaxTotal)
	}
	if student.Scores["Математика"] != 80 {
		t.Fatalf("got score %d, want 80", student.Scores["Математика"])
	}
}

func TestFileExporterRejectsUnsupportedExtension(t *testing.T) {
	exporter := export.NewFileExporter(newExportReport(t))
	exporter.RegisterFormatter(".csv", &export.CSVFormatter{})

	path := filepath.Join(t.TempDir(), "students.txt")
	if err := exporter.ExportStudents(path); err == nil {
		t.Fatal("expected error for unsupported extension")
	}
}
