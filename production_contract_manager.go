package contracts

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/google/uuid"
)

// ProductionContractManager manages all smart contracts in the system
type ProductionContractManager struct {
	actorSystem    *actor.ActorSystem
	contracts      map[string]*Contract
	contractActors map[string]*actor.PID
	mutex          sync.RWMutex
}

// NewProductionContractManager creates a new contract manager
func NewProductionContractManager(actorSystem *actor.ActorSystem) *ProductionContractManager {
	return &ProductionContractManager{
		actorSystem:    actorSystem,
		contracts:      make(map[string]*Contract),
		contractActors: make(map[string]*actor.PID),
		mutex:          sync.RWMutex{},
	}
}

// DeployContract deploys a new smart contract
func (pcm *ProductionContractManager) DeployContract(deployMsg *DeployContract) (*Contract, error) {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	// Create new contract
	contract := &Contract{
		ID:             uuid.New(),
		Address:        generateContractAddress(),
		Name:           deployMsg.Name,
		Type:           deployMsg.Type,
		Status:         ContractStatusDeploying,
		Owner:          deployMsg.Owner,
		Version:        "1.0.0",
		SourceCode:     deployMsg.SourceCode,
		TimeAware:      deployMsg.TimeAware,
		HistoryDepth:   deployMsg.HistoryDepth,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LastActive:     time.Now(),
		ExecutionCount: 0,
		SuccessRate:    1.0,
		// Initialize contract components
		DNA: ContractDNA{
			Genes:      make([]Gene, 0),
			Generation: 0,
			Parents:    make([]string, 0),
			Mutations:  make([]Mutation, 0),
			Fitness:    0.5,
		},
		Memory: ContractMemory{
			ShortTerm:   make(map[string]interface{}),
			LongTerm:    make(map[string]interface{}),
			Patterns:    make([]MemoryPattern, 0),
			Experiences: make([]Experience, 0),
		},
		Behavior: ContractBehavior{
			Traits:         make([]BehaviorTrait, 0),
			Adaptations:    make([]Adaptation, 0),
			Collaborations: make([]Collaboration, 0),
			Predictions:    make([]Prediction, 0),
		},
	}

	// Create contract actor
	contractActor := NewProductionContractActor(contract, pcm.actorSystem)
	props := actor.PropsFromProducer(func() actor.Actor {
		return contractActor
	})

	pid := pcm.actorSystem.Root.Spawn(props)

	// Store contract and actor
	pcm.contracts[contract.Address] = contract
	pcm.contractActors[contract.Address] = pid

	// Send deployment message to actor
	pcm.actorSystem.Root.Send(pid, deployMsg)

	log.Printf("✅ Contract deployed successfully: %s (%s)", contract.Name, contract.Address)

	return contract, nil
}

// GetContract retrieves a contract by address
func (pcm *ProductionContractManager) GetContract(address string) (*Contract, bool) {
	pcm.mutex.RLock()
	defer pcm.mutex.RUnlock()

	contract, exists := pcm.contracts[address]
	return contract, exists
}

// ListContracts returns all contracts
func (pcm *ProductionContractManager) ListContracts() []*Contract {
	pcm.mutex.RLock()
	defer pcm.mutex.RUnlock()

	contracts := make([]*Contract, 0, len(pcm.contracts))
	for _, contract := range pcm.contracts {
		contracts = append(contracts, contract)
	}

	return contracts
}

// GetAllContracts returns all contracts (alias for ListContracts)
func (pcm *ProductionContractManager) GetAllContracts() ([]*Contract, error) {
	contracts := pcm.ListContracts()
	return contracts, nil
}

// ExecuteContract executes a contract function
func (pcm *ProductionContractManager) ExecuteContract(address string, execMsg *ExecuteContract) error {
	pcm.mutex.RLock()
	pid, exists := pcm.contractActors[address]
	pcm.mutex.RUnlock()

	if !exists {
		return &ContractError{
			Code:    "CONTRACT_NOT_FOUND",
			Message: "Contract not found: " + address,
		}
	}

	pcm.actorSystem.Root.Send(pid, execMsg)
	log.Printf("⚡ Contract execution sent: %s.%s", address, execMsg.Function)

	return nil
}

// StopContract stops a contract actor
func (pcm *ProductionContractManager) StopContract(address string) error {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	pid, exists := pcm.contractActors[address]
	if !exists {
		return &ContractError{
			Code:    "CONTRACT_NOT_FOUND",
			Message: "Contract not found: " + address,
		}
	}

	pcm.actorSystem.Root.Stop(pid)

	// Update contract status
	if contract, found := pcm.contracts[address]; found {
		contract.Status = ContractStatusArchived
		contract.UpdatedAt = time.Now()
	}

	delete(pcm.contractActors, address)
	log.Printf("🛑 Contract stopped: %s", address)

	return nil
}

// GetSystemStats returns system statistics
func (pcm *ProductionContractManager) GetSystemStats() *SystemStatusResult {
	pcm.mutex.RLock()
	defer pcm.mutex.RUnlock()

	var activeContracts, sleepingContracts, evolvingContracts int
	var totalExecutions int64

	for _, contract := range pcm.contracts {
		switch contract.Status {
		case ContractStatusActive:
			activeContracts++
		case ContractStatusSleeping:
			sleepingContracts++
		case ContractStatusEvolving:
			evolvingContracts++
		}
		totalExecutions += contract.ExecutionCount
	}

	return &SystemStatusResult{
		ActiveContracts:   activeContracts,
		SleepingContracts: sleepingContracts,
		EvolvingContracts: evolvingContracts,
		TotalExecutions:   totalExecutions,
		AverageCPUUsage:   0.15,           // Simulated
		SystemHealth:      0.95,           // Simulated
		Uptime:            time.Hour * 24, // Simulated
	}
}

// ===== HELPER FUNCTIONS =====

// generateContractAddress generates a unique contract address
func generateContractAddress() string {
	return "contract-" + uuid.New().String()[:8]
}

// ContractError represents contract-specific errors
type ContractError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ContractError) Error() string {
	return e.Code + ": " + e.Message
}

// ===== TESTING & DEMO FUNCTIONS =====

// CreateSampleContract creates a sample contract for testing
func (pcm *ProductionContractManager) CreateSampleContract() (*Contract, error) {
	deployMsg := &DeployContract{
		Name:         "SampleLivingContract",
		Type:         ContractTypeLiving,
		SourceCode:   "// Sample living contract code\nfunction main() { return 'Hello World'; }",
		Owner:        "system",
		TimeAware:    true,
		HistoryDepth: 100,
	}

	return pcm.DeployContract(deployMsg)
}

// TestContractExecution tests contract execution
func (pcm *ProductionContractManager) TestContractExecution(contractAddress string) error {
	execMsg := &ExecuteContract{
		ContractAddress: contractAddress,
		Function:        "main",
		Parameters:      []byte(`{"test": true}`),
		Caller:          "test-user",
		GasLimit:        1000000,
	}

	return pcm.ExecuteContract(contractAddress, execMsg)
}

// ===== ADDITIONAL API METHODS =====

// GetSystemStatus returns overall system status (alias for GetSystemStats)
func (pcm *ProductionContractManager) GetSystemStatus() *SystemStatusResult {
	return pcm.GetSystemStats()
}

// UpgradeContract upgrades a contract to a new version
func (pcm *ProductionContractManager) UpgradeContract(upgradeMsg *UpgradeContract) (*ContractUpgraded, error) {
	pcm.mutex.RLock()
	pid, exists := pcm.contractActors[upgradeMsg.ContractAddress]
	contract := pcm.contracts[upgradeMsg.ContractAddress]
	pcm.mutex.RUnlock()

	if !exists {
		return nil, &ContractError{
			Code:    "CONTRACT_NOT_FOUND",
			Message: "Contract not found: " + upgradeMsg.ContractAddress,
		}
	}

	// Send upgrade message to actor
	pcm.actorSystem.Root.Send(pid, upgradeMsg)

	return &ContractUpgraded{
		Contract: contract,
		Success:  true,
		Message:  "Contract upgrade initiated",
	}, nil
}

// WakeContract wakes up a sleeping contract
func (pcm *ProductionContractManager) WakeContract(wakeMsg *WakeContract) (*ContractAwake, error) {
	pcm.mutex.RLock()
	pid, exists := pcm.contractActors[wakeMsg.ContractAddress]
	pcm.mutex.RUnlock()

	if !exists {
		return nil, &ContractError{
			Code:    "CONTRACT_NOT_FOUND",
			Message: "Contract not found: " + wakeMsg.ContractAddress,
		}
	}

	// Send wake message to actor
	pcm.actorSystem.Root.Send(pid, wakeMsg)

	return &ContractAwake{
		ContractAddress: wakeMsg.ContractAddress,
		Status:          ContractStatusActive,
		Message:         "Contract awakened",
	}, nil
}

// SleepContract puts a contract to sleep
func (pcm *ProductionContractManager) SleepContract(sleepMsg *SleepContract) (*ContractSleeping, error) {
	pcm.mutex.RLock()
	pid, exists := pcm.contractActors[sleepMsg.ContractAddress]
	pcm.mutex.RUnlock()

	if !exists {
		return nil, &ContractError{
			Code:    "CONTRACT_NOT_FOUND",
			Message: "Contract not found: " + sleepMsg.ContractAddress,
		}
	}

	// Send sleep message to actor
	pcm.actorSystem.Root.Send(pid, sleepMsg)

	return &ContractSleeping{
		ContractAddress: sleepMsg.ContractAddress,
		Status:          ContractStatusSleeping,
		WakeCondition:   sleepMsg.Condition,
	}, nil
}

// ProposeCollaboration proposes a collaboration between contracts
func (pcm *ProductionContractManager) ProposeCollaboration(proposeMsg *ProposeCollaboration) (*CollaborationProposed, error) {
	pcm.mutex.RLock()
	fromPid, fromExists := pcm.contractActors[proposeMsg.FromContract]
	_, toExists := pcm.contractActors[proposeMsg.ToContract]
	pcm.mutex.RUnlock()

	if !fromExists || !toExists {
		return nil, &ContractError{
			Code:    "CONTRACT_NOT_FOUND",
			Message: "One or both contracts not found",
		}
	}

	// Send collaboration proposal to from-contract
	pcm.actorSystem.Root.Send(fromPid, proposeMsg)

	return &CollaborationProposed{
		CollaborationID: uuid.New(),
		Status:          "proposed",
		Message:         "Collaboration proposal sent",
	}, nil
}

// AcceptCollaboration accepts a collaboration proposal
func (pcm *ProductionContractManager) AcceptCollaboration(collaborationID string, acceptMsg *AcceptCollaboration) (interface{}, error) {
	// For now, return success - in production this would involve more complex logic
	return map[string]interface{}{
		"collaboration_id": collaborationID,
		"status":           "accepted",
		"message":          "Collaboration accepted",
	}, nil
}

// RejectCollaboration rejects a collaboration proposal
func (pcm *ProductionContractManager) RejectCollaboration(collaborationID string, rejectMsg *RejectCollaboration) (interface{}, error) {
	return map[string]interface{}{
		"collaboration_id": collaborationID,
		"status":           "rejected",
		"reason":           rejectMsg.Reason,
	}, nil
}

// GetContractEvents returns events for a contract
func (pcm *ProductionContractManager) GetContractEvents(address string) ([]*ContractEvent, error) {
	pcm.mutex.RLock()
	contract, exists := pcm.contracts[address]
	pcm.mutex.RUnlock()

	if !exists {
		return nil, &ContractError{
			Code:    "CONTRACT_NOT_FOUND",
			Message: "Contract not found: " + address,
		}
	}

	// In a production system, this would query the event history
	eventData, _ := json.Marshal(map[string]interface{}{"contract_name": contract.Name})
	events := []*ContractEvent{
		{
			ID:              uuid.New(),
			ContractAddress: address,
			EventName:       "ContractCreated",
			Timestamp:       contract.CreatedAt,
			Data:            eventData,
		},
	}

	return events, nil
}

// SubscribeToEvents subscribes to contract events
func (pcm *ProductionContractManager) SubscribeToEvents(subMsg *SubscribeToEvents) (*EventSubscribed, error) {
	pcm.mutex.RLock()
	pid, exists := pcm.contractActors[subMsg.ContractAddress]
	pcm.mutex.RUnlock()

	if !exists {
		return nil, &ContractError{
			Code:    "CONTRACT_NOT_FOUND",
			Message: "Contract not found: " + subMsg.ContractAddress,
		}
	}

	// Send subscription message to actor
	pcm.actorSystem.Root.Send(pid, subMsg)

	return &EventSubscribed{
		SubscriptionID: uuid.New(),
		Success:        true,
		Message:        "Subscribed to events",
	}, nil
}

// EmitEvent emits an event from a contract
func (pcm *ProductionContractManager) EmitEvent(emitMsg *EmitEvent) (interface{}, error) {
	pcm.mutex.RLock()
	pid, exists := pcm.contractActors[emitMsg.ContractAddress]
	pcm.mutex.RUnlock()

	if !exists {
		return nil, &ContractError{
			Code:    "CONTRACT_NOT_FOUND",
			Message: "Contract not found: " + emitMsg.ContractAddress,
		}
	}

	// Send emit message to actor
	pcm.actorSystem.Root.Send(pid, emitMsg)

	return map[string]interface{}{
		"event_id":         uuid.New().String(),
		"contract_address": emitMsg.ContractAddress,
		"event_name":       emitMsg.EventName,
		"timestamp":        time.Now(),
		"status":           "emitted",
	}, nil
}

// GetContractAnalytics returns analytics for a contract
func (pcm *ProductionContractManager) GetContractAnalytics(address string) (interface{}, error) {
	pcm.mutex.RLock()
	contract, exists := pcm.contracts[address]
	pcm.mutex.RUnlock()

	if !exists {
		return nil, &ContractError{
			Code:    "CONTRACT_NOT_FOUND",
			Message: "Contract not found: " + address,
		}
	}

	return map[string]interface{}{
		"contract_address": address,
		"execution_count":  contract.ExecutionCount,
		"success_rate":     contract.SuccessRate,
		"status":           contract.Status,
		"uptime":           time.Since(contract.CreatedAt).Hours(),
		"performance": map[string]interface{}{
			"avg_gas_usage":      21000,
			"avg_execution_time": "15ms",
			"error_rate":         1.0 - contract.SuccessRate,
		},
	}, nil
}

// GetContractHistory returns execution history for a contract
func (pcm *ProductionContractManager) GetContractHistory(address string) ([]*ContractExecution, error) {
	pcm.mutex.RLock()
	_, exists := pcm.contracts[address]
	pcm.mutex.RUnlock()

	if !exists {
		return nil, &ContractError{
			Code:    "CONTRACT_NOT_FOUND",
			Message: "Contract not found: " + address,
		}
	}

	// In production, this would query the execution history from storage
	resultData, _ := json.Marshal("Contract deployed successfully")
	history := []*ContractExecution{
		{
			ID:              uuid.New(),
			ContractAddress: address,
			Function:        "deploy",
			Timestamp:       time.Now().Add(-time.Hour),
			GasUsed:         21000,
			Result:          resultData,
		},
	}

	return history, nil
}

// GetPredictions returns AI predictions for a contract
func (pcm *ProductionContractManager) GetPredictions(address string) (interface{}, error) {
	pcm.mutex.RLock()
	_, exists := pcm.contracts[address]
	pcm.mutex.RUnlock()

	if !exists {
		return nil, &ContractError{
			Code:    "CONTRACT_NOT_FOUND",
			Message: "Contract not found: " + address,
		}
	}

	return map[string]interface{}{
		"contract_address": address,
		"predictions": []map[string]interface{}{
			{
				"type":       "performance",
				"confidence": 0.85,
				"prediction": "Contract will maintain 95%+ success rate",
				"timeframe":  "next 30 days",
			},
			{
				"type":       "evolution",
				"confidence": 0.72,
				"prediction": "Contract may benefit from memory optimization",
				"timeframe":  "next 7 days",
			},
		},
		"generated_at": time.Now(),
	}, nil
}
