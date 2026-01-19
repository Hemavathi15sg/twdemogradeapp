package repository

import (
	"grademanagement-demo/models"
	"sync"
	"time"
)

// GradeRepository manages grade data with thread-safe in-memory storage
type GradeRepository struct {
	grades map[int]models.Grade
	mutex  sync.RWMutex
	nextID int
}

// NewGradeRepository creates a new grade repository
func NewGradeRepository() *GradeRepository {
	return &GradeRepository{
		grades: make(map[int]models.Grade),
		nextID: 1,
	}
}

// Create adds a new grade and returns it with auto-generated ID
func (r *GradeRepository) Create(grade models.Grade) models.Grade {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	grade.ID = r.nextID
	r.nextID++
	grade.CreatedAt = time.Now()
	grade.UpdatedAt = time.Now()
	r.grades[grade.ID] = grade

	return grade
}

// Get retrieves a grade by ID
func (r *GradeRepository) Get(id int) (models.Grade, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	grade, exists := r.grades[id]
	return grade, exists
}

// GetAll retrieves all grades
func (r *GradeRepository) GetAll() []models.Grade {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	grades := make([]models.Grade, 0, len(r.grades))
	for _, grade := range r.grades {
		grades = append(grades, grade)
	}

	return grades
}

// Update modifies an existing grade
func (r *GradeRepository) Update(id int, grade models.Grade) (models.Grade, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.grades[id]; !exists {
		return models.Grade{}, false
	}

	grade.ID = id
	grade.UpdatedAt = time.Now()
	// Preserve CreatedAt from original
	if original, ok := r.grades[id]; ok {
		grade.CreatedAt = original.CreatedAt
	}
	r.grades[id] = grade

	return grade, true
}

// Delete removes a grade by ID
func (r *GradeRepository) Delete(id int) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.grades[id]; !exists {
		return false
	}

	delete(r.grades, id)
	return true
}
