package handlers

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"vibecoders/models"
)

type Handler struct {
	DB *sql.DB
}

// Auth related structs
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Bio        string `json:"bio"`
	LinkedinURL string `json:"linkedin_url"`
	GithubURL  string `json:"github_url"`
	PhotoURL   string `json:"photo_url"`
}

type UpdateUserRequest struct {
	Bio        string `json:"bio"`
	LinkedinURL string `json:"linkedin_url"`
	GithubURL  string `json:"github_url"`
	PhotoURL   string `json:"photo_url"`
}

type AuthResponse struct {
	Token  string       `json:"token"`
	User   *models.User `json:"user"`
}

// RegisterUser handles user registration
func (h *Handler) RegisterUser(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate request
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username and password are required"})
	}

	// Create user
	user, err := models.CreateUser(h.DB, req.Username, req.Password, req.Bio, req.LinkedinURL, req.GithubURL, req.PhotoURL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Generate token
	token := uuid.New().String()

	return c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

// LoginUser handles user authentication
func (h *Handler) LoginUser(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate request
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username and password are required"})
	}

	// Get user
	user, err := models.GetUserByUsername(h.DB, req.Username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Verify password
	if !models.VerifyPassword(user, req.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Generate token
	token := uuid.New().String()

	return c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

// UpdateUserProfile handles updating user profile
func (h *Handler) UpdateUserProfile(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Update user
	user, err := models.UpdateUser(h.DB, userID, req.Bio, req.LinkedinURL, req.GithubURL, req.PhotoURL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// GetUserProfile retrieves user profile
func (h *Handler) GetUserProfile(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}

	// Get user
	user, err := models.GetUserByID(h.DB, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}