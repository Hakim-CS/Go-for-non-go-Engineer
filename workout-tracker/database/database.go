package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Config holds database connection configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DB is our global database connection
// In a real app, you might use dependency injection, but this is simpler for learning
var DB *sql.DB

// Connect establishes a connection to the PostgreSQL database
// This should be called once when the application starts
func Connect(config Config) error {
	// Build the connection string
	// Format: "host=localhost port=5432 user=postgres password=secret dbname=mydb sslmode=disable"
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)
	// Open a connection to the database
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// Verify the connection actually works
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Successfully connected to database")
	return nil
}

// Close closes the database connection
// This should be called when the application shuts down
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// InitializeSchema creates all the necessary database tables
// This is a simple migration system for learning purposes
func InitializeSchema() error {
	log.Println("Creating database schema...")

	// Create users table
	// This stores user account information
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	// Create exercises table
	// This stores our predefined exercise library
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS exercises (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			category VARCHAR(50),
			muscle_group VARCHAR(50)
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating exercises table: %w", err)
	}

	// Create workouts table
	// This stores user-created workout plans
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS workouts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating workouts table: %w", err)
	}

	// Create workout_exercises table
	// This is a join table that connects workouts with exercises
	// and stores the specific details (sets, reps, weight)
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS workout_exercises (
			id SERIAL PRIMARY KEY,
			workout_id INTEGER NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
			exercise_id INTEGER NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
			sets INTEGER NOT NULL,
			reps INTEGER NOT NULL,
			weight DECIMAL(5,2) DEFAULT 0,
			notes TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating workout_exercises table: %w", err)
	}

	// Create schedules table
	// This stores when workouts are scheduled and whether they're completed
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS schedules (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			workout_id INTEGER NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
			scheduled_date TIMESTAMP NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			completed_at TIMESTAMP,
			notes TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating schedules table: %w", err)
	}

	// Create workout_logs table
	// This tracks actual workout performance
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS workout_logs (
			id SERIAL PRIMARY KEY,
			schedule_id INTEGER NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
			exercise_id INTEGER NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
			sets_completed INTEGER NOT NULL,
			reps_completed INTEGER NOT NULL,
			weight_used DECIMAL(5,2) DEFAULT 0,
			duration INTEGER DEFAULT 0,
			notes TEXT,
			logged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating workout_logs table: %w", err)
	}

	log.Println("Database schema created successfully")
	return nil
}
