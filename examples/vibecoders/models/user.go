package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Bio         string    `json:"bio"`
	LinkedInURL string    `json:"linked_in_url"`
	GithubURL   string    `json:"github_url"`
	PhotoURL    string    `json:"photo_url"`
	Password    string    `json:"-"` // Don't include password in JSON responses
	CreatedAt   time.Time `json:"created_at"`
}

func GetTopUsers(db *sql.DB, limit int) ([]User, error) {
	query := `SELECT id, username, bio, linked_in_url, github_url, photo_url, created_at 
              FROM users 
              ORDER BY id 
              LIMIT ?`
	
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		var bio, linkedIn, github sql.NullString
		
		err := rows.Scan(&u.ID, &u.Username, &bio, &linkedIn, &github, &u.PhotoURL, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		if bio.Valid {
			u.Bio = bio.String
		}
		if linkedIn.Valid {
			u.LinkedInURL = linkedIn.String
		}
		if github.Valid {
			u.GithubURL = github.String
		}
		
		users = append(users, u)
	}

	return users, nil
}

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `SELECT id, username, bio, linked_in_url, github_url, photo_url, password, created_at 
              FROM users 
              WHERE username = ?`
	
	var u User
	var bio, linkedIn, github sql.NullString
	
	err := db.QueryRow(query, username).Scan(
		&u.ID, &u.Username, &bio, &linkedIn, &github, &u.PhotoURL, &u.Password, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	if bio.Valid {
		u.Bio = bio.String
	}
	if linkedIn.Valid {
		u.LinkedInURL = linkedIn.String
	}
	if github.Valid {
		u.GithubURL = github.String
	}
	
	return &u, nil
}

func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := `SELECT id, username, bio, linked_in_url, github_url, photo_url, created_at 
              FROM users 
              WHERE id = ?`
	
	var u User
	var bio, linkedIn, github sql.NullString
	
	err := db.QueryRow(query, id).Scan(
		&u.ID, &u.Username, &bio, &linkedIn, &github, &u.PhotoURL, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	if bio.Valid {
		u.Bio = bio.String
	}
	if linkedIn.Valid {
		u.LinkedInURL = linkedIn.String
	}
	if github.Valid {
		u.GithubURL = github.String
	}
	
	return &u, nil
}

func CreateUser(db *sql.DB, username, password, bio, linkedIn, github, photoURL string) error {
	query := `INSERT INTO users (username, password, bio, linked_in_url, github_url, photo_url) 
              VALUES (?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, username, password, bio, linkedIn, github, photoURL)
	return err
}

func UpdateUser(db *sql.DB, id int, bio, linkedIn, github, photoURL string) error {
	query := `UPDATE users 
              SET bio = ?, linked_in_url = ?, github_url = ?, photo_url = ? 
              WHERE id = ?`
	
	_, err := db.Exec(query, bio, linkedIn, github, photoURL, id)
	return err
}