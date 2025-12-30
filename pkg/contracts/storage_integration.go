package contracts

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/google/uuid"
)

// ContractStorageIntegrator manages integration with the document storage system
// This enables contracts to have persistent state stored in the blockchain document system
type ContractStorageIntegrator struct {
	documentStorePID *actor.PID
	system           *actor.ActorSystem

	// Storage management
	stateStorageManager *StateStorageManager
	historyManager      *ContractHistoryManager
	queryInterface      *ContractQueryInterface

	// Caching and performance
	stateCache   map[string]*CachedContractState
	cacheManager *StateCacheManager
}

// StateStorageManager handles contract state persistence
type StateStorageManager struct {
	compressionEnabled bool
	encryptionEnabled  bool
	versioningEnabled  bool

	// State serialization
	serializers map[StateFormat]*StateSerializer
	compressors map[CompressionType]*StateCompressor
}

// ContractHistoryManager manages temporal contract state history
type ContractHistoryManager struct {
	retentionPolicy *HistoryRetentionPolicy
	snapshotManager *StateSnapshotManager
	timeSeriesIndex *TimeSeriesIndex
}

// ContractQueryInterface provides temporal querying capabilities
type ContractQueryInterface struct {
	queryEngine       *TemporalQueryEngine
	indexManager      *ContractIndexManager
	aggregationEngine *StateAggregationEngine
}

// CachedContractState represents cached contract state
type CachedContractState struct {
	ContractAddress string          `json:"contract_address"`
	State           json.RawMessage `json:"state"`
	StateRoot       string          `json:"state_root"`
	Version         int64           `json:"version"`
	CachedAt        time.Time       `json:"cached_at"`
	ExpiresAt       time.Time       `json:"expires_at"`
	AccessCount     int64           `json:"access_count"`
	LastAccessed    time.Time       `json:"last_accessed"`
	Dirty           bool            `json:"dirty"`
}

// StateCacheManager manages contract state caching
type StateCacheManager struct {
	cacheSize      int
	ttl            time.Duration
	evictionPolicy EvictionPolicy
	cacheStats     *CacheStatistics
}

// ContractStateDocument represents a contract state stored as a document
type ContractStateDocument struct {
	// CloudEvents metadata
	SpecVersion     string    `json:"specversion"`
	Type            string    `json:"type"`
	Source          string    `json:"source"`
	ID              string    `json:"id"`
	Time            time.Time `json:"time"`
	Subject         string    `json:"subject"`
	DataContentType string    `json:"datacontenttype"`

	// Contract-specific metadata
	ContractAddress string       `json:"contract_address"`
	ContractType    ContractType `json:"contract_type"`
	StateVersion    int64        `json:"state_version"`
	StateRoot       string       `json:"state_root"`
	PreviousRoot    string       `json:"previous_root,omitempty"`

	// State data (base64 encoded)
	StateData   string          `json:"state_data"`
	StateFormat StateFormat     `json:"state_format"`
	Compression CompressionType `json:"compression,omitempty"`

	// Blockchain context
	BlockHeight     int64  `json:"block_height"`
	TransactionHash string `json:"transaction_hash,omitempty"`

	// Temporal metadata
	ValidFrom  time.Time  `json:"valid_from"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`

	// Change metadata
	ChangeReason string `json:"change_reason"`
	ChangedBy    string `json:"changed_by,omitempty"`

	// Indexing and categorization
	Tags       []string               `json:"tags"`
	Properties map[string]interface{} `json:"properties"`
}

// StateFormat defines different formats for storing contract state
type StateFormat string

const (
	StateFormatJSON     StateFormat = "json"
	StateFormatBinary   StateFormat = "binary"
	StateFormatProtobuf StateFormat = "protobuf"
	StateFormatAvro     StateFormat = "avro"
	StateFormatMsgPack  StateFormat = "msgpack"
)

// CompressionType defines compression algorithms for state data
type CompressionType string

const (
	CompressionNone   CompressionType = "none"
	CompressionGZip   CompressionType = "gzip"
	CompressionSnappy CompressionType = "snappy"
	CompressionLZ4    CompressionType = "lz4"
	CompressionZSTD   CompressionType = "zstd"
)

// EvictionPolicy defines cache eviction strategies
type EvictionPolicy string

const (
	EvictionPolicyLRU  EvictionPolicy = "lru"
	EvictionPolicyLFU  EvictionPolicy = "lfu"
	EvictionPolicyFIFO EvictionPolicy = "fifo"
	EvictionPolicyTTL  EvictionPolicy = "ttl"
)

// HistoryRetentionPolicy defines how long to keep contract state history
type HistoryRetentionPolicy struct {
	MaxAge            time.Duration `json:"max_age"`
	MaxVersions       int           `json:"max_versions"`
	SnapshotInterval  time.Duration `json:"snapshot_interval"`
	CompactionEnabled bool          `json:"compaction_enabled"`
}

// CacheStatistics tracks cache performance
type CacheStatistics struct {
	Hits      int64     `json:"hits"`
	Misses    int64     `json:"misses"`
	Evictions int64     `json:"evictions"`
	HitRate   float64   `json:"hit_rate"`
	Size      int       `json:"size"`
	TotalSize int64     `json:"total_size"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StateSerializer handles different state serialization formats
type StateSerializer struct {
	Format      StateFormat `json:"format"`
	ContentType string      `json:"content_type"`
	Version     string      `json:"version"`
}

// StateCompressor handles state compression
type StateCompressor struct {
	Type       CompressionType `json:"type"`
	Level      int             `json:"level"`
	Dictionary []byte          `json:"dictionary,omitempty"`
}

// TemporalQueryEngine enables time-based queries on contract state
type TemporalQueryEngine struct {
	timeIndexes    map[string]*TimeIndex
	versionIndexes map[string]*VersionIndex
	queryCache     map[string]*QueryResult
}

// TimeIndex indexes contract states by time
type TimeIndex struct {
	ContractAddress string        `json:"contract_address"`
	TimePoints      []TimePoint   `json:"time_points"`
	Resolution      time.Duration `json:"resolution"`
	LastUpdated     time.Time     `json:"last_updated"`
}

// TimePoint represents a point in time with associated state
type TimePoint struct {
	Timestamp  time.Time `json:"timestamp"`
	StateRoot  string    `json:"state_root"`
	Version    int64     `json:"version"`
	DocumentID string    `json:"document_id"`
}

// VersionIndex indexes contract states by version
type VersionIndex struct {
	ContractAddress string         `json:"contract_address"`
	Versions        []VersionEntry `json:"versions"`
	LatestVersion   int64          `json:"latest_version"`
}

// VersionEntry represents a state version
type VersionEntry struct {
	Version    int64     `json:"version"`
	StateRoot  string    `json:"state_root"`
	Timestamp  time.Time `json:"timestamp"`
	DocumentID string    `json:"document_id"`
	ChangeType string    `json:"change_type"`
}

// QueryResult caches temporal query results
type QueryResult struct {
	Query       string            `json:"query"`
	Results     []json.RawMessage `json:"results"`
	GeneratedAt time.Time         `json:"generated_at"`
	ExpiresAt   time.Time         `json:"expires_at"`
}

// NewContractStorageIntegrator creates a new storage integrator
func NewContractStorageIntegrator(documentStorePID *actor.PID, system *actor.ActorSystem) *ContractStorageIntegrator {
	return &ContractStorageIntegrator{
		documentStorePID:    documentStorePID,
		system:              system,
		stateStorageManager: NewStateStorageManager(),
		historyManager:      NewContractHistoryManager(),
		queryInterface:      NewContractQueryInterface(),
		stateCache:          make(map[string]*CachedContractState),
		cacheManager:        NewStateCacheManager(),
	}
}

// SaveContractState persists contract state to document storage
func (csi *ContractStorageIntegrator) SaveContractState(contract *Contract, reason string, changedBy string) error {
	// Serialize contract state
	stateJSON, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("failed to serialize contract state: %w", err)
	}

	// Compress state data if enabled
	compressedData, compression := csi.compressStateData(stateJSON)

	// Encode state data
	encodedData := base64.StdEncoding.EncodeToString(compressedData)

	// Create state document
	stateDoc := &ContractStateDocument{
		SpecVersion:     "1.0",
		Type:            "com.xenese.contract.state",
		Source:          "xenese-dlt://contracts",
		ID:              uuid.New().String(),
		Time:            time.Now(),
		Subject:         contract.Address,
		DataContentType: "application/json",

		ContractAddress: contract.Address,
		ContractType:    contract.Type,
		StateVersion:    csi.nextStateVersion(contract.Address),
		StateRoot:       contract.StateRoot,

		StateData:   encodedData,
		StateFormat: StateFormatJSON,
		Compression: compression,

		BlockHeight: csi.getCurrentBlockHeight(),

		ValidFrom:    time.Now(),
		ChangeReason: reason,
		ChangedBy:    changedBy,

		Tags: []string{
			"contract-state",
			string(contract.Status),
		},
		Properties: map[string]interface{}{
			"contract_name":    contract.Name,
			"contract_owner":   contract.Owner,
			"execution_count":  contract.ExecutionCount,
			"success_rate":     contract.SuccessRate,
			"adaptation_score": contract.AdaptationScore,
		},
	}

	// Store document
	err = csi.storeStateDocument(stateDoc)
	if err != nil {
		return fmt.Errorf("failed to store state document: %w", err)
	}

	stateVersion := stateDoc.StateVersion
	// Update cache
	csi.updateStateCache(contract.Address, stateJSON, contract.StateRoot, stateVersion)

	// Update indexes
	csi.updateIndexes(stateDoc)

	log.Printf(
		"💾 Contract state saved: %s (contract version %s, state version %d)",
		contract.Address,
		contract.Version,
		stateVersion,
	)

	return nil
}

// LoadContractState retrieves contract state from document storage
func (csi *ContractStorageIntegrator) LoadContractState(contractAddress string, version *int64) (*Contract, error) {
	// Check cache first
	if cached := csi.getFromCache(contractAddress, version); cached != nil {
		var contract Contract
		err := json.Unmarshal(cached.State, &contract)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize cached state: %w", err)
		}

		// Update cache statistics
		csi.updateStatsOnHit()
		cached.AccessCount++
		cached.LastAccessed = time.Now()

		return &contract, nil
	}

	// Cache miss - load from document storage
	csi.updateStatsOnMiss()

	// Find the appropriate state document
	documentID, err := csi.findStateDocument(contractAddress, version)
	if err != nil {
		return nil, fmt.Errorf("failed to find state document: %w", err)
	}

	// Retrieve document
	stateDoc, err := csi.retrieveStateDocument(documentID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve state document: %w", err)
	}

	// Decode and decompress state data
	compressedData, err := base64.StdEncoding.DecodeString(stateDoc.StateData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode state data: %w", err)
	}

	stateData, err := csi.decompressStateData(compressedData, stateDoc.Compression)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress state data: %w", err)
	}

	// Deserialize contract
	var contract Contract
	err = json.Unmarshal(stateData, &contract)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize contract state: %w", err)
	}

	// Update cache
	csi.updateStateCache(contractAddress, stateData, stateDoc.StateRoot, stateDoc.StateVersion)

	log.Printf("💾 Contract state loaded: %s (version %d)", contractAddress, stateDoc.StateVersion)

	return &contract, nil
}

// QueryContractHistory performs temporal queries on contract state history
func (csi *ContractStorageIntegrator) QueryContractHistory(contractAddress string, query map[string]interface{}, startTime, endTime time.Time) ([]*ContractStateDocument, error) {
	// Build temporal query
	temporalQuery := map[string]interface{}{
		"contract_address": contractAddress,
		"valid_from": map[string]interface{}{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	// Add additional query criteria, but don't allow overrides of reserved keys
	for key, value := range query {
		if key == "contract_address" || key == "valid_from" {
			continue
		}
		temporalQuery[key] = value
	}

	// Execute query through document storage system
	results, err := csi.executeTemporalQuery(temporalQuery)
	if err != nil {
		return nil, fmt.Errorf("temporal query failed: %w", err)
	}

	log.Printf("💾 Temporal query executed: %s (%d results)", contractAddress, len(results))

	return results, nil
}

// GetContractStateAtTime retrieves contract state as it existed at a specific time
func (csi *ContractStorageIntegrator) GetContractStateAtTime(contractAddress string, timestamp time.Time) (*Contract, error) {
	temporalQuery := map[string]interface{}{
		"contract_address": contractAddress,
		"valid_from": map[string]interface{}{
			"$lte": timestamp,
		},
		"$or": []map[string]interface{}{
			{"valid_until": map[string]interface{}{"$gte": timestamp}},
			{"valid_until": nil},
		},
	}

	results, err := csi.executeTemporalQuery(temporalQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query contract history: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no contract state found at timestamp %v", timestamp)
	}

	// Pick the most recent one before / at timestamp
	stateDoc := results[len(results)-1]

	// Decode / decompress / unmarshal
	compressedData, err := base64.StdEncoding.DecodeString(stateDoc.StateData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode state data: %w", err)
	}

	stateData, err := csi.decompressStateData(compressedData, stateDoc.Compression)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress state data: %w", err)
	}

	var contract Contract
	if err := json.Unmarshal(stateData, &contract); err != nil {
		return nil, fmt.Errorf("failed to deserialize contract state: %w", err)
	}

	log.Printf("💾 Contract state retrieved at time %v: %s", timestamp, contractAddress)
	return &contract, nil
}

// Helper methods

func (csi *ContractStorageIntegrator) compressStateData(data []byte) ([]byte, CompressionType) {
	// TODO: plug in real compression (gzip/zstd/etc.) based on CompressionType
	// Simple compression implementation
	// In production, this would use actual compression libraries
	return data, CompressionNone
}

func (csi *ContractStorageIntegrator) decompressStateData(data []byte, compression CompressionType) ([]byte, error) {
	// TODO: plug in real compression (gzip/zstd/etc.) based on CompressionType
	// Simple decompression implementation
	switch compression {
	case CompressionNone:
		return data, nil
	default:
		return data, nil
	}
}

func (csi *ContractStorageIntegrator) storeStateDocument(stateDoc *ContractStateDocument) error {
	// Integration with document storage system
	// This would call the DocumentActor to store the state document

	docJSON, err := json.Marshal(stateDoc)
	if err != nil {
		return err
	}

	// Send store message to document storage
	ctx := csi.system.Root
	storeMsg := map[string]interface{}{
		"action": "store",
		"data":   docJSON,
		"type":   "contract_state",
	}

	ctx.Send(csi.documentStorePID, storeMsg)

	return nil
}

func (csi *ContractStorageIntegrator) retrieveStateDocument(documentID string) (*ContractStateDocument, error) {
	// Integration with document storage system
	// This would call the DocumentActor to retrieve the state document

	return nil, fmt.Errorf("retrieveStateDocument not implemented (documentID=%s)", documentID)
}

func (csi *ContractStorageIntegrator) findStateDocument(contractAddress string, version *int64) (string, error) {
	// Search for the appropriate state document
	// This would use the document query system

	return "", fmt.Errorf("findStateDocument not implemented (contractAddress=%s, version=%v)", contractAddress, version)
}

func (csi *ContractStorageIntegrator) executeTemporalQuery(query map[string]interface{}) ([]*ContractStateDocument, error) {
	// Execute temporal query through document system
	// This would integrate with the existing document query functionality

	return nil, fmt.Errorf("executeTemporalQuery not implemented (query=%v)", query)
}

func (csi *ContractStorageIntegrator) getCurrentBlockHeight() int64 {
	// Get current block height from ledger
	// This would integrate with the LedgerActor
	return 0
}

func (csi *ContractStorageIntegrator) getFromCache(contractAddress string, version *int64) *CachedContractState {
	// getFromCache returns the cached state for the given contract.
	// NOTE: current implementation caches only the latest version per contract.
	// If version != nil and doesn't match the cached version, this is a miss.
	if cached, exists := csi.stateCache[contractAddress]; exists {
		if version == nil || cached.Version == *version {
			if time.Now().Before(cached.ExpiresAt) {
				return cached
			}
		}
	}
	return nil
}

func (csi *ContractStorageIntegrator) updateStateCache(contractAddress string, state []byte, stateRoot string, version int64) {
	cached := &CachedContractState{
		ContractAddress: contractAddress,
		State:           state,
		StateRoot:       stateRoot,
		Version:         version,
		CachedAt:        time.Now(),
		ExpiresAt:       time.Now().Add(csi.cacheManager.ttl),
		AccessCount:     1,
		LastAccessed:    time.Now(),
		Dirty:           false,
	}

	csi.stateCache[contractAddress] = cached
}

func (csi *ContractStorageIntegrator) updateIndexes(stateDoc *ContractStateDocument) {
	// Update temporal indexes
	// This would update the time and version indexes for faster queries
}

// nextStateVersion returns the next monotonic state version for a contract
func (csi *ContractStorageIntegrator) nextStateVersion(contractAddress string) int64 {
	// Look at cached state if present
	if cached, ok := csi.stateCache[contractAddress]; ok {
		return cached.Version + 1
	}
	// Fallback to 1 if this is the first time
	return 1
}

// updateStatsOnHit updates cache statistics on cache hit
func (csi *ContractStorageIntegrator) updateStatsOnHit() {
	csi.cacheManager.cacheStats.Hits++
	csi.recomputeHitRate()
}

// updateStatsOnMiss updates cache statistics on cache miss
func (csi *ContractStorageIntegrator) updateStatsOnMiss() {
	csi.cacheManager.cacheStats.Misses++
	csi.recomputeHitRate()
}

// recomputeHitRate recalculates cache hit rate
func (csi *ContractStorageIntegrator) recomputeHitRate() {
	total := csi.cacheManager.cacheStats.Hits + csi.cacheManager.cacheStats.Misses
	if total == 0 {
		csi.cacheManager.cacheStats.HitRate = 0
		return
	}
	csi.cacheManager.cacheStats.HitRate = float64(csi.cacheManager.cacheStats.Hits) / float64(total)
	csi.cacheManager.cacheStats.UpdatedAt = time.Now()
}

// Factory functions

func NewStateStorageManager() *StateStorageManager {
	return &StateStorageManager{
		compressionEnabled: true,
		encryptionEnabled:  false,
		versioningEnabled:  true,
		serializers:        make(map[StateFormat]*StateSerializer),     // TODO: register JSON/Protobuf serializers here
		compressors:        make(map[CompressionType]*StateCompressor), // TODO: register compression algorithms here
	}
}

func NewContractHistoryManager() *ContractHistoryManager {
	return &ContractHistoryManager{
		retentionPolicy: &HistoryRetentionPolicy{
			MaxAge:            30 * 24 * time.Hour, // 30 days
			MaxVersions:       1000,                // 1000 versions
			SnapshotInterval:  24 * time.Hour,      // Daily snapshots
			CompactionEnabled: true,
		},
	}
}

func NewContractQueryInterface() *ContractQueryInterface {
	return &ContractQueryInterface{
		queryEngine:       &TemporalQueryEngine{},
		indexManager:      &ContractIndexManager{},
		aggregationEngine: &StateAggregationEngine{},
	}
}

func NewStateCacheManager() *StateCacheManager {
	return &StateCacheManager{
		cacheSize:      1000,
		ttl:            10 * time.Minute,
		evictionPolicy: EvictionPolicyLRU,
		cacheStats: &CacheStatistics{
			UpdatedAt: time.Now(),
		},
	}
}

// Placeholder types for completeness
type StateSnapshotManager struct{}
type TimeSeriesIndex struct{}
type ContractIndexManager struct{}
type StateAggregationEngine struct{}
