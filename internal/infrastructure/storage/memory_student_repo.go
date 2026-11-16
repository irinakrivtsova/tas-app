package storage

import (
	"fmt"
	"sort"

	"tas-app/internal/domain/model"
)

type MemoryStudentRepo struct {
	data map[model.StudentID]*model.Student
}

func NewMemoryStudentRepo() *MemoryStudentRepo {
	return &MemoryStudentRepo{
		data: make(map[model.StudentID]*model.Student),
	}
}

func (r *MemoryStudentRepo) Save(student *model.Student) error {
	r.data[student.ID()] = student
	return nil
}

func (r *MemoryStudentRepo) ExistsByID(id model.StudentID) (bool, error) {
	_, exists := r.data[id]
	return exists, nil
}

func (r *MemoryStudentRepo) GetAll() ([]*model.Student, error) {
	students := make([]*model.Student, 0, len(r.data))

	for _, v := range r.data {
		students = append(students, v)
	}

	sort.Slice(students, func(i, j int) bool {
		return students[i].ID().Value() < students[j].ID().Value()
	})

	return students, nil
}

func (r *MemoryStudentRepo) FindByID(id model.StudentID) (*model.Student, error) {
	student, exists := r.data[id]
	if !exists {
		return nil, fmt.Errorf("student with id %v is not in list", id)
	}

	return student, nil
}
