package repository

import (
	"context"
	"fmt"
	"grademanagement-demo/cache"
	"grademanagement-demo/models"
	"sync"
	"time"
)

// EnrollmentRepository manages enrollment data with thread-safe in-memory storage
type EnrollmentRepository struct {
	enrollments  map[int]models.Enrollment
	mutex        sync.RWMutex
	nextID       int
	cacheManager *cache.Manager
}

// NewEnrollmentRepository creates a new enrollment repository
func NewEnrollmentRepository() *EnrollmentRepository {
	return &EnrollmentRepository{
		enrollments:  make(map[int]models.Enrollment),
		nextID:       1,
		cacheManager: nil,
	}
}

// SetCacheManager sets the cache manager for this repository
func (r *EnrollmentRepository) SetCacheManager(cm *cache.Manager) {
	r.cacheManager = cm
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

	// Update cache
	if r.cacheManager != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("enrollment:%d", enrollment.ID)
		r.cacheManager.Set(ctx, cacheKey, enrollment)
		r.cacheManager.Del(ctx, "enrollments:all")
	}

	return enrollment
}

// Get retrieves an enrollment by ID
func (r *EnrollmentRepository) Get(id int) (models.Enrollment, bool) {
	// Try cache first if available
	if r.cacheManager != nil {
		var cached models.Enrollment
		cacheKey := fmt.Sprintf("enrollment:%d", id)
		if found, err := r.cacheManager.Get(context.Background(), cacheKey, &cached); err == nil && found {
			return cached, true
		}
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	enrollment, exists := r.enrollments[id]

	// Store in cache if found
	if exists && r.cacheManager != nil {
		cacheKey := fmt.Sprintf("enrollment:%d", id)
		r.cacheManager.Set(context.Background(), cacheKey, enrollment)
	}

	return enrollment, exists
}

// GetAll retrieves all enrollments
func (r *EnrollmentRepository) GetAll() []models.Enrollment {
	// Try cache first if available
	if r.cacheManager != nil {
		var cached []models.Enrollment
		if found, err := r.cacheManager.Get(context.Background(), "enrollments:all", &cached); err == nil && found {
			return cached
		}
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	enrollments := make([]models.Enrollment, 0, len(r.enrollments))
	for _, enrollment := range r.enrollments {
		enrollments = append(enrollments, enrollment)
	}

	// Cache the list
	if r.cacheManager != nil {
		r.cacheManager.Set(context.Background(), "enrollments:all", enrollments)
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

	// Update cache
	if r.cacheManager != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("enrollment:%d", id)
		r.cacheManager.Set(ctx, cacheKey, enrollment)
		r.cacheManager.Del(ctx, "enrollments:all")
	}

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

	// Invalidate cache
	if r.cacheManager != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("enrollment:%d", id)
		r.cacheManager.Del(ctx, cacheKey, "enrollments:all")
	}

	return true
}
