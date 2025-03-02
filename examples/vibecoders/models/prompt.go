package models

import (
	"database/sql"
	"strings"
	"time"
)

type Prompt struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}

// GetPromptsByUserID retrieves all prompts for a specific user
func GetPromptsByUserID(db *sql.DB, userID int) ([]Prompt, error) {
	query := `SELECT id, user_id, title, content, tags, created_at 
              FROM prompts 
              WHERE user_id = ? 
              ORDER BY created_at DESC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompts []Prompt
	for rows.Next() {
		var p Prompt
		var tagString sql.NullString

		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &tagString, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		if tagString.Valid && tagString.String != "" {
			p.Tags = strings.Split(tagString.String, ",")
		} else {
			p.Tags = []string{}
		}

		prompts = append(prompts, p)
	}

	if prompts == nil {
		prompts = []Prompt{}
	}

	return prompts, nil
}

// GetPromptByID retrieves a single prompt by ID
func GetPromptByID(db *sql.DB, promptID int) (*Prompt, error) {
	query := `SELECT id, user_id, title, content, tags, created_at 
              FROM prompts 
              WHERE id = ?`

	var p Prompt
	var tagString sql.NullString

	err := db.QueryRow(query, promptID).Scan(
		&p.ID, &p.UserID, &p.Title, &p.Content, &tagString, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if tagString.Valid && tagString.String != "" {
		p.Tags = strings.Split(tagString.String, ",")
	} else {
		p.Tags = []string{}
	}

	return &p, nil
}

// CreatePrompt adds a new prompt to the database
func CreatePrompt(db *sql.DB, userID int, title, content string, tags []string) (int, error) {
	tagString := strings.Join(tags, ",")

	query := `INSERT INTO prompts (user_id, title, content, tags) 
              VALUES (?, ?, ?, ?)`

	result, err := db.Exec(query, userID, title, content, tagString)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// UpdatePrompt modifies an existing prompt
func UpdatePrompt(db *sql.DB, promptID, userID int, title, content string, tags []string) error {
	tagString := strings.Join(tags, ",")

	query := `UPDATE prompts 
              SET title = ?, content = ?, tags = ? 
              WHERE id = ? AND user_id = ?`

	_, err := db.Exec(query, title, content, tagString, promptID, userID)
	return err
}

// DeletePrompt removes a prompt from the database
func DeletePrompt(db *sql.DB, promptID, userID int) error {
	query := `DELETE FROM prompts 
              WHERE id = ? AND user_id = ?`

	_, err := db.Exec(query, promptID, userID)
	return err
}

// GetUserPublicPrompts retrieves public prompts for a user by username
func GetUserPublicPromptsByUsername(db *sql.DB, username string) ([]Prompt, error) {
	query := `SELECT p.id, p.user_id, p.title, p.content, p.tags, p.created_at 
              FROM prompts p
              JOIN users u ON p.user_id = u.id 
              WHERE u.username = ? 
              ORDER BY p.created_at DESC`

	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompts []Prompt
	for rows.Next() {
		var p Prompt
		var tagString sql.NullString

		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &tagString, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		if tagString.Valid && tagString.String != "" {
			p.Tags = strings.Split(tagString.String, ",")
		} else {
			p.Tags = []string{}
		}

		prompts = append(prompts, p)
	}

	if prompts == nil {
		prompts = []Prompt{}
	}

	return prompts, nil
}