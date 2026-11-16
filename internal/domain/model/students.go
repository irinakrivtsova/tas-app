package model

import (
	"fmt"
	"strings"
)

// Value Object
type StudentID struct {
	value int
}

func NewStudentID(value int) (StudentID, error) {
	if value < 0 {
		return StudentID{}, fmt.Errorf("invalid student id: %d", value)
	}
	return StudentID{value: value}, nil
}

func (id StudentID) Value() int {
	return id.value
}

type StudentName struct {
	value string
}

func NewStudentName(value string) (StudentName, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return StudentName{}, fmt.Errorf("student name can not be empty")
	}
	return StudentName{value: value}, nil
}

func (n StudentName) Value() string {
	return n.value
}

func (n StudentName) Valid() bool {
	return n.value != ""
}

// Value Object
type Score struct {
	value int
}

func NewScore(value int, max int) (Score, error) {
	if value < 0 {
		return Score{}, fmt.Errorf("score can not be negative, got: %d", value)
	}
	if value > max {
		return Score{}, fmt.Errorf("score %d exceeds max %d", value, max)
	}
	return Score{value: value}, nil
}

func (s Score) Value() int {
	return s.value
}

// Entity
type Student struct {
	id     StudentID
	name   StudentName
	scores map[CriteriaName]Evaluation
}

func NewStudent(id StudentID, name StudentName) (*Student, error) {
	if !name.Valid() {
		return nil, fmt.Errorf("student name can not be empty")
	}
	return &Student{
		id:     id,
		name:   name,
		scores: make(map[CriteriaName]Evaluation),
	}, nil
}

func (s *Student) GetScore(criteria CriteriaName) (Score, bool) {
	evaluation, exists := s.scores[criteria]
	if !exists {
		return Score{}, false
	}
	return evaluation.Score(), true
}

func (s *Student) GetScores() map[CriteriaName]Score {
	result := make(map[CriteriaName]Score, len(s.scores))
	for k, evaluation := range s.scores {
		result[k] = evaluation.Score()
	}
	return result
}

func (s *Student) ID() StudentID {
	return s.id
}

func (s *Student) Name() StudentName {
	return s.name
}

func (s *Student) Valid() bool {
	return s != nil && s.name.Valid()
}

func (s *Student) AddEvaluation(e Evaluation) error {
	if !s.Valid() {
		return fmt.Errorf("student is nil or invalid")
	}
	if !e.Valid() {
		return fmt.Errorf("evaluation is invalid")
	}
	if s.scores == nil {
		s.scores = make(map[CriteriaName]Evaluation)
	}
	s.scores[e.Criteria()] = e
	return nil
}

func (s *Student) RemoveEvaluation(criteria CriteriaName) error {
	if !s.Valid() {
		return fmt.Errorf("student is nil or invalid")
	}
	if _, exists := s.scores[criteria]; !exists {
		return fmt.Errorf("evaluation for criteria %q not found", criteria.Value())
	}
	delete(s.scores, criteria)
	return nil
}
