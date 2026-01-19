package models

import "time"

// Grade represents a student's grade for a course
type Grade struct {
	ID          int       `json:"id"`
	StudentID   int       `json:"student_id"`
	CourseID    int       `json:"course_id"`
	GradeValue  string    `json:"grade_value"`
	GradePoints float64   `json:"grade_points"`
	GradeDate   time.Time `json:"grade_date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Grade status constants
const (
	GradeStatusDraft     = "draft"
	GradeStatusSubmitted = "submitted"
	GradeStatusFinalized = "finalized"
)
