package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
	
	_ "github.com/lib/pq"
	
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Connect to PostgreSQL
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Database connection failed:", err)
	}
	
	log.Println("Connected to PostgreSQL database")
	
	// Initialize auth service
	authService := service.NewAuthService(db, cfg.SessionTTL, cfg.CleanupInterval)
	
	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	
	// Setup routes
	http.HandleFunc("/v1/register", authHandler.Register)
	http.HandleFunc("/v1/login", authHandler.Login)
	http.HandleFunc("/v1/logout", authHandler.Logout)
	http.HandleFunc("/v1/validate/token", authHandler.ValidateToken)
	http.HandleFunc("/health", authHandler.Health)
	
	// Start server
	log.Printf("Auth service starting on port %s", cfg.Port)
	log.Printf("Session TTL: %v, Cleanup interval: %v", cfg.SessionTTL, cfg.CleanupInterval)
	
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}