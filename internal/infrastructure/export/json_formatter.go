package export

import (
	"encoding/json"
	"io"

	"tas-app/internal/application/usecase"
)

type JSONFormatter struct{}

func (f *JSONFormatter) Write(file io.Writer, students []usecase.StudentDTO, _ []usecase.CriteriaDTO) error {
	// JSON сериализует структуру DTO "как есть"
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Красивое форматирование с отступами
	return encoder.Encode(students)
}
