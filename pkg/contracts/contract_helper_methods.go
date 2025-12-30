package contracts

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// ===== HELPER FUNCTIONS =====

// generateCollaborationID generates a unique collaboration ID
func generateCollaborationID() string {
	return uuid.New().String()
}

// generateExecutionID generates a unique execution ID
func generateExecutionID() string {
	return uuid.New().String()
}

// ===== CONTRACT ACTOR HELPER METHODS =====

// ContractActor is a type alias for ProductionContractActor
type ContractActor = ProductionContractActor

// shouldAutoAcceptCollaboration determines if a collaboration should be auto-accepted
func (ca *ContractActor) shouldAutoAcceptCollaboration(collab *ActiveCollaboration) bool {
	// Auto-accept if the contract has collaborative traits
	for _, trait := range ca.contract.Behavior.Traits {
		if trait.Type == TraitTypeCollaborative && trait.Strength > 0.7 {
			return true
		}
	}

	// Auto-accept data sharing collaborations with trusted partners
	if collab.Type == CollaborationTypeDataSharing {
		// Check if partner has good collaboration history
		for _, pastCollab := range ca.contract.Behavior.Collaborations {
			if pastCollab.PartnerAddress == collab.PartnerAddress && pastCollab.Success {
				return true
			}
		}
	}

	return false
}

// acceptCollaboration processes collaboration acceptance
func (ca *ContractActor) acceptCollaboration(context interface{}, collab *ActiveCollaboration) {
	collab.Status = CollaborationStatusAccepted
	now := time.Now()
	collab.ActivatedAt = &now

	// Add to contract's collaboration history
	ca.contract.Behavior.Collaborations = append(ca.contract.Behavior.Collaborations, Collaboration{
		ID:             collab.ID.String(),
		PartnerAddress: collab.PartnerAddress,
		Type:           collab.Type,
		Success:        true,
		Benefit:        0.5, // Initial benefit estimate
		Timestamp:      now,
	})

	ca.isDirty = true
}

// shouldReactToEvent determines if contract should react to an event
func (ca *ContractActor) shouldReactToEvent(eventType string, eventData json.RawMessage) bool {
	// Check if contract has patterns for this event type
	for _, pattern := range ca.contract.Memory.Patterns {
		if pattern.Type == PatternTypeInteraction {
			// Simple pattern matching - in production, this would be more sophisticated
			return pattern.Confidence > 0.6
		}
	}

	// Default reactive behavior based on contract traits
	for _, trait := range ca.contract.Behavior.Traits {
		if trait.Type == TraitTypeAdaptable && trait.Strength > 0.5 {
			return true
		}
	}

	return false
}

// executeContractLogic executes contract business logic
func (ca *ContractActor) executeContractLogic(context interface{}, execution *ContractExecution) {
	startTime := time.Now()

	// Simulate contract execution
	execution.Status = ExecutionStatusExecuting

	// Add to execution history
	ca.executionHistory = append(ca.executionHistory, execution)

	// Update performance metrics
	ca.executionCount++

	// Simulate successful execution
	execution.Status = ExecutionStatusCompleted
	execution.ExecutionTime = time.Since(startTime)
	execution.GasUsed = int64(rand.Intn(1000) + 500) // Simulate gas usage

	ca.executionCount++
	ca.contract.LastActive = time.Now()
	ca.isDirty = true

	// Learn from this execution
	ca.learnFromExecution(execution)
}

func (ca *ContractActor) notifyEventSubscribers(eventType string, payload json.RawMessage) {
	ca.mutex.RLock()
	defer ca.mutex.RUnlock()

	for _, sub := range ca.eventSubscriptions {
		// If EventTypes not specified, treat as "all"
		if len(sub.EventTypes) == 0 {
			fmt.Printf("Notifying subscriber %v about event: %s\n", sub.SubscriberPID, eventType)
			continue
		}

		for _, t := range sub.EventTypes {
			if t == "*" || t == eventType {
				fmt.Printf("Notifying subscriber %v about event: %s\n", sub.SubscriberPID, eventType)
				break
			}
		}
	}
}

// processExperience processes learning experience data
func (ca *ContractActor) processExperience(experience Experience) {
	ca.contract.Memory.Experiences = append(ca.contract.Memory.Experiences, experience)

	// Trigger pattern recognition if we have enough experiences
	if len(ca.contract.Memory.Experiences) > 10 {
		ca.recognizePatterns()
	}
}

// calculatePredictionConfidence calculates confidence for predictions
func (ca *ContractActor) calculatePredictionConfidence() float64 {
	if ca.executionCount == 0 {
		return 0.5 // Default confidence for new contracts
	}

	successRate := float64(ca.executionCount) / float64(ca.executionCount)
	adaptationScore := ca.contract.AdaptationScore

	// Combine success rate and adaptation score
	confidence := (successRate + adaptationScore) / 2.0

	// Cap confidence between 0.1 and 0.95
	if confidence < 0.1 {
		confidence = 0.1
	}
	if confidence > 0.95 {
		confidence = 0.95
	}

	return confidence
}

// predictExecutionFrequency predicts future execution frequency
func (ca *ContractActor) predictExecutionFrequency(timeHorizon time.Duration) float64 {
	ca.mutex.RLock()
	defer ca.mutex.RUnlock()

	hours := timeHorizon.Hours()
	if hours <= 0 {
		return 0
	}

	if len(ca.executionHistory) < 2 {
		return 1.0
	}

	recent := 0
	cutoff := time.Now().Add(-timeHorizon)
	for _, exec := range ca.executionHistory {
		if exec.Timestamp.After(cutoff) {
			recent++
		}
	}

	return float64(recent) / hours
}

// predictGasUsage predicts future gas usage
func (ca *ContractActor) predictGasUsage(timeHorizon time.Duration) float64 {
	if ca.contract.ExecutionCount == 0 {
		return 1000.0 // Default prediction
	}

	return float64(ca.contract.AverageGasUsed)
}

// predictSuccessRate predicts future success rate
func (ca *ContractActor) predictSuccessRate(timeHorizon time.Duration) float64 {
	return ca.contract.SuccessRate
}

// predictAdaptationScore predicts future adaptation score
func (ca *ContractActor) predictAdaptationScore(timeHorizon time.Duration) float64 {
	return ca.contract.AdaptationScore
}

// calculateExecutionEfficiency calculates execution efficiency metric
func (ca *ContractActor) calculateExecutionEfficiency() float64 {
	ca.mutex.RLock()
	defer ca.mutex.RUnlock()

	if ca.executionCount == 0 || len(ca.executionHistory) == 0 {
		return 0.0
	}

	total := time.Duration(0)
	for _, exec := range ca.executionHistory {
		total += exec.ExecutionTime
	}

	avg := total / time.Duration(len(ca.executionHistory))
	ms := avg.Milliseconds()
	if ms <= 0 {
		return 1.0
	}

	return 1.0 / (1.0 + float64(ms)/1000.0)
}

// calculateErrorRate calculates current error rate
func (ca *ContractActor) calculateErrorRate() float64 {
	if ca.executionCount == 0 {
		return 0.0
	}

	errorCount := ca.executionCount - ca.executionCount
	return float64(errorCount) / float64(ca.executionCount)
}

// calculateAdaptationTrend calculates adaptation trend over time
func (ca *ContractActor) calculateAdaptationTrend() float64 {
	if len(ca.contract.Behavior.Adaptations) < 2 {
		return 0.0
	}

	// Calculate trend from recent adaptations
	recentAdaptations := ca.contract.Behavior.Adaptations
	if len(recentAdaptations) > 10 {
		recentAdaptations = recentAdaptations[len(recentAdaptations)-10:]
	}

	totalImpact := 0.0
	for _, adaptation := range recentAdaptations {
		totalImpact += adaptation.Impact
	}

	return totalImpact / float64(len(recentAdaptations))
}

// identifyPeakUsageTimes identifies patterns in usage timing
func (ca *ContractActor) identifyPeakUsageTimes() interface{} {
	hourCounts := make(map[int]int)

	for _, exec := range ca.executionHistory {
		hour := exec.Timestamp.Hour()
		hourCounts[hour]++
	}

	// Find peak hours
	peakHours := make([]int, 0)
	maxCount := 0

	for hour, count := range hourCounts {
		if count > maxCount {
			maxCount = count
			peakHours = []int{hour}
		} else if count == maxCount {
			peakHours = append(peakHours, hour)
		}
	}

	return peakHours
}

// identifyErrorPatterns identifies common error patterns
func (ca *ContractActor) identifyErrorPatterns() interface{} {
	ca.mutex.RLock()
	defer ca.mutex.RUnlock()

	errorPatterns := make(map[string]int)

	for _, exec := range ca.executionHistory {
		if exec.Status == ExecutionStatusFailed && exec.Error != "" {
			errorPatterns[exec.Error]++
		}
	}

	// Return most common error patterns
	commonErrors := make([]string, 0)
	for errorMsg, count := range errorPatterns {
		if count > 1 { // Only include patterns that occur multiple times
			commonErrors = append(commonErrors, errorMsg)
		}
	}

	return commonErrors
}

// calculateLearningRate calculates current learning rate
func (ca *ContractActor) calculateLearningRate() float64 {
	if !ca.learningEnabled {
		return 0.0
	}

	// Base learning rate
	baseRate := 0.1

	// Adjust based on experience count
	experienceCount := len(ca.contract.Memory.Experiences)
	if experienceCount > 100 {
		// Reduce learning rate as contract gains more experience
		return baseRate * (100.0 / float64(experienceCount))
	}

	return baseRate
}

// learnFromExecution learns from contract execution
func (ca *ContractActor) learnFromExecution(execution *ContractExecution) {
	if !ca.learningEnabled {
		return
	}

	// Create experience from execution
	experience := Experience{
		ID:        uuid.New().String(),
		Context:   execution.Parameters,
		Action:    []byte(fmt.Sprintf(`{"function": "%s"}`, execution.Function)),
		Result:    execution.Result,
		Success:   execution.Status == ExecutionStatusCompleted,
		Timestamp: execution.Timestamp,
	}

	ca.processExperience(experience)
}

// recognizePatterns performs pattern recognition on contract memory
func (ca *ContractActor) recognizePatterns() {
	ca.mutex.Lock()
	defer ca.mutex.Unlock()

	// Simple pattern recognition - in production this would be more sophisticated
	successPatterns := 0
	totalPatterns := len(ca.contract.Memory.Experiences)

	for _, exp := range ca.contract.Memory.Experiences {
		if exp.Success {
			successPatterns++
		}
	}

	if totalPatterns > 0 {
		confidence := float64(successPatterns) / float64(totalPatterns)

		pattern := MemoryPattern{
			ID:         uuid.New().String(),
			Type:       PatternTypeExecution,
			Trigger:    []byte(`{"type": "execution"}`),
			Response:   []byte(`{"optimize": true}`),
			Confidence: confidence,
			Usage:      0,
			LastUsed:   time.Now(),
		}

		ca.contract.Memory.Patterns = append(ca.contract.Memory.Patterns, pattern)
	}
}

// canUpgrade checks if contract can be upgraded
func (ca *ContractActor) canUpgrade(msg *UpgradeContract) bool {
	// Check if contract is not currently evolving
	if ca.contract.Status == ContractStatusEvolving {
		return false
	}

	// Check version compatibility
	if msg.NewVersion <= ca.contract.Version {
		return false
	}

	// Additional upgrade validation would go here
	return true
}

// createContractBackup creates a backup before upgrade
func (ca *ContractActor) createContractBackup() *ContractBackup {
	return &ContractBackup{
		Version:   ca.contract.Version,
		Code:      ca.contract.SourceCode,
		State:     ca.contract.StateRoot,
		Timestamp: time.Now(),
	}
}

// validateContractCode validates new contract code
func (ca *ContractActor) validateContractCode(code string) bool {
	// Basic validation - in production this would include:
	// - Syntax checking
	// - Security analysis
	// - Compatibility verification
	return len(code) > 0 && code != "invalid"
}

// restoreFromBackup restores contract from backup
func (ca *ContractActor) restoreFromBackup(backup *ContractBackup) {
	ca.contract.Version = backup.Version
	ca.contract.SourceCode = backup.Code
	ca.contract.StateRoot = backup.State
	ca.contract.UpdatedAt = time.Now()
}
