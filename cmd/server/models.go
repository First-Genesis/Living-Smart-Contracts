package main

import "time"

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status" example:"healthy"`
	Service   string    `json:"service" example:"living-smart-contracts"`
	Timestamp time.Time `json:"timestamp" example:"2026-04-30T20:08:31Z"`
}

// APIInfoResponse represents the API information response
type APIInfoResponse struct {
	Name        string            `json:"name" example:"Living Smart Contracts"`
	Version     string            `json:"version" example:"1.0.0"`
	Description string            `json:"description" example:"Revolutionary blockchain platform with evolutionary intelligence"`
	Features    []string          `json:"features"`
	Endpoints   map[string]string `json:"endpoints"`
}

// Contract represents a smart contract
type Contract struct {
	ID              string                 `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Address         string                 `json:"address" example:"0xabcdef1234567890abcdef1234567890abcdef12"`
	Name            string                 `json:"name" example:"adaptive_trading_contract"`
	Type            string                 `json:"type" example:"living"`
	Status          string                 `json:"status" example:"active"`
	Owner           string                 `json:"owner" example:"0x1234567890123456789012345678901234567890"`
	Version         string                 `json:"version" example:"1.0.0"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	DNA             ContractDNA            `json:"dna"`
	Memory          ContractMemory         `json:"memory"`
	Collaborations  []Collaboration        `json:"collaborations"`
	Metrics         PerformanceMetrics     `json:"performance_metrics"`
}

// ContractDNA represents the genetic information of a contract
type ContractDNA struct {
	Genes      []Gene    `json:"genes"`
	Generation int       `json:"generation" example:"5"`
	Parents    []string  `json:"parents"`
	Mutations  []string  `json:"mutations"`
	Fitness    float64   `json:"fitness" example:"0.934"`
}

// Gene represents a behavioral trait in contract DNA
type Gene struct {
	ID         string  `json:"id" example:"optimization_001"`
	Type       string  `json:"type" example:"optimization"`
	Dominance  float64 `json:"dominance" example:"0.92"`
	Stability  float64 `json:"stability" example:"0.88"`
	Expression string  `json:"expression" example:"performance_boost"`
}

// ContractMemory represents the memory system of a contract
type ContractMemory struct {
	ShortTerm   map[string]interface{} `json:"short_term"`
	LongTerm    map[string]interface{} `json:"long_term"`
	Patterns    []MemoryPattern        `json:"patterns"`
	Experiences []Experience           `json:"experiences"`
}

// MemoryPattern represents a learned behavioral pattern
type MemoryPattern struct {
	ID         string  `json:"id" example:"pattern_001"`
	Type       string  `json:"type" example:"execution"`
	Confidence float64 `json:"confidence" example:"0.95"`
	Usage      int     `json:"usage" example:"1250"`
}

// Experience represents a historical experience
type Experience struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Context   map[string]interface{} `json:"context"`
	Outcome   string                 `json:"outcome"`
	Learning  string                 `json:"learning"`
}

// Collaboration represents a partnership between contracts
type Collaboration struct {
	ID            string    `json:"id" example:"collab_001"`
	PartnerAddress string   `json:"partner_address" example:"0x9876543210987654321098765432109876543210"`
	Type          string    `json:"type" example:"resource_sharing"`
	Status        string    `json:"status" example:"active"`
	BenefitScore  float64   `json:"benefit_score" example:"0.18"`
	CreatedAt     time.Time `json:"created_at"`
}

// PerformanceMetrics represents contract performance data
type PerformanceMetrics struct {
	ExecutionCount   int64   `json:"execution_count" example:"15420"`
	SuccessRate      float64 `json:"success_rate" example:"0.987"`
	AverageGasUsed   int64   `json:"average_gas_used" example:"45000"`
	AdaptationScore  float64 `json:"adaptation_score" example:"0.923"`
	LastOptimized    time.Time `json:"last_optimized"`
}

// DeployContractRequest represents a contract deployment request
type DeployContractRequest struct {
	Name         string                 `json:"name" example:"adaptive_trading_contract" binding:"required"`
	Type         string                 `json:"type" example:"living" binding:"required"`
	SourceCode   string                 `json:"source_code" example:"contract AdaptiveTrader { ... }" binding:"required"`
	Owner        string                 `json:"owner" example:"0x1234567890123456789012345678901234567890" binding:"required"`
	InitParams   map[string]interface{} `json:"init_params"`
	TimeAware    bool                   `json:"time_aware" example:"true"`
	HistoryDepth int                    `json:"history_depth" example:"1000"`
}

// ExecuteContractRequest represents a contract execution request
type ExecuteContractRequest struct {
	ContractAddress string                 `json:"contract_address" example:"0xabcdef1234567890abcdef1234567890abcdef12" binding:"required"`
	Function        string                 `json:"function" example:"trade" binding:"required"`
	Parameters      map[string]interface{} `json:"parameters" binding:"required"`
	Caller          string                 `json:"caller" example:"0x1234567890123456789012345678901234567890" binding:"required"`
	GasLimit        int64                  `json:"gas_limit" example:"500000"`
}

// EvolutionRequest represents a contract evolution request
type EvolutionRequest struct {
	EvolutionType string                 `json:"evolution_type" example:"optimization" binding:"required"`
	Parameters    map[string]interface{} `json:"parameters"`
}

// CollaborationRequest represents a collaboration proposal request
type CollaborationRequest struct {
	ProposerAddress string                 `json:"proposer_address" example:"0xabcdef1234567890abcdef1234567890abcdef12" binding:"required"`
	TargetAddress   string                 `json:"target_address" example:"0x9876543210987654321098765432109876543210" binding:"required"`
	Type            string                 `json:"collaboration_type" example:"resource_sharing" binding:"required"`
	Terms           map[string]interface{} `json:"terms" binding:"required"`
	ExpectedBenefit float64                `json:"expected_benefit" example:"0.25"`
}

// ContractsResponse represents the contracts list response
type ContractsResponse struct {
	Message   string     `json:"message" example:"Living Smart Contracts API"`
	Status    string     `json:"status" example:"ready"`
	Contracts []Contract `json:"contracts"`
	Total     int        `json:"total" example:"0"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Message string `json:"message" example:"The provided parameters are invalid"`
	Code    int    `json:"code" example:"400"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}
