package seeder

import (
	"database/sql"
	"log"
)

// Exercise represents the exercise data we want to seed
type Exercise struct {
	Name        string
	Description string
	Category    string
	MuscleGroup string
}

// predefinedExercises contains our starter set of exercises
// This gives users a library to choose from when creating workouts
var predefinedExercises = []Exercise{
	// Chest exercises
	{
		Name:        "Bench Press",
		Description: "Lie on a bench and press a barbell or dumbbells upward",
		Category:    "Chest",
		MuscleGroup: "Upper Body",
	},
	{
		Name:        "Push-ups",
		Description: "Classic bodyweight exercise for chest, shoulders, and triceps",
		Category:    "Chest",
		MuscleGroup: "Upper Body",
	},
	{
		Name:        "Dumbbell Flyes",
		Description: "Lying on bench, move dumbbells in an arc motion",
		Category:    "Chest",
		MuscleGroup: "Upper Body",
	},

	// Back exercises
	{
		Name:        "Deadlift",
		Description: "Lift a barbell from the ground to hip level",
		Category:    "Back",
		MuscleGroup: "Full Body",
	},
	{
		Name:        "Pull-ups",
		Description: "Hang from a bar and pull yourself up",
		Category:    "Back",
		MuscleGroup: "Upper Body",
	},
	{
		Name:        "Barbell Rows",
		Description: "Bent over position, pull barbell to your chest",
		Category:    "Back",
		MuscleGroup: "Upper Body",
	},

	// Leg exercises
	{
		Name:        "Squats",
		Description: "Lower your body by bending knees and hips",
		Category:    "Legs",
		MuscleGroup: "Lower Body",
	},
	{
		Name:        "Lunges",
		Description: "Step forward and lower your body until both knees are bent",
		Category:    "Legs",
		MuscleGroup: "Lower Body",
	},
	{
		Name:        "Leg Press",
		Description: "Push weight away using leg press machine",
		Category:    "Legs",
		MuscleGroup: "Lower Body",
	},

	// Shoulder exercises
	{
		Name:        "Shoulder Press",
		Description: "Press dumbbells or barbell overhead",
		Category:    "Shoulders",
		MuscleGroup: "Upper Body",
	},
	{
		Name:        "Lateral Raises",
		Description: "Raise dumbbells to the sides",
		Category:    "Shoulders",
		MuscleGroup: "Upper Body",
	},

	// Arm exercises
	{
		Name:        "Bicep Curls",
		Description: "Curl dumbbells or barbell toward shoulders",
		Category:    "Arms",
		MuscleGroup: "Upper Body",
	},
	{
		Name:        "Tricep Dips",
		Description: "Lower and raise your body using parallel bars",
		Category:    "Arms",
		MuscleGroup: "Upper Body",
	},

	// Core exercises
	{
		Name:        "Plank",
		Description: "Hold a push-up position on your forearms",
		Category:    "Core",
		MuscleGroup: "Core",
	},
	{
		Name:        "Crunches",
		Description: "Lie on back and lift shoulders toward knees",
		Category:    "Core",
		MuscleGroup: "Core",
	},
	{
		Name:        "Russian Twists",
		Description: "Seated position, rotate torso side to side",
		Category:    "Core",
		MuscleGroup: "Core",
	},
}

// SeedExercises populates the exercises table with predefined exercises
// This only adds exercises if the table is empty
func SeedExercises(db *sql.DB) error {
	log.Println("Checking if exercises need to be seeded...")

	// First, check if we already have exercises
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM exercises").Scan(&count)
	if err != nil {
		return err
	}

	// If exercises already exist, skip seeding
	if count > 0 {
		log.Printf("Exercises already seeded (%d exercises found). Skipping...", count)
		return nil
	}

	// Insert each exercise into the database
	log.Println("Seeding exercises...")
	for _, exercise := range predefinedExercises {
		_, err := db.Exec(
			`INSERT INTO exercises (name, description, category, muscle_group) 
			 VALUES ($1, $2, $3, $4)`,
			exercise.Name,
			exercise.Description,
			exercise.Category,
			exercise.MuscleGroup,
		)
		if err != nil {
			return err
		}
	}

	log.Printf("Successfully seeded %d exercises", len(predefinedExercises))
	return nil
}
