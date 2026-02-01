# Setup Instructions

Follow these steps to get the Workout Tracker API running on your machine.

## Prerequisites

1. **Go 1.21 or higher**
   - Download from: https://go.dev/dl/
   - Verify installation: `go version`

2. **PostgreSQL**
   - Download from: https://www.postgresql.org/download/
   - Or use Docker: `docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres`

## Step-by-Step Setup

### 1. Create the Database

Open PostgreSQL command line (psql) and run:

```sql
CREATE DATABASE workout_tracker;
```

Or using command line:
```bash
# Windows (PowerShell)
& "C:\Program Files\PostgreSQL\15\bin\psql.exe" -U postgres -c "CREATE DATABASE workout_tracker;"

# Or if using Docker
docker exec -it postgres psql -U postgres -c "CREATE DATABASE workout_tracker;"
```

### 2. Configure Environment Variables (Optional)

The application works with default values, but you can customize by setting environment variables:

**Windows PowerShell:**
```powershell
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="postgres"
$env:DB_NAME="workout_tracker"
$env:DB_SSLMODE="disable"
$env:SERVER_PORT="8080"
$env:JWT_SECRET="your-secret-key-change-in-production"
$env:JWT_EXPIRATION_HOURS="24"
```

**Or create a `.env` file** (requires additional package):
Copy `.env.example` to `.env` and modify values.

### 3. Install Dependencies

In the `workout-tracker` directory:

```bash
go mod download
```

If you get any errors, try:
```bash
go mod tidy
```

### 4. Run the Application

```bash
go run main.go
```

You should see:
```
Starting Workout Tracker API...
Successfully connected to database
Creating database schema...
Database schema created successfully
Checking if exercises need to be seeded...
Seeding exercises...
Successfully seeded 16 exercises
Server is running on http://localhost:8080
```

### 5. Test the API

Try the health check endpoint:
```bash
curl http://localhost:8080/api/health
```

Should return: `OK`

## Troubleshooting

### "Failed to connect to database"
- Make sure PostgreSQL is running
- Check your database credentials
- Verify the database exists: `psql -U postgres -l`

### "dial tcp: connect: connection refused"
- PostgreSQL service is not running
- **Windows:** Check Services → PostgreSQL should be "Running"
- **Docker:** `docker ps` should show the postgres container

### "go: cannot find module"
- Run `go mod tidy` to fix missing dependencies
- Make sure you're in the correct directory

### Port 8080 already in use
- Change the port: `$env:SERVER_PORT="8081"`
- Or find and stop the process using port 8080

## Project Structure

```
workout-tracker/
├── main.go                    # Application entry point
├── models.go                  # Data structures
├── database.go                # Database connection and schema
├── auth.go                    # JWT and password hashing
├── middleware.go              # Authentication middleware
├── seeder.go                  # Exercise data seeder
├── handlers_auth.go           # Authentication endpoints
├── handlers_exercises.go      # Exercise endpoints
├── handlers_workouts.go       # Workout CRUD endpoints
├── handlers_schedule.go       # Scheduling endpoints
├── handlers_progress.go       # Progress tracking endpoints
├── go.mod                     # Go module definition
├── README.md                  # Project documentation
├── API_EXAMPLES.md           # API usage examples
└── SETUP.md                  # This file
```

## What Happens on First Run?

1. **Database Connection:** App connects to PostgreSQL
2. **Schema Creation:** Creates all necessary tables (users, exercises, workouts, etc.)
3. **Data Seeding:** Populates the exercises table with 16 predefined exercises
4. **Server Start:** API server starts listening on port 8080

## Next Steps

1. **Test the API:** Use the examples in `API_EXAMPLES.md`
2. **Register a user:** POST to `/api/register`
3. **Login:** POST to `/api/login` to get a JWT token
4. **Create workouts:** Use the token to access protected endpoints
5. **Track progress:** Schedule and complete workouts

## Learning Points

As you explore the code, pay attention to:

- **Project Structure:** How files are organized by responsibility
- **Database Operations:** SQL queries and transaction handling
- **JWT Authentication:** Token generation and validation
- **Middleware Pattern:** How authentication is applied to routes
- **Error Handling:** How errors are caught and returned to users
- **REST API Design:** HTTP methods and status codes
- **Go Patterns:** Struct tags, interfaces, error handling

## Development Tips

**Hot Reload:** Install `air` for automatic reloading during development:
```bash
go install github.com/cosmtrek/air@latest
air
```

**Database Inspection:**
```bash
psql -U postgres -d workout_tracker
\dt                    # List tables
SELECT * FROM users;   # View users
SELECT * FROM exercises; # View exercises
```

**Clear Database (start fresh):**
```sql
DROP DATABASE workout_tracker;
CREATE DATABASE workout_tracker;
```

Then restart the application to recreate tables and seed data.
