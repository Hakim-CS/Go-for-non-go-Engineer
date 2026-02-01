# API Testing Examples

This file contains example requests you can use to test the API using tools like Postman, cURL, or any HTTP client.

## 1. Register a New User

**POST** `/api/register`

```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password123"
}
```

**Expected Response (201 Created):**
```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "created_at": "2024-01-15T10:00:00Z"
}
```

## 2. Login

**POST** `/api/login`

```json
{
  "username": "john_doe",
  "password": "secure_password123"
}
```

**Expected Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

**Important:** Save the token! You'll need it for all protected endpoints.

## 3. Get All Exercises

**GET** `/api/exercises`

No authentication required.

**Expected Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Bench Press",
    "description": "Lie on a bench and press a barbell or dumbbells upward",
    "category": "Chest",
    "muscle_group": "Upper Body"
  },
  ...
]
```

## 4. Create a Workout

**POST** `/api/workouts`

**Headers:**
```
Authorization: Bearer YOUR_TOKEN_HERE
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Monday Chest Day",
  "description": "Upper body workout focusing on chest",
  "exercises": [
    {
      "exercise_id": 1,
      "sets": 4,
      "reps": 10,
      "weight": 60.0,
      "notes": "Warm up first"
    },
    {
      "exercise_id": 2,
      "sets": 3,
      "reps": 15,
      "weight": 0,
      "notes": "Bodyweight push-ups"
    }
  ]
}
```

## 5. Get All Your Workouts

**GET** `/api/workouts`

**Headers:**
```
Authorization: Bearer YOUR_TOKEN_HERE
```

## 6. Get a Specific Workout

**GET** `/api/workouts/1`

**Headers:**
```
Authorization: Bearer YOUR_TOKEN_HERE
```

## 7. Update a Workout

**PUT** `/api/workouts/1`

**Headers:**
```
Authorization: Bearer YOUR_TOKEN_HERE
Content-Type: application/json
```

**Body:** (Same format as creating a workout)

## 8. Delete a Workout

**DELETE** `/api/workouts/1`

**Headers:**
```
Authorization: Bearer YOUR_TOKEN_HERE
```

## 9. Schedule a Workout

**POST** `/api/schedule`

**Headers:**
```
Authorization: Bearer YOUR_TOKEN_HERE
Content-Type: application/json
```

**Body:**
```json
{
  "workout_id": 1,
  "scheduled_date": "2024-01-20T09:00:00Z",
  "notes": "Morning workout before work"
}
```

## 10. Get Scheduled Workouts

**GET** `/api/schedule`

**Headers:**
```
Authorization: Bearer YOUR_TOKEN_HERE
```

**Optional Query Parameters:**
- `upcoming=true` - Show only future workouts
- `completed=true` - Show only completed workouts
- `completed=false` - Show only pending workouts

Examples:
- `/api/schedule?upcoming=true`
- `/api/schedule?completed=false`

## 11. Mark Workout as Complete

**POST** `/api/schedule/1/complete`

**Headers:**
```
Authorization: Bearer YOUR_TOKEN_HERE
Content-Type: application/json
```

**Body:**
```json
{
  "notes": "Great workout! Felt strong today.",
  "logs": [
    {
      "exercise_id": 1,
      "sets_completed": 4,
      "reps_completed": 10,
      "weight_used": 60.0,
      "duration": 15,
      "notes": "Last set was challenging"
    },
    {
      "exercise_id": 2,
      "sets_completed": 3,
      "reps_completed": 15,
      "weight_used": 0,
      "duration": 10,
      "notes": "Smooth execution"
    }
  ]
}
```

## 12. Get Progress Report

**GET** `/api/progress`

**Headers:**
```
Authorization: Bearer YOUR_TOKEN_HERE
```

**Expected Response (200 OK):**
```json
{
  "user_id": 1,
  "total_workouts": 25,
  "total_exercises": 150,
  "workouts_this_week": 3,
  "workouts_this_month": 12,
  "most_frequent_exercise": "Bench Press",
  "average_workouts_per_week": 3.5,
  "start_date": "2023-12-01T10:00:00Z"
}
```

## 13. Get Exercise History

**GET** `/api/progress/exercise?exercise_id=1`

**Headers:**
```
Authorization: Bearer YOUR_TOKEN_HERE
```

Shows your performance history for a specific exercise over time.

## cURL Examples

### Register:
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john_doe","email":"john@example.com","password":"secure_password123"}'
```

### Login:
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john_doe","password":"secure_password123"}'
```

### Get Exercises:
```bash
curl http://localhost:8080/api/exercises
```

### Create Workout (replace TOKEN with your actual token):
```bash
curl -X POST http://localhost:8080/api/workouts \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Monday Chest Day","description":"Upper body workout","exercises":[{"exercise_id":1,"sets":4,"reps":10,"weight":60.0}]}'
```

### Get Your Workouts:
```bash
curl http://localhost:8080/api/workouts \
  -H "Authorization: Bearer TOKEN"
```

### Get Progress:
```bash
curl http://localhost:8080/api/progress \
  -H "Authorization: Bearer TOKEN"
```
