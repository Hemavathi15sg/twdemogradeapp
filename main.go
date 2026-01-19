package main

import (
	"context"
	"fmt"
	"grademanagement-demo/cache"
	"grademanagement-demo/handlers"
	"grademanagement-demo/repository"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	// Test Redis connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Println("⚠️  Redis connection failed:", err)
		log.Println("ℹ️  Continuing without cache. Start Redis with: docker-compose up")
		redisClient = nil
	} else {
		log.Println("✅ Redis connected successfully")
	}

	// Initialize cache manager
	var cacheManager *cache.Manager
	if redisClient != nil {
		cacheManager = cache.NewManager(redisClient)
		log.Println("✅ Cache manager initialized")
	}

	// Initialize repositories
	enrollmentRepo := repository.NewEnrollmentRepository()
	gradeRepo := repository.NewGradeRepository()

	// Wire cache manager to repositories
	if cacheManager != nil {
		enrollmentRepo.SetCacheManager(cacheManager)
		gradeRepo.SetCacheManager(cacheManager)
		log.Println("✅ Cache wired to repositories")
	}

	// Initialize handlers
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo)
	gradeHandler := handlers.NewGradeHandler(gradeRepo)

	r := mux.NewRouter()

	// Basic health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Grade Management API - Ready for AI delegation!", "status": "healthy"}`)
	}).Methods("GET")

	// Cache stats endpoint
	if cacheManager != nil {
		r.HandleFunc("/api/cache-stats", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			stats := cacheManager.GetStats()
			w.Header().Set("X-Cache-Stats", cacheManager.LogStats())
			fmt.Fprintf(w, `{"hits":%d,"misses":%d,"updates":%d,"deletes":%d,"hit_rate":%d%%}`,
				stats["hits"], stats["misses"], stats["updates"], stats["deletes"], stats["hit_rate"])
		}).Methods("GET")
		log.Println("📊 Cache stats endpoint available at /api/cache-stats")
	}

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
