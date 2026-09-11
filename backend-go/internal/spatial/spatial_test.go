package spatial

import (
	"metro-synapse/backend/internal/models"
	"testing"
)

func TestGeohashEncodeDecode(t *testing.T) {
	// Bengaluru coordinates (MG Road)
	lat := 12.9716
	lon := 77.5946

	hash := EncodeGeohash(lat, lon, 8)
	if len(hash) != 8 {
		t.Fatalf("Expected geohash length 8, got %d (hash: %s)", len(hash), hash)
	}

	decodedLoc, bbox, err := DecodeGeohash(hash)
	if err != nil {
		t.Fatalf("Failed to decode geohash: %v", err)
	}

	dist := HaversineDistanceMeters(models.Location{Lat: lat, Lon: lon}, decodedLoc)
	if dist > 50.0 { // Precision 8 is within ~38m
		t.Fatalf("Decoded location too far from origin: %f meters", dist)
	}

	if lat < bbox.MinLat || lat > bbox.MaxLat || lon < bbox.MinLon || lon > bbox.MaxLon {
		t.Fatalf("Original point outside decoded bounding box")
	}
}

func TestSpatialIndexRadiusQuery(t *testing.T) {
	idx := NewSpatialIndex(7)

	center := models.Location{Lat: 12.9716, Lon: 77.5946}
	nearNode := models.Location{Lat: 12.9720, Lon: 77.5950} // ~60m away
	farNode := models.Location{Lat: 13.0827, Lon: 80.2707}  // Chennai ~290km away

	idx.Insert(SpatialItem{
		ID:       "node-near",
		Type:     "SENSOR_NODE",
		Location: nearNode,
	})

	idx.Insert(SpatialItem{
		ID:       "node-far",
		Type:     "SENSOR_NODE",
		Location: farNode,
	})

	results := idx.QueryRadius(center, 500.0, "SENSOR_NODE")
	if len(results) != 1 {
		t.Fatalf("Expected 1 result within 500m, got %d", len(results))
	}

	if results[0].Item.ID != "node-near" {
		t.Fatalf("Expected node-near, got %s", results[0].Item.ID)
	}

	nearest := idx.QueryNearest(center, 2, "")
	if len(nearest) != 2 {
		t.Fatalf("Expected 2 nearest items, got %d", len(nearest))
	}
}
