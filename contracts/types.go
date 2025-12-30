package contracts

import (
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

// ContractType represents the different types of innovative contracts
type ContractType string

const (
	ContractTypeLiving     ContractType = "living"     // Persistent actor-based contracts
	ContractTypeTemporal   ContractType = "temporal"   // Time-aware contracts
	ContractTypeMorphic    ContractType = "morphic"    // Self-adapting contracts  
	ContractTypeSymbiotic  ContractType = "symbiotic"  // Inter-dependent contracts
	ContractTypeQuantum    ContractType = "quantum"    // Probabilistic contracts
	ContractTypeMeta       ContractType = "meta"       // Contracts about contracts
)

// ContractStatus represents the lifecycle state of a contract
type ContractStatus string

const (
	ContractStatusDeploying ContractStatus = "deploying"
	ContractStatusActive    ContractStatus = "active"
	ContractStatusSleeping  ContractStatus = "sleeping"  // Living contracts can sleep
	ContractStatusEvolving  ContractStatus = "evolving"  // Morphic contracts evolving
	ContractStatusMating    ContractStatus = "mating"    // Symbiotic contracts collaborating
	ContractStatusArchived  ContractStatus = "archived"
	ContractStatusFailed    ContractStatus = "failed"
)

// Contract represents a smart contract in the Xenese DLT system
type Contract struct {
	ID          uuid.UUID       `json:"id"`
	Address     string          `json:"address"`
	Name        string          `json:"name"`
	Type        ContractType    `json:"type"`
	Status      ContractStatus  `json:"status"`
	Owner       string          `json:"owner"`
	Version     string          `json:"version"`
	
	// Code and execution
	SourceCode  string          `json:"source_code"`
	CompiledPlugin string       `json:"compiled_plugin"` // Path to Go plugin
	ABI         json.RawMessage `json:"abi"`
	
	// Living contract features
	DNA         ContractDNA     `json:"dna"`         // Genetic programming data
	Memory      ContractMemory  `json:"memory"`      // Persistent memory
	Behavior    ContractBehavior `json:"behavior"`   // Learned behaviors
	
	// State management
	StateRoot   string          `json:"state_root"`
	StateSize   int64           `json:"state_size"`
	
	// Temporal features
	TimeAware   bool            `json:"time_aware"`
	HistoryDepth int            `json:"history_depth"` // How far back it can query
	
	// Performance metrics
	ExecutionCount   int64   `json:"execution_count"`
	AverageGasUsed   int64   `json:"average_gas_used"`
	SuccessRate      float64 `json:"success_rate"`
	AdaptationScore  float64 `json:"adaptation_score"`
	
	// Metadata
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	LastActive  time.Time       `json:"last_active"`
}

// ContractDNA represents the genetic programming aspects of contracts
type ContractDNA struct {
	Genes       []Gene          `json:"genes"`        // Behavioral traits
	Generation  int             `json:"generation"`   // Evolution generation
	Parents     []string        `json:"parents"`      // Parent contract addresses
	Mutations   []Mutation      `json:"mutations"`    // Applied mutations
	Fitness     float64         `json:"fitness"`      // Success fitness score
}

// Gene represents a behavioral trait in contract DNA
type Gene struct {
	ID          string          `json:"id"`
	Type        GeneType        `json:"type"`
	Expression  json.RawMessage `json:"expression"`
	Dominance   float64         `json:"dominance"`    // 0.0 to 1.0
	Stability   float64         `json:"stability"`    // Resistance to mutation
}

type GeneType string

const (
	GeneTypeOptimization  GeneType = "optimization"   // Performance optimization traits
	GeneTypeSecurity      GeneType = "security"       // Security enhancement traits  
	GeneTypeCollaboration GeneType = "collaboration"  // Inter-contract cooperation
	GeneTypeAdaptation    GeneType = "adaptation"     // Environmental adaptation
	GeneTypeResilience    GeneType = "resilience"     // Error recovery traits
)

// Mutation represents a genetic change in contract behavior
type Mutation struct {
	ID          string          `json:"id"`
	Type        MutationType    `json:"type"`
	GeneID      string          `json:"gene_id"`
	Change      json.RawMessage `json:"change"`
	Beneficial  *bool           `json:"beneficial"`   // nil if unknown
	AppliedAt   time.Time       `json:"applied_at"`
}

type MutationType string

const (
	MutationTypePointChange MutationType = "point_change"    // Single value change
	MutationTypeInsertion   MutationType = "insertion"       // Add new behavior
	MutationTypeDeletion    MutationType = "deletion"        // Remove behavior
	MutationTypeDuplication MutationType = "duplication"     // Copy existing behavior
	MutationTypeRecombination MutationType = "recombination" // Combine behaviors
)

// ContractMemory represents the persistent memory of a living contract
type ContractMemory struct {
	ShortTerm   map[string]interface{} `json:"short_term"`   // Volatile memory
	LongTerm    map[string]interface{} `json:"long_term"`    // Persistent memory
	Patterns    []MemoryPattern        `json:"patterns"`     // Learned patterns
	Experiences []Experience           `json:"experiences"`  // Historical experiences
}

// MemoryPattern represents a learned behavioral pattern
type MemoryPattern struct {
	ID          string                 `json:"id"`
	Type        PatternType            `json:"type"`
	Trigger     json.RawMessage        `json:"trigger"`      // What activates this pattern
	Response    json.RawMessage        `json:"response"`     // How to respond
	Confidence  float64                `json:"confidence"`   // How reliable this pattern is
	Usage       int                    `json:"usage"`        // How often it's been used
	LastUsed    time.Time              `json:"last_used"`
}

type PatternType string

const (
	PatternTypeExecution    PatternType = "execution"     // Execution optimization patterns
	PatternTypeInteraction  PatternType = "interaction"   // Inter-contract interaction patterns
	PatternTypeError        PatternType = "error"         // Error handling patterns
	PatternTypeResource     PatternType = "resource"      // Resource usage patterns
	PatternTypeTemporal     PatternType = "temporal"      // Time-based patterns
)

// Experience represents a single learning experience
type Experience struct {
	ID          string                 `json:"id"`
	Context     json.RawMessage        `json:"context"`      // Situation context
	Action      json.RawMessage        `json:"action"`       // Action taken
	Result      json.RawMessage        `json:"result"`       // Outcome
	Success     bool                   `json:"success"`      // Whether it was successful
	Timestamp   time.Time              `json:"timestamp"`
}

// ContractBehavior represents learned behaviors of a contract
type ContractBehavior struct {
	Traits          []BehaviorTrait    `json:"traits"`
	Adaptations     []Adaptation       `json:"adaptations"`
	Collaborations  []Collaboration    `json:"collaborations"`
	Predictions     []Prediction       `json:"predictions"`
}

// BehaviorTrait represents a learned behavioral characteristic
type BehaviorTrait struct {
	ID          string          `json:"id"`
	Type        TraitType       `json:"type"`
	Strength    float64         `json:"strength"`     // How strong this trait is
	Stability   float64         `json:"stability"`    // How stable this trait is
	Evidence    []string        `json:"evidence"`     // Experience IDs that support this
}

type TraitType string

const (
	TraitTypeAggressive     TraitType = "aggressive"      // Tends to use more resources
	TraitTypeConservative   TraitType = "conservative"    // Tends to be cautious
	TraitTypeCollaborative  TraitType = "collaborative"   // Works well with others
	TraitTypeIndependent    TraitType = "independent"     // Prefers to work alone
	TraitTypeAdaptable      TraitType = "adaptable"       // Changes behavior quickly
	TraitTypeStable         TraitType = "stable"          // Maintains consistent behavior
)

// Adaptation represents a behavioral adaptation
type Adaptation struct {
	ID          string          `json:"id"`
	Trigger     string          `json:"trigger"`      // What caused this adaptation
	Change      string          `json:"change"`       // What changed
	Impact      float64         `json:"impact"`       // How much it improved performance
	Timestamp   time.Time       `json:"timestamp"`
}

// Collaboration represents inter-contract collaboration
type Collaboration struct {
	ID              string          `json:"id"`
	PartnerAddress  string          `json:"partner_address"`
	Type            CollaborationType `json:"type"`
	Success         bool            `json:"success"`
	Benefit         float64         `json:"benefit"`      // Mutual benefit score
	Timestamp       time.Time       `json:"timestamp"`
}

type CollaborationType string

const (
	CollaborationTypeDataSharing    CollaborationType = "data_sharing"
	CollaborationTypeResourceSharing CollaborationType = "resource_sharing"
	CollaborationTypeCoExecution     CollaborationType = "co_execution"
	CollaborationTypeEventChain      CollaborationType = "event_chain"
	CollaborationTypeMerge          CollaborationType = "merge"        // Temporary merger
)

// Prediction represents a contract's prediction about future events
type Prediction struct {
	ID          string          `json:"id"`
	Type        PredictionType  `json:"type"`
	Prediction  json.RawMessage `json:"prediction"`
	Confidence  float64         `json:"confidence"`
	MadeAt      time.Time       `json:"made_at"`
	ValidUntil  time.Time       `json:"valid_until"`
	Verified    *bool           `json:"verified"`     // nil if not yet verified
}

type PredictionType string

const (
	PredictionTypeExecution     PredictionType = "execution"      // Future execution patterns
	PredictionTypeResource      PredictionType = "resource"       // Resource usage predictions
	PredictionTypeInteraction   PredictionType = "interaction"    // Future interactions
	PredictionTypeMarket        PredictionType = "market"         // Market/economic predictions
	PredictionTypeNetwork       PredictionType = "network"        // Network behavior predictions
)

// ContractExecution represents a contract execution instance
type ContractExecution struct {
	ID              uuid.UUID       `json:"id"`
	ContractAddress string          `json:"contract_address"`
	Function        string          `json:"function"`
	Parameters      json.RawMessage `json:"parameters"`
	Caller          string          `json:"caller"`
	
	// Execution context
	BlockHeight     int64           `json:"block_height"`
	Timestamp       time.Time       `json:"timestamp"`
	GasLimit        int64           `json:"gas_limit"`
	GasUsed         int64           `json:"gas_used"`
	
	// Results
	Status          ExecutionStatus `json:"status"`
	Result          json.RawMessage `json:"result"`
	Events          []ContractEvent `json:"events"`
	StateChanges    []StateChange   `json:"state_changes"`
	Error           string          `json:"error,omitempty"`
	
	// Learning data
	LearningData    json.RawMessage `json:"learning_data,omitempty"`
	AdaptationTriggers []string     `json:"adaptation_triggers,omitempty"`
	
	// Performance metrics
	CPUUsage        float64         `json:"cpu_usage"`
	MemoryUsage     int64           `json:"memory_usage"`
	ExecutionTime   time.Duration   `json:"execution_time"`
}

type ExecutionStatus string

const (
	ExecutionStatusPending    ExecutionStatus = "pending"
	ExecutionStatusExecuting  ExecutionStatus = "executing"
	ExecutionStatusCompleted  ExecutionStatus = "completed"
	ExecutionStatusFailed     ExecutionStatus = "failed"
	ExecutionStatusReverted   ExecutionStatus = "reverted"
	ExecutionStatusTimeout    ExecutionStatus = "timeout"
)

// ContractEvent represents an event emitted by a contract
type ContractEvent struct {
	ID              uuid.UUID       `json:"id"`
	ContractAddress string          `json:"contract_address"`
	EventName       string          `json:"event_name"`
	Data            json.RawMessage `json:"data"`
	Indexed         []string        `json:"indexed"`      // Indexed fields for querying
	Timestamp       time.Time       `json:"timestamp"`
	BlockHeight     int64           `json:"block_height"`
}

// StateChange represents a change in contract state
type StateChange struct {
	Key         string          `json:"key"`
	OldValue    json.RawMessage `json:"old_value"`
	NewValue    json.RawMessage `json:"new_value"`
	Timestamp   time.Time       `json:"timestamp"`
}

// ContractRegistry represents the registry of all contracts
type ContractRegistry struct {
	Contracts       map[string]*Contract    `json:"contracts"`        // Address -> Contract
	Dependencies    map[string][]string     `json:"dependencies"`     // Address -> Dependencies
	Collaborations  map[string][]string     `json:"collaborations"`   // Address -> Collaborators
	Ecosystems      []ContractEcosystem     `json:"ecosystems"`       // Contract ecosystems
	UpdatedAt       time.Time               `json:"updated_at"`
}

// ContractEcosystem represents a group of collaborating contracts
type ContractEcosystem struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Contracts   []string        `json:"contracts"`        // Contract addresses
	Type        EcosystemType   `json:"type"`
	Health      float64         `json:"health"`           // Overall ecosystem health
	Synergy     float64         `json:"synergy"`          // How well contracts work together
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type EcosystemType string

const (
	EcosystemTypeSymbiotic      EcosystemType = "symbiotic"       // Mutually beneficial
	EcosystemTypeCompetitive    EcosystemType = "competitive"     // Competing for resources
	EcosystemTypeParasitic      EcosystemType = "parasitic"       // One benefits, other suffers
	EcosystemTypeCommensalistic EcosystemType = "commensalistic"  // One benefits, other neutral
	EcosystemTypeNeutral        EcosystemType = "neutral"         // No significant interaction
)
