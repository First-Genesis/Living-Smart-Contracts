package contracts

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// LearningEngine implements machine learning for smart contracts
// This enables contracts to learn from their execution history and adapt behavior
type LearningEngine struct {
	contract        *Contract
	patternAnalyzer *PatternAnalyzer
	predictor       *BehaviorPredictor
	optimizer       *PerformanceOptimizer
	mutex           sync.RWMutex

	// Learning configuration
	learningRate        float64
	confidenceThreshold float64
	minSampleSize       int
	maxMemorySize       int

	// Learning state
	totalExperiences int
	learningCycles   int
	lastLearningTime time.Time
}

// AdaptationEngine manages behavioral adaptation based on learning
type AdaptationEngine struct {
	contract          *Contract
	adaptationRules   []*AdaptationRule
	traitEvolution    *TraitEvolution
	environmentSensor *EnvironmentSensor
	mutex             sync.RWMutex

	// Adaptation configuration
	adaptationThreshold   float64
	maxAdaptationsPerHour int
	stabilityRequirement  float64

	// Adaptation state
	totalAdaptations   int
	recentAdaptations  []time.Time
	currentEnvironment *Environment
}

// PatternAnalyzer discovers patterns in contract execution history
type PatternAnalyzer struct {
	patterns            map[PatternType][]*LearnedPattern
	frequencyAnalysis   *FrequencyAnalysis
	sequenceAnalysis    *SequenceAnalysis
	correlationAnalysis *CorrelationAnalysis
}

// LearnedPattern represents a discovered behavioral pattern
type LearnedPattern struct {
	ID             string                 `json:"id"`
	Type           PatternType            `json:"type"`
	Trigger        map[string]interface{} `json:"trigger"`
	Response       map[string]interface{} `json:"response"`
	Confidence     float64                `json:"confidence"`
	SampleCount    int                    `json:"sample_count"`
	SuccessRate    float64                `json:"success_rate"`
	DiscoveredAt   time.Time              `json:"discovered_at"`
	LastReinforced time.Time              `json:"last_reinforced"`
	Stability      float64                `json:"stability"`
}

// BehaviorPredictor predicts future contract behavior
type BehaviorPredictor struct {
	models          map[PredictionType]*PredictionModel
	historicalData  []DataPoint
	predictionCache map[string]*CachedPrediction
}

// PredictionModel represents a trained prediction model
type PredictionModel struct {
	Type        PredictionType      `json:"type"`
	Algorithm   PredictionAlgorithm `json:"algorithm"`
	Parameters  map[string]float64  `json:"parameters"`
	Accuracy    float64             `json:"accuracy"`
	TrainedAt   time.Time           `json:"trained_at"`
	SampleCount int                 `json:"sample_count"`
}

type PredictionAlgorithm string

const (
	AlgorithmLinearRegression   PredictionAlgorithm = "linear_regression"
	AlgorithmDecisionTree       PredictionAlgorithm = "decision_tree"
	AlgorithmNeuralNetwork      PredictionAlgorithm = "neural_network"
	AlgorithmTimeSeriesAnalysis PredictionAlgorithm = "time_series"
	AlgorithmPatternMatching    PredictionAlgorithm = "pattern_matching"
)

// DataPoint represents a single data point for learning
type DataPoint struct {
	Timestamp     time.Time              `json:"timestamp"`
	Input         map[string]interface{} `json:"input"`
	Output        map[string]interface{} `json:"output"`
	Context       map[string]interface{} `json:"context"`
	Success       bool                   `json:"success"`
	GasUsed       int64                  `json:"gas_used"`
	ExecutionTime time.Duration          `json:"execution_time"`
}

// CachedPrediction stores a prediction result
type CachedPrediction struct {
	Prediction interface{} `json:"prediction"`
	Confidence float64     `json:"confidence"`
	CreatedAt  time.Time   `json:"created_at"`
	ValidUntil time.Time   `json:"valid_until"`
}

// AdaptationRule defines rules for behavioral adaptation
type AdaptationRule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Condition   map[string]interface{} `json:"condition"`
	Action      AdaptationAction       `json:"action"`
	Priority    int                    `json:"priority"`
	Enabled     bool                   `json:"enabled"`
	Usage       int                    `json:"usage"`
	SuccessRate float64                `json:"success_rate"`
}

// AdaptationAction defines what adaptation to perform
type AdaptationAction struct {
	Type       AdaptationType         `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
	Magnitude  float64                `json:"magnitude"`
	Reversible bool                   `json:"reversible"`
}

type AdaptationType string

const (
	AdaptationTypeOptimizeGas          AdaptationType = "optimize_gas"
	AdaptationTypeChangeAlgorithm      AdaptationType = "change_algorithm"
	AdaptationTypeAdjustParameters     AdaptationType = "adjust_parameters"
	AdaptationTypeModifyBehavior       AdaptationType = "modify_behavior"
	AdaptationTypeImproveErrorHandling AdaptationType = "improve_error_handling"
	AdaptationTypeEnhanceSecurity      AdaptationType = "enhance_security"
)

// TraitEvolution manages the evolution of contract traits
type TraitEvolution struct {
	traits            map[TraitType]*EvolvingTrait
	evolutionHistory  []TraitChange
	selectionPressure map[TraitType]float64
}

// EvolvingTrait represents a trait that can evolve over time
type EvolvingTrait struct {
	Type          TraitType `json:"type"`
	Value         float64   `json:"value"`
	Trend         float64   `json:"trend"`
	Stability     float64   `json:"stability"`
	Fitness       float64   `json:"fitness"`
	LastChange    time.Time `json:"last_change"`
	ChangeHistory []float64 `json:"change_history"`
}

// TraitChange records a trait evolution event
type TraitChange struct {
	TraitType TraitType `json:"trait_type"`
	OldValue  float64   `json:"old_value"`
	NewValue  float64   `json:"new_value"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
}

// Environment represents the contract's operating environment
type Environment struct {
	NetworkLoad          float64            `json:"network_load"`
	GasPrice             int64              `json:"gas_price"`
	BlockTime            time.Duration      `json:"block_time"`
	PeerContracts        []string           `json:"peer_contracts"`
	ResourceAvailability map[string]float64 `json:"resource_availability"`
	MarketConditions     map[string]float64 `json:"market_conditions"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

// EnvironmentSensor monitors the contract's operating environment
type EnvironmentSensor struct {
	currentEnvironment *Environment
	environmentHistory []*Environment
	changeDetector     *ChangeDetector
	adaptationTriggers map[string]float64
}

// ChangeDetector identifies significant environmental changes
type ChangeDetector struct {
	thresholds     map[string]float64
	changeHistory  []EnvironmentChange
	alertCallbacks []func(EnvironmentChange)
}

// EnvironmentChange represents a detected environmental change
type EnvironmentChange struct {
	Metric       string    `json:"metric"`
	OldValue     float64   `json:"old_value"`
	NewValue     float64   `json:"new_value"`
	ChangeRate   float64   `json:"change_rate"`
	Significance float64   `json:"significance"`
	DetectedAt   time.Time `json:"detected_at"`
}

// NewLearningEngine creates a new learning engine for a contract
func NewLearningEngine(contract *Contract) *LearningEngine {
	return &LearningEngine{
		contract:            contract,
		patternAnalyzer:     NewPatternAnalyzer(),
		predictor:           NewBehaviorPredictor(),
		optimizer:           NewPerformanceOptimizer(),
		learningRate:        0.1,  // 10% learning rate
		confidenceThreshold: 0.7,  // 70% confidence threshold
		minSampleSize:       10,   // Minimum 10 samples for learning
		maxMemorySize:       1000, // Keep last 1000 experiences
		lastLearningTime:    time.Now(),
	}
}

// NewAdaptationEngine creates a new adaptation engine for a contract
func NewAdaptationEngine(contract *Contract) *AdaptationEngine {
	return &AdaptationEngine{
		contract:              contract,
		adaptationRules:       make([]*AdaptationRule, 0),
		traitEvolution:        NewTraitEvolution(),
		environmentSensor:     NewEnvironmentSensor(),
		adaptationThreshold:   0.8, // 80% confidence required for adaptation
		maxAdaptationsPerHour: 5,   // Max 5 adaptations per hour
		stabilityRequirement:  0.6, // 60% stability required
		recentAdaptations:     make([]time.Time, 0),
	}
}

// LearnFromExperiences processes a batch of experiences to discover patterns
func (le *LearningEngine) LearnFromExperiences(experiences []*Experience) error {
	le.mutex.Lock()
	defer le.mutex.Unlock()

	if len(experiences) < le.minSampleSize {
		return fmt.Errorf("insufficient sample size: %d < %d", len(experiences), le.minSampleSize)
	}

	log.Printf("🧠 Learning from %d experiences for contract %s", len(experiences), le.contract.Address)

	// Convert experiences to data points
	dataPoints := le.experiencesToDataPoints(experiences)

	// Discover patterns
	newPatterns, err := le.patternAnalyzer.DiscoverPatterns(dataPoints)
	if err != nil {
		return fmt.Errorf("pattern discovery failed: %w", err)
	}

	// Update prediction models
	err = le.predictor.UpdateModels(dataPoints)
	if err != nil {
		log.Printf("Warning: prediction model update failed: %v", err)
	}

	// Optimize performance based on learned patterns
	optimizations := le.optimizer.GenerateOptimizations(newPatterns, dataPoints)

	// Apply optimizations to contract behavior
	le.applyOptimizations(optimizations)

	// Update contract memory with new experiences
	for _, exp := range experiences {
		le.contract.Memory.Experiences = append(le.contract.Memory.Experiences, *exp)
	}

	// Trim experiences to prevent unbounded memory growth
	le.contract.Memory.Experiences = trimExperiencesValue(le.contract.Memory.Experiences, le.maxMemorySize)

	// Update learning statistics
	le.totalExperiences += len(experiences)
	le.learningCycles++
	le.lastLearningTime = time.Now()

	log.Printf("🎯 Learning completed: discovered %d patterns, %d optimizations", len(newPatterns), len(optimizations))

	return nil
}

// AdaptBehavior applies behavioral adaptations based on current environment
func (ae *AdaptationEngine) AdaptBehavior(environment *Environment) ([]*Adaptation, error) {
	ae.mutex.Lock()
	defer ae.mutex.Unlock()

	// Check adaptation rate limits
	if !ae.canAdapt() {
		return nil, fmt.Errorf("adaptation rate limit exceeded")
	}

	log.Printf("🔄 Adapting behavior for contract %s", ae.contract.Address)

	// Update environment state
	ae.environmentSensor.UpdateEnvironment(environment)

	// Detect environmental changes
	changes := ae.environmentSensor.DetectChanges()

	// Evaluate adaptation rules
	applicableRules := ae.evaluateAdaptationRules(environment, changes)

	// Apply adaptations
	adaptations := make([]*Adaptation, 0)
	for _, rule := range applicableRules {
		adaptation, err := ae.applyAdaptationRule(rule, environment)
		if err != nil {
			log.Printf("Warning: adaptation rule %s failed: %v", rule.ID, err)
			continue
		}

		adaptations = append(adaptations, adaptation)
		ae.recordAdaptation()
	}

	// Evolve traits based on adaptations
	ae.traitEvolution.EvolveTraits(adaptations, environment)

	log.Printf("🔄 Behavior adaptation completed: applied %d adaptations", len(adaptations))

	return adaptations, nil
}

// PredictBehavior predicts future contract behavior
func (le *LearningEngine) PredictBehavior(predictionType PredictionType, ctx map[string]interface{}, timeHorizon time.Duration) (*Prediction, error) {
	cacheKey, err := makePredictionCacheKey(predictionType, ctx, timeHorizon)
	if err != nil {
		return nil, err
	}

	le.mutex.RLock()
	if cached, exists := le.predictor.predictionCache[cacheKey]; exists && time.Now().Before(cached.ValidUntil) {
		le.mutex.RUnlock()
		return cachedToPrediction(predictionType, cached), nil
	}
	model, exists := le.predictor.models[predictionType]
	le.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no prediction model for type: %s", predictionType)
	}

	predictionResult, confidence, err := le.generatePrediction(model, ctx, timeHorizon)
	if err != nil {
		return nil, fmt.Errorf("prediction generation failed: %w", err)
	}

	p := &Prediction{
		ID:         fmt.Sprintf("pred_%d", time.Now().UnixNano()),
		Type:       predictionType,
		Prediction: predictionResult,
		Confidence: confidence,
		MadeAt:     time.Now(),
		ValidUntil: time.Now().Add(timeHorizon),
	}

	le.mutex.Lock()
	le.predictor.predictionCache[cacheKey] = &CachedPrediction{
		Prediction: predictionResult,
		Confidence: confidence,
		CreatedAt:  p.MadeAt,
		ValidUntil: p.ValidUntil,
	}
	le.mutex.Unlock()

	return p, nil
}

// Helper functions and implementations

func makePredictionCacheKey(predType PredictionType, ctx map[string]interface{}, horizon time.Duration) (string, error) {
	blob := map[string]interface{}{
		"type":    predType,
		"horizon": horizon.String(),
		"context": ctx,
	}
	b, err := json.Marshal(blob) // stable ordering for maps in encoding/json is deterministic by key sorting
	if err != nil {
		return "", err
	}
	// hash to keep key small
	sum := sha256.Sum256(b)
	return fmt.Sprintf("%x", sum[:]), nil
}

func cachedToPrediction(predictionType PredictionType, cached *CachedPrediction) *Prediction {
	raw, _ := json.Marshal(cached.Prediction) // ensure valid JSON
	return &Prediction{
		ID:         fmt.Sprintf("cached_%d", time.Now().UnixNano()),
		Type:       predictionType,
		Prediction: raw,
		Confidence: cached.Confidence,
		MadeAt:     cached.CreatedAt,
		ValidUntil: cached.ValidUntil,
	}
}

func trimExperiences(exps []*Experience, max int) []*Experience {
	if max <= 0 || len(exps) <= max {
		return exps
	}
	return exps[len(exps)-max:]
}

func trimExperiencesValue(exps []Experience, max int) []Experience {
	if max <= 0 || len(exps) <= max {
		return exps
	}
	return exps[len(exps)-max:]
}

func (le *LearningEngine) experiencesToDataPoints(experiences []*Experience) []DataPoint {
	out := make([]DataPoint, 0, len(experiences))
	badContext, badAction, badResult := 0, 0, 0

	for _, exp := range experiences {
		var input map[string]interface{}
		var actionCtx map[string]interface{}
		var result map[string]interface{}

		if err := json.Unmarshal(exp.Context, &input); err != nil {
			badContext++
			continue
		}
		if err := json.Unmarshal(exp.Action, &actionCtx); err != nil {
			badAction++
		}
		if err := json.Unmarshal(exp.Result, &result); err != nil {
			badResult++
		}

		dp := DataPoint{
			Timestamp: exp.Timestamp,
			Input:     input,
			Output:    result,
			Context:   actionCtx,
			Success:   exp.Success,
		}

		// Check multiple sources for gas usage and execution time
		if v, ok := result["gas_used"].(float64); ok {
			dp.GasUsed = int64(v)
		} else if v, ok := actionCtx["gas_used"].(float64); ok {
			dp.GasUsed = int64(v)
		} else if v, ok := input["gas_used"].(float64); ok {
			dp.GasUsed = int64(v)
		}

		if v, ok := result["execution_time_ms"].(float64); ok {
			dp.ExecutionTime = time.Duration(v) * time.Millisecond
		} else if v, ok := actionCtx["execution_time_ms"].(float64); ok {
			dp.ExecutionTime = time.Duration(v) * time.Millisecond
		} else if v, ok := input["execution_time_ms"].(float64); ok {
			dp.ExecutionTime = time.Duration(v) * time.Millisecond
		}

		out = append(out, dp)
	}

	if badContext > 0 || badAction > 0 || badResult > 0 {
		log.Printf("Learning decode failures: context=%d action=%d result=%d", badContext, badAction, badResult)
	}

	return out
}

func (le *LearningEngine) applyOptimizations(optimizations []Optimization) {
	// Apply learned optimizations to contract behavior
	for _, opt := range optimizations {
		switch opt.Type {
		case "gas_optimization":
			le.applyGasOptimization(opt)
		case "performance_optimization":
			le.applyPerformanceOptimization(opt)
		case "error_handling":
			le.applyErrorHandlingOptimization(opt)
		}
	}
}

func (ae *AdaptationEngine) canAdapt() bool {
	return len(ae.recentAdaptations) < ae.maxAdaptationsPerHour
}

func (ae *AdaptationEngine) recordAdaptation() {
	now := time.Now()
	oneHourAgo := now.Add(-time.Hour)

	// prune
	pruned := ae.recentAdaptations[:0]
	for _, t := range ae.recentAdaptations {
		if t.After(oneHourAgo) {
			pruned = append(pruned, t)
		}
	}
	ae.recentAdaptations = append(pruned, now)
	ae.totalAdaptations++
}

func (le *LearningEngine) generatePrediction(model *PredictionModel, context map[string]interface{}, timeHorizon time.Duration) (json.RawMessage, float64, error) {
	// Simplified prediction generation
	// In a real implementation, this would use machine learning algorithms

	switch model.Algorithm {
	case AlgorithmLinearRegression:
		return le.linearRegressionPredict(model, context)
	case AlgorithmPatternMatching:
		return le.patternMatchingPredict(model, context)
	default:
		return nil, 0, fmt.Errorf("unsupported algorithm: %s", model.Algorithm)
	}
}

func (le *LearningEngine) linearRegressionPredict(model *PredictionModel, context map[string]interface{}) (json.RawMessage, float64, error) {
	// Simple linear regression prediction
	result := map[string]interface{}{
		"predicted_gas_usage":      50000,
		"predicted_success_rate":   0.95,
		"predicted_execution_time": "100ms",
	}

	resultJSON, _ := json.Marshal(result)
	return resultJSON, model.Accuracy, nil
}

func (le *LearningEngine) patternMatchingPredict(model *PredictionModel, context map[string]interface{}) (json.RawMessage, float64, error) {
	// Pattern-based prediction
	result := map[string]interface{}{
		"matching_pattern":   "high_frequency_trading",
		"predicted_behavior": "aggressive",
		"confidence":         0.8,
	}

	resultJSON, _ := json.Marshal(result)
	return resultJSON, 0.8, nil
}

// Placeholder implementations for supporting types
type Optimization struct {
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
	Impact     float64                `json:"impact"`
}

func (le *LearningEngine) applyGasOptimization(opt Optimization)           {}
func (le *LearningEngine) applyPerformanceOptimization(opt Optimization)   {}
func (le *LearningEngine) applyErrorHandlingOptimization(opt Optimization) {}
func (ae *AdaptationEngine) evaluateAdaptationRules(env *Environment, changes []EnvironmentChange) []*AdaptationRule {
	return nil
}
func (ae *AdaptationEngine) applyAdaptationRule(rule *AdaptationRule, env *Environment) (*Adaptation, error) {
	return nil, nil
}

// Factory functions for supporting components
func NewPatternAnalyzer() *PatternAnalyzer {
	return &PatternAnalyzer{patterns: make(map[PatternType][]*LearnedPattern)}
}
func NewBehaviorPredictor() *BehaviorPredictor {
	return &BehaviorPredictor{models: make(map[PredictionType]*PredictionModel), predictionCache: make(map[string]*CachedPrediction)}
}
func NewPerformanceOptimizer() *PerformanceOptimizer { return &PerformanceOptimizer{} }
func NewTraitEvolution() *TraitEvolution {
	return &TraitEvolution{traits: make(map[TraitType]*EvolvingTrait)}
}
func NewEnvironmentSensor() *EnvironmentSensor { return &EnvironmentSensor{} }

// Stub implementations for pattern analyzer methods
type PerformanceOptimizer struct{}
type FrequencyAnalysis struct{}
type SequenceAnalysis struct{}
type CorrelationAnalysis struct{}

func (pa *PatternAnalyzer) DiscoverPatterns(dataPoints []DataPoint) ([]*LearnedPattern, error) {
	return nil, nil
}
func (bp *BehaviorPredictor) UpdateModels(dataPoints []DataPoint) error { return nil }
func (po *PerformanceOptimizer) GenerateOptimizations(patterns []*LearnedPattern, dataPoints []DataPoint) []Optimization {
	return nil
}
func (te *TraitEvolution) EvolveTraits(adaptations []*Adaptation, environment *Environment) {}
func (es *EnvironmentSensor) UpdateEnvironment(environment *Environment)                    {}
func (es *EnvironmentSensor) DetectChanges() []EnvironmentChange                            { return nil }
