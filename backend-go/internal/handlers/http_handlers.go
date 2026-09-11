package handlers

import (
	"encoding/json"
	"metro-synapse/backend/internal/models"
	"metro-synapse/backend/internal/simulation"
	"net/http"
	"strconv"
)

// APIHandler wraps HTTP handlers with the simulation engine
type APIHandler struct {
	engine *simulation.CityEngine
}

// NewAPIHandler returns a new API handler
func NewAPIHandler(eng *simulation.CityEngine) *APIHandler {
	return &APIHandler{engine: eng}
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// HandleHealth handles healthcheck endpoint
func (h *APIHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"service": "metro-synapse-go-core",
		"version": "1.0.0",
	})
}

// HandleGetState handles full state query
func (h *APIHandler) HandleGetState(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	state := h.engine.GetState()
	json.NewEncoder(w).Encode(state)
}

// HandleDispatchCorridor handles emergency green corridor requests
func (h *APIHandler) HandleDispatchCorridor(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CorridorDispatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Fallback to default ambulance if empty
		req = models.CorridorDispatchRequest{
			VehicleType:    "AMBULANCE_CRITICAL",
			Priority:       1,
			TargetHospital: "Manipal Cardiac Care Centre",
		}
	}

	resp := h.engine.DispatchGreenCorridor(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleEnergyStats handles energy analytics endpoint
func (h *APIHandler) HandleEnergyStats(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	state := h.engine.GetState()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_kinetic_mwh":        state.TotalKineticMWh,
		"total_solar_mwh":          state.TotalSolarMWh,
		"total_power_generated_mw": state.TotalPowerGeneratedMW,
		"total_co2_offset_tonnes":  state.TotalCO2OffsetTonnes,
		"economic_savings_inr":     state.EconomicSavingsINR,
		"active_vehicles":          state.ActiveVehicles,
	})
}

// HandleSpatialQuery performs spatial radius queries
func (h *APIHandler) HandleSpatialQuery(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	radStr := r.URL.Query().Get("radius")

	lat, _ := strconv.ParseFloat(latStr, 64)
	lon, _ := strconv.ParseFloat(lonStr, 64)
	radius, _ := strconv.ParseFloat(radStr, 64)

	if radius <= 0 {
		radius = 1000.0 // 1km default
	}

	state := h.engine.GetState()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"query_loc":     models.Location{Lat: lat, Lon: lon},
		"radius_meters": radius,
		"segments":      state.Segments,
	})
}
