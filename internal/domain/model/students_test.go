package model

import "testing"

func TestNewStudentName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid", input: "Иван Иванов", want: "Иван Иванов"},
		{name: "trimmed", input: "  Иван Иванов  ", want: "Иван Иванов"},
		{name: "empty", input: "", wantErr: true},
		{name: "spaces", input: "   ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStudentName(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Value() != tt.want {
				t.Fatalf("got %q, want %q", got.Value(), tt.want)
			}
		})
	}
}

func TestNewStudentID(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{name: "zero", input: 0},
		{name: "positive", input: 10},
		{name: "negative", input: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStudentID(tt.input)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestNewStudent(t *testing.T) {
	student := mustStudent(t, 1, "Иван Иванов")

	if student.ID().Value() != 1 {
		t.Fatalf("got id %d", student.ID().Value())
	}
	if student.Name().Value() != "Иван Иванов" {
		t.Fatalf("got name %q", student.Name().Value())
	}
}

func TestNewStudentRejectsZeroValueName(t *testing.T) {
	if _, err := NewStudent(StudentID{}, StudentName{}); err == nil {
		t.Fatal("expected error for empty student name")
	}
}

func TestStudentAddEvaluationRejectsZeroValue(t *testing.T) {
	student := mustStudent(t, 1, "Иван Иванов")

	if err := student.AddEvaluation(Evaluation{}); err == nil {
		t.Fatal("expected error for zero value evaluation")
	}
}

func TestStudentAddEvaluationRejectsZeroValueStudent(t *testing.T) {
	criteria := mustCriteria(t, "Математика", 10, 0.5)
	evaluation, err := NewEvaluation(criteria, 5)
	if err != nil {
		t.Fatalf("NewEvaluation: %v", err)
	}

	student := &Student{}
	if student.Valid() {
		t.Fatal("zero value student must be invalid")
	}
	if err := student.AddEvaluation(evaluation); err == nil {
		t.Fatal("expected error for zero value student")
	}
}

func TestStudentAddEvaluationAndGetScore(t *testing.T) {
	criteria := mustCriteria(t, "Математика", 100, 0.5)
	student := mustStudent(t, 1, "Иван Иванов")

	evaluation, err := NewEvaluation(criteria, 80)
	if err != nil {
		t.Fatalf("NewEvaluation: %v", err)
	}
	if err := student.AddEvaluation(evaluation); err != nil {
		t.Fatalf("AddEvaluation: %v", err)
	}

	score, exists := student.GetScore(criteria.Name())
	if !exists {
		t.Fatal("score not found")
	}
	if score.Value() != 80 {
		t.Fatalf("got score %d, want 80", score.Value())
	}
}

func TestStudentRemoveEvaluation(t *testing.T) {
	criteria := mustCriteria(t, "Математика", 100, 0.5)
	student := mustStudent(t, 1, "Иван Иванов")

	evaluation, err := NewEvaluation(criteria, 80)
	if err != nil {
		t.Fatalf("NewEvaluation: %v", err)
	}
	if err := student.AddEvaluation(evaluation); err != nil {
		t.Fatalf("AddEvaluation: %v", err)
	}

	if err := student.RemoveEvaluation(criteria.Name()); err != nil {
		t.Fatalf("RemoveEvaluation: %v", err)
	}
	if _, exists := student.GetScore(criteria.Name()); exists {
		t.Fatal("score still exists")
	}
	if err := student.RemoveEvaluation(criteria.Name()); err == nil {
		t.Fatal("expected error when removing missing evaluation")
	}
}

func TestStudentGetScoresReturnsCopy(t *testing.T) {
	criteria := mustCriteria(t, "Математика", 100, 0.5)
	student := mustStudent(t, 1, "Иван Иванов")

	evaluation, err := NewEvaluation(criteria, 80)
	if err != nil {
		t.Fatalf("NewEvaluation: %v", err)
	}
	if err := student.AddEvaluation(evaluation); err != nil {
		t.Fatalf("AddEvaluation: %v", err)
	}

	snapshot := student.GetScores()
	delete(snapshot, criteria.Name())

	if _, exists := student.GetScore(criteria.Name()); !exists {
		t.Fatal("original student score was modified")
	}
}
