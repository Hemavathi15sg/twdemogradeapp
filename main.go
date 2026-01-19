package main

import (
	"fmt"
	"grademanagement-demo/handlers"
	"grademanagement-demo/repository"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize repositories
	enrollmentRepo := repository.NewEnrollmentRepository()
	gradeRepo := repository.NewGradeRepository()

	// Initialize handlers
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo)
	gradeHandler := handlers.NewGradeHandler(gradeRepo)

	r := mux.NewRouter()

	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API - Ready for AI delegation!", "status": "healthy"}`)
	}).Methods("GET")

	// Enrollment routes
	r.HandleFunc("/api/enrollments", enrollmentHandler.CreateEnrollment).Methods("POST")
	r.HandleFunc("/api/enrollments/{id}", enrollmentHandler.GetEnrollment).Methods("GET")
	r.HandleFunc("/api/enrollments", enrollmentHandler.ListEnrollments).Methods("GET")
	r.HandleFunc("/api/enrollments/{id}", enrollmentHandler.UpdateEnrollment).Methods("PUT")
	r.HandleFunc("/api/enrollments/{id}", enrollmentHandler.DeleteEnrollment).Methods("DELETE")

	// Grade routes
	r.HandleFunc("/api/grades", gradeHandler.CreateGrade).Methods("POST")
	r.HandleFunc("/api/grades/{id}", gradeHandler.GetGrade).Methods("GET")
	r.HandleFunc("/api/grades", gradeHandler.ListGrades).Methods("GET")
	r.HandleFunc("/api/grades/{id}", gradeHandler.UpdateGrade).Methods("PUT")
	r.HandleFunc("/api/grades/{id}", gradeHandler.DeleteGrade).Methods("DELETE")

	port := ":8080"
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Ready for Copilot Agent delegation!")
	fmt.Println("📚 Enrollment API available at /api/enrollments")
	fmt.Println("📊 Grade API available at /api/grades")

	log.Fatal(http.ListenAndServe(port, r))
}
