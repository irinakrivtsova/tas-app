package model

import (
	"fmt"
	"strings"
)

// Value Object
type Weight struct {
	value float64
}

func NewWeight(value float64) (Weight, error) {
	if !(0 <= value && value <= 1) {
		return Weight{}, fmt.Errorf("weight must be between 0 and 1, got: %f", value)
	}
	return Weight{value: value}, nil
}

func (w Weight) Value() float64 {
	return w.value
}

func (w Weight) Valid() bool {
	return 0 <= w.value && w.value <= 1
}

type MaxScore struct {
	value int
}

func NewMaxScore(value int) (MaxScore, error) {
	if value <= 0 {
		return MaxScore{}, fmt.Errorf("max score must be positive, got: %d", value)
	}
	return MaxScore{value: value}, nil
}

func (m MaxScore) Value() int {
	return m.value
}

func (m MaxScore) Valid() bool {
	return m.value > 0
}

type CriteriaName struct {
	value string
}

func NewCriteriaName(value string) (CriteriaName, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return CriteriaName{}, fmt.Errorf("criteria name can not be empty")
	}
	return CriteriaName{value: value}, nil
}

func (n CriteriaName) Value() string {
	return n.value
}

func (n CriteriaName) Valid() bool {
	return n.value != ""
}

// Entity
type Criteria struct {
	name     CriteriaName
	maxScore MaxScore
	weight   Weight
}

func NewCriteria(name CriteriaName, maxScore MaxScore, weight Weight) (*Criteria, error) {
	if !name.Valid() {
		return nil, fmt.Errorf("criteria name can not be empty")
	}
	if !maxScore.Valid() {
		return nil, fmt.Errorf("max score must be positive")
	}
	if !weight.Valid() {
		return nil, fmt.Errorf("weight must be between 0 and 1, got: %f", weight.Value())
	}
	return &Criteria{name: name, maxScore: maxScore, weight: weight}, nil
}

func (c *Criteria) Valid() bool {
	return c != nil && c.name.Valid() && c.maxScore.Valid() && c.weight.Valid()
}

func (c *Criteria) Name() CriteriaName {
	return c.name
}

func (c *Criteria) MaxScore() MaxScore {
	return c.maxScore
}

func (c *Criteria) Weight() Weight {
	return c.weight
}
