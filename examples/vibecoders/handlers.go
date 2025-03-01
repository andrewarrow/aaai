package main

import (
	"net/http"
	"strconv"
	"time"

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
			"error": "Failed to create task",
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
