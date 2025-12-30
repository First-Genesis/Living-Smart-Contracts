package contracts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

// Precompiled regex for contract address validation
var contractAddressRegex = regexp.MustCompile(`^(0x)?[a-fA-F0-9]{40}$`)

// ContractAPIHandlers provides HTTP handlers for smart contract operations
type ContractAPIHandlers struct {
	contractManagerPID *actor.PID
	system             *actor.ActorSystem
}

// NewContractAPIHandlers creates a new contract API handler instance
func NewContractAPIHandlers(contractManagerPID *actor.PID, system *actor.ActorSystem) *ContractAPIHandlers {
	return &ContractAPIHandlers{
		contractManagerPID: contractManagerPID,
		system:             system,
	}
}

// RegisterRoutes adds contract API routes to the HTTP server
func (cah *ContractAPIHandlers) RegisterRoutes(mux *http.ServeMux) {
	// Contract management endpoints (specific routes first)
	mux.HandleFunc("/api/contracts", cah.enableCORS(cah.handleContracts))
	mux.HandleFunc("/api/contracts/deploy", cah.enableCORS(cah.handleDeployContract))
	mux.HandleFunc("/api/contracts/execute", cah.enableCORS(cah.handleExecuteContract))

	// Contract lifecycle endpoints - distinct paths
	mux.HandleFunc("/api/contracts/lifecycle/", cah.enableCORS(cah.handleContractLifecycle))
	mux.HandleFunc("/api/contracts/upgrade/", cah.enableCORS(cah.handleUpgradeContract))
	mux.HandleFunc("/api/contracts/evolve/", cah.enableCORS(cah.handleEvolveContract))

	// Collaboration endpoints
	mux.HandleFunc("/api/contracts/collaborate", cah.enableCORS(cah.handleProposeCollaboration))
	mux.HandleFunc("/api/contracts/collaborations/", cah.enableCORS(cah.handleCollaborationManagement))

	// Query and analytics endpoints
	mux.HandleFunc("/api/contracts/state/", cah.enableCORS(cah.handleQueryContractState))
	mux.HandleFunc("/api/contracts/history/", cah.enableCORS(cah.handleQueryHistory))
	mux.HandleFunc("/api/contracts/predict/", cah.enableCORS(cah.handlePredictBehavior))

	// Ecosystem endpoints
	mux.HandleFunc("/api/contracts/ecosystems", cah.enableCORS(cah.handleEcosystems))
	mux.HandleFunc("/api/contracts/ecosystems/", cah.enableCORS(cah.handleEcosystemByID))

	// System monitoring endpoints
	mux.HandleFunc("/api/contracts/system/status", cah.enableCORS(cah.handleSystemStatus))
	mux.HandleFunc("/api/contracts/system/metrics", cah.enableCORS(cah.handleSystemMetrics))
	mux.HandleFunc("/api/contracts/system/health", cah.enableCORS(cah.handleSystemHealth))

	// Catch-all address handler MUST be registered last to avoid conflicts
	mux.HandleFunc("/api/contracts/", cah.enableCORS(cah.handleContractByAddress))
}

// ===== UTILITY METHODS =====

// getRequestID returns the request id (uses X-Request-Id if provided, else generates one)
func (cah *ContractAPIHandlers) getRequestID(r *http.Request) string {
	if rid := strings.TrimSpace(r.Header.Get("X-Request-Id")); rid != "" {
		return rid
	}
	// Simple generator: timestamp + random suffix (good enough for tracing)
	return fmt.Sprintf("%d-%06d", time.Now().UnixNano(), rand.Intn(1_000_000))
}

// writeJSON writes JSON responses with consistent encoding policy
func (cah *ContractAPIHandlers) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false) // Prevent unnecessary HTML escaping
	_ = enc.Encode(v)
}

// decodeJSON safely decodes JSON with size limits and validation
func (cah *ContractAPIHandlers) decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}, maxBytes int64) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		// Provide more specific error messages for better debugging
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {
		case errors.Is(err, io.EOF):
			return fmt.Errorf("request body is empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return fmt.Errorf("unknown field in JSON: %s", strings.TrimPrefix(err.Error(), "json: unknown field "))
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("invalid JSON syntax at position %d", syntaxErr.Offset)
		case errors.As(err, &unmarshalTypeErr):
			return fmt.Errorf("invalid value for field %s: expected %s", unmarshalTypeErr.Field, unmarshalTypeErr.Type)
		case strings.Contains(err.Error(), "request body too large"):
			return fmt.Errorf("request body too large (max %d bytes)", maxBytes)
		default:
			return fmt.Errorf("invalid JSON: %v", err)
		}
	}

	// Ensure only one JSON value - reject trailing content
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("body must contain a single JSON value")
	}
	return nil
}

// decodeJSONOrRespond decodes JSON and responds with a useful (but safe) message on failure
func (cah *ContractAPIHandlers) decodeJSONOrRespond(w http.ResponseWriter, r *http.Request, dst interface{}, maxBytes int64) bool {
	err := cah.decodeJSON(w, r, dst, maxBytes)
	if err == nil {
		return true
	}

	// Too large
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		cah.writeJSONError(w, r, http.StatusRequestEntityTooLarge, "payload_too_large",
			fmt.Sprintf("Request body too large (max %d bytes)", maxBytes))
		return false
	}

	// Unknown field
	if strings.HasPrefix(err.Error(), "json: unknown field ") {
		cah.writeJSONError(w, r, http.StatusBadRequest, "unknown_field", err.Error())
		return false
	}

	// Syntax errors
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		cah.writeJSONError(w, r, http.StatusBadRequest, "invalid_json",
			fmt.Sprintf("Malformed JSON at position %d", syntaxErr.Offset))
		return false
	}

	// Type errors
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		cah.writeJSONError(w, r, http.StatusBadRequest, "invalid_type",
			fmt.Sprintf("Invalid value for field %q", typeErr.Field))
		return false
	}

	// Trailing JSON or other decode issues
	cah.writeJSONError(w, r, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
	return false
}

// validateContractAddress validates contract address format (Ethereum-style with optional 0x prefix)
func (cah *ContractAPIHandlers) validateContractAddress(address string) error {
	if len(address) == 0 {
		return fmt.Errorf("address cannot be empty")
	}
	// Use precompiled regex for Ethereum-style addresses (40 hex chars with optional 0x prefix)
	if !contractAddressRegex.MatchString(address) {
		return fmt.Errorf("address must be a valid 40-character hex address (with optional 0x prefix)")
	}
	return nil
}

// writeJSONError writes standardized JSON error responses (includes request_id)
func (cah *ContractAPIHandlers) writeJSONError(w http.ResponseWriter, r *http.Request, statusCode int, errorType, message string) {
	cah.writeJSON(w, statusCode, map[string]interface{}{
		"success":    false,
		"error":      errorType,
		"message":    message,
		"request_id": cah.getRequestID(r),
	})
}

// ===== CONTRACT MANAGEMENT ENDPOINTS =====

// handleContracts handles GET /api/contracts and POST /api/contracts
func (cah *ContractAPIHandlers) handleContracts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cah.handleListContracts(w, r)
	case http.MethodPost:
		cah.handleCreateContract(w, r)
	default:
		cah.writeJSONError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

// handleListContracts lists all contracts with optional filtering
func (cah *ContractAPIHandlers) handleListContracts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	query := r.URL.Query()

	listMsg := &ListContracts{
		Limit:  100, // Default limit
		Offset: 0,   // Default offset
	}

	// Parse filters
	if contractType := query.Get("type"); contractType != "" {
		listMsg.Type = ContractType(contractType)
	}

	if status := query.Get("status"); status != "" {
		listMsg.Status = ContractStatus(status)
	}

	if owner := query.Get("owner"); owner != "" {
		listMsg.Owner = owner
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			listMsg.Limit = limit
		}
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			listMsg.Offset = offset
		}
	}

	// Send request to contract manager
	ctx := cah.system.Root
	future := ctx.RequestFuture(cah.contractManagerPID, listMsg, 30*time.Second)

	result, err := future.Result()
	if err != nil {
		rid := cah.getRequestID(r)
		log.Printf("rid=%s Error listing contracts: %v", rid, err)
		cah.writeJSONError(w, r, http.StatusInternalServerError, "internal_error", "Failed to list contracts")
		return
	}

	contractList, ok := result.(*ContractList)
	if !ok {
		rid := cah.getRequestID(r)
		log.Printf("rid=%s Unexpected response type: %T", rid, result)
		cah.writeJSONError(w, r, http.StatusInternalServerError, "internal_error", "Invalid response from contract manager")
		return
	}

	// Return successful response
	cah.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":     true,
		"contracts":   contractList.Contracts,
		"total_count": contractList.TotalCount,
		"has_more":    contractList.HasMore,
		"limit":       listMsg.Limit,
		"offset":      listMsg.Offset,
	})
}

// handleCreateContract creates a new contract (alias for deploy)
func (cah *ContractAPIHandlers) handleCreateContract(w http.ResponseWriter, r *http.Request) {
	cah.handleDeployContract(w, r)
}

// handleDeployContract handles contract deployment
func (cah *ContractAPIHandlers) handleDeployContract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		cah.writeJSONError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	var deployReq DeployContract
	if err := cah.decodeJSON(w, r, &deployReq, 1<<20); err != nil {
		log.Printf("Error decoding deploy request: %v", err)
		cah.writeJSONError(w, r, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}

	// Validate required fields
	if deployReq.Name == "" {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", "Contract name is required")
		return
	}

	if deployReq.Owner == "" {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", "Contract owner is required")
		return
	}

	// Validate owner address format
	if err := cah.validateContractAddress(deployReq.Owner); err != nil {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", fmt.Sprintf("Invalid owner address: %v", err))
		return
	}

	if deployReq.Type == "" {
		deployReq.Type = ContractTypeLiving // Default to living contract
	}

	// Send deployment request to contract manager
	ctx := cah.system.Root
	future := ctx.RequestFuture(cah.contractManagerPID, &deployReq, 60*time.Second)

	result, err := future.Result()
	if err != nil {
		log.Printf("Error deploying contract: %v", err)
		cah.writeJSONError(w, r, http.StatusInternalServerError, "deployment_error", "Contract deployment failed")
		return
	}

	deployResult, ok := result.(*ContractDeployed)
	if !ok {
		log.Printf("Unexpected response type: %T", result)
		cah.writeJSONError(w, r, http.StatusInternalServerError, "internal_error", "Invalid response from contract manager")
		return
	}

	if !deployResult.Success {
		cah.writeJSONError(w, r, http.StatusBadRequest, "deployment_failed", deployResult.Error)
		return
	}

	// Return successful deployment response
	cah.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success":          true,
		"message":          deployResult.Message,
		"contract":         deployResult.Contract,
		"contract_address": deployResult.Contract.Address,
	})
}

// handleContractByAddress handles GET/PUT/DELETE /api/contracts/{address}
func (cah *ContractAPIHandlers) handleContractByAddress(w http.ResponseWriter, r *http.Request) {
	// Extract contract address from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/contracts/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", "Contract address is required")
		return
	}

	address := parts[0]

	// Validate address format
	if err := cah.validateContractAddress(address); err != nil {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", fmt.Sprintf("Invalid contract address: %v", err))
		return
	}

	switch r.Method {
	case http.MethodGet:
		cah.handleGetContract(w, r, address)
	case http.MethodPut:
		cah.handleUpdateContract(w, r, address)
	case http.MethodDelete:
		cah.handleDeleteContract(w, r, address)
	default:
		cah.writeJSONError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

// handleGetContract retrieves a specific contract
func (cah *ContractAPIHandlers) handleGetContract(w http.ResponseWriter, r *http.Request, address string) {
	getMsg := &GetContract{
		ContractAddress: address,
	}

	ctx := cah.system.Root
	future := ctx.RequestFuture(cah.contractManagerPID, getMsg, 30*time.Second)

	result, err := future.Result()
	if err != nil {
		log.Printf("Error getting contract %s: %v", address, err)
		cah.writeJSONError(w, r, http.StatusInternalServerError, "internal_error", "Failed to get contract")
		return
	}

	getResult, ok := result.(*GetContractResult)
	if !ok {
		log.Printf("Unexpected response type: %T", result)
		cah.writeJSONError(w, r, http.StatusInternalServerError, "internal_error", "Invalid response from contract manager")
		return
	}

	if !getResult.Found {
		cah.writeJSONError(w, r, http.StatusNotFound, "not_found", "Contract not found")
		return
	}

	// Return contract information
	cah.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"contract": getResult.Contract,
	})
}

// handleExecuteContract handles contract function execution
func (cah *ContractAPIHandlers) handleExecuteContract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		cah.writeJSONError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	var execReq struct {
		ContractAddress string          `json:"contract_address"`
		Function        string          `json:"function"`
		Parameters      json.RawMessage `json:"parameters"`
		Caller          string          `json:"caller"`
		GasLimit        int64           `json:"gas_limit"`
		Value           int64           `json:"value,omitempty"`
	}

	if err := cah.decodeJSON(w, r, &execReq, 1<<20); err != nil {
		log.Printf("Error decoding execute request: %v", err)
		cah.writeJSONError(w, r, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}

	// Validate required fields
	if execReq.ContractAddress == "" {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", "Contract address is required")
		return
	}

	if err := cah.validateContractAddress(execReq.ContractAddress); err != nil {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", fmt.Sprintf("Invalid contract address: %v", err))
		return
	}

	if execReq.Function == "" {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", "Function name is required")
		return
	}

	if execReq.Caller == "" {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", "Caller address is required")
		return
	}

	if err := cah.validateContractAddress(execReq.Caller); err != nil {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", fmt.Sprintf("Invalid caller address: %v", err))
		return
	}

	if execReq.GasLimit <= 0 {
		execReq.GasLimit = 1000000 // Default gas limit
	} else if execReq.GasLimit > 10000000 {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", "Gas limit exceeds maximum allowed (10M)")
		return
	}

	// Execute the contract function directly (manager will handle routing)
	executeMsg := &ExecuteContract{
		ContractAddress: execReq.ContractAddress,
		Function:        execReq.Function,
		Parameters:      execReq.Parameters,
		Caller:          execReq.Caller,
		GasLimit:        execReq.GasLimit,
		Value:           execReq.Value,
	}

	// Send execute request to contract manager (it will route to the contract actor)
	ctx := cah.system.Root
	execFuture := ctx.RequestFuture(cah.contractManagerPID, executeMsg, 60*time.Second)

	execResult, err := execFuture.Result()
	if err != nil {
		log.Printf("Error executing contract function: %v", err)
		cah.writeJSONError(w, r, http.StatusInternalServerError, "execution_error", "Contract execution failed")
		return
	}

	executionResult, ok := execResult.(*ContractExecuted)
	if !ok {
		log.Printf("Unexpected response type: %T", execResult)
		cah.writeJSONError(w, r, http.StatusInternalServerError, "internal_error", "Invalid response from contract manager")
		return
	}

	if !executionResult.Success {
		// Use stable error codes instead of string matching
		errorCode := "execution_failed"
		statusCode := http.StatusBadRequest

		// Classify errors by content (temporary until we have proper error codes)
		if strings.Contains(executionResult.Error, "out of gas") {
			errorCode = "out_of_gas"
		} else if strings.Contains(executionResult.Error, "unauthorized") {
			errorCode = "unauthorized"
		} else if strings.Contains(executionResult.Error, "revert") {
			errorCode = "contract_revert"
		} else {
			errorCode = "internal_error"
			statusCode = http.StatusInternalServerError
		}

		cah.writeJSONError(w, r, statusCode, errorCode, executionResult.Error)
		return
	}

	// Return successful execution results with safe JSON handling
	// Safely handle result - use RawMessage only if valid JSON, otherwise string
	var resultAny interface{}
	if json.Valid(executionResult.Execution.Result) {
		resultAny = json.RawMessage(executionResult.Execution.Result)
	} else {
		resultAny = string(executionResult.Execution.Result)
	}

	cah.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"execution": executionResult.Execution,
		"gas_used":  executionResult.Execution.GasUsed,
		"result":    resultAny,
	})
}

// ===== TEMPORAL FEATURES =====

// handleQueryHistory handles temporal history queries
func (cah *ContractAPIHandlers) handleQueryHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		cah.writeJSONError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	// Extract contract address from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/contracts/history/")
	if path == "" {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", "Contract address is required")
		return
	}

	// Validate contract address
	if err := cah.validateContractAddress(path); err != nil {
		cah.writeJSONError(w, r, http.StatusBadRequest, "validation_error", fmt.Sprintf("Invalid contract address: %v", err))
		return
	}

	var queryReq struct {
		Query      json.RawMessage `json:"query"`
		StartBlock int64           `json:"start_block,omitempty"`
		EndBlock   int64           `json:"end_block,omitempty"`
		MaxResults int             `json:"max_results,omitempty"`
	}

	if err := cah.decodeJSON(w, r, &queryReq, 1<<20); err != nil {
		log.Printf("Error decoding history query: %v", err)
		cah.writeJSONError(w, r, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}

	historyMsg := &QueryHistory{
		ContractAddress: path,
		Query:           queryReq.Query,
		StartBlock:      queryReq.StartBlock,
		EndBlock:        queryReq.EndBlock,
		MaxResults:      queryReq.MaxResults,
	}

	if historyMsg.MaxResults <= 0 {
		historyMsg.MaxResults = 100 // Default limit
	}

	ctx := cah.system.Root
	future := ctx.RequestFuture(cah.contractManagerPID, historyMsg, 60*time.Second)

	result, err := future.Result()
	if err != nil {
		log.Printf("Error querying history: %v", err)
		cah.writeJSONError(w, r, http.StatusInternalServerError, "query_error", "History query failed")
		return
	}

	historyResult, ok := result.(*HistoryQueryResult)
	if !ok {
		log.Printf("Unexpected response type: %T", result)
		cah.writeJSONError(w, r, http.StatusInternalServerError, "internal_error", "Invalid response from contract manager")
		return
	}

	// Return query results
	cah.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":     historyResult.Success,
		"results":     historyResult.Results,
		"total_count": historyResult.TotalCount,
		"block_range": historyResult.BlockRange,
		"error":       historyResult.Error,
	})
}

// ===== SYSTEM MONITORING =====

// handleSystemStatus handles system status requests
func (cah *ContractAPIHandlers) handleSystemStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		cah.writeJSONError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	statusMsg := &SystemStatus{}

	ctx := cah.system.Root
	future := ctx.RequestFuture(cah.contractManagerPID, statusMsg, 30*time.Second)

	result, err := future.Result()
	if err != nil {
		log.Printf("Error getting system status: %v", err)
		cah.writeJSONError(w, r, http.StatusInternalServerError, "system_error", "Failed to get system status")
		return
	}

	statusResult, ok := result.(*SystemStatusResult)
	if !ok {
		log.Printf("Unexpected response type: %T", result)
		cah.writeJSONError(w, r, http.StatusInternalServerError, "internal_error", "Invalid response from contract manager")
		return
	}

	// Return system status
	cah.writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":            true,
		"active_contracts":   statusResult.ActiveContracts,
		"sleeping_contracts": statusResult.SleepingContracts,
		"evolving_contracts": statusResult.EvolvingContracts,
		"total_executions":   statusResult.TotalExecutions,
		"average_cpu_usage":  statusResult.AverageCPUUsage,
		"system_health":      statusResult.SystemHealth,
		"uptime":             statusResult.Uptime,
	})
}

// handleCORSPreflight is no longer needed - enableCORS handles OPTIONS automatically

// ===== HTTP SERVER TIMEOUT RECOMMENDATIONS =====
// For production deployment, configure http.Server with timeouts:
//
// server := &http.Server{
//     Addr:              ":8080",
//     Handler:           mux,
//     ReadHeaderTimeout: 10 * time.Second,
//     ReadTimeout:       30 * time.Second,
//     WriteTimeout:      30 * time.Second,
//     IdleTimeout:       120 * time.Second,
// }
//
// This prevents slow clients from tying up connections indefinitely.

// enableCORS adds CORS headers to responses with security best practices
func (cah *ContractAPIHandlers) enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add Vary header to prevent cache mixing
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Restrict to known domains in production
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Placeholder implementations for remaining handlers
func (cah *ContractAPIHandlers) handleContractLifecycle(w http.ResponseWriter, r *http.Request) {
	// Implementation for wake/sleep contract operations
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "Contract lifecycle operations not yet implemented")
}

func (cah *ContractAPIHandlers) handleUpgradeContract(w http.ResponseWriter, r *http.Request) {
	// Implementation for contract upgrades
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "Contract upgrade operations not yet implemented")
}

func (cah *ContractAPIHandlers) handleEvolveContract(w http.ResponseWriter, r *http.Request) {
	// Implementation for contract evolution
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "Contract evolution operations not yet implemented")
}

func (cah *ContractAPIHandlers) handleProposeCollaboration(w http.ResponseWriter, r *http.Request) {
	// Implementation for contract collaboration
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "Contract collaboration not yet implemented")
}

func (cah *ContractAPIHandlers) handleCollaborationManagement(w http.ResponseWriter, r *http.Request) {
	// Implementation for collaboration management
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "Collaboration management not yet implemented")
}

func (cah *ContractAPIHandlers) handleQueryContractState(w http.ResponseWriter, r *http.Request) {
	// Implementation for state queries
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "Contract state queries not yet implemented")
}

func (cah *ContractAPIHandlers) handlePredictBehavior(w http.ResponseWriter, r *http.Request) {
	// Implementation for behavior prediction
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "Behavior prediction not yet implemented")
}

func (cah *ContractAPIHandlers) handleEcosystems(w http.ResponseWriter, r *http.Request) {
	// Implementation for ecosystem management
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "Ecosystem management not yet implemented")
}

func (cah *ContractAPIHandlers) handleEcosystemByID(w http.ResponseWriter, r *http.Request) {
	// Implementation for specific ecosystem operations
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "Ecosystem operations not yet implemented")
}

func (cah *ContractAPIHandlers) handleSystemMetrics(w http.ResponseWriter, r *http.Request) {
	// Implementation for system metrics
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "System metrics not yet implemented")
}

func (cah *ContractAPIHandlers) handleSystemHealth(w http.ResponseWriter, r *http.Request) {
	// Implementation for system health
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "System health monitoring not yet implemented")
}

func (cah *ContractAPIHandlers) handleUpdateContract(w http.ResponseWriter, r *http.Request, address string) {
	// Implementation for contract updates
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "Contract updates not yet implemented")
}

func (cah *ContractAPIHandlers) handleDeleteContract(w http.ResponseWriter, r *http.Request, address string) {
	// Implementation for contract deletion
	cah.writeJSONError(w, r, http.StatusNotImplemented, "not_implemented", "Contract deletion not yet implemented")
}
