package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"tas-app/internal/application/usecase"
)

type Option func(*config)

type config struct {
	reader    io.Reader
	writer    io.Writer
	errWriter io.Writer
	deps      Dependencies
}

func WithReader(r io.Reader) Option {
	return func(c *config) {
		c.reader = r
	}
}

func WithWriter(w io.Writer) Option {
	return func(c *config) {
		c.writer = w
	}
}

func WithErrorWriter(w io.Writer) Option {
	return func(c *config) {
		c.errWriter = w
	}
}

func WithAddStudent(a StudentAdder) Option {
	return func(c *config) {
		c.deps.AddStudent = a
	}
}

func WithAddCriteria(a CriteriaAdder) Option {
	return func(c *config) {
		c.deps.AddCriteria = a
	}
}

func WithEvaluate(e StudentEvaluator) Option {
	return func(c *config) {
		c.deps.Evaluate = e
	}
}

func WithReport(r usecase.ReportReader) Option {
	return func(c *config) {
		c.deps.GetReport = r
	}
}

func WithRemoveEvaluation(e EvaluationRemover) Option {
	return func(c *config) {
		c.deps.RemoveEvaluation = e
	}
}

func WithExporter(e StudentsExporter) Option {
	return func(c *config) {
		c.deps.ExportStudents = e
	}
}

func NewCLI(opts ...Option) (*CLI, error) {
	cfg := config{
		reader:    os.Stdin,
		writer:    os.Stdout,
		errWriter: os.Stderr,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.reader == nil {
		cfg.reader = os.Stdin
	}
	if cfg.writer == nil {
		cfg.writer = os.Stdout
	}
	if cfg.errWriter == nil {
		cfg.errWriter = os.Stderr
	}

	if err := cfg.deps.validate(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(cfg.reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	return &CLI{
		deps:      cfg.deps,
		reader:    scanner,
		writer:    cfg.writer,
		errWriter: cfg.errWriter,
	}, nil
}

func (d Dependencies) validate() error {
	var missing []string
	if d.AddStudent == nil {
		missing = append(missing, "AddStudent")
	}
	if d.AddCriteria == nil {
		missing = append(missing, "AddCriteria")
	}
	if d.Evaluate == nil {
		missing = append(missing, "Evaluate")
	}
	if d.RemoveEvaluation == nil {
		missing = append(missing, "RemoveEvaluation")
	}
	if d.GetReport == nil {
		missing = append(missing, "GetReport")
	}
	if d.ExportStudents == nil {
		missing = append(missing, "ExportStudents")
	}
	if len(missing) > 0 {
		return fmt.Errorf("cli: missing dependencies: %s", strings.Join(missing, ", "))
	}
	return nil
}
