package models

import "time"

// Location defines 3D geospatial coordinates
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Alt float64 `json:"alt,omitempty"`
}

// BoundingBox represents a geographic search boundary
type BoundingBox struct {
	MinLat float64 `json:"min_lat"`
	MinLon float64 `json:"min_lon"`
	MaxLat float64 `json:"max_lat"`
	MaxLon float64 `json:"max_lon"`
}

// RoadSegment represents an arterial or smart highway segment with kinetic tiles
type RoadSegment struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	StartLoc           Location  `json:"start_loc"`
	EndLoc             Location  `json:"end_loc"`
	Lanes              int       `json:"lanes"`
	SurfaceType        string    `json:"surface_type"`        // "piezo_solar_matrix", "asphalt_hybrid", "permeable_biopolymer"
	ConditionScore     float64   `json:"condition_score"`     // 0.0 - 100.0 (structural health)
	AvgSpeedKmh        float64   `json:"avg_speed_kmh"`
	VehicleDensity     float64   `json:"vehicle_density"`     // vehicles/km/lane
	KineticMWhToday    float64   `json:"kinetic_mwh_today"`   // Energy generated today
	SolarMWhToday      float64   `json:"solar_mwh_today"`     // Solar captured today
	HeatIslandDegC     float64   `json:"heat_island_deg_c"`   // Surface micro-climate delta
	IsGreenCorridor    bool      `json:"is_green_corridor"`   // True when emergency corridor active
	LastMaintained     time.Time `json:"last_maintained"`
}

// EnergyTile represents an individual smart piezo-photovoltaic road tile
type EnergyTile struct {
	ID               string    `json:"id"`
	SegmentID        string    `json:"segment_id"`
	Location         Location  `json:"location"`
	PiezoOutputWatts float64   `json:"piezo_output_watts"`
	SolarOutputWatts float64   `json:"solar_output_watts"`
	StressPressure   float64   `json:"stress_pressure_kpa"`
	TemperatureDegC  float64   `json:"temperature_deg_c"`
	Status           string    `json:"status"` // "OPTIMAL", "SURGE", "DEGRADED", "MAINTENANCE_REQUIRED"
	UpdatedAt        time.Time `json:"updated_at"`
}

// EmergencyVehicle represents ambulances, fire trucks or disaster response units
type EmergencyVehicle struct {
	ID               string     `json:"id"`
	Callsign         string     `json:"callsign"`
	VehicleType      string     `json:"vehicle_type"` // "AMBULANCE_CRITICAL", "FIRE_ENGINE", "DISASTER_RESCUE"
	Priority         int        `json:"priority"`     // 1 (Highest) to 5
	CurrentLoc       Location   `json:"current_loc"`
	DestinationLoc   Location   `json:"destination_loc"`
	DestinationName  string     `json:"destination_name"`
	CurrentSpeedKmh  float64    `json:"current_speed_kmh"`
	Status           string     `json:"status"` // "DISPATCHED", "CORRIDOR_ACTIVE", "ARRIVED", "IDLE"
	AssignedSegments []string   `json:"assigned_segments"`
	EstimatedArrival time.Time  `json:"estimated_arrival"`
	TimeSavedSeconds float64    `json:"time_saved_seconds"`
}

// TelemetryPacket represents high-frequency IoT streaming telemetry
type TelemetryPacket struct {
	NodeID          string    `json:"node_id"`
	Timestamp       time.Time `json:"timestamp"`
	VehicleCount    int       `json:"vehicle_count"`
	KineticJoules   float64   `json:"kinetic_joules"`
	SolarWatts      float64   `json:"solar_watts"`
	AvgVelocityKmh  float64   `json:"avg_velocity_kmh"`
	AcousticDecibels float64  `json:"acoustic_decibels"`
	MicroVibrationHz float64  `json:"micro_vibration_hz"`
	AirQualityAQI   int       `json:"air_quality_aqi"`
	SurfaceTempDegC float64   `json:"surface_temp_deg_c"`
}

// AnomalyAlert represents an AI-detected road defect, congestion wave or disaster alert
type AnomalyAlert struct {
	ID              string    `json:"id"`
	Timestamp       time.Time `json:"timestamp"`
	SegmentID       string    `json:"segment_id"`
	Severity        string    `json:"severity"`        // "LOW", "MEDIUM", "HIGH", "CRITICAL"
	AnomalyType     string    `json:"anomaly_type"`    // "POTHOLE_EARLY_STAGE", "SUBSIDENCE", "HEAT_ANOMALY", "CORRIDOR_BLOCKAGE"
	ConfidenceScore float64   `json:"confidence_score"`
	Description     string    `json:"description"`
	SuggestedAction string    `json:"suggested_action"`
	GeminiReasoning string    `json:"gemini_reasoning,omitempty"`
}

// UrbanGridState represents the complete real-time digital twin state of the city
type UrbanGridState struct {
	CityID               string              `json:"city_id"`
	CityName             string              `json:"city_name"`
	Timestamp            time.Time           `json:"timestamp"`
	ActiveVehicles       int                 `json:"active_vehicles"`
	TotalKineticMWh      float64             `json:"total_kinetic_mwh"`
	TotalSolarMWh        float64             `json:"total_solar_mwh"`
	TotalPowerGeneratedMW float64            `json:"total_power_generated_mw"`
	TotalCO2OffsetTonnes float64             `json:"total_co2_offset_tonnes"`
	EconomicSavingsINR   float64             `json:"economic_savings_inr"`
	ActiveCorridorsCount int                 `json:"active_corridors_count"`
	AvgResponseTimeMin   float64             `json:"avg_response_time_min"`
	Segments             []RoadSegment       `json:"segments"`
	ActiveEmergencies    []EmergencyVehicle  `json:"active_emergencies"`
	RecentAnomalies      []AnomalyAlert      `json:"recent_anomalies"`
}

// CorridorDispatchRequest holds payload to initiate an emergency green corridor
type CorridorDispatchRequest struct {
	VehicleType     string   `json:"vehicle_type"`
	Priority        int      `json:"priority"`
	StartSegmentID  string   `json:"start_segment_id"`
	TargetHospital  string   `json:"target_hospital"`
	DestinationLoc  Location `json:"destination_loc"`
}

// CorridorDispatchResponse holds response for green corridor activation
type CorridorDispatchResponse struct {
	DispatchID        string    `json:"dispatch_id"`
	VehicleID         string    `json:"vehicle_id"`
	AssignedSegments  []string  `json:"assigned_segments"`
	EstimatedArrival  time.Time `json:"estimated_arrival"`
	TimeSavedSeconds  float64   `json:"time_saved_seconds"`
	SignalPhaseWaveID string    `json:"signal_phase_wave_id"`
	Status            string    `json:"status"`
}
