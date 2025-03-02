package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"vibecoders/models"

	"github.com/labstack/echo/v4"
)

type AdminUpdateUserRequest struct {
	Username    string `json:"username"`
	Fullname    string `json:"fullname"`
	Bio         string `json:"bio"`
	LinkedInURL string `json:"linked_in_url"`
	GithubURL   string `json:"github_url"`
	PhotoURL    string `json:"photo_url"`
	IsAdmin     bool   `json:"is_admin"`
}

// Check if the current user is an admin
func IsAdmin(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
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

			if !user.IsAdmin {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Administrator access required"})
			}

			return next(c)
		}
	}
}

// Get all users with pagination
func GetAllUsers(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse pagination parameters
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page < 1 {
			page = 1
		}
		
		pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
		if err != nil || pageSize < 1 || pageSize > 100 {
			pageSize = 10 // Default page size
		}

		// Get users for the requested page
		users, err := models.GetAllUsers(db, page, pageSize)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch users"})
		}

		// Get total user count for pagination
		total, err := models.GetTotalUserCount(db)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get user count"})
		}

		// Calculate total pages
		totalPages := (total + pageSize - 1) / pageSize

		return c.JSON(http.StatusOK, map[string]interface{}{
			"users": users,
			"pagination": map[string]interface{}{
				"total":      total,
				"page":       page,
				"pageSize":   pageSize,
				"totalPages": totalPages,
			},
		})
	}
}

// Get a specific user by ID
func GetUserByID(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		user, err := models.GetUserByID(db, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch user"})
		}

		return c.JSON(http.StatusOK, user)
	}
}

// Update a user (admin can update any field)
func UpdateUserAsAdmin(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		// Check if user exists
		_, err = models.GetUserByID(db, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch user"})
		}

		var req AdminUpdateUserRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// If changing username, check that it's not already taken
		if req.Username != "" {
			currentUser, err := models.GetUserByID(db, userID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch user"})
			}

			if currentUser.Username != req.Username {
				// Check if new username is already taken
				_, err = models.GetUserByUsername(db, req.Username)
				if err == nil {
					return c.JSON(http.StatusConflict, map[string]string{"error": "Username already exists"})
				} else if err != sql.ErrNoRows {
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
				}
			}
		}

		// Update user with admin privileges
		err = models.UpdateUserAdmin(db, userID, req.Username, req.Fullname, req.Bio, req.LinkedInURL, req.GithubURL, req.PhotoURL, req.IsAdmin)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update user"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User updated successfully"})
	}
}

// Delete a user
func DeleteUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		// Check if user exists
		_, err = models.GetUserByID(db, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch user"})
		}

		// Delete user
		err = models.DeleteUser(db, userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete user"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
	}
}