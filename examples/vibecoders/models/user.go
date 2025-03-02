package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	Bio        string    `json:"bio"`
	LinkedinURL string   `json:"linkedin_url"`
	GithubURL  string    `json:"github_url"`
	PhotoURL   string    `json:"photo_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// InitDB sets up the database connection and creates tables if needed
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./vibecoders.db")
	if err != nil {
		return nil, err
	}

	// Create users table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		bio TEXT,
		linkedin_url TEXT,
		github_url TEXT,
		photo_url TEXT,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateUser creates a new user in the database
func CreateUser(db *sql.DB, username, password, bio, linkedinURL, githubURL, photoURL string) (*User, error) {
	// Check if username already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create a new user
	user := &User{
		ID:          uuid.New().String(),
		Username:    username,
		Password:    string(hashedPassword),
		Bio:         bio,
		LinkedinURL: linkedinURL,
		GithubURL:   githubURL,
		PhotoURL:    photoURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Insert the user into the database
	_, err = db.Exec(
		"INSERT INTO users (id, username, password, bio, linkedin_url, github_url, photo_url, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.ID, user.Username, user.Password, user.Bio, user.LinkedinURL, user.GithubURL, user.PhotoURL, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	user := &User{}
	err := db.QueryRow(
		"SELECT id, username, password, bio, linkedin_url, github_url, photo_url, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Bio, &user.LinkedinURL, &user.GithubURL, &user.PhotoURL, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(db *sql.DB, id string) (*User, error) {
	user := &User{}
	err := db.QueryRow(
		"SELECT id, username, password, bio, linkedin_url, github_url, photo_url, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Bio, &user.LinkedinURL, &user.GithubURL, &user.PhotoURL, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// UpdateUser updates user information
func UpdateUser(db *sql.DB, id, bio, linkedinURL, githubURL, photoURL string) (*User, error) {
	// Get the current user
	user, err := GetUserByID(db, id)
	if err != nil {
		return nil, err
	}

	// Update the user fields
	user.Bio = bio
	user.LinkedinURL = linkedinURL
	user.GithubURL = githubURL
	user.PhotoURL = photoURL
	user.UpdatedAt = time.Now()

	// Update the user in the database
	_, err = db.Exec(
		"UPDATE users SET bio = ?, linkedin_url = ?, github_url = ?, photo_url = ?, updated_at = ? WHERE id = ?",
		user.Bio, user.LinkedinURL, user.GithubURL, user.PhotoURL, user.UpdatedAt, user.ID,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Authentication functions
func VerifyPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// UpdatePassword updates a user's password
func UpdatePassword(db *sql.DB, id, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE users SET password = ?, updated_at = ? WHERE id = ?", 
		string(hashedPassword), time.Now(), id)
	return err
}