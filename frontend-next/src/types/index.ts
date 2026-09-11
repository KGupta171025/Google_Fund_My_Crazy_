export interface Location {
  lat: number;
  lon: number;
  alt?: number;
}

export interface RoadSegment {
  id: string;
  name: string;
  start_loc: Location;
  end_loc: Location;
  lanes: number;
  surface_type: string;
  condition_score: number;
  avg_speed_kmh: number;
  vehicle_density: number;
  kinetic_mwh_today: number;
  solar_mwh_today: number;
  heat_island_deg_c: number;
  is_green_corridor: boolean;
  last_maintained: string;
}

export interface EnergyTile {
  id: string;
  segment_id: string;
  location: Location;
  piezo_output_watts: number;
  solar_output_watts: number;
  stress_pressure_kpa: number;
  temperature_deg_c: number;
  status: "OPTIMAL" | "SURGE" | "DEGRADED" | "MAINTENANCE_REQUIRED";
  updated_at: string;
}

export interface EmergencyVehicle {
  id: string;
  callsign: string;
  vehicle_type: string;
  priority: number;
  current_loc: Location;
  destination_loc: Location;
  destination_name: string;
  current_speed_kmh: number;
  status: "DISPATCHED" | "CORRIDOR_ACTIVE" | "ARRIVED" | "IDLE";
  assigned_segments: string[];
  estimated_arrival: string;
  time_saved_seconds: number;
}

export interface AnomalyAlert {
  id: string;
  timestamp: string;
  segment_id: string;
  severity: "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  anomaly_type: string;
  confidence_score: number;
  description: string;
  suggested_action: string;
  gemini_reasoning?: string;
}

export interface UrbanGridState {
  city_id: string;
  city_name: string;
  timestamp: string;
  active_vehicles: number;
  total_kinetic_mwh: number;
  total_solar_mwh: number;
  total_power_generated_mw: number;
  total_co2_offset_tonnes: number;
  economic_savings_inr: number;
  active_corridors_count: number;
  avg_response_time_min: number;
  segments: RoadSegment[];
  active_emergencies: EmergencyVehicle[];
  recent_anomalies: AnomalyAlert[];
}

export interface CorridorDispatchRequest {
  vehicle_type: string;
  priority: number;
  start_segment_id: string;
  target_hospital: string;
  destination_loc?: Location;
}

export interface CorridorDispatchResponse {
  dispatch_id: string;
  vehicle_id: string;
  assigned_segments: string[];
  estimated_arrival: string;
  time_saved_seconds: number;
  signal_phase_wave_id: string;
  status: string;
}

export interface AIPredictionResponse {
  congestion_probability_30m: number;
  pothole_degradation_risk: number;
  optimal_green_corridor_speed_kmh: number;
  structural_status: string;
}

export interface GeminiReasoningResponse {
  assessment_summary: string;
  structural_risk_level: "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  green_corridor_recommendation: string;
  energy_harvesting_optimization: string;
  autonomous_drone_dispatch: boolean;
  engine_mode?: string;
}
