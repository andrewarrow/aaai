package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"vibecoders/models"

	"github.com/labstack/echo/v4"
)

type PromptRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// GetUserPrompts retrieves all prompts for the current user
func GetUserPrompts(db *sql.DB) echo.HandlerFunc {
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

		prompts, err := models.GetPromptsByUserID(db, userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch prompts"})
		}

		return c.JSON(http.StatusOK, prompts)
	}
}

// CreatePrompt adds a new prompt for the current user
func CreatePrompt(db *sql.DB) echo.HandlerFunc {
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

		var req PromptRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Validate required fields
		if req.Title == "" || req.Content == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title and content are required"})
		}

		// Create the prompt
		promptID, err := models.CreatePrompt(db, userID, req.Title, req.Content, req.Tags)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create prompt"})
		}

		prompt, err := models.GetPromptByID(db, promptID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch created prompt"})
		}

		return c.JSON(http.StatusCreated, prompt)
	}
}

// UpdatePrompt modifies an existing prompt
func UpdatePrompt(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get prompt ID from URL parameter
		promptIDStr := c.Param("id")
		promptID, err := strconv.Atoi(promptIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid prompt ID"})
		}

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

		// Check if prompt exists and belongs to the user
		prompt, err := models.GetPromptByID(db, promptID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Prompt not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch prompt"})
		}

		if prompt.UserID != userID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "You don't have permission to update this prompt"})
		}

		var req PromptRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Validate required fields
		if req.Title == "" || req.Content == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title and content are required"})
		}

		// Update the prompt
		err = models.UpdatePrompt(db, promptID, userID, req.Title, req.Content, req.Tags)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update prompt"})
		}

		updatedPrompt, err := models.GetPromptByID(db, promptID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch updated prompt"})
		}

		return c.JSON(http.StatusOK, updatedPrompt)
	}
}

// DeletePrompt removes a prompt
func DeletePrompt(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get prompt ID from URL parameter
		promptIDStr := c.Param("id")
		promptID, err := strconv.Atoi(promptIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid prompt ID"})
		}

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

		// Check if prompt exists and belongs to the user
		prompt, err := models.GetPromptByID(db, promptID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Prompt not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch prompt"})
		}

		if prompt.UserID != userID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "You don't have permission to delete this prompt"})
		}

		// Delete the prompt
		err = models.DeletePrompt(db, promptID, userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete prompt"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Prompt deleted successfully"})
	}
}

// GetUserPublicPrompts retrieves public prompts for a specific user by username
func GetUserPublicPrompts(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")
		if username == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username is required"})
		}

		prompts, err := models.GetUserPublicPromptsByUsername(db, username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch prompts"})
		}

		return c.JSON(http.StatusOK, prompts)
	}
}