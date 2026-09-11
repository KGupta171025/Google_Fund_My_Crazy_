package spatial

import (
	"math"
	"metro-synapse/backend/internal/models"
	"sort"
	"sync"
)

// SpatialItem wraps any entity with location and arbitrary metadata
type SpatialItem struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"` // "ROAD_SEGMENT", "ENERGY_TILE", "EMERGENCY_VEHICLE", "SENSOR_NODE"
	Location models.Location `json:"location"`
	Data     interface{}     `json:"data"`
}

// SpatialSearchResult holds search result with exact distance
type SpatialSearchResult struct {
	Item           SpatialItem `json:"item"`
	DistanceMeters float64     `json:"distance_meters"`
}

// SpatialIndex is a high-speed concurrent spatial grid index
type SpatialIndex struct {
	mu         sync.RWMutex
	precision  int
	buckets    map[string][]SpatialItem
	itemLookup map[string]string // ID -> Geohash bucket key
}

// NewSpatialIndex creates a new spatial search engine
func NewSpatialIndex(precision int) *SpatialIndex {
	if precision <= 0 || precision > 12 {
		precision = 7 // ~150m cell size
	}
	return &SpatialIndex{
		precision:  precision,
		buckets:    make(map[string][]SpatialItem),
		itemLookup: make(map[string]string),
	}
}

// Insert adds or updates an item in the spatial index
func (si *SpatialIndex) Insert(item SpatialItem) {
	si.mu.Lock()
	defer si.mu.Unlock()

	// If item already exists, remove it from old bucket
	if oldHash, exists := si.itemLookup[item.ID]; exists {
		si.removeFromBucket(oldHash, item.ID)
	}

	hash := EncodeGeohash(item.Location.Lat, item.Location.Lon, si.precision)
	si.buckets[hash] = append(si.buckets[hash], item)
	si.itemLookup[item.ID] = hash
}

// Remove removes an item by its ID
func (si *SpatialIndex) Remove(id string) bool {
	si.mu.Lock()
	defer si.mu.Unlock()

	hash, exists := si.itemLookup[id]
	if !exists {
		return false
	}

	si.removeFromBucket(hash, id)
	delete(si.itemLookup, id)
	return true
}

func (si *SpatialIndex) removeFromBucket(hash, id string) {
	items := si.buckets[hash]
	for i, item := range items {
		if item.ID == id {
			si.buckets[hash] = append(items[:i], items[i+1:]...)
			break
		}
	}
	if len(si.buckets[hash]) == 0 {
		delete(si.buckets, hash)
	}
}

// QueryRadius finds all items within radiusMeters from target location
func (si *SpatialIndex) QueryRadius(target models.Location, radiusMeters float64, filterType string) []SpatialSearchResult {
	si.mu.RLock()
	defer si.mu.RUnlock()

	var results []SpatialSearchResult

	// Collect candidate items from buckets
	for _, items := range si.buckets {
		for _, item := range items {
			if filterType != "" && item.Type != filterType {
				continue
			}
			dist := HaversineDistanceMeters(target, item.Location)
			if dist <= radiusMeters {
				results = append(results, SpatialSearchResult{
					Item:           item,
					DistanceMeters: dist,
				})
			}
		}
	}

	// Sort results by proximity
	sort.Slice(results, func(i, j int) bool {
		return results[i].DistanceMeters < results[j].DistanceMeters
	})

	return results
}

// QueryNearest finds K nearest items to target location
func (si *SpatialIndex) QueryNearest(target models.Location, k int, filterType string) []SpatialSearchResult {
	if k <= 0 {
		return nil
	}

	all := si.QueryRadius(target, math.MaxFloat64, filterType)
	if len(all) > k {
		return all[:k]
	}
	return all
}

// Count returns total indexed items
func (si *SpatialIndex) Count() int {
	si.mu.RLock()
	defer si.mu.RUnlock()
	return len(si.itemLookup)
}
