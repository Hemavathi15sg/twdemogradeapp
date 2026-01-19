---
applyTo: '**'
---
# GitHub Copilot Agent Instructions for Enrollment Feature Demo

## 🎯 CRITICAL: Follow These Rules Exactly

This instruction file ensures consistent, predictable behavior for demo purposes. **Every feature must follow this structure exactly.**

---

## 📁 MANDATORY Folder Structure

**ALWAYS create these folders if they don't exist:**

```
project-root/
├── models/           # Data structures and domain entities
├── handlers/         # HTTP handlers (controllers)
├── repository/       # Data access layer
└── main.go          # Entry point with routes only
```

**❌ NEVER put handler logic, models, or repository code directly in main.go**
**❌ NEVER create files directly in root (except main.go, go.mod, go.sum)**

---

## 🔢 Data Type Standards (MANDATORY)

### ID Fields - ALWAYS Use `int`
```go
// ✅ CORRECT - Use simple integers
type Enrollment struct {
    ID        int    `json:"id"`
    StudentID int    `json:"student_id"`
    CourseID  int    `json:"course_id"`
    // ... other fields
}

// ❌ WRONG - Never use UUID or string IDs for demo
type Enrollment struct {
    ID        string    `json:"id"`  // NO!
    StudentID uuid.UUID `json:"student_id"`  // NO!
}
```

### Date/Time Fields - Use `time.Time` with JSON tags
```go
EnrollmentDate time.Time `json:"enrollment_date"`
CreatedAt      time.Time `json:"created_at"`
UpdatedAt      time.Time `json:"updated_at"`
```

### Status Fields - Use string constants
```go
const (
    StatusPending   = "pending"
    StatusActive    = "active"
    StatusCompleted = "completed"
)

type Enrollment struct {
    Status string `json:"status"`
}
```

---

## 📝 File Organization Rules

### 1. models/enrollment.go
**Purpose:** Define data structures ONLY
**Contains:**
- Struct definition with JSON tags (snake_case)
- Status constants
- NO business logic
- NO HTTP handling

```go
package models

import "time"

type Enrollment struct {
    ID             int       `json:"id"`
    StudentID      int       `json:"student_id"`
    CourseID       int       `json:"course_id"`
    EnrollmentDate time.Time `json:"enrollment_date"`
    Status         string    `json:"status"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

const (
    StatusPending   = "pending"
    StatusActive    = "active"
    StatusCompleted = "completed"
)
```

---

### 2. repository/enrollment_repository.go
**Purpose:** Data access and storage ONLY
**Contains:**
- In-memory storage with map[int]models.Enrollment
- CRUD operations: Create, Get, GetAll, Update, Delete
- Thread safety with sync.RWMutex
- Auto-increment ID logic

**Required Methods:**
```go
type EnrollmentRepository struct {
    enrollments map[int]models.Enrollment
    mutex       sync.RWMutex
    nextID      int
}

func NewEnrollmentRepository() *EnrollmentRepository
func (r *EnrollmentRepository) Create(enrollment models.Enrollment) models.Enrollment
func (r *EnrollmentRepository) Get(id int) (models.Enrollment, bool)
func (r *EnrollmentRepository) GetAll() []models.Enrollment
func (r *EnrollmentRepository) Update(id int, enrollment models.Enrollment) (models.Enrollment, bool)
func (r *EnrollmentRepository) Delete(id int) bool
```

**Auto-increment ID logic:**
```go
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
```

---

### 3. handlers/enrollment_handler.go
**Purpose:** HTTP request/response handling ONLY
**Contains:**
- All 5 CRUD handlers
- Input validation
- Error responses with proper HTTP codes
- JSON encoding/decoding

**Required Handler Structure:**
```go
type EnrollmentHandler struct {
    repo *repository.EnrollmentRepository
}

func NewEnrollmentHandler(repo *repository.EnrollmentRepository) *EnrollmentHandler {
    return &EnrollmentHandler{repo: repo}
}

// CreateEnrollment - POST /api/enrollments
func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request)

// GetEnrollment - GET /api/enrollments/{id}
func (h *EnrollmentHandler) GetEnrollment(w http.ResponseWriter, r *http.Request)

// ListEnrollments - GET /api/enrollments
func (h *EnrollmentHandler) ListEnrollments(w http.ResponseWriter, r *http.Request)

// UpdateEnrollment - PUT /api/enrollments/{id}
func (h *EnrollmentHandler) UpdateEnrollment(w http.ResponseWriter, r *http.Request)

// DeleteEnrollment - DELETE /api/enrollments/{id}
func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request)
```

---

### 4. main.go
**Purpose:** Entry point and route registration ONLY
**Contains:**
- Import statements
- Route setup
- Server start
- NO handler logic
- NO business logic

```go
package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "yourproject/models"
    "yourproject/handlers"
    "yourproject/repository"
)

func main() {
    // Initialize repository
    enrollmentRepo := repository.NewEnrollmentRepository()
    
    // Initialize handlers
    enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo)
    
    // Setup routes
    r := mux.NewRouter()
    
    // Enrollment routes
    r.HandleFunc("/api/enrollments", enrollmentHandler.CreateEnrollment).Methods("POST")
    r.HandleFunc("/api/enrollments/{id}", enrollmentHandler.GetEnrollment).Methods("GET")
    r.HandleFunc("/api/enrollments", enrollmentHandler.ListEnrollments).Methods("GET")
    r.HandleFunc("/api/enrollments/{id}", enrollmentHandler.UpdateEnrollment).Methods("PUT")
    r.HandleFunc("/api/enrollments/{id}", enrollmentHandler.DeleteEnrollment).Methods("DELETE")
    
    // Start server
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
```

---

## 🧪 Validation Rules (MANDATORY)

### Status Validation
```go
func validateStatus(status string) error {
    if status != models.StatusPending && 
       status != models.StatusActive && 
       status != models.StatusCompleted {
        return fmt.Errorf("status must be one of: pending, active, completed")
    }
    return nil
}
```

### Required Fields Validation
```go
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
```

---

## 📡 HTTP Response Standards

### Success Responses
```go
// 200 OK - GET, PUT
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(enrollment)

// 201 Created - POST
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(enrollment)

// 204 No Content - DELETE
w.WriteHeader(http.StatusNoContent)
```

### Error Responses (MUST be consistent)
```go
// 400 Bad Request - Validation errors
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusBadRequest)
json.NewEncoder(w).Encode(map[string]string{"error": "status must be one of: pending, active, completed"})

// 404 Not Found - Resource doesn't exist
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusNotFound)
json.NewEncoder(w).Encode(map[string]string{"error": "enrollment not found"})

// 500 Internal Server Error - Unexpected errors
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusInternalServerError)
json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
```

---

## 🔧 Testing Commands (For Demo Verification)

### These commands MUST work after implementation:

```powershell
# 1. Create enrollment (should return ID: 1)
Invoke-RestMethod -Uri http://localhost:8080/api/enrollments -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"student_id":42,"course_id":101,"status":"pending"}'

# Expected response:
# {
#   "id": 1,
#   "student_id": 42,
#   "course_id": 101,
#   "status": "pending",
#   "enrollment_date": "2026-01-19T10:30:00Z",
#   "created_at": "2026-01-19T10:30:00Z",
#   "updated_at": "2026-01-19T10:30:00Z"
# }

# 2. Get enrollment by ID
Invoke-RestMethod -Uri http://localhost:8080/api/enrollments/1

# 3. List all enrollments
Invoke-RestMethod -Uri http://localhost:8080/api/enrollments

# 4. Update enrollment
Invoke-RestMethod -Uri http://localhost:8080/api/enrollments/1 -Method Put -Headers @{"Content-Type"="application/json"} -Body '{"student_id":42,"course_id":101,"status":"completed"}'

# 5. Test validation (should FAIL with 400)
Invoke-RestMethod -Uri http://localhost:8080/api/enrollments -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"student_id":42,"course_id":101,"status":"invalid"}'

# Expected error:
# Invoke-RestMethod : status must be one of: pending, active, completed

# 6. Delete enrollment
Invoke-RestMethod -Uri http://localhost:8080/api/enrollments/1 -Method Delete
```

---

## 🚨 Common Mistakes to AVOID

### ❌ DON'T: Put everything in main.go
```go
// ❌ WRONG
func main() {
    type Enrollment struct { ... }  // NO!
    func createHandler(...) { ... }  // NO!
}
```

### ❌ DON'T: Use string or UUID for IDs
```go
// ❌ WRONG
type Enrollment struct {
    ID string `json:"id"`  // Use int!
}
```

### ❌ DON'T: Create files in root directory
```
project-root/
├── enrollment.go        # ❌ NO!
├── enrollment_handler.go  # ❌ NO!
```

### ❌ DON'T: Inconsistent error messages
```go
// ❌ WRONG
return errors.New("bad status")  // Too vague
return errors.New("Status must be one of: pending, active, or completed")  // Inconsistent format
```

---

## ✅ Success Criteria Checklist

Before marking implementation complete, verify:

- [ ] Folder structure exists: `models/`, `handlers/`, `repository/`
- [ ] All IDs are `int` type (not string/UUID)
- [ ] Status constants defined in models package
- [ ] Repository has auto-increment ID logic starting from 1
- [ ] All 5 CRUD handlers implemented in separate handler file
- [ ] Routes registered in main.go with `/api` prefix
- [ ] Status validation rejects invalid values with proper error message
- [ ] All test commands execute successfully
- [ ] First created enrollment has ID: 1
- [ ] Error responses include descriptive messages
- [ ] Thread safety implemented with mutex
- [ ] JSON tags use snake_case convention

---

## 🤖 Agent Prompt Template

**When delegating to GitHub Copilot Coding Agent, use this prompt:**

```
@github #github-pull-request_copilot-coding-agent 

Create enrollment feature following COPILOT_AGENT_INSTRUCTIONS.md:

CRITICAL REQUIREMENTS:
1. Create folder structure: models/, handlers/, repository/
2. Use int for ALL IDs (not UUID/string)
3. models/enrollment.go: struct with int ID, student_id, course_id, status (string), enrollment_date, created_at, updated_at; add status constants
4. repository/enrollment_repository.go: in-memory storage with auto-increment int IDs starting from 1, thread-safe with mutex
5. handlers/enrollment_handler.go: 5 CRUD handlers with validation
6. main.go: ONLY routes and server setup, no handler logic
7. Status validation: must be "pending", "active", or "completed"
8. Routes under /api prefix
9. Proper HTTP codes: 200, 201, 204, 400, 404, 500
10. Descriptive error messages in JSON format

Verify these test commands work:
- Create: POST /api/enrollments (returns ID: 1)
- Get: GET /api/enrollments/1
- List: GET /api/enrollments
- Update: PUT /api/enrollments/1
- Delete: DELETE /api/enrollments/1
- Validation: Invalid status returns 400

Follow instruction file exactly. Do not deviate from folder structure or data types.
```

---

## 📞 Quick Reference Card

**For Demo Script Usage:**

1. **Before agent call:** Ensure COPILOT_AGENT_INSTRUCTIONS.md exists in repo
2. **Agent prompt:** Reference instruction file explicitly
3. **After implementation:** Run test commands to verify
4. **First enrollment ID:** MUST be 1 (not 0, not UUID)
5. **Invalid status test:** MUST return 400 with clear message
6. **Folder structure:** MUST have models/, handlers/, repository/

**If agent deviates:**
- Stop and reference this instruction file
- Regenerate with explicit folder structure requirement
- Verify ID data type is int

---

## 🎯 Demo Success Pattern

```
1. Agent receives prompt with instruction file reference
2. Agent creates folder structure
3. Agent implements with int IDs
4. Test commands execute successfully
5. First enrollment has ID: 1
6. Status validation works correctly
7. Demo proceeds smoothly ✅
```

---

**This instruction file is your source of truth for consistent, predictable demo behavior. Reference it in EVERY agent prompt for enrollment feature work.**
