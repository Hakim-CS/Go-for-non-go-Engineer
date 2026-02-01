# Learning Guide - Understanding the Code

This guide helps you understand what each part of the code does and why it's structured this way.

## 📁 File Overview

### Core Application Files

**`main.go`** - The Entry Point
- Starts the application
- Loads configuration
- Connects to database
- Sets up routes
- Starts the HTTP server

**Key Learning Points:**
- How Go applications are initialized
- Environment variable handling
- Dependency injection basics
- HTTP server setup with Gorilla Mux

---

**`models.go`** - Data Structures
- Defines all the data types (User, Exercise, Workout, etc.)
- Uses struct tags for JSON serialization
- Shows relationships between entities

**Key Learning Points:**
- Go structs and fields
- JSON struct tags (`json:"field_name"`)
- How to model database entities
- Understanding `-` tag to exclude fields from JSON

---

**`database.go`** - Database Layer
- Handles PostgreSQL connection
- Creates database schema (tables)
- Manages the database lifecycle

**Key Learning Points:**
- SQL database connections in Go
- `database/sql` package usage
- How to execute SQL queries
- Database schema management
- Foreign key relationships
- CASCADE deletion

---

**`auth.go`** - Security & Authentication
- Password hashing with bcrypt
- JWT token generation
- Token validation

**Key Learning Points:**
- Password security (never store plain text!)
- How bcrypt works (salting and hashing)
- JWT structure (header, payload, signature)
- Token expiration
- Claims (data stored in tokens)

---

**`middleware.go`** - Request Processing
- Authentication middleware (checks JWT tokens)
- CORS middleware (allows cross-origin requests)
- Context usage for passing data between handlers

**Key Learning Points:**
- Middleware pattern in Go
- HTTP middleware chaining
- Request context for user data
- CORS headers
- Authorization vs Authentication

---

**`seeder.go`** - Initial Data
- Populates the database with predefined exercises
- Runs only once (checks if data already exists)

**Key Learning Points:**
- Database seeding concept
- Preventing duplicate data
- Batch insertions

---

### Handler Files (API Endpoints)

**`handlers_auth.go`** - User Management
- Register new users
- Login existing users
- Get current user info

**Key Learning Points:**
- HTTP request/response handling
- JSON parsing (`json.Decoder`)
- Database queries and inserts
- Error handling and HTTP status codes
- Password validation

---

**`handlers_exercises.go`** - Exercise Library
- List all available exercises
- Simple read-only endpoint

**Key Learning Points:**
- SQL SELECT queries
- Iterating database results with `rows.Next()`
- Returning JSON arrays
- Public vs protected endpoints

---

**`handlers_workouts.go`** - Workout CRUD
- Create, Read, Update, Delete workouts
- Manage workout exercises
- Uses database transactions

**Key Learning Points:**
- CRUD operations in Go
- Database transactions (BEGIN, COMMIT, ROLLBACK)
- Why transactions are important (data consistency)
- URL parameter extraction (mux.Vars)
- Complex database queries with JOINs
- Helper functions to reduce code duplication

---

**`handlers_schedule.go`** - Workout Scheduling
- Schedule workouts for specific dates
- Mark workouts as completed
- Log workout performance

**Key Learning Points:**
- Working with time in Go
- Date parsing (RFC3339 format)
- Nullable database fields (sql.NullTime)
- Query parameters for filtering
- Authorization (ensure users can only modify their own data)

---

**`handlers_progress.go`** - Analytics
- Generate progress reports
- View exercise history
- Calculate statistics

**Key Learning Points:**
- Aggregate SQL queries (COUNT, MIN, MAX)
- Date arithmetic in SQL
- Complex data analysis
- Performance tracking over time

---

## 🎯 Key Go Concepts Used

### 1. **Structs and Methods**
```go
type User struct {
    ID       int
    Username string
}
```
Structs are like classes in other languages but simpler.

### 2. **Pointers**
```go
func getUser(id int) (*User, error)
```
The `*` means "pointer to" - it allows functions to return nil for "not found" cases.

### 3. **Error Handling**
```go
if err != nil {
    return err
}
```
Go uses explicit error handling instead of exceptions.

### 4. **JSON Tags**
```go
type User struct {
    ID int `json:"id"`
}
```
Tags control how structs are converted to/from JSON.

### 5. **Defer**
```go
defer rows.Close()
defer tx.Rollback()
```
Defer ensures cleanup happens when function exits.

### 6. **Interfaces** (implicit)
```go
http.HandlerFunc
```
HTTP handlers implement interfaces without explicitly declaring it.

### 7. **Package Organization**
All files are in `package main` because this is an executable program.

---

## 🔍 Common Patterns

### Pattern 1: Get User from Context
```go
claims := middleware.GetUserFromContext(r)
if claims == nil {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
}
```
**Why?** Authentication middleware puts user info in context, handlers extract it.

### Pattern 2: Parse JSON Request
```go
var req CreateWorkoutRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    http.Error(w, "Invalid request body", http.StatusBadRequest)
    return
}
```
**Why?** Converts JSON from client into Go structs.

### Pattern 3: Database Transaction
```go
tx, err := database.DB.Begin()
defer tx.Rollback()
// ... do work ...
tx.Commit()
```
**Why?** Ensures all changes succeed or none do (atomicity).

### Pattern 4: Return JSON
```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(data)
```
**Why?** Converts Go structs to JSON and sends to client.

---

## 🚀 How Requests Flow

1. **Client sends request** → `POST /api/workouts`
2. **Router matches route** → Finds the handler function
3. **Middleware runs** → Checks JWT token, adds user to context
4. **Handler executes** → `CreateWorkout` function
5. **Parse request** → Extract JSON from body
6. **Validate data** → Check required fields
7. **Database operation** → Insert workout and exercises
8. **Return response** → Send JSON back to client

---

## 📊 Database Relationships

```
users
  ↓ (has many)
workouts
  ↓ (has many)
workout_exercises → exercises
  ↓ (scheduled in)
schedules
  ↓ (logged as)
workout_logs → exercises
```

---

## 💡 Study Suggestions

### Week 1: Understanding Structure
- Read through `main.go` and understand the startup flow
- Look at `models.go` to understand the data structure
- Trace one request from start to finish

### Week 2: Database & Authentication
- Study `database.go` - how tables are created
- Understand `auth.go` - password hashing and JWT
- Try modifying the exercise seeder

### Week 3: HTTP Handlers
- Pick one handler file and understand every line
- Try adding a new endpoint
- Practice reading error messages

### Week 4: Advanced Features
- Study transaction usage in workout creation
- Understand the progress report calculations
- Try adding a new feature (e.g., workout sharing)

---

## 🛠️ Modification Ideas

### Easy:
1. Add more exercises to the seeder
2. Add a "favorite" flag to workouts
3. Add user profile fields (age, weight, height)

### Medium:
1. Add workout categories/tags
2. Implement workout templates
3. Add pagination to workout lists
4. Add search functionality for exercises

### Challenging:
1. Add personal records tracking (highest weight for each exercise)
2. Implement workout streak tracking
3. Add social features (follow users, share workouts)
4. Generate graphical progress reports

---

## 🐛 Common Issues & Solutions

### Issue: "dial tcp: connection refused"
**Solution:** PostgreSQL isn't running. Start it or check connection settings.

### Issue: "duplicate key violates unique constraint"
**Solution:** Username/email already exists. Use different credentials.

### Issue: "invalid or expired token"
**Solution:** Token has expired (24 hours). Login again to get a new token.

### Issue: "Workout not found" when it exists
**Solution:** Check if the workout belongs to the current user (authorization).

---

## 📚 Resources for Learning

### Go Basics:
- Tour of Go: https://go.dev/tour/
- Go by Example: https://gobyexample.com/

### Database:
- PostgreSQL Tutorial: https://www.postgresqltutorial.com/
- SQL Practice: https://sqlbolt.com/

### HTTP & REST APIs:
- HTTP Status Codes: https://httpstatuses.com/
- REST API Design: https://restfulapi.net/

### Security:
- JWT.io (JWT debugger): https://jwt.io/
- OWASP Guidelines: https://owasp.org/

---

## ✅ Code Quality Checklist

When reviewing or modifying code:

- [ ] Are errors handled properly?
- [ ] Are HTTP status codes correct?
- [ ] Is user authorization checked?
- [ ] Are database resources closed (defer)?
- [ ] Are SQL queries safe from injection?
- [ ] Are passwords hashed (never plain text)?
- [ ] Are tokens validated?
- [ ] Is JSON properly encoded/decoded?
- [ ] Are database transactions used when needed?
- [ ] Are comments clear and helpful?

---

## 🎓 What You've Learned

By studying this project, you'll understand:

✅ RESTful API design principles
✅ JWT authentication implementation
✅ Database design and relationships
✅ SQL queries and transactions
✅ Go HTTP server setup
✅ Middleware patterns
✅ Error handling in Go
✅ Security best practices
✅ Project structure and organization
✅ Testing APIs

This is production-quality code that demonstrates real-world patterns used in professional Go development!
