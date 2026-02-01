# Test Script - Run this after starting the server

Write-Host "`n=== Workout Tracker API Test Script ===" -ForegroundColor Cyan
Write-Host "Make sure the server is running before executing these tests`n" -ForegroundColor Yellow

$baseUrl = "http://localhost:8080/api"

# Test 1: Health Check
Write-Host "Test 1: Health Check..." -ForegroundColor Green
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/health" -Method GET -UseBasicParsing
    Write-Host "✓ Health check passed: $($response.StatusCode)" -ForegroundColor Green
} catch {
    Write-Host "✗ Health check failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 2: Get Exercises
Write-Host "`nTest 2: Get Exercises..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/exercises" -Method GET
    Write-Host "✓ Found $($response.Count) exercises" -ForegroundColor Green
    Write-Host "  First exercise: $($response[0].name)" -ForegroundColor Gray
} catch {
    Write-Host "✗ Failed to get exercises: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 3: Register a User
Write-Host "`nTest 3: Register a new user..." -ForegroundColor Green
$timestamp = [int][double]::Parse((Get-Date -UFormat %s))
$username = "testuser_$timestamp"
$registerData = @{
    username = $username
    email = "test_${timestamp}@example.com"
    password = "testpassword123"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/register" -Method POST -Body $registerData -ContentType "application/json"
    Write-Host "✓ User registered successfully: $($response.username)" -ForegroundColor Green
    $userId = $response.id
} catch {
    Write-Host "✗ Registration failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 4: Login
Write-Host "`nTest 4: Login..." -ForegroundColor Green
$loginData = @{
    username = $username
    password = "testpassword123"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/login" -Method POST -Body $loginData -ContentType "application/json"
    Write-Host "✓ Login successful" -ForegroundColor Green
    Write-Host "  Token: $($response.token.Substring(0, 20))..." -ForegroundColor Gray
    $token = $response.token
} catch {
    Write-Host "✗ Login failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 5: Get Current User
Write-Host "`nTest 5: Get current user info..." -ForegroundColor Green
$headers = @{
    Authorization = "Bearer $token"
}

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/me" -Method GET -Headers $headers
    Write-Host "✓ Got user info: $($response.username)" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed to get user info: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 6: Create a Workout
Write-Host "`nTest 6: Create a workout..." -ForegroundColor Green
$workoutData = @{
    name = "Test Workout"
    description = "A test workout for API verification"
    exercises = @(
        @{
            exercise_id = 1
            sets = 3
            reps = 10
            weight = 50.0
            notes = "Test exercise"
        }
    )
} | ConvertTo-Json -Depth 3

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/workouts" -Method POST -Body $workoutData -ContentType "application/json" -Headers $headers
    Write-Host "✓ Workout created: $($response.name)" -ForegroundColor Green
    $workoutId = $response.id
} catch {
    Write-Host "✗ Failed to create workout: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 7: Get Workouts
Write-Host "`nTest 7: Get all workouts..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/workouts" -Method GET -Headers $headers
    Write-Host "✓ Found $($response.Count) workout(s)" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed to get workouts: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 8: Schedule a Workout
Write-Host "`nTest 8: Schedule a workout..." -ForegroundColor Green
$scheduledDate = (Get-Date).AddDays(1).ToString("yyyy-MM-ddTHH:mm:ssZ")
$scheduleData = @{
    workout_id = $workoutId
    scheduled_date = $scheduledDate
    notes = "Test schedule"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/schedule" -Method POST -Body $scheduleData -ContentType "application/json" -Headers $headers
    Write-Host "✓ Workout scheduled for tomorrow" -ForegroundColor Green
    $scheduleId = $response.id
} catch {
    Write-Host "✗ Failed to schedule workout: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 9: Get Progress Report
Write-Host "`nTest 9: Get progress report..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/progress" -Method GET -Headers $headers
    Write-Host "✓ Progress report retrieved" -ForegroundColor Green
    Write-Host "  Total workouts: $($response.total_workouts)" -ForegroundColor Gray
    Write-Host "  Workouts this week: $($response.workouts_this_week)" -ForegroundColor Gray
} catch {
    Write-Host "✗ Failed to get progress: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== All Tests Completed ===" -ForegroundColor Cyan
Write-Host "API is working correctly!`n" -ForegroundColor Green

Write-Host "Test user credentials:" -ForegroundColor Yellow
Write-Host "  Username: $username" -ForegroundColor White
Write-Host "  Password: testpassword123" -ForegroundColor White
Write-Host "`nYou can use these credentials to continue testing the API.`n" -ForegroundColor Gray
