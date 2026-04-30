package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/First-Genesis/Living-Smart-Contracts/cmd/server/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

//go:generate swag init

// Living Smart Contracts API Server with Swagger documentation
func main() {
	// Setup HTTP server with Swagger-documented endpoints
	mux := http.NewServeMux()

	// Swagger UI endpoint
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Health check endpoints
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/health", healthHandler)

	// API info endpoint
	mux.HandleFunc("/api/info", apiInfoHandler)

	// Contract management endpoints
	mux.HandleFunc("/api/contracts", contractsHandler)
	mux.HandleFunc("/api/contracts/deploy", deployContractHandler)
	mux.HandleFunc("/api/contracts/execute", executeContractHandler)

	// Contract lifecycle endpoints
	mux.HandleFunc("/api/contracts/", contractDetailHandler)

	// Evolution endpoints
	mux.HandleFunc("/api/evolution/trigger", triggerEvolutionHandler)

	// Collaboration endpoints
	mux.HandleFunc("/api/collaborations/propose", proposeCollaborationHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      enableCORS(mux),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Living Smart Contracts server starting on port %s", port)
		log.Printf("📊 Swagger UI available at: http://localhost:%s/swagger/", port)
		log.Printf("💚 Health check available at: http://localhost:%s/health", port)
		log.Printf("📋 API info available at: http://localhost:%s/api/info", port)
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

// enableCORS adds CORS headers to all responses
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// healthHandler handles health check requests
// @Summary Health Check
// @Description Get the health status of the Living Smart Contracts service
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Status:    "healthy",
		Service:   "living-smart-contracts",
		Timestamp: time.Now().UTC(),
	}

	json.NewEncoder(w).Encode(response)
}

// apiInfoHandler handles API information requests
// @Summary API Information
// @Description Get information about the Living Smart Contracts API
// @Tags Health
// @Produce json
// @Success 200 {object} APIInfoResponse
// @Router /api/info [get]
func apiInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := APIInfoResponse{
		Name:        "Living Smart Contracts",
		Version:     "1.0.0",
		Description: "Revolutionary blockchain platform with evolutionary intelligence",
		Features: []string{
			"Evolutionary Intelligence",
			"Living Contract Types",
			"Machine Learning Integration",
			"Inter-Contract Collaboration",
			"High-Performance Actor System",
		},
		Endpoints: map[string]string{
			"health":        "/health",
			"info":          "/api/info",
			"contracts":     "/api/contracts",
			"deploy":        "/api/contracts/deploy",
			"execute":       "/api/contracts/execute",
			"evolution":     "/api/evolution/trigger",
			"collaboration": "/api/collaborations/propose",
			"swagger":       "/swagger/",
		},
	}

	json.NewEncoder(w).Encode(response)
}

// contractsHandler handles contract listing requests
// @Summary List Contracts
// @Description Get a list of all deployed smart contracts
// @Tags Contracts
// @Produce json
// @Success 200 {object} ContractsResponse
// @Router /api/contracts [get]
func contractsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Mock contracts for demonstration
	mockContracts := []Contract{
		{
			ID:        "550e8400-e29b-41d4-a716-446655440000",
			Address:   "0xabcdef1234567890abcdef1234567890abcdef12",
			Name:      "adaptive_trading_contract",
			Type:      "living",
			Status:    "active",
			Owner:     "0x1234567890123456789012345678901234567890",
			Version:   "1.0.0",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
			DNA: ContractDNA{
				Generation: 5,
				Fitness:    0.934,
				Genes: []Gene{
					{
						ID:        "optimization_001",
						Type:      "optimization",
						Dominance: 0.92,
						Stability: 0.88,
					},
				},
			},
			Metrics: PerformanceMetrics{
				ExecutionCount:  15420,
				SuccessRate:     0.987,
				AverageGasUsed:  45000,
				AdaptationScore: 0.923,
			},
		},
	}

	response := ContractsResponse{
		Message:   "Living Smart Contracts API",
		Status:    "ready",
		Contracts: mockContracts,
		Total:     len(mockContracts),
	}

	json.NewEncoder(w).Encode(response)
}

// deployContractHandler handles contract deployment requests
// @Summary Deploy Contract
// @Description Deploy a new living smart contract
// @Tags Contracts
// @Accept json
// @Produce json
// @Param request body DeployContractRequest true "Contract deployment request"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/contracts/deploy [post]
func deployContractHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req DeployContractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Invalid JSON",
			Message: "Failed to parse request body",
			Code:    400,
		})
		return
	}

	// Mock deployment logic
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("Contract '%s' deployed successfully", req.Name),
		Data: map[string]interface{}{
			"contract_id":      "550e8400-e29b-41d4-a716-446655440001",
			"contract_address": "0x" + strings.Repeat("ab", 20),
			"status":           "deploying",
		},
	})
}

// executeContractHandler handles contract execution requests
// @Summary Execute Contract
// @Description Execute a function on a deployed smart contract
// @Tags Contracts
// @Accept json
// @Produce json
// @Param request body ExecuteContractRequest true "Contract execution request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/contracts/execute [post]
func executeContractHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req ExecuteContractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Invalid JSON",
			Message: "Failed to parse request body",
			Code:    400,
		})
		return
	}

	// Mock execution logic
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("Function '%s' executed successfully", req.Function),
		Data: map[string]interface{}{
			"execution_id": "exec_" + fmt.Sprintf("%d", time.Now().Unix()),
			"gas_used":     42000,
			"result":       "success",
		},
	})
}

// contractDetailHandler handles individual contract detail requests
// @Summary Get Contract Details
// @Description Get detailed information about a specific contract
// @Tags Contracts
// @Produce json
// @Param address path string true "Contract Address"
// @Success 200 {object} Contract
// @Failure 404 {object} ErrorResponse
// @Router /api/contracts/{address} [get]
func contractDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract contract address from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/contracts/")
	if path == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Contract not found",
			Message: "Contract address is required",
			Code:    404,
		})
		return
	}

	// Mock contract details
	contract := Contract{
		ID:        "550e8400-e29b-41d4-a716-446655440000",
		Address:   path,
		Name:      "adaptive_trading_contract",
		Type:      "living",
		Status:    "active",
		Owner:     "0x1234567890123456789012345678901234567890",
		Version:   "1.0.0",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	json.NewEncoder(w).Encode(contract)
}

// triggerEvolutionHandler handles evolution trigger requests
// @Summary Trigger Evolution
// @Description Trigger evolutionary process for a contract
// @Tags Evolution
// @Accept json
// @Produce json
// @Param address path string true "Contract Address"
// @Param request body EvolutionRequest true "Evolution request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/evolution/trigger [post]
func triggerEvolutionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req EvolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Invalid JSON",
			Message: "Failed to parse request body",
			Code:    400,
		})
		return
	}

	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("Evolution '%s' triggered successfully", req.EvolutionType),
		Data: map[string]interface{}{
			"evolution_id":   "evo_" + fmt.Sprintf("%d", time.Now().Unix()),
			"estimated_time": "300s",
			"evolution_type": req.EvolutionType,
		},
	})
}

// proposeCollaborationHandler handles collaboration proposal requests
// @Summary Propose Collaboration
// @Description Propose a collaboration between two contracts
// @Tags Collaboration
// @Accept json
// @Produce json
// @Param request body CollaborationRequest true "Collaboration proposal request"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/collaborations/propose [post]
func proposeCollaborationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req CollaborationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Invalid JSON",
			Message: "Failed to parse request body",
			Code:    400,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Message: "Collaboration proposal submitted successfully",
		Data: map[string]interface{}{
			"collaboration_id": "collab_" + fmt.Sprintf("%d", time.Now().Unix()),
			"status":           "pending",
			"type":             req.Type,
		},
	})
}
