package container

import (
	"fmt"

	"tas-app/internal/application/usecase"
	"tas-app/internal/infrastructure/cli"
	"tas-app/internal/infrastructure/export"
)

func newFileExporter(reader usecase.ReportReader) *export.FileExporter {
	exporter := export.NewFileExporter(reader)
	exporter.RegisterFormatter(".csv", &export.CSVFormatter{})
	exporter.RegisterFormatter(".json", &export.JSONFormatter{})
	return exporter
}

func newCLI(
	addStudent cli.StudentAdder,
	addCriteria cli.CriteriaAdder,
	evaluate cli.StudentEvaluator,
	removeEvaluation cli.EvaluationRemover,
	report usecase.ReportReader,
	exporter cli.StudentsExporter,
) (*cli.CLI, error) {
	return cli.NewCLI(
		cli.WithAddStudent(addStudent),
		cli.WithAddCriteria(addCriteria),
		cli.WithEvaluate(evaluate),
		cli.WithRemoveEvaluation(removeEvaluation),
		cli.WithReport(report),
		cli.WithExporter(exporter),
	)
}

func BuildCLI() (*cli.CLI, error) {
	c, err := New()
	if err != nil {
		return nil, err
	}

	var terminalCLI *cli.CLI
	if err := c.Invoke(func(resolvedCLI *cli.CLI) {
		terminalCLI = resolvedCLI
	}); err != nil {
		return nil, fmt.Errorf("build CLI: %w", err)
	}

	return terminalCLI, nil
}
