package usecase

import (
	"fmt"

	"tas-app/internal/application/ports"
	"tas-app/internal/domain/model"
)

type AddStudentCommand struct {
	ID   int
	Name string
}

type AddStudentUseCase struct {
	studentsRepo ports.StudentSaver
}

func NewAddStudentUseCase(repo ports.StudentSaver) *AddStudentUseCase {
	newUseCase := AddStudentUseCase{studentsRepo: repo}
	return &newUseCase
}

func (c *AddStudentUseCase) Execute(cmd AddStudentCommand) error {
	studentID, err := model.NewStudentID(cmd.ID)
	if err != nil {
		return err
	}

	name, err := model.NewStudentName(cmd.Name)
	if err != nil {
		return err
	}

	exists, err := c.studentsRepo.ExistsByID(studentID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("student %s with id %d already in list", name.Value(), cmd.ID)
	}

	student, err := model.NewStudent(studentID, name)
	if err != nil {
		return err
	}

	return c.studentsRepo.Save(student)
}
