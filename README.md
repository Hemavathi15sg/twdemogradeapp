# Grade Management API - Demo Application

## 🚀 Overview

A demonstration Grade Management API built with Go, featuring both **Enrollment** and **Grade** CRUD operations. This application showcases clean architecture patterns and is designed for AI agent delegation demos.

## ✨ Features

### 📚 Enrollment Management
- Complete CRUD operations for student enrollments
- Status tracking (pending, active, completed)
- Thread-safe in-memory storage
- REST API at `/api/enrollments`

### 📊 Grade Management
- Complete CRUD operations for student grades
- Grade tracking with letter grades and points
- Status workflow (draft, submitted, finalized)
- Thread-safe in-memory storage
- REST API at `/api/grades`

## 🏗️ Architecture

The application follows a clean, layered architecture:

```
twdemogradeapp/
├── models/              # Data structures and domain entities
│   ├── enrollment.go
│   └── grade.go
├── repository/          # Data access layer (in-memory storage)
│   ├── enrollment_repository.go
│   └── grade_repository.go
├── handlers/            # HTTP request handlers
│   ├── enrollment_handler.go
│   └── grade_handler.go
└── main.go             # Application entry point and routing
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Docker (optional, for Redis)

### Run the Application

```bash
# Build and run
go build -o gradeapp
./gradeapp

# Or run directly
go run main.go
```

The server will start on `http://localhost:8080`

### Quick Test

```bash
# Health check
curl http://localhost:8080

# Create an enrollment
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{"student_id":42,"course_id":101,"status":"pending"}'

# Create a grade
curl -X POST http://localhost:8080/api/grades \
  -H "Content-Type: application/json" \
  -d '{"student_id":101,"course_id":201,"grade_value":"A","grade_points":4.0,"status":"draft"}'
```

## 📚 API Documentation

### Enrollment API

Base path: `/api/enrollments`

- **POST** `/api/enrollments` - Create enrollment
- **GET** `/api/enrollments/{id}` - Get enrollment by ID
- **GET** `/api/enrollments` - List all enrollments
- **PUT** `/api/enrollments/{id}` - Update enrollment
- **DELETE** `/api/enrollments/{id}` - Delete enrollment

**Status values:** `pending`, `active`, `completed`

### Grade API

Base path: `/api/grades`

- **POST** `/api/grades` - Create grade
- **GET** `/api/grades/{id}` - Get grade by ID
- **GET** `/api/grades` - List all grades
- **PUT** `/api/grades/{id}` - Update grade
- **DELETE** `/api/grades/{id}` - Delete grade

**Status values:** `draft`, `submitted`, `finalized`

📖 **[View detailed Grade API documentation](GRADE_FEATURE.md)**

## 🧪 Testing

### Using curl

```bash
# Test enrollments
curl -X POST http://localhost:8080/api/enrollments \
  -H "Content-Type: application/json" \
  -d '{"student_id":42,"course_id":101,"status":"pending"}'

curl http://localhost:8080/api/enrollments/1

# Test grades
curl -X POST http://localhost:8080/api/grades \
  -H "Content-Type: application/json" \
  -d '{"student_id":101,"course_id":201,"grade_value":"A","grade_points":4.0,"status":"draft"}'

curl http://localhost:8080/api/grades/1
```

### Using PowerShell

```powershell
# Test enrollments
Invoke-RestMethod -Uri http://localhost:8080/api/enrollments -Method Post `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"student_id":42,"course_id":101,"status":"pending"}'

# Test grades
Invoke-RestMethod -Uri http://localhost:8080/api/grades -Method Post `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"student_id":101,"course_id":201,"grade_value":"A","grade_points":4.0,"status":"draft"}'
```

## 🔒 Validation

Both features include comprehensive validation:

- **Required fields:** All required fields must be provided
- **Status validation:** Only predefined status values are accepted
- **ID validation:** IDs must be valid integers
- **Error messages:** Clear, descriptive error messages for all validation failures

## 🛠️ Technical Details

### Data Storage
- In-memory storage using Go maps
- Thread-safe with `sync.RWMutex`
- Auto-increment integer IDs starting from 1

### HTTP Status Codes
- `200 OK` - Successful GET/PUT
- `201 Created` - Successful POST
- `204 No Content` - Successful DELETE
- `400 Bad Request` - Validation errors
- `404 Not Found` - Resource not found

### Dependencies
- [gorilla/mux](https://github.com/gorilla/mux) - HTTP routing

## 🎯 AI Agent Delegation Demo

### Session 1 Plan

**Act 1: CRUD Boilerplate** → **Cloud Coding Agent**
- ✅ Complete enrollment CRUD API
- ✅ Complete grade CRUD API
- ✅ Validation and error handling
- ✅ Repository pattern implementation

**Act 2: Performance Caching** → **Local Agent** 
- Redis integration
- Cache-aside pattern
- Performance optimization

**Act 3: Quality Gates** → **Local Agent**
- OpenAPI specification
- Contract validation
- Integration testing

**Act 4: Documentation** → **Background Agent**
- Godoc comments
- README generation
- API documentation

### ⚙️ Copilot Agent Setup Tips

**For Cloud Coding Agent:**
- Ensure you're on feature branch before prompting
- Agent will create PR against current branch's upstream
- Use detailed Jira-style requirements

**For Local Agent:**
- Great for iterative improvements
- Faster for testing and validation
- Perfect for incremental changes

**For Background Agent:**
- Ideal for documentation tasks
- Non-blocking workflow
- Can work while you demo other features

## 📋 Project Status

- ✅ Enrollment feature - Complete
- ✅ Grade feature - Complete
- ⏳ Redis caching - Pending
- ⏳ OpenAPI specification - Pending
- ⏳ Integration tests - Pending

---

**Ready for 30-minute AI delegation demo!** 🤖✨