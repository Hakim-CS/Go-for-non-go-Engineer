package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"workout-tracker/auth"
	"workout-tracker/database"
	"workout-tracker/handlers"
	"workout-tracker/middleware"
	"workout-tracker/seeder"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	//gorilla mux is a powerful URL router and dispatcher for golang
)

func main() {
	log.Println("Starting Workout Tracker API...")

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration from environment variables
	// In production, you'd use a proper config library
	config := loadConfig()

	// what is the database itself here?
	//answer: The database here refers to the PostgreSQL database that the application connects to for storing and retrieving workout-related data.
	//is that act like a func object,class or package?
	// answer: In this context, "database" is a package that contains functions and types related to database operations, such as connecting to the database, initializing the schema, and managing the database connection.
	//where is the database package defined?
	// answer: The database package is defined in the "workout-tracker/database" directory of the project.

	// Connect to the database
	err = database.Connect(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize database schema (create tables if they don't exist)
	err = database.InitializeSchema()
	if err != nil {
		log.Fatal("Failed to initialize schema:", err)
	}

	// Seed exercises into the database
	err = seeder.SeedExercises(database.DB)
	if err != nil {
		log.Fatal("Failed to seed exercises:", err)
	}

	//mux: NewRouter creates a new router instance to handle incoming HTTP requests.
	//CORS: Cross-Origin Resource Sharing
	//it means that the server allows requests from different origins (domains) to access its resources.
	// This is important for web applications that may be hosted on different domains than the API server.
	//for example, a frontend application running on "http://localhost:3000" might need to make requests to an API server running on "http://localhost:8080".
	// Without proper CORS configuration, the browser would block these requests for security reasons.
	// By implementing CORS middleware, the server can specify which origins are allowed to access its resources, what HTTP methods are permitted, and which headers can be included in the requests.
	// This enables seamless communication between the frontend and backend, enhancing the overall user experience.

	//middleware : CORS middleware is applied to handle cross-origin requests.
	//it's job is to add the necessary CORS headers to the HTTP responses.
	//it act like a function that wraps around the actual request handler to add CORS support.

	// Create a new router
	router := mux.NewRouter()

	// Apply CORS middleware to all routes
	//what actually is happening here?
	// answer: Here, we are applying the CORS middleware to all routes in the router.
	// The middleware function wraps around the actual request handlers to add CORS headers to the HTTP responses.
	// This ensures that all incoming requests are processed with CORS support, allowing cross-origin requests to be handled properly.
	// explaining the code line by line:
	// router.Use(func(next http.Handler) http.Handler { ... }): This line adds a middleware function to the router.
	// The middleware function takes the next handler in the chain as an argument and returns a new handler.
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { ... }): This line defines the new handler function that will be executed for each incoming request.
	// middleware.CORSMiddleware(next.ServeHTTP)(w, r): This line calls the CORS middleware function, passing in the next handler's ServeHTTP method.
	// The CORS middleware processes the request and adds the necessary CORS headers to the response before passing control to the next handler.
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middleware.CORSMiddleware(next.ServeHTTP)(w, r)
		})
	})

	// Set up API routes
	setupRoutes(router)

	// Serve static files (HTML, CSS, JS)
	// The PathPrefix tells the router to handle all URLs starting with /static
	// StripPrefix removes /static from the URL before looking for files
	// http.FileServer serves files from the ./static directory
	fileServer := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// important: this route must be defined after the API routes to avoid conflicts
	// Serve the web UI at the root path
	// Redirect root to static index page
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/index.html", http.StatusFound)
	})

	// Start the server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on http://localhost:%s", port)
	log.Printf("Web UI available at http://localhost:%s/static/index.html", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// setupRoutes configures all the API endpoints
func setupRoutes(router *mux.Router) {
	// API prefix
	//why we use PathPrefix here?
	//  We use PathPrefix to group all API routes under the /api path.

	api := router.PathPrefix("/api").Subrouter()

	// Public routes (no authentication required)
	api.HandleFunc("/register", handlers.Register).Methods("POST")
	api.HandleFunc("/login", handlers.Login).Methods("POST")
	api.HandleFunc("/exercises", handlers.GetExercises).Methods("GET")

	// Protected routes (authentication required)
	// These routes use the AuthMiddleware to verify JWT tokens

	// User routes
	api.HandleFunc("/me", middleware.AuthMiddleware(handlers.GetCurrentUser)).Methods("GET")

	// Workout routes
	api.HandleFunc("/workouts", middleware.AuthMiddleware(handlers.CreateWorkout)).Methods("POST")
	api.HandleFunc("/workouts", middleware.AuthMiddleware(handlers.GetWorkouts)).Methods("GET")
	api.HandleFunc("/workouts/{id}", middleware.AuthMiddleware(handlers.GetWorkout)).Methods("GET")
	api.HandleFunc("/workouts/{id}", middleware.AuthMiddleware(handlers.UpdateWorkout)).Methods("PUT")
	api.HandleFunc("/workouts/{id}", middleware.AuthMiddleware(handlers.DeleteWorkout)).Methods("DELETE")

	// Schedule routes
	api.HandleFunc("/schedule", middleware.AuthMiddleware(handlers.CreateSchedule)).Methods("POST")
	api.HandleFunc("/schedule", middleware.AuthMiddleware(handlers.GetSchedules)).Methods("GET")
	api.HandleFunc("/schedule/{id}/complete", middleware.AuthMiddleware(handlers.CompleteSchedule)).Methods("POST")
	api.HandleFunc("/schedule/{id}", middleware.AuthMiddleware(handlers.DeleteSchedule)).Methods("DELETE")

	// Progress routes
	api.HandleFunc("/progress", middleware.AuthMiddleware(handlers.GetProgress)).Methods("GET")
	api.HandleFunc("/progress/exercise", middleware.AuthMiddleware(handlers.GetExerciseHistory)).Methods("GET")

	// Health check endpoint
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
}

// loadConfig loads configuration from environment variables
// This provides default values if environment variables are not set
//Go doesn't automatically read .env files - you need a library like godotenv for that
//That's why you had to manually set

func loadConfig() database.Config {
	// Get database configuration
	host := getEnv("DB_HOST", "localhost")
	portStr := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "workout")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Convert port string to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Invalid DB_PORT, using default: 5432")
		port = 5432
	}

	// Set JWT configuration
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-this-in-production")
	auth.JWTSecret = []byte(jwtSecret)

	jwtExpirationStr := getEnv("JWT_EXPIRATION_HOURS", "24")
	jwtExpiration, err := strconv.Atoi(jwtExpirationStr)
	if err != nil {
		log.Printf("Invalid JWT_EXPIRATION_HOURS, using default: 24")
		jwtExpiration = 24
	}
	auth.TokenExpirationHours = jwtExpiration

	return database.Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
		SSLMode:  sslmode,
	}
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
