package main

import (
	"database/sql"
	"embed"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// Task represents a task in our application
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// User represents a user in our application
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // omitempty to not send in responses
}

// LoginCredentials for login requests
type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//go:embed templates/*.html static/*.css static/*.js
var templateFS embed.FS

// Template renderer
type TemplateRenderer struct {
	templates *template.Template
}

// Implement echo.Renderer interface
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Initialize templates
func initTemplates() *TemplateRenderer {
	// Get templates from embedded filesystem
	templatesContent, err := fs.Sub(templateFS, "templates")
	if err != nil {
		log.Fatal(err)
	}

	return &TemplateRenderer{
		templates: template.Must(template.ParseFS(templatesContent, "*.html")),
	}
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tasks table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createUsersTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

// Setup static file server from embedded files
func setupStaticFiles(e *echo.Echo) {
	// Get the static files from embedded filesystem
	staticFiles, err := fs.Sub(templateFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	// Create a filesystem HTTP handler
	fileServer := http.FileServer(http.FS(staticFiles))

	// Register the handler for the /static/ path
	staticHandler := echo.WrapHandler(http.StripPrefix("/static/", fileServer))
	e.GET("/static/*", staticHandler)
}

func main() {
	// Initialize database
	initDB()
	defer db.Close()

	// Create a new Echo instance
	e := echo.New()

	// Set renderer
	renderer := initTemplates()
	e.Renderer = renderer

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5: true,
	}))
	e.Use(middleware.CORS())

	// Page Routes
	e.GET("/", dashboard)

	// Task API Routes
	e.GET("/tasks", getTasks)
	e.POST("/tasks", createTask)
	e.PUT("/tasks/:id", updateTask)
	e.DELETE("/tasks/:id", deleteTask)

	// User API Routes
	e.POST("/api/register", registerUser)
	e.POST("/api/login", loginUser)

	// Setup static file serving
	setupStaticFiles(e)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler for the dashboard page
func dashboard(c echo.Context) error {
	rows, err := db.Query("SELECT id, title, completed, created_at FROM tasks ORDER BY created_at DESC")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch tasks",
		})
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
			continue
		}
		tasks = append(tasks, t)
	}

	data := map[string]interface{}{
		"PageTitle": "Task Manager Dashboard",
		"Tasks":     tasks,
		"TaskCount": len(tasks),
	}
	return c.Render(http.StatusOK, "dashboard.html", data)
}

// Handler to get all tasks

// Handler to get all tasks
func getTasks(c echo.Context) error {
	rows, err := db.Query("SELECT id, title, completed, created_at FROM tasks")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch tasks",
		})
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
			continue
		}
		tasks = append(tasks, t)
	}

	return c.JSON(http.StatusOK, tasks)
}
