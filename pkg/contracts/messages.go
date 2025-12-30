package contracts

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ===== CONTRACT MANAGEMENT MESSAGES =====

// Background task tick messages for actor-based scheduling
type MaintenanceTick struct{}
type HealthTick struct{}
type PerfTick struct{}
type EcosystemTick struct{}

// DeployContract message to deploy a new smart contract
type DeployContract struct {
	Name         string          `json:"name"`
	Type         ContractType    `json:"type"`
	SourceCode   string          `json:"source_code"`
	Owner        string          `json:"owner"`
	InitParams   json.RawMessage `json:"init_params,omitempty"`
	TimeAware    bool            `json:"time_aware"`
	HistoryDepth int             `json:"history_depth"`
}

// ContractDeployed response to successful contract deployment
type ContractDeployed struct {
	Contract *Contract `json:"contract"`
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Error    string    `json:"error,omitempty"`
}

// ExecuteContract message to execute a contract function
type ExecuteContract struct {
	ContractAddress string          `json:"contract_address"`
	Function        string          `json:"function"`
	Parameters      json.RawMessage `json:"parameters"`
	Caller          string          `json:"caller"`
	GasLimit        int64           `json:"gas_limit"`
	Value           int64           `json:"value,omitempty"`
}

// ContractExecuted response to contract execution
type ContractExecuted struct {
	Execution *ContractExecution `json:"execution"`
	Success   bool               `json:"success"`
	Message   string             `json:"message"`
	Error     string             `json:"error,omitempty"`
}

// UpgradeContract message to upgrade/evolve a contract
type UpgradeContract struct {
	ContractAddress string     `json:"contract_address"`
	NewVersion      string     `json:"new_version"`
	SourceCode      string     `json:"source_code,omitempty"`
	MigrationScript string     `json:"migration_script,omitempty"`
	Mutations       []Mutation `json:"mutations,omitempty"`
}

// ContractUpgraded response to contract upgrade
type ContractUpgraded struct {
	Contract *Contract `json:"contract"`
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Error    string    `json:"error,omitempty"`
}

// ===== LIVING CONTRACT MESSAGES =====

// WakeContract message to wake up a sleeping contract
type WakeContract struct {
	ContractAddress string          `json:"contract_address"`
	Reason          string          `json:"reason"`
	Context         json.RawMessage `json:"context,omitempty"`
}

// ContractAwake response when contract wakes up
type ContractAwake struct {
	ContractAddress string         `json:"contract_address"`
	Status          ContractStatus `json:"status"`
	Message         string         `json:"message"`
}

// SleepContract message to put contract to sleep
type SleepContract struct {
	ContractAddress string          `json:"contract_address"`
	Duration        time.Duration   `json:"duration,omitempty"`
	Condition       json.RawMessage `json:"condition,omitempty"`
}

// ContractSleeping response when contract goes to sleep
type ContractSleeping struct {
	ContractAddress string          `json:"contract_address"`
	Status          ContractStatus  `json:"status"`
	WakeCondition   json.RawMessage `json:"wake_condition,omitempty"`
}

// ===== EVOLUTIONARY MESSAGES =====

// TriggerEvolution message to start contract evolution
type TriggerEvolution struct {
	ContractAddress string          `json:"contract_address"`
	EvolutionType   EvolutionType   `json:"evolution_type"`
	Parameters      json.RawMessage `json:"parameters,omitempty"`
}

type EvolutionType string

const (
	EvolutionTypeMutation      EvolutionType = "mutation"      // Random beneficial change
	EvolutionTypeOptimization  EvolutionType = "optimization"  // Performance optimization
	EvolutionTypeAdaptation    EvolutionType = "adaptation"    // Environmental adaptation
	EvolutionTypeCollaboration EvolutionType = "collaboration" // Improve collaboration
	EvolutionTypeBreeding      EvolutionType = "breeding"      // Genetic combination
)

// EvolutionStarted response when evolution begins
type EvolutionStarted struct {
	ContractAddress string        `json:"contract_address"`
	EvolutionID     uuid.UUID     `json:"evolution_id"`
	EstimatedTime   time.Duration `json:"estimated_time"`
	Success         bool          `json:"success"`
}

// EvolutionCompleted message when evolution finishes
type EvolutionCompleted struct {
	ContractAddress    string    `json:"contract_address"`
	EvolutionID        uuid.UUID `json:"evolution_id"`
	Success            bool      `json:"success"`
	Changes            []string  `json:"changes"`
	FitnessImprovement float64   `json:"fitness_improvement"`
}

// ===== COLLABORATION MESSAGES =====

// ProposeCollaboration message to propose inter-contract collaboration
type ProposeCollaboration struct {
	FromContract string            `json:"from_contract"`
	ToContract   string            `json:"to_contract"`
	Type         CollaborationType `json:"type"`
	Proposal     json.RawMessage   `json:"proposal"`
	Duration     time.Duration     `json:"duration,omitempty"`
	Terms        json.RawMessage   `json:"terms,omitempty"`
}

// CollaborationProposed represents a collaboration proposal response
type CollaborationProposed struct {
	CollaborationID uuid.UUID `json:"collaboration_id"`
	Status          string    `json:"status"`
	Message         string    `json:"message"`
}

// CollaborationAccepted represents a collaboration acceptance response
type CollaborationAccepted struct {
	CollaborationID uuid.UUID `json:"collaboration_id"`
	Success         bool      `json:"success"`
	Message         string    `json:"message"`
}

// CollaborationRejected represents a collaboration rejection response
type CollaborationRejected struct {
	CollaborationID uuid.UUID `json:"collaboration_id"`
	Success         bool      `json:"success"`
	Message         string    `json:"message"`
	Reason          string    `json:"reason"`
}

// AcceptCollaboration message to accept a collaboration
type AcceptCollaboration struct {
	CollaborationID uuid.UUID       `json:"collaboration_id"`
	Terms           json.RawMessage `json:"terms,omitempty"`
}

// RejectCollaboration message to reject a collaboration
type RejectCollaboration struct {
	CollaborationID uuid.UUID `json:"collaboration_id"`
	Reason          string    `json:"reason"`
}

// CollaborationStatusChanged message when collaboration status changes
type CollaborationStatusChanged struct {
	CollaborationID uuid.UUID           `json:"collaboration_id"`
	Status          CollaborationStatus `json:"status"`
	Participants    []string            `json:"participants"`
}

type CollaborationStatus string

const (
	CollaborationStatusProposed   CollaborationStatus = "proposed"
	CollaborationStatusAccepted   CollaborationStatus = "accepted"
	CollaborationStatusActive     CollaborationStatus = "active"
	CollaborationStatusCompleted  CollaborationStatus = "completed"
	CollaborationStatusRejected   CollaborationStatus = "rejected"
	CollaborationStatusTerminated CollaborationStatus = "terminated"
)

// ===== TEMPORAL MESSAGES =====

// QueryHistory message to query historical blockchain state
type QueryHistory struct {
	ContractAddress string          `json:"contract_address"`
	Query           json.RawMessage `json:"query"`
	StartBlock      int64           `json:"start_block,omitempty"`
	EndBlock        int64           `json:"end_block,omitempty"`
	MaxResults      int             `json:"max_results,omitempty"`
}

// HistoryQueryResult response to historical query
type HistoryQueryResult struct {
	Results    []json.RawMessage `json:"results"`
	TotalCount int               `json:"total_count"`
	BlockRange string            `json:"block_range"`
	Success    bool              `json:"success"`
	Error      string            `json:"error,omitempty"`
}

// PredictFuture message for temporal prediction
type PredictFuture struct {
	ContractAddress string          `json:"contract_address"`
	PredictionType  PredictionType  `json:"prediction_type"`
	Context         json.RawMessage `json:"context"`
	TimeHorizon     time.Duration   `json:"time_horizon"`
}

// PredictionMade response to prediction request
type PredictionMade struct {
	Prediction *Prediction `json:"prediction"`
	Success    bool        `json:"success"`
	Confidence float64     `json:"confidence"`
	Error      string      `json:"error,omitempty"`
}

// ===== EVENT INTEGRATION MESSAGES =====

// SubscribeToEvents message to subscribe to event patterns
type SubscribeToEvents struct {
	ContractAddress string          `json:"contract_address"`
	EventPattern    json.RawMessage `json:"event_pattern"`
	Callback        string          `json:"callback"` // Function to call
	FilterCriteria  json.RawMessage `json:"filter_criteria,omitempty"`
}

// EventSubscribed response to event subscription
type EventSubscribed struct {
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Success        bool      `json:"success"`
	Message        string    `json:"message"`
}

// EventTriggered message when subscribed event occurs
type EventTriggered struct {
	SubscriptionID  uuid.UUID       `json:"subscription_id"`
	ContractAddress string          `json:"contract_address"` // subscriber
	EventType       string          `json:"event_type"`
	Emitter         string          `json:"emitter,omitempty"`
	Event           json.RawMessage `json:"event"`
	Timestamp       time.Time       `json:"timestamp"`
	Callback        string          `json:"callback"`
}

// EmitEvent message to emit a contract event
type EmitEvent struct {
	ContractAddress string          `json:"contract_address"`
	EventName       string          `json:"event_name"`
	Data            json.RawMessage `json:"data"`
	Indexed         []string        `json:"indexed"`
}

// EventEmitted response to event emission
type EventEmitted struct {
	Event   *ContractEvent `json:"event"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
}

// ===== LEARNING AND ADAPTATION MESSAGES =====

// LearnFromExperience message to add learning data
type LearnFromExperience struct {
	ContractAddress string          `json:"contract_address"`
	Experience      *Experience     `json:"experience"`
	Context         json.RawMessage `json:"context,omitempty"`
}

// LearningProcessed response to learning input
type LearningProcessed struct {
	Success         bool     `json:"success"`
	PatternsFound   []string `json:"patterns_found"`
	AdaptationsMade []string `json:"adaptations_made"`
	Message         string   `json:"message"`
}

// AnalyzeBehavior message to analyze contract behavior
type AnalyzeBehavior struct {
	ContractAddress string        `json:"contract_address"`
	TimeWindow      time.Duration `json:"time_window"`
	Metrics         []string      `json:"metrics"`
}

// BehaviorAnalysis response to behavior analysis
type BehaviorAnalysis struct {
	Analysis        json.RawMessage `json:"analysis"`
	Recommendations []string        `json:"recommendations"`
	Traits          []BehaviorTrait `json:"traits"`
	Success         bool            `json:"success"`
}

// ===== STATE MANAGEMENT MESSAGES =====

// QueryContractState message to query contract state
type QueryContractState struct {
	ContractAddress string          `json:"contract_address"`
	StateQuery      json.RawMessage `json:"state_query"`
	IncludeHistory  bool            `json:"include_history"`
}

// ContractStateResult response to state query
type ContractStateResult struct {
	State        json.RawMessage `json:"state"`
	StateRoot    string          `json:"state_root"`
	LastModified time.Time       `json:"last_modified"`
	Success      bool            `json:"success"`
	Error        string          `json:"error,omitempty"`
}

// UpdateContractState message to update contract state
type UpdateContractState struct {
	ContractAddress string        `json:"contract_address"`
	StateChanges    []StateChange `json:"state_changes"`
	Reason          string        `json:"reason"`
}

// ContractStateUpdated response to state update
type ContractStateUpdated struct {
	NewStateRoot string        `json:"new_state_root"`
	Changes      []StateChange `json:"changes"`
	Success      bool          `json:"success"`
	Message      string        `json:"message"`
}

// ===== REGISTRY MESSAGES =====

// RegisterContract message to register a new contract
type RegisterContract struct {
	Contract *Contract `json:"contract"`
}

// ContractRegistered response to contract registration
type ContractRegistered struct {
	ContractAddress string `json:"contract_address"`
	Success         bool   `json:"success"`
	Message         string `json:"message"`
}

// GetContract message to retrieve contract information
type GetContract struct {
	ContractAddress string `json:"contract_address"`
}

// GetContractResult response to contract retrieval
type GetContractResult struct {
	Contract *Contract `json:"contract"`
	Found    bool      `json:"found"`
}

// ListContracts message to list contracts with filters
type ListContracts struct {
	Type   ContractType   `json:"type,omitempty"`
	Status ContractStatus `json:"status,omitempty"`
	Owner  string         `json:"owner,omitempty"`
	Limit  int            `json:"limit,omitempty"`
	Offset int            `json:"offset,omitempty"`
}

// ContractList response to contract listing
type ContractList struct {
	Contracts  []*Contract `json:"contracts"`
	TotalCount int         `json:"total_count"`
	HasMore    bool        `json:"has_more"`
}

// ===== ECOSYSTEM MESSAGES =====

// CreateEcosystem message to create a contract ecosystem
type CreateEcosystem struct {
	Name      string          `json:"name"`
	Type      EcosystemType   `json:"type"`
	Contracts []string        `json:"contracts"`
	Rules     json.RawMessage `json:"rules,omitempty"`
}

// EcosystemCreated response to ecosystem creation
type EcosystemCreated struct {
	Ecosystem *ContractEcosystem `json:"ecosystem"`
	Success   bool               `json:"success"`
	Message   string             `json:"message"`
}

// AnalyzeEcosystem message to analyze ecosystem health
type AnalyzeEcosystem struct {
	EcosystemID string   `json:"ecosystem_id"`
	Metrics     []string `json:"metrics"`
}

// EcosystemAnalysis response to ecosystem analysis
type EcosystemAnalysis struct {
	Health          float64  `json:"health"`
	Synergy         float64  `json:"synergy"`
	Bottlenecks     []string `json:"bottlenecks"`
	Recommendations []string `json:"recommendations"`
	Success         bool     `json:"success"`
}

// ===== SYSTEM MESSAGES =====

// SystemStatus message to check system health
type SystemStatus struct{}

// SystemStatusResult response to system status check
type SystemStatusResult struct {
	ActiveContracts   int           `json:"active_contracts"`
	SleepingContracts int           `json:"sleeping_contracts"`
	EvolvingContracts int           `json:"evolving_contracts"`
	TotalExecutions   int64         `json:"total_executions"`
	AverageCPUUsage   float64       `json:"average_cpu_usage"`
	SystemHealth      float64       `json:"system_health"`
	Uptime            time.Duration `json:"uptime"`
}
