package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo/v4"
)

// Handler to create a new task
func createTask(c echo.Context) error {
	task := new(Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate task
	if task.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Title cannot be empty",
		})
	}

	// Set creation time
	task.CreatedAt = time.Now()

	// Insert into database
	result, err := db.Exec("INSERT INTO tasks (title, completed) VALUES (?, ?)",
		task.Title, task.Completed)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create task: " + err.Error(),
		})
	}

	id, _ := result.LastInsertId()
	task.ID = int(id)

	return c.JSON(http.StatusCreated, task)
}

// Handler to update an existing task
func updateTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	task := new(Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	_, err := db.Exec("UPDATE tasks SET title = ?, completed = ? WHERE id = ?",
		task.Title, task.Completed, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, task)
}

// Handler to delete a task
func deleteTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return c.NoContent(http.StatusNoContent)
}

// Handler to register a new user
func registerUser(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate user data
	if user.Username == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Username and password cannot be empty",
		})
	}

	// Check if username already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
	}

	if count > 0 {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Username already exists",
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to hash password",
		})
	}

	// Insert user into database
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)",
		user.Username, hashedPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}

// Handler to login
func loginUser(c echo.Context) error {
	credentials := new(LoginCredentials)
	if err := c.Bind(credentials); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Check if user exists in DB and password matches
	var user User
	var hashedPassword string
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?",
		credentials.Username).Scan(&user.ID, &user.Username, &hashedPassword)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Generate a new session UUID
	sessionID := uuid.New().String()

	// Store the session in the database
	_, err = db.Exec("INSERT INTO auth_sessions (uuid, user_id) VALUES (?, ?)",
		sessionID, user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
	}

	// Set the session cookie
	cookie := new(http.Cookie)
	cookie.Name = "auth_session"
	cookie.Value = sessionID
	cookie.Path = "/"
	cookie.HttpOnly = true
	// In production, you'd want to set cookie.Secure = true and use HTTPS
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// Handler to logout user
func logoutUser(c echo.Context) error {
	// Get the session cookie
	cookie, err := c.Cookie("auth_session")
	if err == nil {
		// Delete the session from the database
		_, err = db.Exec("DELETE FROM auth_sessions WHERE uuid = ?", cookie.Value)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete session"})
		}

		// Expire the cookie
		cookie.Value = ""
		cookie.Path = "/"
		cookie.Expires = time.Now().Add(-1 * time.Hour) // Expire the cookie
		c.SetCookie(cookie)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// Handler to get current user information
func getCurrentUser(c echo.Context) error {
	// Get the session cookie
	cookie, err := c.Cookie("auth_session")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"authenticated": false,
		})
	}

	// Look up the session in the database
	var userID int
	var username string
	err = db.QueryRow(`
		SELECT u.id, u.username
		FROM auth_sessions a
		JOIN users u ON a.user_id = u.id
		WHERE a.uuid = ?
	`, cookie.Value).Scan(&userID, &username)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"authenticated": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"id":            userID,
		"username":      username,
	})
}
