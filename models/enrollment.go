package models

import "time"

// Enrollment represents a student enrollment in a course
type Enrollment struct {
	ID             int       `json:"id"`
	StudentID      int       `json:"student_id"`
	CourseID       int       `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Status constants
const (
	StatusPending   = "pending"
	StatusActive    = "active"
	StatusCompleted = "completed"
)
