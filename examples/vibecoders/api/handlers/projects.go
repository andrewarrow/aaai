package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"vibecoders/models"

	"github.com/labstack/echo/v4"
)

type ProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	GithubURL   string `json:"github_url"`
	WebsiteURL  string `json:"website_url"`
	ImageURL1   string `json:"image_url1"`
	ImageURL2   string `json:"image_url2"`
	ImageURL3   string `json:"image_url3"`
}

// GetUserProjects retrieves all projects for the current user
func GetUserProjects(db *sql.DB) echo.HandlerFunc {
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

		projects, err := models.GetProjectsByUserID(db, userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch projects"})
		}

		return c.JSON(http.StatusOK, projects)
	}
}

// CreateProject adds a new project for the current user
func CreateProject(db *sql.DB) echo.HandlerFunc {
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

		var req ProjectRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Validate required fields
		if req.Title == "" || req.Description == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title and description are required"})
		}

		// Create the project
		projectID, err := models.CreateProject(db, userID, req.Title, req.Description, req.GithubURL, req.WebsiteURL,
			req.ImageURL1, req.ImageURL2, req.ImageURL3)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create project"})
		}

		project, err := models.GetProjectByID(db, projectID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch created project"})
		}

		return c.JSON(http.StatusCreated, project)
	}
}

// UpdateProject modifies an existing project
func UpdateProject(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get project ID from URL parameter
		projectIDStr := c.Param("id")
		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
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

		// Check if project exists and belongs to the user
		project, err := models.GetProjectByID(db, projectID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Project not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch project"})
		}

		if project.UserID != userID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "You don't have permission to update this project"})
		}

		var req ProjectRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Validate required fields
		if req.Title == "" || req.Description == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title and description are required"})
		}

		// Update the project
		err = models.UpdateProject(db, projectID, userID, req.Title, req.Description, req.GithubURL, req.WebsiteURL,
			req.ImageURL1, req.ImageURL2, req.ImageURL3)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update project"})
		}

		updatedProject, err := models.GetProjectByID(db, projectID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch updated project"})
		}

		return c.JSON(http.StatusOK, updatedProject)
	}
}

// DeleteProject removes a project
func DeleteProject(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get project ID from URL parameter
		projectIDStr := c.Param("id")
		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
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

		// Check if project exists and belongs to the user
		project, err := models.GetProjectByID(db, projectID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Project not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch project"})
		}

		if project.UserID != userID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "You don't have permission to delete this project"})
		}

		// Delete the project
		err = models.DeleteProject(db, projectID, userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete project"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Project deleted successfully"})
	}
}

// GetUserPublicProjects retrieves public projects for a specific user by username
func GetUserPublicProjects(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")
		if username == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username is required"})
		}

		projects, err := models.GetUserPublicProjectsByUsername(db, username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch projects"})
		}

		return c.JSON(http.StatusOK, projects)
	}
}