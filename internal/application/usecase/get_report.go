package usecase

import (
	"tas-app/internal/application/ports"
	"tas-app/internal/domain/service"
)

type CriteriaDTO struct {
	Name     string
	MaxScore int
	Weight   float64
}

type StudentDTO struct {
	ID       int
	Name     string
	Scores   map[string]int
	Total    float64
	MaxTotal float64
}

type ReportReader interface {
	GetStudents() ([]StudentDTO, error)
	GetCriteria() ([]CriteriaDTO, error)
}

type GetReportUseCase struct {
	studentRepo  ports.StudentProvider
	criteriaRepo ports.CriteriaProvider
	gradingSvc   *service.GradingService
}

var _ ReportReader = (*GetReportUseCase)(nil)

func NewGetReportUseCase(sr ports.StudentProvider, cr ports.CriteriaProvider, gs *service.GradingService) *GetReportUseCase {
	return &GetReportUseCase{
		studentRepo:  sr,
		criteriaRepo: cr,
		gradingSvc:   gs,
	}
}

func (uc *GetReportUseCase) GetCriteria() ([]CriteriaDTO, error) {
	domainCriteria, err := uc.criteriaRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var dtos []CriteriaDTO
	for _, c := range domainCriteria {
		dtos = append(dtos, CriteriaDTO{
			Name:     c.Name().Value(),
			MaxScore: c.MaxScore().Value(),
			Weight:   c.Weight().Value(),
		})
	}

	return dtos, nil
}

func (uc *GetReportUseCase) GetStudents() ([]StudentDTO, error) {
	students, err := uc.studentRepo.GetAll()
	if err != nil {
		return nil, err
	}

	criteria, err := uc.criteriaRepo.GetAll()
	if err != nil {
		return nil, err
	}

	maxTotal := uc.gradingSvc.CalculateMaxTotal(criteria)

	var dtos []StudentDTO
	for _, s := range students {
		total := uc.gradingSvc.CalculateFinalGrade(s, criteria)

		scoresMap := make(map[string]int)

		for _, c := range criteria {
			if score, exists := s.GetScore(c.Name()); exists {
				scoresMap[c.Name().Value()] = score.Value()
			}
		}

		dtos = append(dtos, StudentDTO{
			ID:       s.ID().Value(),
			Name:     s.Name().Value(),
			Scores:   scoresMap,
			Total:    total,
			MaxTotal: maxTotal,
		})
	}

	return dtos, nil
}
