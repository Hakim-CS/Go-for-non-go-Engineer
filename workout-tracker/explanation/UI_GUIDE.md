# Simple UI Guide

## 📁 File Structure

```
static/
├── index.html          # Home page (landing page)
├── login.html          # Login form
├── register.html       # Registration form
├── dashboard.html      # Main app interface (after login)
├── css/
│   └── style.css      # All styling
└── js/
    ├── auth.js        # Login & Register logic
    └── dashboard.js   # Dashboard functionality
```

---

## 🎯 How It Works

### **1. Landing Page (index.html)**
- First page users see
- Has two buttons: Login and Register
- Simple welcome screen

### **2. Authentication (login.html & register.html)**
- Forms to collect user credentials
- `auth.js` handles sending data to your Go API
- On success, saves JWT token to `localStorage`
- Redirects to dashboard

### **3. Dashboard (dashboard.html)**
- Shows after successful login
- Protected: redirects to login if no token
- Displays:
  - Available exercises
  - User's workouts
  - Form to create new workouts

---

## 🔑 Key Concepts

### **localStorage**
```javascript
localStorage.setItem('token', data.token);  // Save
localStorage.getItem('token');              // Read
localStorage.removeItem('token');           // Delete
```
- Stores JWT token in browser
- Persists even after closing browser
- Used for authentication

### **fetch API**
```javascript
fetch('http://localhost:8080/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
})
```
- Makes HTTP requests to your Go backend
- Returns promises (async/await)
- Handles JSON data

### **Authorization Header**
```javascript
'Authorization': `Bearer ${token}`
```
- Sends JWT token with protected requests
- Backend verifies token
- Grants access to user data

---

## 🚀 How to Use

### **1. Start Your Go Server**
```powershell
go run main.go
```

### **2. Open Browser**
```
http://localhost:8080
```

### **3. Try It Out**
1. Click "Register"
2. Create an account
3. You'll be redirected to dashboard
4. See exercises
5. Create a workout

---

## 📝 File Explanations

### **auth.js**
- Handles login and register forms
- Submits data to `/api/login` and `/api/register`
- Saves token when successful
- Shows error messages

### **dashboard.js**
- Checks if user is logged in (has token)
- Loads exercises from `/api/exercises`
- Loads workouts from `/api/workouts`
- Creates new workouts
- Includes `fetchWithAuth()` helper for authenticated requests

### **style.css**
- Modern, clean design
- Purple gradient background
- Responsive forms
- Styled buttons and cards

---

## 🔒 Authentication Flow

```
1. User registers/logs in
   ↓
2. Backend returns JWT token
   ↓
3. Frontend saves token to localStorage
   ↓
4. All protected API calls include token
   ↓
5. Backend verifies token
   ↓
6. Returns user-specific data
```

---

## 🛠️ How Go Serves Static Files

In `main.go`:
```go
// Serve files from ./static directory
fileServer := http.FileServer(http.Dir("./static"))
router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
```

**What this does:**
- `/static/index.html` → serves `./static/index.html`
- `/static/css/style.css` → serves `./static/css/style.css`
- `/static/js/auth.js` → serves `./static/js/auth.js`

**StripPrefix:**
- Removes `/static/` from URL before looking for file
- Request: `/static/index.html` → File: `./static/index.html`

---

## 🎨 Customization

### **Change Colors**
Edit `style.css`:
```css
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
/* Change to your preferred gradient */
```

### **Add More Features**
1. Create new HTML page
2. Add JavaScript file for logic
3. Link in main.go (or access via `/static/`)

### **Modify Forms**
Edit HTML files to add/remove fields, then update corresponding JavaScript

---

## 🐛 Troubleshooting

### **404 Not Found**
- Make sure server is running: `go run main.go`
- Check URL starts with `/static/`

### **CORS Errors**
- Your Go backend already has CORS enabled
- Check browser console for specific errors

### **Token Invalid**
- Logout and login again
- Token might be expired (24 hours default)

### **Can't Connect to API**
- Verify Go server is running on port 8080
- Check API_URL in JavaScript files matches your server

---

## 📚 What You Learned

1. **Frontend talks to Backend via API**
   - JavaScript fetch → Go handlers → Database

2. **Authentication with JWT**
   - Token stored in browser
   - Sent with each request
   - Backend verifies

3. **Static File Serving**
   - Go can serve HTML/CSS/JS
   - FileServer handles file requests

4. **Separation of Concerns**
   - HTML: Structure
   - CSS: Styling
   - JavaScript: Logic
   - Go: API/Database

---

## 🎯 Next Steps

- Add more features (edit workout, delete workout)
- Add schedule page
- Add progress tracking page
- Improve styling
- Add form validation
- Add loading spinners

Keep it simple and build one feature at a time!
