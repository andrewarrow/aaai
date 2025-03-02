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

//go:embed db/init.sql
var initSQLContent embed.FS

//go:embed static/dist
var staticContent embed.FS

func main() {
	// Initialize database
	os.Remove("./vibecoders.db") // Remove existing database during development
	db, err := sql.Open("sqlite3", "./vibecoders.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Execute database initialization script
	initSQL, err := initSQLContent.ReadFile("db/init.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(initSQL))
	if err != nil {
		log.Fatal(err)
	}

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

	assetHandler := http.FileServer(http.FS(staticFS))

	e.GET("/", func(c echo.Context) error {
		indexHTML, err := fs.ReadFile(staticFS, "index.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error reading index.html")
		}

		return c.Blob(http.StatusOK, "text/html", indexHTML)
	})

	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
