package repository

import (
	"context"
	"fmt"
	"grademanagement-demo/cache"
	"grademanagement-demo/models"
	"sync"
	"time"
)

// GradeRepository manages grade data with thread-safe in-memory storage
type GradeRepository struct {
	grades       map[int]models.Grade
	mutex        sync.RWMutex
	nextID       int
	cacheManager *cache.Manager
}

// NewGradeRepository creates a new grade repository
func NewGradeRepository() *GradeRepository {
	return &GradeRepository{
		grades:       make(map[int]models.Grade),
		nextID:       1,
		cacheManager: nil,
	}
}

// SetCacheManager sets the cache manager for this repository
func (r *GradeRepository) SetCacheManager(cm *cache.Manager) {
	r.cacheManager = cm
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

	// Update cache
	if r.cacheManager != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("grade:%d", grade.ID)
		r.cacheManager.Set(ctx, cacheKey, grade)
		r.cacheManager.Del(ctx, "grades:all")
	}

	return grade
}

// Get retrieves a grade by ID
func (r *GradeRepository) Get(id int) (models.Grade, bool) {
	// Try cache first if available
	if r.cacheManager != nil {
		var cached models.Grade
		cacheKey := fmt.Sprintf("grade:%d", id)
		if found, err := r.cacheManager.Get(context.Background(), cacheKey, &cached); err == nil && found {
			return cached, true
		}
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	grade, exists := r.grades[id]

	// Store in cache if found
	if exists && r.cacheManager != nil {
		cacheKey := fmt.Sprintf("grade:%d", id)
		r.cacheManager.Set(context.Background(), cacheKey, grade)
	}

	return grade, exists
}

// GetAll retrieves all grades
func (r *GradeRepository) GetAll() []models.Grade {
	// Try cache first if available
	if r.cacheManager != nil {
		var cached []models.Grade
		if found, err := r.cacheManager.Get(context.Background(), "grades:all", &cached); err == nil && found {
			return cached
		}
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	grades := make([]models.Grade, 0, len(r.grades))
	for _, grade := range r.grades {
		grades = append(grades, grade)
	}

	// Cache the list
	if r.cacheManager != nil {
		r.cacheManager.Set(context.Background(), "grades:all", grades)
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

	// Update cache
	if r.cacheManager != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("grade:%d", id)
		r.cacheManager.Set(ctx, cacheKey, grade)
		r.cacheManager.Del(ctx, "grades:all")
	}

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

	// Invalidate cache
	if r.cacheManager != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("grade:%d", id)
		r.cacheManager.Del(ctx, cacheKey, "grades:all")
	}

	return true
}
