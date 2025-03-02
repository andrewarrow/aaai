package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"vibecoders/api/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed static/dist
var staticContent embed.FS

func main() {
	// Initialize database connection
	// Database is expected to be migrated using Flyway before server startup
	dbPath := "./vibecoders.db"

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatalf("Database %s not found. Please run migrations first with: npm run db:migrate", dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to database")

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Serve static files from embedded filesystem
	staticFS, err := fs.Sub(staticContent, "static/dist")
	if err != nil {
		log.Fatal(err)
	}

	// API routes
	api := e.Group("/api")
	api.POST("/login", handlers.Login(db))
	api.DELETE("/logout", handlers.Logout(db))
	api.POST("/register", handlers.Register(db))
	api.PATCH("/user", handlers.UpdateUser(db))
	api.GET("/homepage-users", handlers.GetHomepageUsers(db))
	api.GET("/user", handlers.GetCurrentUser(db))
	api.GET("/users/:username", handlers.GetPublicUserByUsername(db))
	
	// Prompt routes
	api.GET("/prompts", handlers.GetUserPrompts(db))
	api.POST("/prompts", handlers.CreatePrompt(db))
	api.PUT("/prompts/:id", handlers.UpdatePrompt(db))
	api.DELETE("/prompts/:id", handlers.DeletePrompt(db))
	api.GET("/users/:username/prompts", handlers.GetUserPublicPrompts(db))
	
	// Project routes
	api.GET("/projects", handlers.GetUserProjects(db))
	api.POST("/projects", handlers.CreateProject(db))
	api.PUT("/projects/:id", handlers.UpdateProject(db))
	api.DELETE("/projects/:id", handlers.DeleteProject(db))
	api.GET("/users/:username/projects", handlers.GetUserPublicProjects(db))

	// Admin API routes with admin middleware
	adminMiddleware := handlers.IsAdmin(db)
	admin := api.Group("/admin", adminMiddleware)
	admin.GET("/users", handlers.GetAllUsers(db))
	admin.GET("/users/:id", handlers.GetUserByID(db))
	admin.PUT("/users/:id", handlers.UpdateUserAsAdmin(db))
	admin.DELETE("/users/:id", handlers.DeleteUser(db))

	assetHandler := http.FileServer(http.FS(staticFS))

	// Function to serve the SPA index.html for any frontend routes
	serveSPA := func(c echo.Context) error {
		indexHTML, err := fs.ReadFile(staticFS, "index.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error reading index.html")
		}

		return c.Blob(http.StatusOK, "text/html", indexHTML)
	}

	// Serve SPA routes
	e.GET("/", serveSPA)
	e.GET("/login", serveSPA)
	e.GET("/register", serveSPA)
	e.GET("/profile", serveSPA)
	e.GET("/users/:username", serveSPA)
	e.GET("/admin", serveSPA)
	e.GET("/admin/*", serveSPA)

	// Serve static assets
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
