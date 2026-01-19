# Grade Feature Documentation

## Overview

The Grade feature provides a complete CRUD API for managing student grades in courses. It follows the same architectural pattern as the Enrollment feature.

## Architecture

The implementation follows a clean, layered architecture:

```
models/grade.go           # Data structures and constants
repository/grade_repository.go  # Data access layer with in-memory storage
handlers/grade_handler.go       # HTTP request/response handlers
main.go                   # Route registration
```

## Data Model

### Grade Structure

```go
type Grade struct {
    ID          int       `json:"id"`           // Auto-incremented, starts from 1
    StudentID   int       `json:"student_id"`   // Required
    CourseID    int       `json:"course_id"`    // Required
    GradeValue  string    `json:"grade_value"`  // Required (e.g., "A", "B+", "C")
    GradePoints float64   `json:"grade_points"` // Grade point value
    GradeDate   time.Time `json:"grade_date"`   // Auto-set to current time if not provided
    Status      string    `json:"status"`       // Required: draft, submitted, or finalized
    CreatedAt   time.Time `json:"created_at"`   // Auto-set on creation
    UpdatedAt   time.Time `json:"updated_at"`   // Auto-updated on modification
}
```

### Status Constants

- `draft` - Grade is being drafted
- `submitted` - Grade has been submitted
- `finalized` - Grade has been finalized

## API Endpoints

All grade endpoints are under the `/api/grades` prefix.

### 1. Create Grade

**POST** `/api/grades`

Creates a new grade record.

**Request Body:**
```json
{
  "student_id": 101,
  "course_id": 201,
  "grade_value": "A",
  "grade_points": 4.0,
  "status": "draft"
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "student_id": 101,
  "course_id": 201,
  "grade_value": "A",
  "grade_points": 4.0,
  "grade_date": "2026-01-19T10:30:00Z",
  "status": "draft",
  "created_at": "2026-01-19T10:30:00Z",
  "updated_at": "2026-01-19T10:30:00Z"
}
```

### 2. Get Grade by ID

**GET** `/api/grades/{id}`

Retrieves a specific grade by its ID.

**Response:** `200 OK`
```json
{
  "id": 1,
  "student_id": 101,
  "course_id": 201,
  "grade_value": "A",
  "grade_points": 4.0,
  "grade_date": "2026-01-19T10:30:00Z",
  "status": "draft",
  "created_at": "2026-01-19T10:30:00Z",
  "updated_at": "2026-01-19T10:30:00Z"
}
```

**Error Response:** `404 Not Found`
```json
{
  "error": "grade not found"
}
```

### 3. List All Grades

**GET** `/api/grades`

Retrieves all grade records.

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "student_id": 101,
    "course_id": 201,
    "grade_value": "A",
    "grade_points": 4.0,
    "grade_date": "2026-01-19T10:30:00Z",
    "status": "draft",
    "created_at": "2026-01-19T10:30:00Z",
    "updated_at": "2026-01-19T10:30:00Z"
  }
]
```

### 4. Update Grade

**PUT** `/api/grades/{id}`

Updates an existing grade record.

**Request Body:**
```json
{
  "student_id": 101,
  "course_id": 201,
  "grade_value": "A+",
  "grade_points": 4.0,
  "status": "finalized"
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "student_id": 101,
  "course_id": 201,
  "grade_value": "A+",
  "grade_points": 4.0,
  "grade_date": "2026-01-19T10:30:00Z",
  "status": "finalized",
  "created_at": "2026-01-19T10:30:00Z",
  "updated_at": "2026-01-19T10:35:00Z"
}
```

**Error Response:** `404 Not Found`
```json
{
  "error": "grade not found"
}
```

### 5. Delete Grade

**DELETE** `/api/grades/{id}`

Deletes a grade record.

**Response:** `204 No Content`

**Error Response:** `404 Not Found`
```json
{
  "error": "grade not found"
}
```

## Validation Rules

### Required Fields

- `student_id` - Must be a non-zero integer
- `course_id` - Must be a non-zero integer
- `grade_value` - Must be a non-empty string
- `status` - Must be one of: `draft`, `submitted`, `finalized`

### Status Validation

The status field must be one of the predefined constants. Any other value will result in a `400 Bad Request` error:

```json
{
  "error": "status must be one of: draft, submitted, finalized"
}
```

### Validation Examples

**Missing student_id:**
```bash
curl -X POST http://localhost:8080/api/grades \
  -H "Content-Type: application/json" \
  -d '{"course_id":201,"grade_value":"A","status":"draft"}'
```
Response:
```json
{
  "error": "student_id is required"
}
```

**Invalid status:**
```bash
curl -X POST http://localhost:8080/api/grades \
  -H "Content-Type: application/json" \
  -d '{"student_id":101,"course_id":201,"grade_value":"A","status":"invalid"}'
```
Response:
```json
{
  "error": "status must be one of: draft, submitted, finalized"
}
```

## HTTP Status Codes

- `200 OK` - Successful GET or PUT request
- `201 Created` - Successful POST request
- `204 No Content` - Successful DELETE request
- `400 Bad Request` - Validation error or invalid request body
- `404 Not Found` - Grade not found

## Testing

### Using curl

```bash
# Create a grade
curl -X POST http://localhost:8080/api/grades \
  -H "Content-Type: application/json" \
  -d '{"student_id":101,"course_id":201,"grade_value":"A","grade_points":4.0,"status":"draft"}'

# Get a grade
curl http://localhost:8080/api/grades/1

# List all grades
curl http://localhost:8080/api/grades

# Update a grade
curl -X PUT http://localhost:8080/api/grades/1 \
  -H "Content-Type: application/json" \
  -d '{"student_id":101,"course_id":201,"grade_value":"A+","grade_points":4.0,"status":"finalized"}'

# Delete a grade
curl -X DELETE http://localhost:8080/api/grades/1
```

### Using PowerShell

```powershell
# Create a grade
Invoke-RestMethod -Uri http://localhost:8080/api/grades -Method Post `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"student_id":101,"course_id":201,"grade_value":"A","grade_points":4.0,"status":"draft"}'

# Get a grade
Invoke-RestMethod -Uri http://localhost:8080/api/grades/1

# List all grades
Invoke-RestMethod -Uri http://localhost:8080/api/grades

# Update a grade
Invoke-RestMethod -Uri http://localhost:8080/api/grades/1 -Method Put `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"student_id":101,"course_id":201,"grade_value":"A+","grade_points":4.0,"status":"finalized"}'

# Delete a grade
Invoke-RestMethod -Uri http://localhost:8080/api/grades/1 -Method Delete
```

## Implementation Details

### Thread Safety

The grade repository uses `sync.RWMutex` to ensure thread-safe concurrent access:
- Read operations (Get, GetAll) use read locks
- Write operations (Create, Update, Delete) use write locks

### Auto-increment IDs

Grade IDs are automatically assigned starting from 1 and increment sequentially for each new grade created.

### Timestamp Management

- `created_at` is set automatically when a grade is created
- `updated_at` is updated automatically on every modification
- `created_at` is preserved during updates

## Comparison with Enrollment Feature

Both features follow identical architectural patterns:

| Aspect | Enrollment | Grade |
|--------|-----------|-------|
| Route Prefix | `/api/enrollments` | `/api/grades` |
| Status Values | pending, active, completed | draft, submitted, finalized |
| Key Fields | enrollment_date | grade_value, grade_points, grade_date |
| CRUD Operations | ✅ All 5 | ✅ All 5 |
| Validation | ✅ Comprehensive | ✅ Comprehensive |
| Thread Safety | ✅ Mutex-based | ✅ Mutex-based |
| Auto-increment IDs | ✅ Starting from 1 | ✅ Starting from 1 |
