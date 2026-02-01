package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"workout-tracker/database"
	"workout-tracker/middleware"
	"workout-tracker/models"

	"github.com/gorilla/mux"
)

// CreateScheduleRequest represents the data needed to schedule a workout
type CreateScheduleRequest struct {
	WorkoutID     int    `json:"workout_id"`
	ScheduledDate string `json:"scheduled_date"` // Format: "2024-01-15T10:00:00Z"
	Notes         string `json:"notes"`
}

// CompleteScheduleRequest represents data for marking a workout as complete
type CompleteScheduleRequest struct {
	Notes string              `json:"notes"`
	Logs  []WorkoutLogRequest `json:"logs"` // Optional: detailed exercise logs
}

// WorkoutLogRequest represents a logged exercise
type WorkoutLogRequest struct {
	ExerciseID    int     `json:"exercise_id"`
	SetsCompleted int     `json:"sets_completed"`
	RepsCompleted int     `json:"reps_completed"`
	WeightUsed    float64 `json:"weight_used"`
	Duration      int     `json:"duration"` // in minutes
	Notes         string  `json:"notes"`
}

// CreateSchedule schedules a workout for a specific date
func CreateSchedule(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body
	var req CreateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse the scheduled date
	scheduledDate, err := time.Parse(time.RFC3339, req.ScheduledDate)
	if err != nil {
		http.Error(w, "Invalid date format. Use RFC3339 format like: 2024-01-15T10:00:00Z", http.StatusBadRequest)
		return
	}

	// Verify the workout belongs to this user
	var workoutUserID int
	err = database.DB.QueryRow(`SELECT user_id FROM workouts WHERE id = $1`, req.WorkoutID).Scan(&workoutUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Workout not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error verifying workout", http.StatusInternalServerError)
		return
	}

	if workoutUserID != claims.UserID {
		http.Error(w, "You don't have permission to schedule this workout", http.StatusForbidden)
		return
	}

	// Insert the schedule
	var scheduleID int
	err = database.DB.QueryRow(
		`INSERT INTO schedules (user_id, workout_id, scheduled_date, completed, notes) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		claims.UserID, req.WorkoutID, scheduledDate, false, req.Notes,
	).Scan(&scheduleID)

	if err != nil {
		http.Error(w, "Error creating schedule", http.StatusInternalServerError)
		return
	}

	// Fetch and return the created schedule
	schedule, err := getScheduleByID(scheduleID)
	if err != nil {
		http.Error(w, "Error fetching created schedule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schedule)
}

// GetSchedules returns all scheduled workouts for the logged-in user
func GetSchedules(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Optional query parameters for filtering
	upcoming := r.URL.Query().Get("upcoming")   // "true" to show only future workouts
	completed := r.URL.Query().Get("completed") // "true" or "false" to filter by completion

	// Build the query
	query := `
		SELECT s.id, s.user_id, s.workout_id, s.scheduled_date, s.completed, s.completed_at, s.notes,
		       w.id, w.user_id, w.name, w.description, w.created_at, w.updated_at
		FROM schedules s
		JOIN workouts w ON s.workout_id = w.id
		WHERE s.user_id = $1
	`

	// Add filters
	args := []interface{}{claims.UserID}
	argCount := 1

	if upcoming == "true" {
		argCount++
		query += ` AND s.scheduled_date >= NOW()`
	}

	if completed == "true" {
		query += ` AND s.completed = true`
	} else if completed == "false" {
		query += ` AND s.completed = false`
	}

	query += ` ORDER BY s.scheduled_date ASC`

	// Execute the query
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Error fetching schedules", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	schedules := []models.Schedule{}
	for rows.Next() {
		var schedule models.Schedule
		var workout models.Workout
		var completedAt sql.NullTime

		err := rows.Scan(
			&schedule.ID,
			&schedule.UserID,
			&schedule.WorkoutID,
			&schedule.ScheduledDate,
			&schedule.Completed,
			&completedAt,
			&schedule.Notes,
			&workout.ID,
			&workout.UserID,
			&workout.Name,
			&workout.Description,
			&workout.CreatedAt,
			&workout.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "Error scanning schedule", http.StatusInternalServerError)
			return
		}

		// Handle nullable completed_at
		if completedAt.Valid {
			schedule.CompletedAt = &completedAt.Time
		}

		schedule.Workout = &workout
		schedules = append(schedules, schedule)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)
}

// CompleteSchedule marks a scheduled workout as completed
func CompleteSchedule(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get schedule ID from URL
	vars := mux.Vars(r)
	scheduleID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var req CompleteScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Start a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Mark the schedule as completed
	completedAt := time.Now()
	result, err := tx.Exec(
		`UPDATE schedules SET completed = true, completed_at = $1, notes = $2 
		 WHERE id = $3 AND user_id = $4`,
		completedAt, req.Notes, scheduleID, claims.UserID,
	)
	if err != nil {
		http.Error(w, "Error completing schedule", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Schedule not found", http.StatusNotFound)
		return
	}

	// If logs are provided, insert them
	for _, log := range req.Logs {
		_, err = tx.Exec(
			`INSERT INTO workout_logs (schedule_id, exercise_id, sets_completed, reps_completed, weight_used, duration, notes, logged_at) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			scheduleID, log.ExerciseID, log.SetsCompleted, log.RepsCompleted, log.WeightUsed, log.Duration, log.Notes, time.Now(),
		)
		if err != nil {
			http.Error(w, "Error logging workout details", http.StatusInternalServerError)
			return
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	// Fetch and return the updated schedule
	schedule, err := getScheduleByID(scheduleID)
	if err != nil {
		http.Error(w, "Error fetching updated schedule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}

// DeleteSchedule deletes a scheduled workout
func DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get schedule ID from URL
	vars := mux.Vars(r)
	scheduleID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	// Delete the schedule
	result, err := database.DB.Exec(
		`DELETE FROM schedules WHERE id = $1 AND user_id = $2`,
		scheduleID, claims.UserID,
	)
	if err != nil {
		http.Error(w, "Error deleting schedule", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Schedule not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function to get a schedule by ID
func getScheduleByID(scheduleID int) (*models.Schedule, error) {
	var schedule models.Schedule
	var workout models.Workout
	var completedAt sql.NullTime

	err := database.DB.QueryRow(`
		SELECT s.id, s.user_id, s.workout_id, s.scheduled_date, s.completed, s.completed_at, s.notes,
		       w.id, w.user_id, w.name, w.description, w.created_at, w.updated_at
		FROM schedules s
		JOIN workouts w ON s.workout_id = w.id
		WHERE s.id = $1
	`, scheduleID).Scan(
		&schedule.ID,
		&schedule.UserID,
		&schedule.WorkoutID,
		&schedule.ScheduledDate,
		&schedule.Completed,
		&completedAt,
		&schedule.Notes,
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.Description,
		&workout.CreatedAt,
		&workout.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		schedule.CompletedAt = &completedAt.Time
	}

	schedule.Workout = &workout
	return &schedule, nil
}
