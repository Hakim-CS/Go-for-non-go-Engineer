
# Docker Setup Guide

This guide explains how to run the Workout Tracker application with Docker or locally on your machine.

## Prerequisites

- **For Docker**: Docker and Docker Compose installed
- **For Local**: Go 1.21+, PostgreSQL installed locally

---

## Option 1: Running with Docker (Recommended)

Docker runs both PostgreSQL and the application in containers, making it easy to start and stop everything together.

### Start the Application

```powershell
# Build and start all services (PostgreSQL + App)
docker-compose up --build
```

The application will be available at `http://localhost:8080`

### Useful Docker Commands

```powershell
# Start in detached mode (background)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Stop and remove volumes (deletes database data)
docker-compose down -v

# Rebuild the app after code changes
docker-compose up --build app
```

### Access PostgreSQL in Docker

```powershell
# Connect to PostgreSQL container
docker exec -it workout-tracker-db psql -U postgres -d workout_tracker
```

---

## Option 2: Running Locally (PostgreSQL on your machine)

### Step 1: Create the Database

Connect to your local PostgreSQL and create the database:

```powershell
# Using psql command line
psql -U postgres

# In psql, run:
CREATE DATABASE workout_tracker;
\q
```

Or use pgAdmin or any PostgreSQL GUI tool to create a database named `workout_tracker`.

### Step 2: Set Environment Variables

Create a `.env` file in the project root (copy from `.env.example`):

```powershell
Copy-Item .env.example .env
```

Edit `.env` to match your local PostgreSQL settings:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=workout_tracker
DB_SSLMODE=disable
```

### Step 3: Load Environment Variables (PowerShell)

```powershell
# Load environment variables from .env file
Get-Content .env | ForEach-Object {
    if ($_ -match '^([^=]+)=(.*)$') {
        [Environment]::SetEnvironmentVariable($matches[1], $matches[2], "Process")
    }
}
```

### Step 4: Run the Application

```powershell
# Download dependencies
go mod download

# Run the application
go run main.go
```

The application will be available at `http://localhost:8080`

---

## Option 3: Docker PostgreSQL + Local Go App

Run PostgreSQL in Docker but develop the Go app locally.

### Step 1: Start Only PostgreSQL

```powershell
# Start only the database
docker-compose up postgres
```

### Step 2: Run the App Locally

Follow steps 2-4 from Option 2, but use these environment variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=workout_tracker
DB_SSLMODE=disable
```

Then run:

```powershell
# Load environment variables
Get-Content .env | ForEach-Object {
    if ($_ -match '^([^=]+)=(.*)$') {
        [Environment]::SetEnvironmentVariable($matches[1], $matches[2], "Process")
    }
}

# Run the app
go run main.go
```

---

## Testing the API

After starting the application, test it:

```powershell
# Health check
curl http://localhost:8080/api/health

# Register a user
curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d '{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\"}'
```

Or use the provided test script:

```powershell
.\explanation\test_api.ps1
```

---

## Troubleshooting

### Port Already in Use

If port 5432 or 8080 is already in use:

**For Docker:**
```powershell
# Stop conflicting services
docker-compose down

# Or change ports in docker-compose.yml
# Change "5432:5432" to "5433:5432" for PostgreSQL
# Change "8080:8080" to "8081:8080" for the app
```

**For Local:**
```powershell
# Find process using port 5432 (PostgreSQL)
netstat -ano | findstr :5432

# Find process using port 8080 (App)
netstat -ano | findstr :8080

# Stop the process by PID
taskkill /PID <process_id> /F
```

### Database Connection Error

If you see "Failed to connect to database":

1. **Docker**: Make sure PostgreSQL container is healthy:
   ```powershell
   docker-compose ps
   ```

2. **Local**: Verify PostgreSQL is running:
   ```powershell
   # Check PostgreSQL service
   Get-Service -Name postgresql*
   ```

3. Check your credentials in `.env` or environment variables

### Cannot Find go.sum

```powershell
# Initialize go.sum
go mod tidy
```

---

## Configuration Priority

The application reads configuration in this order:

1. Environment variables (highest priority)
2. Default values in code (fallback)

So you can override any setting by setting environment variables.
