package usecase

import (
	"fmt"

	"tas-app/internal/application/ports"
	"tas-app/internal/domain/model"
)

type AddCriteriaCommand struct {
	Name        string
	MaxScoreInt int
	WeightFloat float64
}

type AddCriteriaUseCase struct {
	criteriaRepo ports.CriteriaSaver
}

func NewAddCriteriaUseCase(repo ports.CriteriaSaver) *AddCriteriaUseCase {
	return &AddCriteriaUseCase{criteriaRepo: repo}
}

func (c *AddCriteriaUseCase) Execute(cmd AddCriteriaCommand) error {
	criteriaName, err := model.NewCriteriaName(cmd.Name)
	if err != nil {
		return err
	}

	exists, err := c.criteriaRepo.ExistsByName(criteriaName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("criteria %s already in list", criteriaName.Value())
	}

	score, err := model.NewMaxScore(cmd.MaxScoreInt)
	if err != nil {
		return err
	}

	weight, err := model.NewWeight(cmd.WeightFloat)
	if err != nil {
		return err
	}

	newCrit, err := model.NewCriteria(criteriaName, score, weight)
	if err != nil {
		return err
	}
	return c.criteriaRepo.Save(newCrit)
}
