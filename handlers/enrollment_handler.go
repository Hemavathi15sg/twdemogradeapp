package handlers

import (
	"encoding/json"
	"fmt"
	"grademanagement-demo/middleware"
	"grademanagement-demo/models"
	"grademanagement-demo/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// EnrollmentHandler handles enrollment-related HTTP requests
type EnrollmentHandler struct {
	repo *repository.EnrollmentRepository
}

// NewEnrollmentHandler creates a new enrollment handler
func NewEnrollmentHandler(repo *repository.EnrollmentRepository) *EnrollmentHandler {
	return &EnrollmentHandler{repo: repo}
}

// validateStatus checks if the status is valid
func validateStatus(status string) error {
	if status != models.StatusPending &&
		status != models.StatusActive &&
		status != models.StatusCompleted {
		return fmt.Errorf("status must be one of: pending, active, completed")
	}
	return nil
}

// validateEnrollment validates required fields
func validateEnrollment(enrollment *models.Enrollment) error {
	if enrollment.StudentID == 0 {
		return fmt.Errorf("student_id is required")
	}
	if enrollment.CourseID == 0 {
		return fmt.Errorf("course_id is required")
	}
	if enrollment.Status == "" {
		return fmt.Errorf("status is required")
	}
	return validateStatus(enrollment.Status)
}

// CreateEnrollment handles POST /api/enrollments
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var enrollment models.Enrollment

	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	// Set enrollment date if not provided
	if enrollment.EnrollmentDate.IsZero() {
		enrollment.EnrollmentDate = time.Now()
	}

	if err := validateEnrollment(&enrollment); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	created := h.repo.Create(enrollment)

	w.Header().Set("Content-Type", "application/json")
	middleware.SetCacheHeader(w, middleware.CacheUpdate, 5*time.Minute)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetEnrollment handles GET /api/enrollments/{id}
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid enrollment id"})
		return
	}

	enrollment, exists := h.repo.Get(id)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "enrollment not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	middleware.SetCacheHeader(w, middleware.CacheHit, 5*time.Minute)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(enrollment)
}

// ListEnrollments handles GET /api/enrollments
func (h *EnrollmentHandler) ListEnrollments(w http.ResponseWriter, r *http.Request) {
	enrollments := h.repo.GetAll()

	middleware.SetCacheHeader(w, middleware.CacheHit, 5*time.Minute)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(enrollments)
}

// UpdateEnrollment handles PUT /api/enrollments/{id}
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid enrollment id"})
		return
	}

	var enrollment models.Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if err := validateEnrollment(&enrollment); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	updated, exists := h.repo.Update(id, enrollment)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "enrollment not found"})
		return
	}

	middleware.SetCacheHeader(w, middleware.CacheUpdate, 5*time.Minute)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

// DeleteEnrollment handles DELETE /api/enrollments/{id}
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid enrollment id"})
		return
	}

	if !h.repo.Delete(id) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "enrollment not found"})
		return
	}
	middleware.SetCacheHeader(w, middleware.CacheDelete)

	w.WriteHeader(http.StatusNoContent)
}
