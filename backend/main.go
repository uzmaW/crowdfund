package main

import (
	"context"
	"crowdfund/backend/handlers"
	"crowdfund/backend/middlewares"
	"crowdfund/backend/models"
	"crowdfund/backend/services"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pst "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// getEnvOrDefault returns the value of the environment variable if it exists, or the default value if it doesn't
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		getEnvOrDefault("POSTGRES_PORT", "5434"),
		os.Getenv("POSTGRES_DB"),
	)
	log.Println(dbURL)
	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database for migrations: %v", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	dsn := dbURL //os.Getenv("DATABASE_URL")
	db, err := gorm.Open(pst.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	userService := services.NewUserService(db)
	projectService := services.NewProjectService(db)
	emailService := services.NewEmailService()
	cacheService := services.NewCacheService()

	// Worker Pool Setup
	donationTasks := make(chan models.Donation, 100) // Buffered channel
	var donationWg sync.WaitGroup
	numDonationWorkers := 5
	for i := 1; i <= numDonationWorkers; i++ {
		donationWg.Add(1)
		go services.DonationWorker(i, donationTasks, &donationWg, db, emailService)
	}
	donationService := services.NewDonationService(db, emailService, donationTasks)

	r := gin.Default()

	userHandlers := handlers.NewUserHandlers(userService, cacheService)
	projectHandlers := handlers.NewProjectHandlers(projectService, cacheService)
	donationHandlers := handlers.NewDonationHandlers(donationService)
	passHandlers := handlers.PassHandlers{}

	r.POST("/users/register", userHandlers.Register)
	r.POST("/users/login", userHandlers.Login)
	r.GET("/api/users/profile", middlewares.AuthMiddleware(), userHandlers.Profile)

	r.POST("/api/projects", middlewares.AuthMiddleware(), projectHandlers.CreateProject)
	r.GET("/api/projects/:id", projectHandlers.GetProject)
	r.PUT("/api/projects/:id", middlewares.AuthMiddleware(), projectHandlers.UpdateProject)
	r.DELETE("/api/projects/:id", middlewares.AuthMiddleware(), projectHandlers.DeleteProject)
	r.GET("/api/projects", projectHandlers.ListProjects)

	r.POST("/api/projects/:id/donations", middlewares.AuthMiddleware(), donationHandlers.CreateDonation)
	r.GET("/api/projects/:id/donations", middlewares.AuthMiddleware(), donationHandlers.GetDonationsByProjectID)
	r.POST("/password", passHandlers.GetHashForPass)
	
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	close(donationTasks) // Signal workers to stop
	donationWg.Wait()    // Wait for workers to finish

	log.Println("Server exiting")
}
