package usecase

import (
	"fmt"

	"tas-app/internal/application/ports"
	"tas-app/internal/domain/model"
)

type StudentEvaluationCommand struct {
	StudentID    int
	CriteriaName string
	Score        int
}

type EvaluateStudentUseCase struct {
	studentRepo  ports.StudentReadSaver
	criteriaRepo ports.CriteriaReader
}

func NewEvaluateStudentUseCase(students ports.StudentReadSaver, criteriaReader ports.CriteriaReader) *EvaluateStudentUseCase {
	return &EvaluateStudentUseCase{studentRepo: students, criteriaRepo: criteriaReader}
}

func (c *EvaluateStudentUseCase) Execute(cmd StudentEvaluationCommand) error {
	studentID, err := model.NewStudentID(cmd.StudentID)
	if err != nil {
		return err
	}

	student, err := c.studentRepo.FindByID(studentID)
	if err != nil {
		return err
	}

	criteriaName, err := model.NewCriteriaName(cmd.CriteriaName)
	if err != nil {
		return err
	}

	criteria, err := c.criteriaRepo.FindByName(criteriaName)
	if err != nil {
		return err
	}

	evaluation, err := model.NewEvaluation(criteria, cmd.Score)
	if err != nil {
		return err
	}

	if err := student.AddEvaluation(evaluation); err != nil {
		return err
	}

	if err := c.studentRepo.Save(student); err != nil {
		return fmt.Errorf("can not save new score: %w", err)
	}

	return nil
}
