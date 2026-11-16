package cli

import (
	"bytes"
	"strings"
	"testing"

	"tas-app/internal/application/usecase"
)

type stubStudentAdder struct{}

func (stubStudentAdder) Execute(usecase.AddStudentCommand) error { return nil }

type stubCriteriaAdder struct{}

func (stubCriteriaAdder) Execute(usecase.AddCriteriaCommand) error { return nil }

type stubStudentEvaluator struct{}

func (stubStudentEvaluator) Execute(usecase.StudentEvaluationCommand) error { return nil }

type stubEvaluationRemover struct{}

func (stubEvaluationRemover) Execute(usecase.RemoveEvaluationCommand) error { return nil }

type stubReportReader struct{}

func (stubReportReader) GetStudents() ([]usecase.StudentDTO, error) { return nil, nil }

func (stubReportReader) GetCriteria() ([]usecase.CriteriaDTO, error) { return nil, nil }

type stubExporter struct{}

func (stubExporter) ExportStudents(string) error { return nil }

func TestNewCLIRequiresDependencies(t *testing.T) {
	_, err := NewCLI()
	if err == nil {
		t.Fatal("expected error")
	}

	for _, name := range []string{
		"AddStudent",
		"AddCriteria",
		"Evaluate",
		"RemoveEvaluation",
		"GetReport",
		"ExportStudents",
	} {
		if !strings.Contains(err.Error(), name) {
			t.Fatalf("error %q does not mention %s", err.Error(), name)
		}
	}
}

func TestNewCLIWithDependencies(t *testing.T) {
	cli, err := NewCLI(
		WithAddStudent(stubStudentAdder{}),
		WithAddCriteria(stubCriteriaAdder{}),
		WithEvaluate(stubStudentEvaluator{}),
		WithRemoveEvaluation(stubEvaluationRemover{}),
		WithReport(stubReportReader{}),
		WithExporter(stubExporter{}),
	)
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}
	if cli == nil {
		t.Fatal("expected CLI")
	}
}

func TestCLIRunExit(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cli, err := NewCLI(
		WithReader(strings.NewReader("0\n")),
		WithWriter(&stdout),
		WithErrorWriter(&stderr),
		WithAddStudent(stubStudentAdder{}),
		WithAddCriteria(stubCriteriaAdder{}),
		WithEvaluate(stubStudentEvaluator{}),
		WithRemoveEvaluation(stubEvaluationRemover{}),
		WithReport(stubReportReader{}),
		WithExporter(stubExporter{}),
	)
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	cli.Run()

	if !strings.Contains(stdout.String(), "до встречи!") {
		t.Fatalf("stdout does not contain exit message: %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}
