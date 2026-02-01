# 🚀 Getting Started Checklist

Follow this checklist to get your Workout Tracker API up and running!

## ☐ Prerequisites Check

- [ ] Go 1.21 or higher installed
  - Check: `go version`
  - Download: https://go.dev/dl/

- [ ] PostgreSQL installed and running
  - Check: `psql --version`
  - Download: https://www.postgresql.org/download/
  - Or use Docker: `docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres`

- [ ] Git installed (optional, for version control)
  - Check: `git --version`

## ☐ Database Setup

- [ ] Start PostgreSQL service
  - Windows: Check Services app
  - Docker: `docker start postgres`

- [ ] Create database
  ```sql
  CREATE DATABASE workout_tracker;
  ```
  
  Or via command line:
  ```powershell
  & "C:\Program Files\PostgreSQL\15\bin\psql.exe" -U postgres -c "CREATE DATABASE workout_tracker;"
  ```

- [ ] Verify database was created
  ```sql
  \l
  -- or
  SELECT datname FROM pg_database WHERE datname = 'workout_tracker';
  ```

## ☐ Project Setup

- [ ] Navigate to project directory
  ```powershell
  cd "d:\Go\Go for non-go Engineer\workout-tracker"
  ```

- [ ] Install Go dependencies
  ```bash
  go mod download
  ```

- [ ] Verify all files are present
  - main.go
  - models.go
  - database.go
  - auth.go
  - middleware.go
  - seeder.go
  - handlers_*.go files
  - go.mod

## ☐ Configuration (Optional)

- [ ] Review default configuration in main.go
  - Database: localhost:5432
  - User: postgres
  - Password: postgres
  - Database: workout_tracker

- [ ] Set custom environment variables if needed
  ```powershell
  $env:DB_HOST="localhost"
  $env:DB_PASSWORD="your_password"
  $env:SERVER_PORT="8080"
  ```

## ☐ First Run

- [ ] Start the application
  ```bash
  go run main.go
  ```

- [ ] Look for success messages:
  ```
  ✓ Successfully connected to database
  ✓ Database schema created successfully
  ✓ Successfully seeded 16 exercises
  ✓ Server is running on http://localhost:8080
  ```

- [ ] If errors occur, see Troubleshooting section below

## ☐ Verify Installation

- [ ] Test health endpoint
  ```bash
  curl http://localhost:8080/api/health
  ```
  Expected: `OK`

- [ ] Test exercises endpoint
  ```bash
  curl http://localhost:8080/api/exercises
  ```
  Expected: JSON array with 16 exercises

- [ ] Run automated test script
  ```powershell
  .\test_api.ps1
  ```
  Expected: All tests pass ✓

## ☐ First API Calls

- [ ] Register a user
  ```bash
  curl -X POST http://localhost:8080/api/register \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
  ```

- [ ] Login to get token
  ```bash
  curl -X POST http://localhost:8080/api/login \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","password":"password123"}'
  ```

- [ ] Save the token from response

- [ ] Create your first workout
  ```bash
  curl -X POST http://localhost:8080/api/workouts \
    -H "Authorization: Bearer YOUR_TOKEN_HERE" \
    -H "Content-Type: application/json" \
    -d '{"name":"First Workout","exercises":[{"exercise_id":1,"sets":3,"reps":10,"weight":50}]}'
  ```

## ☐ Explore the Code

- [ ] Read PROJECT_SUMMARY.md for overview
- [ ] Read LEARNING_GUIDE.md for code explanations
- [ ] Open main.go and trace the startup flow
- [ ] Review models.go to understand data structures
- [ ] Study one handler file (start with handlers_auth.go)

## ☐ Next Steps

- [ ] Review API_EXAMPLES.md for all endpoints
- [ ] Try creating multiple workouts
- [ ] Schedule some workouts
- [ ] Mark workouts as complete
- [ ] Check your progress report
- [ ] Experiment with modifying the code

## 🐛 Troubleshooting

### Problem: "Failed to connect to database"
**Solutions:**
- [ ] Ensure PostgreSQL is running
- [ ] Check connection credentials
- [ ] Verify database exists: `psql -U postgres -l`
- [ ] Check firewall/port 5432

### Problem: "Port 8080 already in use"
**Solutions:**
- [ ] Change port: `$env:SERVER_PORT="8081"`
- [ ] Find process using port: `netstat -ano | findstr :8080`
- [ ] Stop the conflicting process

### Problem: "go: cannot find module"
**Solutions:**
- [ ] Run `go mod tidy`
- [ ] Delete go.sum and run `go mod download`
- [ ] Check internet connection

### Problem: "invalid or expired token"
**Solutions:**
- [ ] Login again to get fresh token
- [ ] Check token is being sent in Authorization header
- [ ] Verify format: `Bearer YOUR_TOKEN_HERE`

### Problem: "Workout not found" but it exists
**Solutions:**
- [ ] Check if workout belongs to current user
- [ ] Verify you're using correct workout ID
- [ ] Check authentication token is valid

## 📚 Documentation Reference

- **README.md** - Project overview
- **SETUP.md** - Detailed setup guide
- **API_EXAMPLES.md** - Complete API reference
- **LEARNING_GUIDE.md** - Code explanations
- **PROJECT_SUMMARY.md** - Feature summary

## ✅ Success Indicators

You're ready to go when:
- ✓ Server starts without errors
- ✓ Health endpoint returns OK
- ✓ Exercises endpoint returns 16 items
- ✓ You can register and login
- ✓ You can create workouts with auth token
- ✓ Test script passes all tests

## 🎉 You're All Set!

Start experimenting with:
- Creating different workout combinations
- Scheduling workouts for different dates
- Tracking your progress over time
- Modifying the code to add features
- Understanding how each component works

**Happy Learning! 🚀**

---

**Need Help?**
- Review error messages carefully
- Check the documentation files
- Ensure all prerequisites are met
- Verify database connection
- Test with simple curl commands first
