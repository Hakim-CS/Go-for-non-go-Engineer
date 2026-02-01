# 🏋️ Workout Tracker API

A complete backend API for tracking workouts, managing exercise plans, and monitoring fitness progress. Built with Go, PostgreSQL, and JWT authentication.

## ✨ Features

- ✅ User registration and JWT authentication
- ✅ Predefined exercise library with 16 exercises across 6 categories
- ✅ Create and manage custom workouts (full CRUD)
- ✅ Schedule workouts for specific dates and times
- ✅ Track workout completion with detailed performance logs
- ✅ Generate progress reports and analytics
- ✅ Exercise history tracking over time
- ✅ Secure password hashing with bcrypt
- ✅ Authorization middleware for protected routes

## 🛠️ Tech Stack

- **Go 1.21** - Modern, efficient programming language
- **PostgreSQL** - Robust relational database
- **JWT** - Token-based authentication
- **Gorilla Mux** - Flexible HTTP router
- **bcrypt** - Secure password hashing

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 12 or higher

### 1. Clone and Navigate
```bash
cd workout-tracker
```

### 2. Setup Database
```sql
CREATE DATABASE workout_tracker;
```

### 3. Install Dependencies
```bash
go mod download
```

### 4. Run the Application
```bash
go run main.go
```

The server will start on `http://localhost:8080`

### 5. Test the API
```powershell
# Windows PowerShell
.\test_api.ps1

# Or manually test health check
curl http://localhost:8080/api/health
```

📖 **For detailed setup instructions, see [SETUP.md](SETUP.md)**

## 📡 API Endpoints

### Authentication (Public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/register` | Register new user |
| POST | `/api/login` | Login and get JWT token |

### Exercises (Public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/exercises` | List all exercises |

### Workouts (Protected - Requires JWT)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/workouts` | Create new workout |
| GET | `/api/workouts` | List user's workouts |
| GET | `/api/workouts/{id}` | Get workout details |
| PUT | `/api/workouts/{id}` | Update workout |
| DELETE | `/api/workouts/{id}` | Delete workout |

### Schedule (Protected - Requires JWT)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/schedule` | Schedule a workout |
| GET | `/api/schedule` | Get scheduled workouts |
| POST | `/api/schedule/{id}/complete` | Mark workout as completed |
| DELETE | `/api/schedule/{id}` | Delete scheduled workout |

### Progress (Protected - Requires JWT)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/progress` | Get progress report |
| GET | `/api/progress/exercise?exercise_id={id}` | Get exercise history |

📖 **For detailed API examples with request/response samples, see [API_EXAMPLES.md](API_EXAMPLES.md)**

## 📊 Example Usage

### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"secure123"}'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secure123"}'
```

### 3. Create a Workout (with JWT token)
```bash
curl -X POST http://localhost:8080/api/workouts \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Morning Routine",
    "description": "Quick morning workout",
    "exercises": [
      {"exercise_id": 1, "sets": 3, "reps": 10, "weight": 50}
    ]
  }'
```

## 🗂️ Project Structure

```
workout-tracker/
├── main.go                    # Application entry point & routing
├── models.go                  # Data structures
├── database.go               # Database connection & schema
├── auth.go                   # JWT & password security
├── middleware.go             # Authentication middleware
├── seeder.go                 # Exercise data seeding
├── handlers_*.go             # API endpoint handlers
├── go.mod                    # Go dependencies
├── README.md                 # This file
├── SETUP.md                  # Detailed setup guide
├── API_EXAMPLES.md          # Complete API examples
├── LEARNING_GUIDE.md        # Code learning guide
├── PROJECT_SUMMARY.md       # Project overview
└── test_api.ps1             # Automated test script
```

## 🎓 Learning Resources

This project is designed for **intermediate Go learners**. It demonstrates:

- RESTful API design principles
- JWT authentication and authorization
- Database operations with PostgreSQL
- SQL transactions for data integrity
- Middleware patterns
- Error handling in Go
- Project structure and organization
- Security best practices

**📚 New to this codebase?** Start with [LEARNING_GUIDE.md](LEARNING_GUIDE.md) for a comprehensive code walkthrough.

## 🧪 Testing

### Automated Testing
Run the PowerShell test script:
```powershell
.\test_api.ps1
```

This will test:
- Health check
- Exercise retrieval  
- User registration
- Login authentication
- Workout creation
- Scheduling
- Progress tracking

### Manual Testing
Use the examples in [API_EXAMPLES.md](API_EXAMPLES.md) with:
- cURL (command line)
- Postman (GUI)
- Any HTTP client

## 🔒 Security Features

- **Password Security**: bcrypt hashing with salt
- **Authentication**: JWT tokens with expiration
- **Authorization**: Middleware validates user access
- **SQL Injection Prevention**: Parameterized queries
- **CORS Configuration**: Controlled cross-origin access

## 📝 Configuration

### Environment Variables (Optional)
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=workout_tracker
DB_SSLMODE=disable
SERVER_PORT=8080
JWT_SECRET=your-secret-key
JWT_EXPIRATION_HOURS=24
```

Default values are provided, so the application works out of the box.

## 💡 Key Features in Detail

### Exercise Library
16 predefined exercises across categories:
- **Chest**: Bench Press, Push-ups, Dumbbell Flyes
- **Back**: Deadlift, Pull-ups, Barbell Rows
- **Legs**: Squats, Lunges, Leg Press
- **Shoulders**: Shoulder Press, Lateral Raises
- **Arms**: Bicep Curls, Tricep Dips
- **Core**: Plank, Crunches, Russian Twists

### Workout Management
- Create workouts with multiple exercises
- Specify sets, reps, and weight for each exercise
- Add notes and descriptions
- Update and delete workouts
- Full CRUD operations

### Progress Tracking
- Total workouts completed
- Weekly and monthly statistics
- Most frequent exercises
- Average workout frequency
- Exercise performance history
- Trend analysis over time

## 🛠️ Development

### Code Style
- Clear, descriptive variable names
- Extensive comments explaining logic
- Consistent error handling
- Proper HTTP status codes
- Clean function signatures

### Database Schema
6 tables with proper relationships:
- `users` - User accounts
- `exercises` - Exercise library
- `workouts` - Workout plans
- `workout_exercises` - Join table
- `schedules` - Scheduled workouts
- `workout_logs` - Performance logs

## 📦 Dependencies

```go
require (
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/gorilla/mux v1.8.1
    github.com/lib/pq v1.10.9
    golang.org/x/crypto v0.18.0
)
```

## 🐛 Troubleshooting

**Database Connection Issues?**
- Ensure PostgreSQL is running
- Check credentials in environment variables
- Verify database exists: `psql -U postgres -l`

**Port Already in Use?**
- Change port: `$env:SERVER_PORT="8081"`
- Or stop process using port 8080

**Token Expired?**
- Tokens expire after 24 hours
- Login again to get a new token

📖 **See [SETUP.md](SETUP.md) for more troubleshooting tips**
