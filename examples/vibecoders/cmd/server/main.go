package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"vibecoders/api/handlers"
	"vibecoders/models"
)

func main() {
	// Initialize the database
	db, err := models.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize handler
	h := &handlers.Handler{
		DB: db,
	}

	// Static files
	e.Static("/", "static")

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.File("templates/index.html")
	})

	// API routes
	api := e.Group("/api")
	
	// Auth routes
	api.POST("/register", h.RegisterUser)
	api.POST("/login", h.LoginUser)
	
	// User routes
	api.GET("/users/:id", h.GetUserProfile)
	api.PUT("/users/:id", h.UpdateUserProfile)

	// Start server
	log.Println("Starting server on :3000")
	if err := e.Start(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}