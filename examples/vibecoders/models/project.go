package models

import (
	"database/sql"
	"time"
)

type Project struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	GithubURL   string    `json:"github_url"`
	WebsiteURL  string    `json:"website_url"`
	ImageURL1   string    `json:"image_url1"`
	ImageURL2   string    `json:"image_url2"`
	ImageURL3   string    `json:"image_url3"`
	CreatedAt   time.Time `json:"created_at"`
}

// GetProjectsByUserID retrieves all projects for a specific user
func GetProjectsByUserID(db *sql.DB, userID int) ([]Project, error) {
	query := `SELECT id, user_id, title, description, github_url, website_url, 
                  image_url1, image_url2, image_url3, created_at 
              FROM projects 
              WHERE user_id = ? 
              ORDER BY created_at DESC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		var githubURL, websiteURL, imageURL1, imageURL2, imageURL3 sql.NullString

		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Description, &githubURL, &websiteURL,
			&imageURL1, &imageURL2, &imageURL3, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		if githubURL.Valid {
			p.GithubURL = githubURL.String
		}
		if websiteURL.Valid {
			p.WebsiteURL = websiteURL.String
		}
		if imageURL1.Valid {
			p.ImageURL1 = imageURL1.String
		}
		if imageURL2.Valid {
			p.ImageURL2 = imageURL2.String
		}
		if imageURL3.Valid {
			p.ImageURL3 = imageURL3.String
		}

		projects = append(projects, p)
	}

	if projects == nil {
		projects = []Project{}
	}

	return projects, nil
}

// GetProjectByID retrieves a single project by ID
func GetProjectByID(db *sql.DB, projectID int) (*Project, error) {
	query := `SELECT id, user_id, title, description, github_url, website_url, 
                 image_url1, image_url2, image_url3, created_at 
              FROM projects 
              WHERE id = ?`

	var p Project
	var githubURL, websiteURL, imageURL1, imageURL2, imageURL3 sql.NullString

	err := db.QueryRow(query, projectID).Scan(
		&p.ID, &p.UserID, &p.Title, &p.Description, &githubURL, &websiteURL,
		&imageURL1, &imageURL2, &imageURL3, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if githubURL.Valid {
		p.GithubURL = githubURL.String
	}
	if websiteURL.Valid {
		p.WebsiteURL = websiteURL.String
	}
	if imageURL1.Valid {
		p.ImageURL1 = imageURL1.String
	}
	if imageURL2.Valid {
		p.ImageURL2 = imageURL2.String
	}
	if imageURL3.Valid {
		p.ImageURL3 = imageURL3.String
	}

	return &p, nil
}

// CreateProject adds a new project to the database
func CreateProject(db *sql.DB, userID int, title, description, githubURL, websiteURL,
	imageURL1, imageURL2, imageURL3 string) (int, error) {

	query := `INSERT INTO projects (user_id, title, description, github_url, website_url, 
                               image_url1, image_url2, image_url3) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, userID, title, description, githubURL, websiteURL,
		imageURL1, imageURL2, imageURL3)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// UpdateProject modifies an existing project
func UpdateProject(db *sql.DB, projectID, userID int, title, description, githubURL, websiteURL,
	imageURL1, imageURL2, imageURL3 string) error {

	query := `UPDATE projects 
              SET title = ?, description = ?, github_url = ?, website_url = ?, 
                  image_url1 = ?, image_url2 = ?, image_url3 = ? 
              WHERE id = ? AND user_id = ?`

	_, err := db.Exec(query, title, description, githubURL, websiteURL,
		imageURL1, imageURL2, imageURL3, projectID, userID)
	return err
}

// DeleteProject removes a project from the database
func DeleteProject(db *sql.DB, projectID, userID int) error {
	query := `DELETE FROM projects 
              WHERE id = ? AND user_id = ?`

	_, err := db.Exec(query, projectID, userID)
	return err
}

// GetUserPublicProjectsByUsername retrieves public projects for a user by username
func GetUserPublicProjectsByUsername(db *sql.DB, username string) ([]Project, error) {
	query := `SELECT p.id, p.user_id, p.title, p.description, p.github_url, p.website_url, 
                  p.image_url1, p.image_url2, p.image_url3, p.created_at 
              FROM projects p
              JOIN users u ON p.user_id = u.id 
              WHERE u.username = ? 
              ORDER BY p.created_at DESC`

	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		var githubURL, websiteURL, imageURL1, imageURL2, imageURL3 sql.NullString

		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Description, &githubURL, &websiteURL,
			&imageURL1, &imageURL2, &imageURL3, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		if githubURL.Valid {
			p.GithubURL = githubURL.String
		}
		if websiteURL.Valid {
			p.WebsiteURL = websiteURL.String
		}
		if imageURL1.Valid {
			p.ImageURL1 = imageURL1.String
		}
		if imageURL2.Valid {
			p.ImageURL2 = imageURL2.String
		}
		if imageURL3.Valid {
			p.ImageURL3 = imageURL3.String
		}

		projects = append(projects, p)
	}

	if projects == nil {
		projects = []Project{}
	}

	return projects, nil
}