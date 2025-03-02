package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Fullname    string    `json:"fullname"`
	Bio         string    `json:"bio"`
	LinkedInURL string    `json:"linked_in_url"`
	GithubURL   string    `json:"github_url"`
	PhotoURL    string    `json:"photo_url"`
	Password    string    `json:"-"` // Don't include password in JSON responses
	CreatedAt   time.Time `json:"created_at"`
	IsAdmin     bool      `json:"is_admin"`
}

func GetTopUsers(db *sql.DB, limit int) ([]User, error) {
	query := `SELECT id, username, fullname, bio, linked_in_url, github_url, photo_url, created_at, is_admin 
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
		var bio, linkedIn, github, fullname sql.NullString
		
		err := rows.Scan(&u.ID, &u.Username, &fullname, &bio, &linkedIn, &github, &u.PhotoURL, &u.CreatedAt, &u.IsAdmin)
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
		if fullname.Valid {
			u.Fullname = fullname.String
		}
		
		users = append(users, u)
	}

	return users, nil
}

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `SELECT id, username, fullname, bio, linked_in_url, github_url, photo_url, password, created_at, is_admin 
              FROM users 
              WHERE username = ?`
	
	var u User
	var bio, linkedIn, github, fullname sql.NullString
	
	err := db.QueryRow(query, username).Scan(
		&u.ID, &u.Username, &fullname, &bio, &linkedIn, &github, &u.PhotoURL, &u.Password, &u.CreatedAt, &u.IsAdmin,
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
	if fullname.Valid {
		u.Fullname = fullname.String
	}
	
	return &u, nil
}

func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := `SELECT id, username, fullname, bio, linked_in_url, github_url, photo_url, created_at, is_admin 
              FROM users 
              WHERE id = ?`
	
	var u User
	var bio, linkedIn, github, fullname sql.NullString
	
	err := db.QueryRow(query, id).Scan(
		&u.ID, &u.Username, &fullname, &bio, &linkedIn, &github, &u.PhotoURL, &u.CreatedAt, &u.IsAdmin,
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
	if fullname.Valid {
		u.Fullname = fullname.String
	}
	
	return &u, nil
}

func CreateUser(db *sql.DB, username, password, fullname, bio, linkedIn, github, photoURL string) error {
	query := `INSERT INTO users (username, password, fullname, bio, linked_in_url, github_url, photo_url) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, username, password, fullname, bio, linkedIn, github, photoURL)
	return err
}

func UpdateUser(db *sql.DB, id int, fullname, bio, linkedIn, github, photoURL string) error {
	query := `UPDATE users 
              SET fullname = ?, bio = ?, linked_in_url = ?, github_url = ?, photo_url = ? 
              WHERE id = ?`
	
	_, err := db.Exec(query, fullname, bio, linkedIn, github, photoURL, id)
	return err
}

// GetAllUsers returns a paginated list of users sorted by created_at in descending order
func GetAllUsers(db *sql.DB, page, pageSize int) ([]User, error) {
	offset := (page - 1) * pageSize
	
	query := `SELECT id, username, fullname, bio, linked_in_url, github_url, photo_url, created_at, is_admin 
              FROM users 
              ORDER BY created_at DESC 
              LIMIT ? OFFSET ?`
	
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		var bio, linkedIn, github, fullname sql.NullString
		
		err := rows.Scan(&u.ID, &u.Username, &fullname, &bio, &linkedIn, &github, &u.PhotoURL, &u.CreatedAt, &u.IsAdmin)
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
		if fullname.Valid {
			u.Fullname = fullname.String
		}
		
		users = append(users, u)
	}

	return users, nil
}

// GetTotalUserCount returns the total number of users in the database
func GetTotalUserCount(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM users`
	
	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// UpdateUserAdmin allows an admin to update all fields of a user, including admin status
func UpdateUserAdmin(db *sql.DB, id int, username, fullname, bio, linkedIn, github, photoURL string, isAdmin bool) error {
	query := `UPDATE users 
              SET username = ?, fullname = ?, bio = ?, linked_in_url = ?, github_url = ?, photo_url = ?, is_admin = ?
              WHERE id = ?`
	
	_, err := db.Exec(query, username, fullname, bio, linkedIn, github, photoURL, isAdmin, id)
	return err
}

// DeleteUser deletes a user from the database
func DeleteUser(db *sql.DB, id int) error {
	// First delete any sessions for this user
	sessionQuery := `DELETE FROM sessions WHERE user_id = ?`
	_, err := db.Exec(sessionQuery, id)
	if err != nil {
		return err
	}
	
	// Then delete the user
	userQuery := `DELETE FROM users WHERE id = ?`
	_, err = db.Exec(userQuery, id)
	return err
}