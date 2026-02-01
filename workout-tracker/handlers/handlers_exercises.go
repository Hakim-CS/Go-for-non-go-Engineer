package handlers

import (
	"encoding/json"
	"net/http"
	"workout-tracker/database"
	"workout-tracker/models"
)

// GetExercises returns all available exercises from the library
// This is a public endpoint - no authentication required
func GetExercises(w http.ResponseWriter, r *http.Request) {
	// Query all exercises from the database
	rows, err := database.DB.Query(`
		SELECT id, name, description, category, muscle_group 
		FROM exercises 
		ORDER BY category, name
	`)
	if err != nil {
		http.Error(w, "Error fetching exercises", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Build a slice to hold all exercises
	exercises := []models.Exercise{}

	// Loop through each row and scan into an Exercise struct
	for rows.Next() {
		var exercise models.Exercise
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category,
			&exercise.MuscleGroup,
		)
		if err != nil {
			http.Error(w, "Error scanning exercise", http.StatusInternalServerError)
			return
		}
		exercises = append(exercises, exercise)
	}

	// Check for any errors during iteration
	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating exercises", http.StatusInternalServerError)
		return
	}

	// Return the exercises as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercises)
}
