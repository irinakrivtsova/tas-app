package export

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"tas-app/internal/application/usecase"
)

// Strategy
type Formatter interface {
	Write(file io.Writer, students []usecase.StudentDTO, criteria []usecase.CriteriaDTO) error
}

type FileExporter struct {
	reportReader usecase.ReportReader
	formatters   map[string]Formatter
}

func NewFileExporter(reportReader usecase.ReportReader) *FileExporter {
	return &FileExporter{
		reportReader: reportReader,
		formatters:   make(map[string]Formatter),
	}
}

func (e *FileExporter) RegisterFormatter(extension string, f Formatter) {
	e.formatters[extension] = f
}

// Template Method
func (e *FileExporter) ExportStudents(path string) error {
	ext := filepath.Ext(path)
	formatter, exists := e.formatters[ext]
	if !exists {
		return fmt.Errorf("неподдерживаемый формат файла: %s", ext)
	}

	students, err := e.reportReader.GetStudents()
	if err != nil {
		return err
	}
	criteria, err := e.reportReader.GetCriteria()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return fmt.Errorf("не удалось создать временный файл: %w", err)
	}
	tmpName := tmp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = tmp.Close()
			_ = os.Remove(tmpName)
		}
	}()

	if err := tmp.Chmod(0o644); err != nil {
		return fmt.Errorf("не удалось выставить права на файл: %w", err)
	}

	if err := formatter.Write(tmp, students, criteria); err != nil {
		return fmt.Errorf("ошибка записи данных: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		return fmt.Errorf("не удалось синхронизировать файл: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("не удалось закрыть файл: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("не удалось сохранить файл: %w", err)
	}

	cleanup = false
	return nil
}
