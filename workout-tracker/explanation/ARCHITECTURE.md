# 🏗️ System Architecture

This document explains how the Workout Tracker API is structured and how components interact.

## 📐 High-Level Architecture

```
┌─────────────┐
│   Client    │ (Browser, Mobile App, Postman, etc.)
└──────┬──────┘
       │ HTTP Requests (JSON)
       ▼
┌─────────────────────────────────────────┐
│         HTTP Server (Gorilla Mux)       │
│  ┌──────────────────────────────────┐   │
│  │   CORS Middleware                │   │
│  └──────────┬───────────────────────┘   │
│             ▼                            │
│  ┌──────────────────────────────────┐   │
│  │   Auth Middleware (JWT Check)    │   │
│  └──────────┬───────────────────────┘   │
│             ▼                            │
│  ┌──────────────────────────────────┐   │
│  │   Route Handlers                 │   │
│  │  - Auth Handlers                 │   │
│  │  - Workout Handlers              │   │
│  │  - Schedule Handlers             │   │
│  │  - Progress Handlers             │   │
│  └──────────┬───────────────────────┘   │
└─────────────┼───────────────────────────┘
              ▼
   ┌──────────────────┐
   │  Database Layer  │
   │   (PostgreSQL)   │
   └──────────────────┘
```

## 🔄 Request Flow

### 1. Public Endpoint (No Auth Required)

```
GET /api/exercises

Client Request
    ↓
CORS Middleware (adds headers)
    ↓
Route Handler (handlers_exercises.go)
    ↓
Database Query (SELECT * FROM exercises)
    ↓
JSON Response
    ↓
Client Receives Data
```

### 2. Protected Endpoint (Auth Required)

```
POST /api/workouts
Authorization: Bearer <token>

Client Request with Token
    ↓
CORS Middleware
    ↓
Auth Middleware
    ├─ Parse JWT token
    ├─ Validate token
    ├─ Extract user info
    └─ Add to request context
    ↓
Route Handler (handlers_workouts.go)
    ├─ Get user from context
    ├─ Parse request body
    ├─ Validate data
    └─ Begin transaction
        ├─ Insert workout
        ├─ Insert exercises
        └─ Commit transaction
    ↓
JSON Response
    ↓
Client Receives Data
```

## 📦 Component Details

### Main Application (main.go)

```
main()
  ├─ Load Configuration
  ├─ Connect to Database
  ├─ Initialize Schema (create tables)
  ├─ Seed Exercises
  ├─ Setup Routes
  └─ Start HTTP Server
```

**Responsibilities:**
- Application bootstrap
- Configuration management
- Route registration
- Server lifecycle

### Database Layer (database.go)

```
Database Package
  ├─ Connect() - Establish connection
  ├─ Close() - Cleanup connection
  └─ InitializeSchema() - Create tables
```

**Tables Created:**
1. users
2. exercises
3. workouts
4. workout_exercises
5. schedules
6. workout_logs

### Authentication (auth.go)

```
Auth Package
  ├─ HashPassword() - bcrypt hashing
  ├─ CheckPassword() - verify password
  ├─ GenerateToken() - create JWT
  └─ ValidateToken() - verify JWT
```

**Token Structure:**
```json
{
  "user_id": 1,
  "username": "john",
  "exp": 1234567890, // expiration time
  "iat": 1234567890 // issued at time
}
```

### Middleware (middleware.go)

```
Middleware Package
  ├─ AuthMiddleware() - JWT validation
  ├─ GetUserFromContext() - extract user
  └─ CORSMiddleware() - cross-origin headers
```

**Middleware Chain:**
```
Request → CORS → Auth → Handler → Response
```

### Handlers (handlers_*.go)

Each handler file handles a specific domain:

**handlers_auth.go**
- Register()
- Login()
- GetCurrentUser()

**handlers_exercises.go**
- GetExercises()

**handlers_workouts.go**
- CreateWorkout()
- GetWorkouts()
- GetWorkout()
- UpdateWorkout()
- DeleteWorkout()

**handlers_schedule.go**
- CreateSchedule()
- GetSchedules()
- CompleteSchedule()
- DeleteSchedule()

**handlers_progress.go**
- GetProgress()
- GetExerciseHistory()

## 🗄️ Database Schema

```
┌─────────────┐
│    users    │
│─────────────│
│ id (PK)     │───┐
│ username    │   │
│ email       │   │
│ password_hash│  │
│ created_at  │   │
└─────────────┘   │
                  │
        ┌─────────┴─────────┐
        │                   │
        ▼                   ▼
┌─────────────┐    ┌─────────────┐
│  workouts   │    │  schedules  │
│─────────────│    │─────────────│
│ id (PK)     │◄───│ id (PK)     │
│ user_id (FK)│    │ user_id (FK)│
│ name        │    │ workout_id(FK)
│ description │    │ scheduled_date
│ created_at  │    │ completed   │
│ updated_at  │    │ completed_at│
└──────┬──────┘    └──────┬──────┘
       │                  │
       │                  │
       ▼                  ▼
┌─────────────────┐ ┌─────────────┐
│workout_exercises│ │workout_logs │
│─────────────────│ │─────────────│
│ id (PK)         │ │ id (PK)     │
│ workout_id (FK) │ │ schedule_id │
│ exercise_id (FK)│ │ exercise_id │
│ sets            │ │ sets_completed
│ reps            │ │ reps_completed
│ weight          │ │ weight_used │
│ notes           │ │ duration    │
└────────┬────────┘ │ notes       │
         │          │ logged_at   │
         │          └─────────────┘
         ▼
    ┌─────────────┐
    │  exercises  │
    │─────────────│
    │ id (PK)     │
    │ name        │
    │ description │
    │ category    │
    │ muscle_group│
    └─────────────┘
```

**Relationships:**
- 1 User → Many Workouts
- 1 User → Many Schedules
- 1 Workout → Many Workout_Exercises
- 1 Exercise → Many Workout_Exercises
- 1 Schedule → Many Workout_Logs

## 🔐 Authentication Flow

### Registration
```
1. User submits credentials
   ↓
2. Server validates input
   ↓
3. Password hashed with bcrypt
   ↓
4. User saved to database
   ↓
5. User object returned (no password)
```

### Login
```
1. User submits credentials
   ↓
2. Server finds user by username
   ↓
3. Password verified with bcrypt
   ↓
4. JWT token generated
   ↓
5. Token + user object returned
```

### Protected Request
```
1. Client sends request with token
   Authorization: Bearer <token>
   ↓
2. Auth middleware extracts token
   ↓
3. Token validated and decoded
   ↓
4. User info added to context
   ↓
5. Handler accesses user from context
   ↓
6. Handler performs authorized action
```

## 📊 Data Flow Examples

### Creating a Workout

```
Client: POST /api/workouts
Body: {
  name: "Chest Day",
  exercises: [
    {exercise_id: 1, sets: 3, reps: 10}
  ]
}

Handler receives request
  ↓
Extract user_id from context (from JWT)
  ↓
Begin database transaction
  ↓
Insert workout record
  workouts table: {user_id, name, description}
  Returns: workout_id
  ↓
For each exercise:
  Insert workout_exercise record
    workout_exercises table: {workout_id, exercise_id, sets, reps, weight}
  ↓
Commit transaction
  ↓
Fetch complete workout (with exercises joined)
  ↓
Return JSON response
```

### Completing a Workout

```
Client: POST /api/schedule/{id}/complete
Body: {
  notes: "Great workout!",
  logs: [{exercise_id, sets, reps, weight}]
}

Handler receives request
  ↓
Extract user_id from context
  ↓
Begin transaction
  ↓
Update schedule:
  SET completed = true
  SET completed_at = NOW()
  SET notes = "Great workout!"
  WHERE id = {id} AND user_id = {user_id}
  ↓
For each log:
  Insert workout_log record
    workout_logs table: {schedule_id, exercise_id, sets, reps, weight, logged_at}
  ↓
Commit transaction
  ↓
Return updated schedule
```

### Getting Progress Report

```
Client: GET /api/progress

Handler receives request
  ↓
Extract user_id from context
  ↓
Run multiple queries:
  1. COUNT completed schedules (total workouts)
  2. COUNT workout_logs (total exercises)
  3. COUNT where completed_at >= 7 days ago
  4. COUNT where completed_at >= 30 days ago
  5. Find most frequent exercise (GROUP BY + ORDER BY)
  6. MIN(completed_at) for start date
  ↓
Calculate average workouts per week
  ↓
Build response object
  ↓
Return JSON
```

## 🎯 Key Design Patterns

### 1. **Middleware Pattern**
```go
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Do work before handler
        next(w, r)
        // Do work after handler
    }
}
```

### 2. **Repository Pattern** (Implicit)
```go
// Handlers interact with database
// Could be extracted to repository layer for cleaner separation
```

### 3. **Transaction Pattern**
```go
tx, _ := db.Begin()
defer tx.Rollback() // Rollback if anything fails

// Multiple operations
tx.Exec(...)
tx.Exec(...)

tx.Commit() // Commit if all succeed
```

### 4. **Context Pattern**
```go
// Middleware adds data to context
ctx := context.WithValue(r.Context(), UserKey, claims)

// Handler retrieves data from context
claims := r.Context().Value(UserKey).(*Claims)
```

## 🔄 Error Handling Strategy

```
Error occurs
  ↓
Log error details (for debugging)
  ↓
Return appropriate HTTP status:
  - 400 Bad Request (invalid input)
  - 401 Unauthorized (auth failed)
  - 403 Forbidden (no permission)
  - 404 Not Found (resource missing)
  - 500 Internal Server Error (server issue)
  ↓
Return error message to client
  http.Error(w, "Error message", statusCode)
```

## 📈 Scalability Considerations

### Current Design:
- Single database connection
- Synchronous request handling
- In-memory token validation

### Future Improvements:
- Connection pooling
- Caching layer (Redis)
- Rate limiting
- Load balancing
- Horizontal scaling
- Microservices architecture

## 🎓 Learning Path Through Code

1. **Start**: main.go (entry point)
2. **Routes**: See how endpoints are registered
3. **Middleware**: Understand request flow
4. **Handlers**: See business logic
5. **Database**: Understand data operations
6. **Models**: Learn data structures
7. **Auth**: Study security implementation

---

This architecture provides a solid foundation for learning Go web development while maintaining production-ready patterns and practices! 🚀
