package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Simple HTTP server for Living Smart Contracts
func main() {
	// Setup HTTP server with basic endpoints
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "healthy",
			"service":   "living-smart-contracts",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API info endpoint
	mux.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":        "Living Smart Contracts",
			"version":     "1.0.0",
			"description": "Revolutionary blockchain platform with evolutionary intelligence",
			"features": []string{
				"Evolutionary Intelligence",
				"Living Contract Types",
				"Machine Learning Integration",
				"Inter-Contract Collaboration",
				"High-Performance Actor System",
			},
			"endpoints": map[string]string{
				"health":    "/health",
				"info":      "/api/info",
				"contracts": "/api/contracts",
			},
		})
	})

	// Placeholder contracts endpoint
	mux.HandleFunc("/api/contracts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "Living Smart Contracts API",
			"status":    "ready",
			"contracts": []map[string]interface{}{},
		})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Living Smart Contracts server starting on port %s", port)
		log.Printf("Health check available at: http://localhost:%s/health", port)
		log.Printf("API info available at: http://localhost:%s/api/info", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	fmt.Println("Living Smart Contracts server stopped")
}
