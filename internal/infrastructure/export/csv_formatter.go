package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"tas-app/internal/application/usecase"
)

type CSVFormatter struct{}

func (f *CSVFormatter) Write(file io.Writer, students []usecase.StudentDTO, criteria []usecase.CriteriaDTO) error {
	writer := csv.NewWriter(file)

	// 1. Пишем заголовки
	headers := []string{"ID", "Имя"}
	for _, c := range criteria {
		headers = append(headers, c.Name)
	}
	headers = append(headers, "Итог", "Макс. Итог")

	if err := writer.Write(headers); err != nil {
		return err
	}

	// 2. Пишем строки
	for _, s := range students {
		row := []string{strconv.Itoa(s.ID), s.Name}
		for _, c := range criteria {
			if score, ok := s.Scores[c.Name]; ok {
				row = append(row, strconv.Itoa(score))
			} else {
				row = append(row, "-")
			}
		}
		row = append(row, fmt.Sprintf("%.2f", s.Total), fmt.Sprintf("%.2f", s.MaxTotal))

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return nil
}
