package service

import (
	"tas-app/internal/domain/model"
)

type GradingService struct{}

func (s *GradingService) CalculateFinalGrade(student *model.Student, criteriaList []*model.Criteria) float64 {
	var finalGrade float64

	for _, criteria := range criteriaList {
		if score, exists := student.GetScore(criteria.Name()); exists {
			finalGrade += float64(score.Value()) * criteria.Weight().Value()
		}
	}
	return finalGrade
}

func (s *GradingService) CalculateMaxTotal(criteriaList []*model.Criteria) float64 {
	var maxTotal float64
	for _, c := range criteriaList {
		maxTotal += float64(c.MaxScore().Value()) * c.Weight().Value()
	}
	return maxTotal
}
