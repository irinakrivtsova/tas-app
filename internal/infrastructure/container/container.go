package container

import (
	"fmt"

	"go.uber.org/dig"

	"tas-app/internal/application/ports"
	"tas-app/internal/application/usecase"
	"tas-app/internal/domain/service"
	"tas-app/internal/infrastructure/cli"
	"tas-app/internal/infrastructure/storage"
)

type provider struct {
	name        string
	constructor any
	options     []dig.ProvideOption
}

func New() (*dig.Container, error) {
	c := dig.New()

	providers := []provider{
		{
			name:        "student repository",
			constructor: storage.NewMemoryStudentRepo,
			options: []dig.ProvideOption{
				dig.As(
					new(ports.StudentSaver),
					new(ports.StudentReadSaver),
					new(ports.StudentProvider),
				),
			},
		},
		{
			name:        "criteria repository",
			constructor: storage.NewMemoryCriteriaRepo,
			options: []dig.ProvideOption{
				dig.As(
					new(ports.CriteriaSaver),
					new(ports.CriteriaReader),
					new(ports.CriteriaProvider),
				),
			},
		},
		{
			name:        "grading service",
			constructor: func() *service.GradingService { return &service.GradingService{} },
		},
		{
			name:        "add student use case",
			constructor: usecase.NewAddStudentUseCase,
			options:     []dig.ProvideOption{dig.As(new(cli.StudentAdder))},
		},
		{
			name:        "add criteria use case",
			constructor: usecase.NewAddCriteriaUseCase,
			options:     []dig.ProvideOption{dig.As(new(cli.CriteriaAdder))},
		},
		{
			name:        "evaluate student use case",
			constructor: usecase.NewEvaluateStudentUseCase,
			options:     []dig.ProvideOption{dig.As(new(cli.StudentEvaluator))},
		},
		{
			name:        "remove evaluation use case",
			constructor: usecase.NewRemoveEvaluationUseCase,
			options:     []dig.ProvideOption{dig.As(new(cli.EvaluationRemover))},
		},
		{
			name:        "get report use case",
			constructor: usecase.NewGetReportUseCase,
			options:     []dig.ProvideOption{dig.As(new(usecase.ReportReader))},
		},
		{
			name:        "file exporter",
			constructor: newFileExporter,
			options:     []dig.ProvideOption{dig.As(new(cli.StudentsExporter))},
		},
		{
			name:        "CLI",
			constructor: newCLI,
		},
	}

	for _, p := range providers {
		if err := c.Provide(p.constructor, p.options...); err != nil {
			return nil, fmt.Errorf("provide %s: %w", p.name, err)
		}
	}

	return c, nil
}
