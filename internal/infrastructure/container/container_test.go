package container_test

import (
	"testing"

	"tas-app/internal/application/usecase"
	"tas-app/internal/infrastructure/cli"
	"tas-app/internal/infrastructure/container"
)

func TestBuildCLI(t *testing.T) {
	terminalCLI, err := container.BuildCLI()
	if err != nil {
		t.Fatalf("BuildCLI: %v", err)
	}
	if terminalCLI == nil {
		t.Fatal("expected CLI")
	}
}

func TestNewResolvesAllDependencies(t *testing.T) {
	c, err := container.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	err = c.Invoke(func(
		_ cli.StudentAdder,
		_ cli.CriteriaAdder,
		_ cli.StudentEvaluator,
		_ cli.EvaluationRemover,
		_ usecase.ReportReader,
		_ cli.StudentsExporter,
		_ *cli.CLI,
	) {
	})
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
}
