package usecase

import (
	"tas-app/internal/application/ports"
	"tas-app/internal/domain/model"
)

type RemoveEvaluationCommand struct {
	StudentID    int
	CriteriaName string
}

type RemoveEvaluationUseCase struct {
	students ports.StudentReadSaver
	criteria ports.CriteriaReader
}

func NewRemoveEvaluationUseCase(
	students ports.StudentReadSaver,
	criteria ports.CriteriaReader,
) *RemoveEvaluationUseCase {
	return &RemoveEvaluationUseCase{
		students: students,
		criteria: criteria,
	}
}

func (uc *RemoveEvaluationUseCase) Execute(cmd RemoveEvaluationCommand) error {
	studentID, err := model.NewStudentID(cmd.StudentID)
	if err != nil {
		return err
	}

	criteriaName, err := model.NewCriteriaName(cmd.CriteriaName)
	if err != nil {
		return err
	}

	student, err := uc.students.FindByID(studentID)
	if err != nil {
		return err
	}

	if _, err := uc.criteria.FindByName(criteriaName); err != nil {
		return err
	}

	if err := student.RemoveEvaluation(criteriaName); err != nil {
		return err
	}

	return uc.students.Save(student)
}
