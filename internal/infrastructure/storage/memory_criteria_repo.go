package storage

import (
	"fmt"
	"sort"

	"tas-app/internal/domain/model"
)

type MemoryCriteriaRepo struct {
	data map[model.CriteriaName]*model.Criteria
}

func NewMemoryCriteriaRepo() *MemoryCriteriaRepo {
	return &MemoryCriteriaRepo{
		data: make(map[model.CriteriaName]*model.Criteria),
	}
}

func (r *MemoryCriteriaRepo) Save(criteria *model.Criteria) error {
	r.data[criteria.Name()] = criteria
	return nil
}

func (r *MemoryCriteriaRepo) ExistsByName(name model.CriteriaName) (bool, error) {
	_, exists := r.data[name]
	return exists, nil
}

func (r *MemoryCriteriaRepo) FindByName(name model.CriteriaName) (*model.Criteria, error) {
	criteria, exists := r.data[name]
	if !exists {
		return nil, fmt.Errorf("criteria %s is not in list", name.Value())
	}
	return criteria, nil
}

func (r *MemoryCriteriaRepo) GetAll() ([]*model.Criteria, error) {
	criteria := make([]*model.Criteria, 0, len(r.data))
	for _, v := range r.data {
		criteria = append(criteria, v)
	}

	sort.Slice(criteria, func(i, j int) bool {
		return criteria[i].Name().Value() < criteria[j].Name().Value()
	})

	return criteria, nil
}
