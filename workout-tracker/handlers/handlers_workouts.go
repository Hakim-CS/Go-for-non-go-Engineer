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

// CreateWorkoutRequest represents the data needed to create a workout
type CreateWorkoutRequest struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	Exercises   []CreateWorkoutExerciseRequest `json:"exercises"`
}

// CreateWorkoutExerciseRequest represents an exercise in a workout
type CreateWorkoutExerciseRequest struct {
	ExerciseID int     `json:"exercise_id"`
	Sets       int     `json:"sets"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
	Notes      string  `json:"notes"`
}

// CreateWorkout creates a new workout for the logged-in user
func CreateWorkout(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body
	var req CreateWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Name == "" {
		http.Error(w, "Workout name is required", http.StatusBadRequest)
		return
	}

	// Start a database transaction
	// This ensures that either all changes succeed or none do
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}
	// If anything goes wrong, rollback the transaction
	defer tx.Rollback()

	// Insert the workout
	var workoutID int
	err = tx.QueryRow(
		`INSERT INTO workouts (user_id, name, description, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		claims.UserID, req.Name, req.Description, time.Now(), time.Now(),
	).Scan(&workoutID)

	if err != nil {
		http.Error(w, "Error creating workout", http.StatusInternalServerError)
		return
	}

	// Insert each exercise in the workout
	for _, exercise := range req.Exercises {
		_, err = tx.Exec(
			`INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight, notes) 
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			workoutID, exercise.ExerciseID, exercise.Sets, exercise.Reps, exercise.Weight, exercise.Notes,
		)
		if err != nil {
			http.Error(w, "Error adding exercise to workout", http.StatusInternalServerError)
			return
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	// Fetch and return the complete workout
	workout, err := getWorkoutByID(workoutID, claims.UserID)
	if err != nil {
		http.Error(w, "Error fetching created workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workout)
}

// GetWorkouts returns all workouts for the logged-in user
func GetWorkouts(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Query workouts for this user
	rows, err := database.DB.Query(`
		SELECT id, user_id, name, description, created_at, updated_at 
		FROM workouts 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`, claims.UserID)

	if err != nil {
		http.Error(w, "Error fetching workouts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	workouts := []models.Workout{}
	for rows.Next() {
		var workout models.Workout
		err := rows.Scan(
			&workout.ID,
			&workout.UserID,
			&workout.Name,
			&workout.Description,
			&workout.CreatedAt,
			&workout.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "Error scanning workout", http.StatusInternalServerError)
			return
		}

		// Load exercises for this workout
		exercises, err := getWorkoutExercises(workout.ID)
		if err != nil {
			http.Error(w, "Error loading workout exercises", http.StatusInternalServerError)
			return
		}
		workout.Exercises = exercises

		workouts = append(workouts, workout)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

// GetWorkout returns a specific workout by ID
func GetWorkout(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get workout ID from URL
	vars := mux.Vars(r)
	workoutID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	// Fetch the workout
	workout, err := getWorkoutByID(workoutID, claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Workout not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error fetching workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// UpdateWorkout updates an existing workout
func UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get workout ID from URL
	vars := mux.Vars(r)
	workoutID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var req CreateWorkoutRequest
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

	// Update the workout
	result, err := tx.Exec(
		`UPDATE workouts SET name = $1, description = $2, updated_at = $3 
		 WHERE id = $4 AND user_id = $5`,
		req.Name, req.Description, time.Now(), workoutID, claims.UserID,
	)
	if err != nil {
		http.Error(w, "Error updating workout", http.StatusInternalServerError)
		return
	}

	// Check if the workout was found
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	// Delete existing exercises
	_, err = tx.Exec(`DELETE FROM workout_exercises WHERE workout_id = $1`, workoutID)
	if err != nil {
		http.Error(w, "Error removing old exercises", http.StatusInternalServerError)
		return
	}

	// Insert new exercises
	for _, exercise := range req.Exercises {
		_, err = tx.Exec(
			`INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight, notes) 
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			workoutID, exercise.ExerciseID, exercise.Sets, exercise.Reps, exercise.Weight, exercise.Notes,
		)
		if err != nil {
			http.Error(w, "Error adding exercise to workout", http.StatusInternalServerError)
			return
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	// Fetch and return the updated workout
	workout, err := getWorkoutByID(workoutID, claims.UserID)
	if err != nil {
		http.Error(w, "Error fetching updated workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// DeleteWorkout deletes a workout
func DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get workout ID from URL
	vars := mux.Vars(r)
	workoutID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	// Delete the workout (cascades to workout_exercises due to foreign key)
	result, err := database.DB.Exec(
		`DELETE FROM workouts WHERE id = $1 AND user_id = $2`,
		workoutID, claims.UserID,
	)
	if err != nil {
		http.Error(w, "Error deleting workout", http.StatusInternalServerError)
		return
	}

	// Check if the workout was found
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function to get a workout by ID
func getWorkoutByID(workoutID, userID int) (*models.Workout, error) {
	var workout models.Workout
	err := database.DB.QueryRow(`
		SELECT id, user_id, name, description, created_at, updated_at 
		FROM workouts 
		WHERE id = $1 AND user_id = $2
	`, workoutID, userID).Scan(
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

	// Load exercises
	exercises, err := getWorkoutExercises(workoutID)
	if err != nil {
		return nil, err
	}
	workout.Exercises = exercises

	return &workout, nil
}

// Helper function to get exercises for a workout
func getWorkoutExercises(workoutID int) ([]models.WorkoutExercise, error) {
	rows, err := database.DB.Query(`
		SELECT we.id, we.workout_id, we.exercise_id, we.sets, we.reps, we.weight, we.notes,
		       e.id, e.name, e.description, e.category, e.muscle_group
		FROM workout_exercises we
		JOIN exercises e ON we.exercise_id = e.id
		WHERE we.workout_id = $1
	`, workoutID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := []models.WorkoutExercise{}
	for rows.Next() {
		var we models.WorkoutExercise
		var exercise models.Exercise

		err := rows.Scan(
			&we.ID,
			&we.WorkoutID,
			&we.ExerciseID,
			&we.Sets,
			&we.Reps,
			&we.Weight,
			&we.Notes,
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category,
			&exercise.MuscleGroup,
		)
		if err != nil {
			return nil, err
		}

		we.Exercise = &exercise
		exercises = append(exercises, we)
	}

	return exercises, nil
}
