package model

import "fmt"

type Evaluation struct {
	criteria CriteriaName
	score    Score
}

func NewEvaluation(criteria *Criteria, value int) (Evaluation, error) {
	if !criteria.Valid() {
		return Evaluation{}, fmt.Errorf("criteria is nil or invalid")
	}

	score, err := NewScore(value, criteria.MaxScore().Value())
	if err != nil {
		return Evaluation{}, err
	}

	return Evaluation{
		criteria: criteria.Name(),
		score:    score,
	}, nil
}

func (e Evaluation) Criteria() CriteriaName {
	return e.criteria
}

func (e Evaluation) Score() Score {
	return e.score
}

func (e Evaluation) Valid() bool {
	return e.criteria.Valid()
}
