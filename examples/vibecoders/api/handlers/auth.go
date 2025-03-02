package handlers

import (
	"database/sql"
	"net/http"

	"vibecoders/models"

	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Bio             string `json:"bio"`
	LinkedInURL     string `json:"linked_in_url"`
	GithubURL       string `json:"github_url"`
	PhotoURL        string `json:"photo_url"`
}

type UpdateUserRequest struct {
	Bio         string `json:"bio"`
	LinkedInURL string `json:"linked_in_url"`
	GithubURL   string `json:"github_url"`
	PhotoURL    string `json:"photo_url"`
}

func Login(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		user, err := models.GetUserByUsername(db, req.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}

		// In a real app, we would use password hashing
		if user.Password != req.Password {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}

		// Create session
		token, err := models.CreateSession(db, user.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create session"})
		}

		// Set cookie
		cookie := new(http.Cookie)
		cookie.Name = "session_token"
		cookie.Value = token
		cookie.Path = "/"
		cookie.HttpOnly = true
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Login successful",
			"user": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
			},
		})
	}
}

func Logout(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusOK, map[string]string{"message": "Already logged out"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
		}

		// Delete session from database
		if err := models.DeleteSession(db, cookie.Value); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete session"})
		}

		// Clear cookie
		cookie = new(http.Cookie)
		cookie.Name = "session_token"
		cookie.Value = ""
		cookie.Path = "/"
		cookie.MaxAge = -1
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
	}
}

func Register(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req RegisterRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Validate passwords match
		if req.Password != req.ConfirmPassword {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Passwords do not match"})
		}

		// Check if username already exists
		_, err := models.GetUserByUsername(db, req.Username)
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Username already exists"})
		} else if err != sql.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}

		// Create user
		err = models.CreateUser(db, req.Username, req.Password, "", req.Bio, req.LinkedInURL, req.GithubURL, req.PhotoURL)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create user"})
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully"})
	}
}

func UpdateUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get user ID from session
		cookie, err := c.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Not logged in"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
		}

		userID, err := models.GetUserIDByToken(db, cookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid session"})
		}

		var req UpdateUserRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Update user
		err = models.UpdateUser(db, userID, "", req.Bio, req.LinkedInURL, req.GithubURL, req.PhotoURL)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update user"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User updated successfully"})
	}
}

func GetHomepageUsers(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := models.GetTopUsers(db, 3)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch users"})
		}

		return c.JSON(http.StatusOK, users)
	}
}

func GetCurrentUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Not logged in"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
		}

		userID, err := models.GetUserIDByToken(db, cookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid session"})
		}

		user, err := models.GetUserByID(db, userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch user"})
		}

		return c.JSON(http.StatusOK, user)
	}
}

func GetPublicUserByUsername(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")
		if username == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username is required"})
		}

		user, err := models.GetUserByUsername(db, username)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch user"})
		}

		// Remove sensitive information for public profile
		user.Password = ""
		
		return c.JSON(http.StatusOK, user)
	}
}
