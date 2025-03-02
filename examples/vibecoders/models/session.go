package models

import (
	"database/sql"
	"github.com/google/uuid"
)

type Session struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}

func CreateSession(db *sql.DB, userID int) (string, error) {
	token := uuid.New().String()
	
	query := `INSERT INTO sessions (user_id, token) VALUES (?, ?)`
	_, err := db.Exec(query, userID, token)
	if err != nil {
		return "", err
	}
	
	return token, nil
}

func GetUserIDByToken(db *sql.DB, token string) (int, error) {
	query := `SELECT user_id FROM sessions WHERE token = ?`
	
	var userID int
	err := db.QueryRow(query, token).Scan(&userID)
	if err != nil {
		return 0, err
	}
	
	return userID, nil
}

func DeleteSession(db *sql.DB, token string) error {
	query := `DELETE FROM sessions WHERE token = ?`
	_, err := db.Exec(query, token)
	return err
}