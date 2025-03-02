package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"vibecoders.com/api/handlers"
)

func main() {
	// Initialize database
	os.Remove("./vibecoders.db") // Remove existing database during development
	db, err := sql.Open("sqlite3", "./vibecoders.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Execute database initialization script
	initDB, err := os.ReadFile("./db/init.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(initDB))
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Serve static files
	e.Static("/", "static/dist")

	// API routes
	api := e.Group("/api")
	api.POST("/login", handlers.Login(db))
	api.DELETE("/logout", handlers.Logout(db))
	api.POST("/register", handlers.Register(db))
	api.PATCH("/user", handlers.UpdateUser(db))
	api.GET("/homepage-users", handlers.GetHomepageUsers(db))
	api.GET("/user", handlers.GetCurrentUser(db))

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
