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
	// Initialize repository
	enrollmentRepo := repository.NewEnrollmentRepository()

	// Initialize handlers
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo)

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

	port := ":8080"
	fmt.Printf("🚀 Grade Management API starting on port %s\n", port)
	fmt.Println("📋 Ready for Copilot Agent delegation!")
	fmt.Println("📚 Enrollment API available at /api/enrollments")

	log.Fatal(http.ListenAndServe(port, r))
}
