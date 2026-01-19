package repository

import (
	"grademanagement-demo/models"
	"sync"
	"time"
)

// EnrollmentRepository manages enrollment data with thread-safe in-memory storage
type EnrollmentRepository struct {
	enrollments map[int]models.Enrollment
	mutex       sync.RWMutex
	nextID      int
}

// NewEnrollmentRepository creates a new enrollment repository
func NewEnrollmentRepository() *EnrollmentRepository {
	return &EnrollmentRepository{
		enrollments: make(map[int]models.Enrollment),
		nextID:      1,
	}
}

// Create adds a new enrollment and returns it with auto-generated ID
func (r *EnrollmentRepository) Create(enrollment models.Enrollment) models.Enrollment {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	enrollment.ID = r.nextID
	r.nextID++
	enrollment.CreatedAt = time.Now()
	enrollment.UpdatedAt = time.Now()
	r.enrollments[enrollment.ID] = enrollment

	return enrollment
}

// Get retrieves an enrollment by ID
func (r *EnrollmentRepository) Get(id int) (models.Enrollment, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	enrollment, exists := r.enrollments[id]
	return enrollment, exists
}

// GetAll retrieves all enrollments
func (r *EnrollmentRepository) GetAll() []models.Enrollment {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	enrollments := make([]models.Enrollment, 0, len(r.enrollments))
	for _, enrollment := range r.enrollments {
		enrollments = append(enrollments, enrollment)
	}

	return enrollments
}

// Update modifies an existing enrollment
func (r *EnrollmentRepository) Update(id int, enrollment models.Enrollment) (models.Enrollment, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.enrollments[id]; !exists {
		return models.Enrollment{}, false
	}

	enrollment.ID = id
	enrollment.UpdatedAt = time.Now()
	// Preserve CreatedAt from original
	if original, ok := r.enrollments[id]; ok {
		enrollment.CreatedAt = original.CreatedAt
	}
	r.enrollments[id] = enrollment

	return enrollment, true
}

// Delete removes an enrollment by ID
func (r *EnrollmentRepository) Delete(id int) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.enrollments[id]; !exists {
		return false
	}

	delete(r.enrollments, id)
	return true
}
