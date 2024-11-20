package main

import (
	"context"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	// "whoami-go/src/docs" // import the generated docs package
)

// User represents a user structure
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Response struct {
	Message string `json:"message"`
}

// Database pool
var db *pgxpool.Pool

// Initialize the database connection
func initDB() {
	// Connection URL for the Neon database
	databaseURL := "postgresql://whoamidb_owner:o4t1DRpQJBGh@ep-snowy-pond-a2ywwuj6.eu-central-1.aws.neon.tech/whoamidb?sslmode=require"

	var err error
	db, err = pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}
	log.Println("Connected to the database successfully!")
}

// Close the database connection when the app shuts down
func closeDB() {
	if db != nil {
		db.Close()
		log.Println("Database connection closed.")
	}
}

// @title Go Fiber API with Swagger
// @version 1.0
// @description This is a simple API with Swagger documentation using Fiber in Go.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:3000
// @BasePath /api
func main() {
	// Initialize the database
	initDB()
	defer closeDB()

	// Create a new Fiber app
	app := fiber.New()

	// Apply CORS middleware to allow cross-origin requests
	app.Use(func(c *fiber.Ctx) error {
		// CORS headers for frontend (Svelte in this case)
		c.Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Method() == "OPTIONS" {
			// If it's an OPTIONS request, just return a 200 status
			return c.SendStatus(fiber.StatusOK)
		}
		return c.Next()
	})

	// Define a simple GET endpoint
	// @Summary Greet the user
	// @Description A simple endpoint to return a greeting
	// @Produce json
	// @Success 200 {string} string "Hello, World!"
	// @Router / [get]
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Define an API endpoint
	// @Summary Returns API status
	// @Description Returns a success message
	// @Produce json
	// @Success 200 {object} map[string]string{"status": "success", "message": "Welcome to the API"}
	// @Router /api [get]
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to the API",
		})
	})

	// Define a dynamic route
	// @Summary Get user by ID
	// @Description Fetch user details using user ID
	// @Param id path int true "User ID"
	// @Produce json
	// @Success 200 {object} map[string]string{"user_id": "1"}
	// @Router /api/user/{id} [get]
	app.Get("/api/user/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.JSON(fiber.Map{
			"user_id": id,
		})
	})

	// POST endpoint: Create a new user
	// @Summary Create a new user
	// @Description Create a new user in the system
	// @Accept json
	// @Produce json
	// @Param user body User true "User Data"
	// @Success 201 {object} map[string]interface{}{"status": "success", "message": "User created successfully", "user": User{}}
	// @Router /api/user [post]
	app.Post("/api/user", func(c *fiber.Ctx) error {
		// Parse the incoming JSON into a User struct
		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid JSON data",
			})
		}

		// Insert user into the database
		query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
		var newUserID int
		err := db.QueryRow(context.Background(), query, user.Name, user.Email).Scan(&newUserID)
		if err != nil {
			log.Printf("Failed to insert user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to create user",
			})
		}

		// Respond with the created user data
		user.ID = newUserID
		log.Printf("User created: %+v", user)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":  "success",
			"message": "User created successfully",
			"user":    user,
		})
	})

	// Swagger UI endpoint
	// @Router /swagger/*any [get]
	// app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Start the server on port 3000
	port := ":3000"
	log.Printf("Server running on http://localhost%s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
