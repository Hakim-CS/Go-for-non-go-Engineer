package models

import "time"

// User represents a registered user in the system
// Each user has their own workouts and progress tracking
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // "-" means this field won't be included in JSON responses
	CreatedAt    time.Time `json:"created_at"`
}

// Exercise represents a single exercise from our predefined library
// Examples: Push-ups, Squats, Bench Press, etc.
type Exercise struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`     // e.g., "Chest", "Legs", "Back", "Arms"
	MuscleGroup string `json:"muscle_group"` // e.g., "Upper Body", "Lower Body", "Core"
}

// Workout represents a collection of exercises that a user plans to do
// Think of it as a "workout plan" or "routine"
type Workout struct {
	ID          int               `json:"id"`
	UserID      int               `json:"user_id"`
	Name        string            `json:"name"`        // e.g., "Morning Chest Day", "Leg Workout"
	Description string            `json:"description"` // Optional notes about the workout
	Exercises   []WorkoutExercise `json:"exercises"`   // The exercises included in this workout
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// WorkoutExercise links an exercise to a workout with specific details
// This is the "join" between workouts and exercises, with extra info
type WorkoutExercise struct {
	ID         int       `json:"id"`
	WorkoutID  int       `json:"workout_id"`
	ExerciseID int       `json:"exercise_id"`
	Sets       int       `json:"sets"`               // Number of sets to perform
	Reps       int       `json:"reps"`               // Number of repetitions per set
	Weight     float64   `json:"weight"`             // Weight in kg or lbs
	Notes      string    `json:"notes"`              // Any additional notes
	Exercise   *Exercise `json:"exercise,omitempty"` // The actual exercise details (populated when needed)
}

// Schedule represents when a user plans to do a specific workout
// This connects workouts to calendar dates
type Schedule struct {
	ID            int        `json:"id"`
	UserID        int        `json:"user_id"`
	WorkoutID     int        `json:"workout_id"`
	ScheduledDate time.Time  `json:"scheduled_date"`    // When the workout is planned
	Completed     bool       `json:"completed"`         // Whether the workout was done
	CompletedAt   *time.Time `json:"completed_at"`      // When it was actually completed (null if not done)
	Notes         string     `json:"notes"`             // Post-workout notes
	Workout       *Workout   `json:"workout,omitempty"` // The workout details (populated when needed)
}

// WorkoutLog tracks the actual performance of a workout
// This records what the user actually did (may differ from the plan)
type WorkoutLog struct {
	ID            int       `json:"id"`
	ScheduleID    int       `json:"schedule_id"`
	ExerciseID    int       `json:"exercise_id"`
	SetsCompleted int       `json:"sets_completed"`
	RepsCompleted int       `json:"reps_completed"`
	WeightUsed    float64   `json:"weight_used"`
	Duration      int       `json:"duration"` // Duration in minutes
	Notes         string    `json:"notes"`
	LoggedAt      time.Time `json:"logged_at"`
}

// ProgressReport represents aggregated statistics for a user
// This helps users see their improvement over time
type ProgressReport struct {
	UserID                 int       `json:"user_id"`
	TotalWorkouts          int       `json:"total_workouts"`            // Total workouts completed
	TotalExercises         int       `json:"total_exercises"`           // Total exercises performed
	WorkoutsThisWeek       int       `json:"workouts_this_week"`        // Workouts in the last 7 days
	WorkoutsThisMonth      int       `json:"workouts_this_month"`       // Workouts in the last 30 days
	MostFrequentExercise   string    `json:"most_frequent_exercise"`    // Most commonly done exercise
	AverageWorkoutsPerWeek float64   `json:"average_workouts_per_week"` // Overall average
	StartDate              time.Time `json:"start_date"`                // When user started tracking
}
