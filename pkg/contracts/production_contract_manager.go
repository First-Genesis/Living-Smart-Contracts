package contracts

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/google/uuid"
)

// ProductionContractManager manages all smart contracts in the system
type ProductionContractManager struct {
	actorSystem    *actor.ActorSystem
	contracts      map[string]*Contract // stores immutable summary copies
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

// DeployContract deploys a new smart contract with proper safety measures
func (pcm *ProductionContractManager) DeployContract(deployMsg *DeployContract) (*Contract, error) {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	// Create new contract with full initialization
	contract := &Contract{
		ID:           uuid.New(),
		Address:      "contract_" + uuid.New().String()[:8],
		Name:         deployMsg.Name,
		Type:         deployMsg.Type,
		Status:       ContractStatusDeploying,
		Owner:        deployMsg.Owner,
		Version:      "1.0.0",
		SourceCode:   deployMsg.SourceCode,
		TimeAware:    deployMsg.TimeAware,
		HistoryDepth: deployMsg.HistoryDepth,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActive:   time.Now(),
		SuccessRate:  1.0,
		DNA: ContractDNA{
			Generation: 0,
			Fitness:    1.0,
			Mutations:  make([]Mutation, 0),
		},
		Memory: ContractMemory{
			Experiences: make([]Experience, 0),
			Patterns:    make([]MemoryPattern, 0),
		},
		Behavior: ContractBehavior{
			Traits: make([]BehaviorTrait, 0),
		},
	}

	// Create and start contract actor with named spawn
	contractActor := NewProductionContractActor(contract)
	props := actor.PropsFromProducer(func() actor.Actor {
		return contractActor
	})

	pid, err := pcm.actorSystem.Root.SpawnNamed(props, contract.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to spawn contract actor: %w", err)
	}

	// Store an IMMUTABLE SUMMARY COPY (avoid pointer sharing with actor)
	summary := &Contract{
		ID:        contract.ID,
		Address:   contract.Address,
		Name:      contract.Name,
		Type:      contract.Type,
		Status:    contract.Status,
		Owner:     contract.Owner,
		Version:   contract.Version,
		CreatedAt: contract.CreatedAt,
		UpdatedAt: contract.UpdatedAt,
		TimeAware: contract.TimeAware,
		// Exclude mutable fields like SourceCode/DNA/Memory/Behavior
	}

	// Store summary and PID
	pcm.contracts[contract.Address] = summary
	pcm.contractActors[contract.Address] = pid

	// Don't send deployMsg to actor - it already starts with initialized state
	log.Printf("✅ Contract deployed successfully: %s (%s)", contract.Name, contract.Address)
	return contract, nil
}

// getPID safely retrieves a contract actor PID with proper error handling
func (pcm *ProductionContractManager) getPID(address string) (*actor.PID, error) {
	pcm.mutex.RLock()
	pid, exists := pcm.contractActors[address]
	pcm.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("contract not found: %s", address)
	}
	return pid, nil
}

// GetContract retrieves a contract by address (returns immutable summary)
func (pcm *ProductionContractManager) GetContract(address string) (*Contract, bool) {
	pcm.mutex.RLock()
	defer pcm.mutex.RUnlock()

	contract, exists := pcm.contracts[address]
	return contract, exists
}

// ListContracts returns all contracts (immutable summaries)
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

// WakeContract wakes a sleeping contract using RequestFuture for real confirmation
func (pcm *ProductionContractManager) WakeContract(address string) (*ContractAwake, error) {
	pid, err := pcm.getPID(address)
	if err != nil {
		return nil, err
	}

	wakeMsg := &WakeContract{ContractAddress: address}
	fut := pcm.actorSystem.Root.RequestFuture(pid, wakeMsg, 5*time.Second)
	res, err := fut.Result()
	if err != nil {
		return nil, fmt.Errorf("wake contract failed: %w", err)
	}

	msg, ok := res.(*ContractAwake)
	if !ok {
		return nil, fmt.Errorf("unexpected wake response type")
	}

	// Update summary status after successful wake
	pcm.mutex.Lock()
	if s, ok := pcm.contracts[address]; ok {
		s.Status = msg.Status
		s.UpdatedAt = time.Now()
	}
	pcm.mutex.Unlock()

	return msg, nil
}

// SleepContract puts a contract to sleep using RequestFuture for real confirmation
func (pcm *ProductionContractManager) SleepContract(address string) (*ContractSleeping, error) {
	pid, err := pcm.getPID(address)
	if err != nil {
		return nil, err
	}

	sleepMsg := &SleepContract{ContractAddress: address}
	fut := pcm.actorSystem.Root.RequestFuture(pid, sleepMsg, 5*time.Second)
	res, err := fut.Result()
	if err != nil {
		return nil, fmt.Errorf("sleep contract failed: %w", err)
	}

	msg, ok := res.(*ContractSleeping)
	if !ok {
		return nil, fmt.Errorf("unexpected sleep response type")
	}

	// Update summary status after successful sleep
	pcm.mutex.Lock()
	if s, ok := pcm.contracts[address]; ok {
		s.Status = msg.Status
		s.UpdatedAt = time.Now()
	}
	pcm.mutex.Unlock()

	return msg, nil
}

// ExecuteContract executes a contract function using RequestFuture for confirmation
func (pcm *ProductionContractManager) ExecuteContract(address string, execMsg *ExecuteContract) error {
	pid, err := pcm.getPID(address)
	if err != nil {
		return err
	}

	fut := pcm.actorSystem.Root.RequestFuture(pid, execMsg, 5*time.Second)
	res, err := fut.Result()
	if err != nil {
		return fmt.Errorf("execute contract failed: %w", err)
	}

	execRes, ok := res.(*ContractExecuted)
	if !ok {
		return fmt.Errorf("unexpected execute response type: %T", res)
	}
	if !execRes.Success {
		return fmt.Errorf("execution failed: %s", execRes.Error)
	}

	// Update summary LastActive after successful execution
	pcm.mutex.Lock()
	if s, ok := pcm.contracts[address]; ok {
		s.LastActive = time.Now()
		s.UpdatedAt = time.Now()
		// Promote from Deploying to Active on first successful execution
		if s.Status == ContractStatusDeploying {
			s.Status = ContractStatusActive
		}
	}
	pcm.mutex.Unlock()

	log.Printf("⚡ Contract execution completed: %s.%s", address, execMsg.Function)
	return nil
}

// StopContract stops a contract actor safely and archives the summary
func (pcm *ProductionContractManager) StopContract(address string) error {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	pid, exists := pcm.contractActors[address]
	if !exists {
		return fmt.Errorf("contract not found: %s", address)
	}

	// Stop the actor
	pcm.actorSystem.Root.Stop(pid)

	// Keep immutable summary, mark archived
	if summary, ok := pcm.contracts[address]; ok {
		summary.Status = ContractStatusArchived
		summary.UpdatedAt = time.Now()
	}

	// Remove only the actor PID, keep the archived summary
	delete(pcm.contractActors, address)

	log.Printf("🛑 Contract stopped and archived: %s", address)
	return nil
}

// GetSystemStats returns system statistics from immutable summaries
func (pcm *ProductionContractManager) GetSystemStats() *SystemStatusResult {
	pcm.mutex.RLock()
	defer pcm.mutex.RUnlock()

	totalContracts := len(pcm.contracts)
	activeContracts := 0
	sleepingContracts := 0

	for _, contract := range pcm.contracts {
		switch contract.Status {
		case ContractStatusActive:
			activeContracts++
		case ContractStatusSleeping:
			sleepingContracts++
		}
	}

	systemHealth := 1.0
	if totalContracts > 0 {
		systemHealth = float64(activeContracts) / float64(totalContracts)
	}

	return &SystemStatusResult{
		ActiveContracts:   activeContracts,
		SleepingContracts: sleepingContracts,
		SystemHealth:      systemHealth,
		Uptime:            24 * time.Hour, // placeholder
	}
}
