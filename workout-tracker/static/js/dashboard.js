// API Base URL
const API_URL = 'http://localhost:8080/api';

// Check if user is logged in
const token = localStorage.getItem('token');
if (!token) {
    window.location.href = '/static/login.html';
}

// Display username
document.getElementById('username').textContent = localStorage.getItem('username');

// Logout function
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    window.location.href = '/static/login.html';
}

// Helper function to make authenticated requests
async function fetchWithAuth(url, options = {}) {
    const token = localStorage.getItem('token');
    
    options.headers = {
        ...options.headers,
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
    };
    
    const response = await fetch(url, options);
    
    // If unauthorized, redirect to login
    if (response.status === 401) {
        logout();
    }
    
    return response;
}

// Load exercises
async function loadExercises() {
    try {
        //await : wait for the fetch to complete
        //what if we don't use await here?
        //answer: we would get a Promise object instead of the actual response data
        //promise: an object representing the eventual completion or failure of an asynchronous operation
        const response = await fetch(`${API_URL}/exercises`);
        const data = await response.json();
        
        if (response.ok) {
            // API returns array directly, not wrapped in {exercises: [...]}
            const exercises = Array.isArray(data) ? data : (data.exercises || []);
            displayExercises(exercises);
            populateExerciseDropdown(exercises);
        }
    } catch (error) {
        console.error('Error loading exercises:', error);
        showMessage('Failed to load exercises', 'error');
    }
}

// Display exercises
function displayExercises(exercises) {
    const container = document.getElementById('exercisesList');
    
    if (!exercises || exercises.length === 0) {
        container.innerHTML = '<p>No exercises available</p>';
        return;
    }
    
    container.innerHTML = exercises.map(ex => `
        <div class="exercise-item">
            <h3>${ex.name}</h3>
            <p><strong>Muscle:</strong> ${ex.muscle_group}</p>
            <p>${ex.description}</p>
        </div>
    `).join('');
}

// Populate exercise dropdown
function populateExerciseDropdown(exercises) {
    const select = document.getElementById('workoutExercise');
    
    if (!exercises || exercises.length === 0) {
        select.innerHTML = '<option value="">No exercises available</option>';
        return;
    }
    
    select.innerHTML = exercises.map(ex => 
        `<option value="${ex.id}">${ex.name}</option>`
    ).join('');
}

// Load workouts
async function loadWorkouts() {
    try {
        const response = await fetchWithAuth(`${API_URL}/workouts`);
        const data = await response.json();
        
        if (response.ok) {
            displayWorkouts(data.workouts);
        }
    } catch (error) {
        showMessage('Failed to load workouts', 'error');
    }
}

// Display workouts
function displayWorkouts(workouts) {
    const container = document.getElementById('workoutsList');
    
    if (!workouts || workouts.length === 0) {
        container.innerHTML = '<p>No workouts yet. Create your first workout!</p>';
        return;
    }
    
    container.innerHTML = workouts.map(workout => `
        <div class="workout-item">
            <h3>${workout.name}</h3>
            <p><strong>Exercise:</strong> ${workout.exercise_name}</p>
            <p><strong>Sets:</strong> ${workout.sets} | <strong>Reps:</strong> ${workout.reps}</p>
            <p><small>Created: ${new Date(workout.created_at).toLocaleDateString()}</small></p>
        </div>
    `).join('');
}

// Create workout form
document.getElementById('createWorkoutForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const name = document.getElementById('workoutName').value;
    const exercise_id = parseInt(document.getElementById('workoutExercise').value);
    const sets = parseInt(document.getElementById('sets').value);
    const reps = parseInt(document.getElementById('reps').value);
    
    try {
        const response = await fetchWithAuth(`${API_URL}/workouts`, {
            method: 'POST',
            body: JSON.stringify({ name, exercise_id, sets, reps })
        });
        
        if (response.ok) {
            showMessage('Workout created successfully!', 'success');
            document.getElementById('createWorkoutForm').reset();
            toggleForm('workoutForm');
            loadWorkouts();
        } else {
            const data = await response.json();
            showMessage(data.error || 'Failed to create workout', 'error');
        }
    } catch (error) {
        showMessage('Failed to create workout', 'error');
    }
});

// Toggle form visibility
function toggleForm(formId) {
    const form = document.getElementById(formId);
    form.style.display = form.style.display === 'none' ? 'block' : 'none';
}

// Show message
function showMessage(text, type) {
    const messageEl = document.getElementById('message');
    messageEl.textContent = text;
    messageEl.className = `message ${type}`;
    
    // Clear message after 3 seconds
    setTimeout(() => {
        messageEl.textContent = '';
        messageEl.className = 'message';
    }, 3000);
}

// Load data when page loads
loadExercises();
loadWorkouts();
