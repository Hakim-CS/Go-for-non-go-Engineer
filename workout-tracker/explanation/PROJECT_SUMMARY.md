# 🏋️ Workout Tracker API - Complete Implementation

A professional-grade RESTful API for tracking workouts, built with Go, PostgreSQL, and JWT authentication.

## 🎯 What's Been Built

### ✅ Complete Feature Set

1. **User Authentication System**
   - User registration with secure password hashing (bcrypt)
   - Login with JWT token generation
   - Protected routes with authentication middleware
   - Token expiration and validation

2. **Exercise Library**
   - 16 predefined exercises across multiple categories
   - Categories: Chest, Back, Legs, Shoulders, Arms, Core
   - Automatically seeded on first run
   - Read-only public access

3. **Workout Management (Full CRUD)**
   - Create custom workouts with multiple exercises
   - Specify sets, reps, and weight for each exercise
   - Update existing workouts
   - Delete workouts
   - List all user workouts
   - View detailed workout information

4. **Workout Scheduling**
   - Schedule workouts for specific dates and times
   - Filter by upcoming, completed, or pending
   - Mark workouts as complete with performance logs
   - Track actual performance vs planned workout

5. **Progress Tracking & Analytics**
   - Total workouts completed
   - Workouts this week/month
   - Most frequently performed exercise
   - Average workouts per week
   - Exercise history over time
   - Performance trends

## 📂 Project Structure

```
workout-tracker/
├── main.go                      # Application entry point & routing
├── models.go                    # Data structures for all entities
├── database.go                  # Database connection & schema
├── auth.go                      # JWT & password security
├── middleware.go                # Authentication & CORS middleware
├── seeder.go                    # Exercise data seeding
├── handlers_auth.go            # Registration & login endpoints
├── handlers_exercises.go       # Exercise listing endpoint
├── handlers_workouts.go        # Workout CRUD endpoints
├── handlers_schedule.go        # Scheduling endpoints
├── handlers_progress.go        # Progress & analytics endpoints
├── go.mod                      # Go dependencies
├── .env.example                # Environment variable template
├── config.yaml                 # Configuration file
├── .gitignore                  # Git ignore rules
├── README.md                   # Project overview
├── SETUP.md                    # Setup instructions
├── API_EXAMPLES.md            # API usage examples
├── LEARNING_GUIDE.md          # Code learning guide
└── test_api.ps1               # Automated API test script
```

## 🔧 Technologies Used

- **Go 1.21** - Programming language
- **PostgreSQL** - Relational database
- **JWT** (golang-jwt/jwt) - Token-based authentication
- **Gorilla Mux** - HTTP routing
- **bcrypt** - Password hashing
- **lib/pq** - PostgreSQL driver

## 🚀 Quick Start

### 1. Install Dependencies
```bash
go mod download
```

### 2. Setup Database
```sql
CREATE DATABASE workout_tracker;
```

### 3. Run the Application
```bash
go run main.go
```

### 4. Test the API
```powershell
.\test_api.ps1
```

## 📡 API Endpoints

### Public Endpoints
- `POST /api/register` - Register new user
- `POST /api/login` - Login and get token
- `GET /api/exercises` - List all exercises
- `GET /api/health` - Health check

### Protected Endpoints (Require JWT Token)
- `GET /api/me` - Get current user info
- `POST /api/workouts` - Create workout
- `GET /api/workouts` - List workouts
- `GET /api/workouts/{id}` - Get workout details
- `PUT /api/workouts/{id}` - Update workout
- `DELETE /api/workouts/{id}` - Delete workout
- `POST /api/schedule` - Schedule workout
- `GET /api/schedule` - Get scheduled workouts
- `POST /api/schedule/{id}/complete` - Mark as complete
- `DELETE /api/schedule/{id}` - Delete schedule
- `GET /api/progress` - Get progress report
- `GET /api/progress/exercise` - Get exercise history

## 🎓 Learning Value

### Intermediate Go Concepts Demonstrated

1. **Project Organization**
   - Clean separation of concerns
   - Logical file structure
   - Package management

2. **HTTP & REST**
   - RESTful API design
   - HTTP methods and status codes
   - Request/response handling
   - JSON encoding/decoding

3. **Database Operations**
   - Connection management
   - CRUD operations
   - SQL transactions
   - Foreign key relationships
   - Query optimization

4. **Security**
   - Password hashing (bcrypt)
   - JWT token generation & validation
   - Authorization middleware
   - Secure credential handling

5. **Go Patterns**
   - Error handling
   - Middleware pattern
   - Context usage
   - Defer statements
   - Struct tags
   - Pointer receivers

6. **API Design**
   - Authentication flow
   - Protected routes
   - Query parameters
   - URL parameters
   - Response formatting

## 📊 Database Schema

**6 Tables:**
- `users` - User accounts
- `exercises` - Exercise library
- `workouts` - User workout plans
- `workout_exercises` - Exercises in workouts
- `schedules` - Scheduled workouts
- `workout_logs` - Performance tracking

**Relationships:**
- Users → Workouts (one-to-many)
- Workouts → Exercises (many-to-many through workout_exercises)
- Users → Schedules (one-to-many)
- Schedules → Workout Logs (one-to-many)

## 🧪 Testing

**Automated Test Script:**
```powershell
.\test_api.ps1
```

Tests all major functionality:
- Health check
- Exercise retrieval
- User registration
- Login
- Workout creation
- Scheduling
- Progress tracking

**Manual Testing:**
See `API_EXAMPLES.md` for detailed cURL and JSON examples.

## 📝 Code Comments

Every file includes detailed comments explaining:
- What each function does
- Why certain approaches are used
- How components interact
- Key concepts being demonstrated

## 🎯 Next Steps for Learning

1. **Run and Test**: Start the server and test all endpoints
2. **Read the Code**: Start with `main.go` and follow the flow
3. **Study Patterns**: Notice repeated patterns across handlers
4. **Modify It**: Try adding new features or endpoints
5. **Debug It**: Use breakpoints to understand execution flow
6. **Extend It**: Add features like workout sharing or social aspects

## 💡 Modification Ideas

### Beginner:
- Add more exercises to the seeder
- Change token expiration time
- Add user profile fields

### Intermediate:
- Add pagination to lists
- Implement search functionality
- Add workout categories/tags
- Create workout templates

### Advanced:
- Personal records tracking
- Achievement system
- Social features (follow users)
- Graphical progress reports
- Real-time notifications

## 📖 Documentation Files

- **README.md** - Project overview and features
- **SETUP.md** - Detailed setup instructions
- **API_EXAMPLES.md** - Complete API usage guide
- **LEARNING_GUIDE.md** - Code explanation and learning path

## ✨ Key Features

✅ Clean, readable code with extensive comments
✅ Production-ready patterns and practices
✅ Secure authentication and authorization
✅ Complete CRUD operations
✅ Database transactions for data integrity
✅ Proper error handling
✅ RESTful API design
✅ Middleware architecture
✅ Relationship modeling
✅ Progress tracking and analytics

## 🔒 Security Features

- Passwords hashed with bcrypt (never stored as plain text)
- JWT tokens with expiration
- Authorization checks on all protected routes
- SQL injection prevention (parameterized queries)
- CORS configuration for API security

## 🎉 Project Highlights

This is a **complete, working application** that demonstrates:
- Professional Go development practices
- Real-world API architecture
- Secure authentication patterns
- Database design and relationships
- Clean code principles
- Comprehensive documentation

Perfect for learning, showcasing skills, or as a foundation for more complex projects!

---

**Built for intermediate Go learners** - Clear, understandable, and well-documented code that demonstrates real-world development patterns. 🚀
