package handlers

import (
	"encoding/json"
	"fmt"
	"grademanagement-demo/models"
	"grademanagement-demo/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// GradeHandler handles grade-related HTTP requests
type GradeHandler struct {
	repo *repository.GradeRepository
}

// NewGradeHandler creates a new grade handler
func NewGradeHandler(repo *repository.GradeRepository) *GradeHandler {
	return &GradeHandler{repo: repo}
}

// validateGradeStatus checks if the grade status is valid
func validateGradeStatus(status string) error {
	if status != models.GradeStatusDraft &&
		status != models.GradeStatusSubmitted &&
		status != models.GradeStatusFinalized {
		return fmt.Errorf("status must be one of: draft, submitted, finalized")
	}
	return nil
}

// validateGrade validates required fields
func validateGrade(grade *models.Grade) error {
	if grade.StudentID == 0 {
		return fmt.Errorf("student_id is required")
	}
	if grade.CourseID == 0 {
		return fmt.Errorf("course_id is required")
	}
	if grade.GradeValue == "" {
		return fmt.Errorf("grade_value is required")
	}
	if grade.Status == "" {
		return fmt.Errorf("status is required")
	}
	return validateGradeStatus(grade.Status)
}

// CreateGrade handles POST /api/grades
func (h *GradeHandler) CreateGrade(w http.ResponseWriter, r *http.Request) {
	var grade models.Grade

	if err := json.NewDecoder(r.Body).Decode(&grade); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	// Set grade date if not provided
	if grade.GradeDate.IsZero() {
		grade.GradeDate = time.Now()
	}

	if err := validateGrade(&grade); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	created := h.repo.Create(grade)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetGrade handles GET /api/grades/{id}
func (h *GradeHandler) GetGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid grade id"})
		return
	}

	grade, exists := h.repo.Get(id)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "grade not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(grade)
}

// ListGrades handles GET /api/grades
func (h *GradeHandler) ListGrades(w http.ResponseWriter, r *http.Request) {
	grades := h.repo.GetAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(grades)
}

// UpdateGrade handles PUT /api/grades/{id}
func (h *GradeHandler) UpdateGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid grade id"})
		return
	}

	var grade models.Grade
	if err := json.NewDecoder(r.Body).Decode(&grade); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if err := validateGrade(&grade); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	updated, exists := h.repo.Update(id, grade)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "grade not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

// DeleteGrade handles DELETE /api/grades/{id}
func (h *GradeHandler) DeleteGrade(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid grade id"})
		return
	}

	if !h.repo.Delete(id) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "grade not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
