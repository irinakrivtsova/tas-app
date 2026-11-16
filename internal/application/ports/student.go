package ports

import "tas-app/internal/domain/model"

type StudentSaver interface {
	Save(*model.Student) error
	ExistsByID(model.StudentID) (bool, error)
}

type StudentReadSaver interface {
	FindByID(model.StudentID) (*model.Student, error)
	Save(*model.Student) error
}

type StudentProvider interface {
	GetAll() ([]*model.Student, error)
}
