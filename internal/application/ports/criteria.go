package ports

import "tas-app/internal/domain/model"

type CriteriaSaver interface {
	Save(*model.Criteria) error
	ExistsByName(model.CriteriaName) (bool, error)
}

type CriteriaReader interface {
	FindByName(model.CriteriaName) (*model.Criteria, error)
}

type CriteriaProvider interface {
	GetAll() ([]*model.Criteria, error)
}
