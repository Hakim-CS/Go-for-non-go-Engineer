package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"workout-tracker/database"
	"workout-tracker/middleware"
	"workout-tracker/models"
)

// GetProgress returns a progress report for the logged-in user
// This includes statistics about workouts completed, frequency, and trends
func GetProgress(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Initialize the progress report
	report := models.ProgressReport{
		UserID: claims.UserID,
	}

	// Get total completed workouts
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM schedules WHERE user_id = $1 AND completed = true
	`, claims.UserID).Scan(&report.TotalWorkouts)
	if err != nil {
		http.Error(w, "Error fetching total workouts", http.StatusInternalServerError)
		return
	}

	// Get total exercises performed (from workout logs)
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM workout_logs wl
		JOIN schedules s ON wl.schedule_id = s.id
		WHERE s.user_id = $1
	`, claims.UserID).Scan(&report.TotalExercises)
	if err != nil {
		http.Error(w, "Error fetching total exercises", http.StatusInternalServerError)
		return
	}

	// Get workouts completed in the last 7 days
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM schedules 
		WHERE user_id = $1 AND completed = true AND completed_at >= $2
	`, claims.UserID, sevenDaysAgo).Scan(&report.WorkoutsThisWeek)
	if err != nil {
		http.Error(w, "Error fetching weekly workouts", http.StatusInternalServerError)
		return
	}

	// Get workouts completed in the last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM schedules 
		WHERE user_id = $1 AND completed = true AND completed_at >= $2
	`, claims.UserID, thirtyDaysAgo).Scan(&report.WorkoutsThisMonth)
	if err != nil {
		http.Error(w, "Error fetching monthly workouts", http.StatusInternalServerError)
		return
	}

	// Get the most frequently performed exercise
	var exerciseName sql.NullString
	err = database.DB.QueryRow(`
		SELECT e.name 
		FROM workout_logs wl
		JOIN exercises e ON wl.exercise_id = e.id
		JOIN schedules s ON wl.schedule_id = s.id
		WHERE s.user_id = $1
		GROUP BY e.name
		ORDER BY COUNT(*) DESC
		LIMIT 1
	`, claims.UserID).Scan(&exerciseName)

	// If no exercises logged yet, this is okay
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Error fetching most frequent exercise", http.StatusInternalServerError)
		return
	}

	if exerciseName.Valid {
		report.MostFrequentExercise = exerciseName.String
	} else {
		report.MostFrequentExercise = "None yet"
	}

	// Get the user's start date (when they first completed a workout)
	var startDate sql.NullTime
	err = database.DB.QueryRow(`
		SELECT MIN(completed_at) FROM schedules 
		WHERE user_id = $1 AND completed = true
	`, claims.UserID).Scan(&startDate)

	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Error fetching start date", http.StatusInternalServerError)
		return
	}

	if startDate.Valid {
		report.StartDate = startDate.Time

		// Calculate average workouts per week
		weeksSinceStart := time.Since(report.StartDate).Hours() / 24 / 7
		if weeksSinceStart > 0 {
			report.AverageWorkoutsPerWeek = float64(report.TotalWorkouts) / weeksSinceStart
		}
	} else {
		// If no workouts completed yet, use current time as start date
		report.StartDate = time.Now()
		report.AverageWorkoutsPerWeek = 0
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetExerciseHistory returns the history of a specific exercise for the user
// This shows how the user's performance has changed over time
func GetExerciseHistory(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get exercise ID from query parameter
	exerciseIDStr := r.URL.Query().Get("exercise_id")
	if exerciseIDStr == "" {
		http.Error(w, "exercise_id parameter is required", http.StatusBadRequest)
		return
	}

	// Query the workout logs for this exercise
	rows, err := database.DB.Query(`
		SELECT wl.sets_completed, wl.reps_completed, wl.weight_used, wl.duration, wl.notes, wl.logged_at,
		       e.name, e.description, e.category
		FROM workout_logs wl
		JOIN exercises e ON wl.exercise_id = e.id
		JOIN schedules s ON wl.schedule_id = s.id
		WHERE s.user_id = $1 AND wl.exercise_id = $2
		ORDER BY wl.logged_at DESC
		LIMIT 50
	`, claims.UserID, exerciseIDStr)

	if err != nil {
		http.Error(w, "Error fetching exercise history", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a custom response structure
	type ExerciseHistoryEntry struct {
		SetsCompleted int       `json:"sets_completed"`
		RepsCompleted int       `json:"reps_completed"`
		WeightUsed    float64   `json:"weight_used"`
		Duration      int       `json:"duration"`
		Notes         string    `json:"notes"`
		LoggedAt      time.Time `json:"logged_at"`
		ExerciseName  string    `json:"exercise_name"`
		Category      string    `json:"category"`
	}

	history := []ExerciseHistoryEntry{}
	for rows.Next() {
		var entry ExerciseHistoryEntry
		var description string // We query it but don't use it in the response

		err := rows.Scan(
			&entry.SetsCompleted,
			&entry.RepsCompleted,
			&entry.WeightUsed,
			&entry.Duration,
			&entry.Notes,
			&entry.LoggedAt,
			&entry.ExerciseName,
			&description,
			&entry.Category,
		)
		if err != nil {
			http.Error(w, "Error scanning history entry", http.StatusInternalServerError)
			return
		}

		history = append(history, entry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
