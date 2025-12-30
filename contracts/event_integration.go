package contracts

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/google/uuid"
)

// ContractEventIntegrator manages integration between contracts and the event streaming system
type ContractEventIntegrator struct {
	eventStreamPID    *actor.PID
	registryPID       *actor.PID
	contractManager   *actor.PID
	system           *actor.ActorSystem
	
	// Event mapping and routing
	eventSubscriptions map[uuid.UUID]*ContractEventSubscription
	contractEventSinks map[string]*ContractEventSink  // contract_address -> sink
	
	// Event transformation and enrichment
	eventTransformer  *ContractEventTransformer
	eventEnricher    *ContractEventEnricher
}

// ContractEventSubscription represents a contract's subscription to events
type ContractEventSubscription struct {
	ID              uuid.UUID       `json:"id"`
	ContractAddress string          `json:"contract_address"`
	EventPattern    json.RawMessage `json:"event_pattern"`
	FilterCriteria  json.RawMessage `json:"filter_criteria"`
	Callback        string          `json:"callback"`
	Priority        int             `json:"priority"`
	Active          bool            `json:"active"`
	CreatedAt       time.Time       `json:"created_at"`
	LastTriggered   *time.Time      `json:"last_triggered,omitempty"`
	TriggerCount    int64           `json:"trigger_count"`
}

// ContractEventSink represents a dedicated sink for contract events
type ContractEventSink struct {
	ID              uuid.UUID       `json:"id"`
	ContractAddress string          `json:"contract_address"`
	SinkType        string          `json:"sink_type"`
	Configuration   json.RawMessage `json:"configuration"`
	Filter          json.RawMessage `json:"filter"`
	BufferSize      int             `json:"buffer_size"`
	FlushInterval   time.Duration   `json:"flush_interval"`
	Active          bool            `json:"active"`
	CreatedAt       time.Time       `json:"created_at"`
	EventCount      int64           `json:"event_count"`
	LastActivity    time.Time       `json:"last_activity"`
}

// ContractEventTransformer transforms contract events to different formats
type ContractEventTransformer struct {
	transformationRules map[string]*EventTransformationRule
	formatters         map[EventFormat]*EventFormatter
}

// EventTransformationRule defines how to transform events
type EventTransformationRule struct {
	ID            string                 `json:"id"`
	SourceFormat  EventFormat            `json:"source_format"`
	TargetFormat  EventFormat            `json:"target_format"`
	Mapping       map[string]interface{} `json:"mapping"`
	Conditions    json.RawMessage        `json:"conditions"`
	Priority      int                    `json:"priority"`
	Active        bool                   `json:"active"`
}

// EventFormat represents different event formats
type EventFormat string

const (
	EventFormatContractNative EventFormat = "contract_native"  // Native contract events
	EventFormatCloudEvents    EventFormat = "cloudevents"      // CloudEvents 1.0 format
	EventFormatLedgerEvents   EventFormat = "ledger_events"    // Native ledger events
	EventFormatWebhook        EventFormat = "webhook"          // Webhook format
	EventFormatKafka          EventFormat = "kafka"            // Kafka message format
	EventFormatJSON           EventFormat = "json"             // Generic JSON
)

// EventFormatter formats events for different outputs
type EventFormatter struct {
	Format      EventFormat `json:"format"`
	Template    string      `json:"template"`
	ContentType string      `json:"content_type"`
	Compression string      `json:"compression,omitempty"`
}

// ContractEventEnricher enriches events with additional context
type ContractEventEnricher struct {
	enrichmentSources map[string]*EnrichmentSource
	contextProviders  []ContextProvider
}

// EnrichmentSource provides additional data for event enrichment
type EnrichmentSource struct {
	ID          string                 `json:"id"`
	Type        EnrichmentType         `json:"type"`
	Source      string                 `json:"source"`
	Fields      []string               `json:"fields"`
	CacheTime   time.Duration          `json:"cache_time"`
	LastUpdated time.Time              `json:"last_updated"`
	Active      bool                   `json:"active"`
}

// EnrichmentType defines types of enrichment
type EnrichmentType string

const (
	EnrichmentTypeBlockchain  EnrichmentType = "blockchain"   // Blockchain state data
	EnrichmentTypeContract    EnrichmentType = "contract"     // Contract state data  
	EnrichmentTypeNetwork     EnrichmentType = "network"      // Network information
	EnrichmentTypeMarket      EnrichmentType = "market"       // Market data
	EnrichmentTypeExternal    EnrichmentType = "external"     // External API data
)

// ContextProvider interface for providing enrichment context
type ContextProvider interface {
	GetContext(eventType string, contractAddress string) (map[string]interface{}, error)
	GetCacheKey(eventType string, contractAddress string) string
	GetCacheDuration() time.Duration
}

// ContractEvent types specific to contract operations
type ContractEventType string

const (
	ContractEventDeployed        ContractEventType = "contract.deployed"
	ContractEventExecuted        ContractEventType = "contract.executed"
	ContractEventEvolved         ContractEventType = "contract.evolved"
	ContractEventCollaborated    ContractEventType = "contract.collaborated"
	ContractEventSlept          ContractEventType = "contract.slept"
	ContractEventAwoke          ContractEventType = "contract.awoke"
	ContractEventPredicted      ContractEventType = "contract.predicted"
	ContractEventLearned        ContractEventType = "contract.learned"
	ContractEventAdapted        ContractEventType = "contract.adapted"
	ContractEventStateChanged   ContractEventType = "contract.state_changed"
	ContractEventErrorOccurred  ContractEventType = "contract.error_occurred"
)

// StandardContractEvent represents a standardized contract event
type StandardContractEvent struct {
	// CloudEvents-compatible fields
	SpecVersion     string                 `json:"specversion"`
	Type            string                 `json:"type"`
	Source          string                 `json:"source"`
	ID              string                 `json:"id"`
	Time            time.Time              `json:"time"`
	Subject         string                 `json:"subject,omitempty"`
	DataContentType string                 `json:"datacontenttype,omitempty"`
	
	// Contract-specific fields
	ContractAddress string                 `json:"contract_address"`
	ContractType    ContractType           `json:"contract_type"`
	Function        string                 `json:"function,omitempty"`
	Caller          string                 `json:"caller,omitempty"`
	
	// Event data
	Data            json.RawMessage        `json:"data"`
	
	// Execution context
	BlockHeight     int64                  `json:"block_height,omitempty"`
	GasUsed         int64                  `json:"gas_used,omitempty"`
	ExecutionTime   time.Duration          `json:"execution_time,omitempty"`
	Success         bool                   `json:"success"`
	
	// Learning and adaptation context
	LearningData    json.RawMessage        `json:"learning_data,omitempty"`
	AdaptationData  json.RawMessage        `json:"adaptation_data,omitempty"`
	
	// Enrichment data
	EnrichmentData  map[string]interface{} `json:"enrichment_data,omitempty"`
	
	// Metadata
	Tags            []string               `json:"tags,omitempty"`
	Priority        int                    `json:"priority"`
	Indexed         []string               `json:"indexed,omitempty"`
}

// NewContractEventIntegrator creates a new event integrator
func NewContractEventIntegrator(eventStreamPID, registryPID, contractManagerPID *actor.PID, system *actor.ActorSystem) *ContractEventIntegrator {
	return &ContractEventIntegrator{
		eventStreamPID:     eventStreamPID,
		registryPID:        registryPID,
		contractManager:    contractManagerPID,
		system:            system,
		eventSubscriptions: make(map[uuid.UUID]*ContractEventSubscription),
		contractEventSinks: make(map[string]*ContractEventSink),
		eventTransformer:   NewContractEventTransformer(),
		eventEnricher:     NewContractEventEnricher(),
	}
}

// EmitContractEvent emits a contract event to the event streaming system
func (cei *ContractEventIntegrator) EmitContractEvent(contractAddress string, eventType ContractEventType, data json.RawMessage, context map[string]interface{}) error {
	// Create standardized contract event
	event := &StandardContractEvent{
		SpecVersion:     "1.0",
		Type:            string(eventType),
		Source:          fmt.Sprintf("xenese-dlt://contracts/%s", contractAddress),
		ID:              uuid.New().String(),
		Time:            time.Now(),
		Subject:         contractAddress,
		DataContentType: "application/json",
		ContractAddress: contractAddress,
		Data:            data,
		Success:         true,
		Priority:        1,
		Tags:            []string{"contract", "smart-contract"},
	}
	
	// Add execution context if available
	if context != nil {
		if blockHeight, ok := context["block_height"].(int64); ok {
			event.BlockHeight = blockHeight
		}
		if gasUsed, ok := context["gas_used"].(int64); ok {
			event.GasUsed = gasUsed
		}
		if executionTime, ok := context["execution_time"].(time.Duration); ok {
			event.ExecutionTime = executionTime
		}
		if function, ok := context["function"].(string); ok {
			event.Function = function
		}
		if caller, ok := context["caller"].(string); ok {
			event.Caller = caller
		}
		if contractType, ok := context["contract_type"].(ContractType); ok {
			event.ContractType = contractType
		}
		if success, ok := context["success"].(bool); ok {
			event.Success = success
		}
	}
	
	// Enrich the event with additional context
	enrichedEvent, err := cei.eventEnricher.EnrichEvent(event)
	if err != nil {
		log.Printf("Warning: Failed to enrich contract event: %v", err)
		enrichedEvent = event
	}
	
	// Transform event to different formats if needed
	transformedEvents, err := cei.eventTransformer.TransformEvent(enrichedEvent)
	if err != nil {
		log.Printf("Warning: Failed to transform contract event: %v", err)
		transformedEvents = []*StandardContractEvent{enrichedEvent}
	}
	
	// Send events to appropriate sinks
	for _, transformedEvent := range transformedEvents {
		err := cei.sendToEventSystem(transformedEvent)
		if err != nil {
			log.Printf("Warning: Failed to send contract event to event system: %v", err)
		}
	}
	
	// Trigger any subscribed contracts
	cei.triggerSubscriptions(enrichedEvent)
	
	log.Printf("📡 Contract event emitted: %s from %s", eventType, contractAddress)
	
	return nil
}

// SubscribeToEvents creates an event subscription for a contract
func (cei *ContractEventIntegrator) SubscribeToEvents(contractAddress string, eventPattern json.RawMessage, callback string, filterCriteria json.RawMessage) (uuid.UUID, error) {
	subscription := &ContractEventSubscription{
		ID:              uuid.New(),
		ContractAddress: contractAddress,
		EventPattern:    eventPattern,
		FilterCriteria:  filterCriteria,
		Callback:        callback,
		Priority:        1,
		Active:          true,
		CreatedAt:       time.Now(),
		TriggerCount:    0,
	}
	
	cei.eventSubscriptions[subscription.ID] = subscription
	
	log.Printf("📡 Contract %s subscribed to events with callback %s", contractAddress, callback)
	
	return subscription.ID, nil
}

// UnsubscribeFromEvents removes an event subscription
func (cei *ContractEventIntegrator) UnsubscribeFromEvents(subscriptionID uuid.UUID) error {
	if subscription, exists := cei.eventSubscriptions[subscriptionID]; exists {
		delete(cei.eventSubscriptions, subscriptionID)
		log.Printf("📡 Contract %s unsubscribed from events", subscription.ContractAddress)
		return nil
	}
	
	return fmt.Errorf("subscription not found: %s", subscriptionID)
}

// CreateContractEventSink creates a dedicated event sink for a contract
func (cei *ContractEventIntegrator) CreateContractEventSink(contractAddress, sinkType string, configuration json.RawMessage) (*ContractEventSink, error) {
	sink := &ContractEventSink{
		ID:              uuid.New(),
		ContractAddress: contractAddress,
		SinkType:        sinkType,
		Configuration:   configuration,
		BufferSize:      100,
		FlushInterval:   5 * time.Second,
		Active:          true,
		CreatedAt:       time.Now(),
		EventCount:      0,
		LastActivity:    time.Now(),
	}
	
	cei.contractEventSinks[contractAddress] = sink
	
	// Create the actual sink in the event system
	err := cei.createEventSystemSink(sink)
	if err != nil {
		return nil, fmt.Errorf("failed to create event system sink: %w", err)
	}
	
	log.Printf("📡 Created contract event sink: %s for contract %s", sinkType, contractAddress)
	
	return sink, nil
}

// Helper methods

func (cei *ContractEventIntegrator) sendToEventSystem(event *StandardContractEvent) error {
	// Convert StandardContractEvent to native LedgerEvent format
	// This integration point connects to the existing event streaming system
	
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	
	// Create event message for the event stream system
	eventMessage := map[string]interface{}{
		"type":      event.Type,
		"source":    event.Source,
		"data":      eventJSON,
		"timestamp": event.Time,
	}
	
	// Send to event stream PID (this would be the actual integration)
	ctx := cei.system.Root
	ctx.Send(cei.eventStreamPID, eventMessage)
	
	return nil
}

func (cei *ContractEventIntegrator) triggerSubscriptions(event *StandardContractEvent) {
	for _, subscription := range cei.eventSubscriptions {
		if subscription.Active && cei.matchesEventPattern(event, subscription.EventPattern) {
			cei.triggerSubscription(subscription, event)
		}
	}
}

func (cei *ContractEventIntegrator) triggerSubscription(subscription *ContractEventSubscription, event *StandardContractEvent) {
	// Create event triggered message
	triggeredMsg := &EventTriggered{
		SubscriptionID:  subscription.ID,
		ContractAddress: subscription.ContractAddress,
		Event:          event.Data,
		Timestamp:      time.Now(),
		Callback:       subscription.Callback,
	}
	
	// Send to contract manager to route to the specific contract
	ctx := cei.system.Root
	ctx.Send(cei.contractManager, triggeredMsg)
	
	// Update subscription statistics
	subscription.TriggerCount++
	now := time.Now()
	subscription.LastTriggered = &now
	
	log.Printf("📡 Triggered event subscription %s for contract %s", subscription.ID, subscription.ContractAddress)
}

func (cei *ContractEventIntegrator) matchesEventPattern(event *StandardContractEvent, pattern json.RawMessage) bool {
	// Simple pattern matching implementation
	// In a full implementation, this would support complex pattern matching
	
	if len(pattern) == 0 {
		return true // Empty pattern matches all events
	}
	
	var patternMap map[string]interface{}
	if err := json.Unmarshal(pattern, &patternMap); err != nil {
		return false
	}
	
	// Match event type if specified
	if eventType, ok := patternMap["type"].(string); ok {
		if event.Type != eventType {
			return false
		}
	}
	
	// Match contract address if specified
	if contractAddr, ok := patternMap["contract_address"].(string); ok {
		if event.ContractAddress != contractAddr {
			return false
		}
	}
	
	// Match function if specified
	if function, ok := patternMap["function"].(string); ok {
		if event.Function != function {
			return false
		}
	}
	
	return true
}

func (cei *ContractEventIntegrator) createEventSystemSink(sink *ContractEventSink) error {
	// Integration with existing event system sink creation
	// This would call the existing RegistryActor to create a sink
	
	sinkConfig := map[string]interface{}{
		"name":     fmt.Sprintf("contract_%s_sink", sink.ContractAddress),
		"type":     sink.SinkType,
		"config":   sink.Configuration,
		"filter":   sink.Filter,
		"enabled":  sink.Active,
	}
	
	// Send sink creation message to registry
	ctx := cei.system.Root
	ctx.Send(cei.registryPID, sinkConfig)
	
	return nil
}

// Factory functions for supporting components

func NewContractEventTransformer() *ContractEventTransformer {
	return &ContractEventTransformer{
		transformationRules: make(map[string]*EventTransformationRule),
		formatters:         make(map[EventFormat]*EventFormatter),
	}
}

func NewContractEventEnricher() *ContractEventEnricher {
	return &ContractEventEnricher{
		enrichmentSources: make(map[string]*EnrichmentSource),
		contextProviders:  make([]ContextProvider, 0),
	}
}

// Placeholder implementations for transformer and enricher

func (cet *ContractEventTransformer) TransformEvent(event *StandardContractEvent) ([]*StandardContractEvent, error) {
	// Simple implementation - just return the original event
	// In a full implementation, this would apply transformation rules
	return []*StandardContractEvent{event}, nil
}

func (cee *ContractEventEnricher) EnrichEvent(event *StandardContractEvent) (*StandardContractEvent, error) {
	// Simple implementation - add basic enrichment
	if event.EnrichmentData == nil {
		event.EnrichmentData = make(map[string]interface{})
	}
	
	event.EnrichmentData["enriched_at"] = time.Now()
	event.EnrichmentData["enricher_version"] = "1.0.0"
	
	return event, nil
}
