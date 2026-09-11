package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"metro-synapse/backend/internal/broker"
	"metro-synapse/backend/internal/models"
	"metro-synapse/backend/internal/spatial"
	"sync"
	"time"
)

// CityEngine coordinates live urban digital twin telemetry, kinetic energy generation and emergency green corridors
type CityEngine struct {
	mu                sync.RWMutex
	cityID            string
	cityName          string
	spatialIndex      *spatial.SpatialIndex
	broker            *broker.EventBroker
	segments          map[string]*models.RoadSegment
	energyTiles       map[string]*models.EnergyTile
	emergencies       map[string]*models.EmergencyVehicle
	anomalies         []models.AnomalyAlert
	totalKineticMWh   float64
	totalSolarMWh     float64
	totalCO2OffsetKg  float64
	totalSavingsINR   float64
	activeVehicles    int
	stopChan          chan struct{}
}

// NewCityEngine initializes the urban simulation engine
func NewCityEngine(cityID, cityName string, b *broker.EventBroker) *CityEngine {
	engine := &CityEngine{
		cityID:          cityID,
		cityName:        cityName,
		spatialIndex:    spatial.NewSpatialIndex(7),
		broker:          b,
		segments:        make(map[string]*models.RoadSegment),
		energyTiles:     make(map[string]*models.EnergyTile),
		emergencies:     make(map[string]*models.EmergencyVehicle),
		anomalies:       make([]models.AnomalyAlert, 0),
		totalKineticMWh: 142.85,
		totalSolarMWh:   86.40,
		totalCO2OffsetKg: 178500.0,
		totalSavingsINR: 3250000.0,
		activeVehicles:  14850,
		stopChan:        make(chan struct{}),
	}

	engine.seedArterialNetwork()
	return engine
}

// seedArterialNetwork populates the city with real-world inspired smart road corridors
func (ce *CityEngine) seedArterialNetwork() {
	routes := []struct {
		id       string
		name     string
		startLoc models.Location
		endLoc   models.Location
		lanes    int
	}{
		{
			id:       "seg-koramangala-ecity",
			name:     "Hosur Elevated Smart Highway (Outer Ring Link)",
			startLoc: models.Location{Lat: 12.9279, Lon: 77.6271},
			endLoc:   models.Location{Lat: 12.8452, Lon: 77.6602},
			lanes:    6,
		},
		{
			id:       "seg-mgroad-indiranagar",
			name:     "MG Road Arterial Kinetic Axis",
			startLoc: models.Location{Lat: 12.9756, Lon: 77.6066},
			endLoc:   models.Location{Lat: 12.9784, Lon: 77.6408},
			lanes:    4,
		},
		{
			id:       "seg-hebbal-airport",
			name:     "Hebbal Express Kinetic & Solar Corridor",
			startLoc: models.Location{Lat: 13.0358, Lon: 77.5970},
			endLoc:   models.Location{Lat: 13.1986, Lon: 77.7066},
			lanes:    8,
		},
		{
			id:       "seg-whitefield-marathahalli",
			name:     "Whitefield Cyber Tech Spine",
			startLoc: models.Location{Lat: 12.9698, Lon: 77.7500},
			endLoc:   models.Location{Lat: 12.9591, Lon: 77.6974},
			lanes:    6,
		},
	}

	for _, r := range routes {
		seg := &models.RoadSegment{
			ID:              r.id,
			Name:            r.name,
			StartLoc:        r.startLoc,
			EndLoc:          r.endLoc,
			Lanes:           r.lanes,
			SurfaceType:     "piezo_solar_matrix",
			ConditionScore:  94.5 + rand.Float64()*5.0,
			AvgSpeedKmh:     42.0 + rand.Float64()*15.0,
			VehicleDensity:  35.0 + rand.Float64()*20.0,
			KineticMWhToday: 24.5 + rand.Float64()*10.0,
			SolarMWhToday:   18.2 + rand.Float64()*6.0,
			HeatIslandDegC:  -2.4, // Surface cooler due to energy harvesting & biopolymer matrix
			IsGreenCorridor: false,
			LastMaintained:  time.Now().AddDate(0, -1, -10),
		}
		ce.segments[seg.ID] = seg

		// Register in spatial index
		ce.spatialIndex.Insert(spatial.SpatialItem{
			ID:       seg.ID,
			Type:     "ROAD_SEGMENT",
			Location: seg.StartLoc,
			Data:     seg,
		})

		// Seed sample smart tiles along this segment
		for i := 0; i < 5; i++ {
			ratio := float64(i) / 4.0
			tileLat := seg.StartLoc.Lat + ratio*(seg.EndLoc.Lat-seg.StartLoc.Lat)
			tileLon := seg.StartLoc.Lon + ratio*(seg.EndLoc.Lon-seg.StartLoc.Lon)
			tileID := fmt.Sprintf("tile-%s-%d", seg.ID, i+1)

			tile := &models.EnergyTile{
				ID:               tileID,
				SegmentID:        seg.ID,
				Location:         models.Location{Lat: tileLat, Lon: tileLon},
				PiezoOutputWatts: 420.0 + rand.Float64()*180.0,
				SolarOutputWatts: 280.0 + rand.Float64()*70.0,
				StressPressure:   120.0 + rand.Float64()*30.0,
				TemperatureDegC:  28.5 + rand.Float64()*4.0,
				Status:           "OPTIMAL",
				UpdatedAt:        time.Now(),
			}
			ce.energyTiles[tile.ID] = tile

			ce.spatialIndex.Insert(spatial.SpatialItem{
				ID:       tile.ID,
				Type:     "ENERGY_TILE",
				Location: tile.Location,
				Data:     tile,
			})
		}
	}

	// Add initial anomaly alert
	ce.anomalies = append(ce.anomalies, models.AnomalyAlert{
		ID:              "alert-anom-01",
		Timestamp:       time.Now().Add(-12 * time.Minute),
		SegmentID:       "seg-whitefield-marathahalli",
		Severity:        "MEDIUM",
		AnomalyType:     "POTHOLE_EARLY_STAGE",
		ConfidenceScore: 0.942,
		Description:     "Micro-acoustic vibration signature indicates subsurface acoustic resonance deviation at km 3.4.",
		SuggestedAction: "Deploy autonomous thermal resin patch drone within 48 hours.",
		GeminiReasoning: "Gemini Vision analysis of drone feed confirms 4cm micro-fissure forming along outer wheel track.",
	})
}

// Start launches the simulation ticker loop
func (ce *CityEngine) Start() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				ce.tick()
			case <-ce.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop terminates the simulation
func (ce *CityEngine) Stop() {
	close(ce.stopChan)
}

// tick calculates kinetic power generation, updates active emergencies and broadcasts state
func (ce *CityEngine) tick() {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	deltaVehicles := rand.Intn(21) - 10
	ce.activeVehicles += deltaVehicles
	if ce.activeVehicles < 5000 {
		ce.activeVehicles = 5000
	}

	// Calculate incremental kinetic & solar generation (in MWh)
	// Energy = sum of piezo generation across active tiles
	incrementalMWh := 0.00045 + (float64(ce.activeVehicles)/10000.0)*0.00025
	ce.totalKineticMWh += incrementalMWh
	ce.totalSolarMWh += incrementalMWh * 0.65
	ce.totalCO2OffsetKg += incrementalMWh * 780.0 // ~780kg CO2 offset per MWh
	ce.totalSavingsINR += incrementalMWh * 7500.0 // ~₹7.5/kWh value

	// Update tiles with dynamic oscillation
	for _, tile := range ce.energyTiles {
		tile.PiezoOutputWatts = 400.0 + rand.Float64()*250.0
		tile.SolarOutputWatts = 260.0 + rand.Float64()*90.0
		tile.UpdatedAt = time.Now()
	}

	// Update active emergencies (simulated movement towards destination)
	for id, em := range ce.emergencies {
		if em.Status == "CORRIDOR_ACTIVE" {
			em.TimeSavedSeconds += 1.0
			// Progress location slightly
			dLat := (em.DestinationLoc.Lat - em.CurrentLoc.Lat) * 0.05
			dLon := (em.DestinationLoc.Lon - em.CurrentLoc.Lon) * 0.05
			em.CurrentLoc.Lat += dLat
			em.CurrentLoc.Lon += dLon
			em.CurrentSpeedKmh = 72.0 + rand.Float64()*8.0 // High speed on cleared green corridor

			dist := spatial.HaversineDistanceMeters(em.CurrentLoc, em.DestinationLoc)
			if dist < 80.0 {
				em.Status = "ARRIVED"
				// Clear green corridor flag on segments
				for _, segID := range em.AssignedSegments {
					if s, exists := ce.segments[segID]; exists {
						s.IsGreenCorridor = false
					}
				}
			}
		} else if em.Status == "ARRIVED" {
			// Auto clear after 10s
			if time.Since(em.EstimatedArrival) > 10*time.Second {
				delete(ce.emergencies, id)
			}
		}
	}

	// Broadcast full state snapshot to broker
	state := ce.getSnapshotUnsafe()
	ce.broker.Publish("city:state", state)
}

// DispatchGreenCorridor initiates an emergency wave clearance along optimal segments
func (ce *CityEngine) DispatchGreenCorridor(req models.CorridorDispatchRequest) models.CorridorDispatchResponse {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	vehicleID := fmt.Sprintf("em-%s-%d", req.VehicleType, rand.Intn(9000)+1000)
	dispatchID := fmt.Sprintf("disp-%d", time.Now().Unix())

	// Find optimal road segments for corridor
	assignedSegments := []string{"seg-mgroad-indiranagar", "seg-koramangala-ecity"}
	if req.StartSegmentID != "" {
		assignedSegments = []string{req.StartSegmentID}
	}

	// Activate green corridor on segments
	for _, segID := range assignedSegments {
		if s, ok := ce.segments[segID]; ok {
			s.IsGreenCorridor = true
		}
	}

	destLoc := req.DestinationLoc
	if destLoc.Lat == 0 && destLoc.Lon == 0 {
		// Default to Manipal / Narayana Hospital location
		destLoc = models.Location{Lat: 12.9592, Lon: 77.6534}
	}

	startLoc := models.Location{Lat: 12.9756, Lon: 77.6066}
	if s, ok := ce.segments[assignedSegments[0]]; ok {
		startLoc = s.StartLoc
	}

	distMeters := spatial.HaversineDistanceMeters(startLoc, destLoc)
	standardTransitMinutes := (distMeters / 1000.0) / 20.0 * 60.0 // at 20 km/h gridlock
	clearedTransitMinutes := (distMeters / 1000.0) / 70.0 * 60.0  // at 70 km/h green corridor
	timeSavedSeconds := (standardTransitMinutes - clearedTransitMinutes) * 60.0

	arrival := time.Now().Add(time.Duration(clearedTransitMinutes * float64(time.Minute)))

	ev := &models.EmergencyVehicle{
		ID:               vehicleID,
		Callsign:         fmt.Sprintf("PRANA-LIFE-%d", rand.Intn(900)+100),
		VehicleType:      req.VehicleType,
		Priority:         req.Priority,
		CurrentLoc:       startLoc,
		DestinationLoc:   destLoc,
		DestinationName:  req.TargetHospital,
		CurrentSpeedKmh:  68.0,
		Status:           "CORRIDOR_ACTIVE",
		AssignedSegments: assignedSegments,
		EstimatedArrival: arrival,
		TimeSavedSeconds: timeSavedSeconds,
	}

	ce.emergencies[vehicleID] = ev

	resp := models.CorridorDispatchResponse{
		DispatchID:        dispatchID,
		VehicleID:         vehicleID,
		AssignedSegments:  assignedSegments,
		EstimatedArrival:  arrival,
		TimeSavedSeconds:  timeSavedSeconds,
		SignalPhaseWaveID: fmt.Sprintf("WAVE-PHASE-%x", rand.Intn(0xFFFFFF)),
		Status:            "ACTIVE_GREEN_CORRIDOR",
	}

	// Broadcast alert
	ce.broker.Publish("emergency:dispatch", resp)
	return resp
}

// GetState returns thread-safe snapshot of city state
func (ce *CityEngine) GetState() models.UrbanGridState {
	ce.mu.RLock()
	defer ce.mu.RUnlock()
	return ce.getSnapshotUnsafe()
}

func (ce *CityEngine) getSnapshotUnsafe() models.UrbanGridState {
	segmentList := make([]models.RoadSegment, 0, len(ce.segments))
	for _, s := range ce.segments {
		segmentList = append(segmentList, *s)
	}

	emergencyList := make([]models.EmergencyVehicle, 0, len(ce.emergencies))
	for _, e := range ce.emergencies {
		emergencyList = append(emergencyList, *e)
	}

	activeCorridors := 0
	for _, s := range ce.segments {
		if s.IsGreenCorridor {
			activeCorridors++
		}
	}

	// Calculate instantaneous MW power
	totalPowerMW := (ce.totalKineticMWh + ce.totalSolarMWh) * 0.12

	return models.UrbanGridState{
		CityID:                ce.cityID,
		CityName:              ce.cityName,
		Timestamp:             time.Now().UTC(),
		ActiveVehicles:        ce.activeVehicles,
		TotalKineticMWh:       math.Round(ce.totalKineticMWh*100) / 100,
		TotalSolarMWh:         math.Round(ce.totalSolarMWh*100) / 100,
		TotalPowerGeneratedMW: math.Round(totalPowerMW*100) / 100,
		TotalCO2OffsetTonnes:  math.Round((ce.totalCO2OffsetKg/1000.0)*100) / 100,
		EconomicSavingsINR:    math.Round(ce.totalSavingsINR),
		ActiveCorridorsCount:  activeCorridors,
		AvgResponseTimeMin:    4.2, // Down from 28.5 mins in traditional gridlock
		Segments:              segmentList,
		ActiveEmergencies:     emergencyList,
		RecentAnomalies:       ce.anomalies,
	}
}
