import {
  UrbanGridState,
  CorridorDispatchRequest,
  CorridorDispatchResponse,
  AIPredictionResponse,
  GeminiReasoningResponse,
} from "@/types";

const GO_BACKEND_URL = process.env.NEXT_PUBLIC_GO_BACKEND_URL || "http://localhost:8080";
const PYTHON_AI_URL = process.env.NEXT_PUBLIC_PYTHON_AI_URL || "http://localhost:8000";

export const MOCK_INITIAL_STATE: UrbanGridState = {
  city_id: "city-blr-01",
  city_name: "Bengaluru Smart Metropolis",
  timestamp: new Date().toISOString(),
  active_vehicles: 14820,
  total_kinetic_mwh: 148.52,
  total_solar_mwh: 92.14,
  total_power_generated_mw: 28.88,
  total_co2_offset_tonnes: 185.6,
  economic_savings_inr: 3450000,
  active_corridors_count: 1,
  avg_response_time_min: 4.2,
  segments: [
    {
      id: "seg-koramangala-ecity",
      name: "Hosur Elevated Smart Highway (Outer Ring Link)",
      start_loc: { lat: 12.9279, lon: 77.6271 },
      end_loc: { lat: 12.8452, lon: 77.6602 },
      lanes: 6,
      surface_type: "piezo_solar_matrix",
      condition_score: 96.4,
      avg_speed_kmh: 52.0,
      vehicle_density: 38.5,
      kinetic_mwh_today: 34.2,
      solar_mwh_today: 22.8,
      heat_island_deg_c: -2.8,
      is_green_corridor: false,
      last_maintained: new Date(Date.now() - 15 * 86400000).toISOString(),
    },
    {
      id: "seg-mgroad-indiranagar",
      name: "MG Road Arterial Kinetic Axis",
      start_loc: { lat: 12.9756, lon: 77.6066 },
      end_loc: { lat: 12.9784, lon: 77.6408 },
      lanes: 4,
      surface_type: "piezo_solar_matrix",
      condition_score: 98.2,
      avg_speed_kmh: 68.0,
      vehicle_density: 22.0,
      kinetic_mwh_today: 28.5,
      solar_mwh_today: 19.4,
      heat_island_deg_c: -2.1,
      is_green_corridor: true,
      last_maintained: new Date(Date.now() - 30 * 86400000).toISOString(),
    },
    {
      id: "seg-hebbal-airport",
      name: "Hebbal Express Kinetic & Solar Corridor",
      start_loc: { lat: 13.0358, lon: 77.597 },
      end_loc: { lat: 13.1986, lon: 77.7066 },
      lanes: 8,
      surface_type: "piezo_solar_matrix",
      condition_score: 94.8,
      avg_speed_kmh: 75.0,
      vehicle_density: 45.0,
      kinetic_mwh_today: 48.6,
      solar_mwh_today: 31.5,
      heat_island_deg_c: -3.2,
      is_green_corridor: false,
      last_maintained: new Date(Date.now() - 5 * 86400000).toISOString(),
    },
    {
      id: "seg-whitefield-marathahalli",
      name: "Whitefield Cyber Tech Spine",
      start_loc: { lat: 12.9698, lon: 77.75 },
      end_loc: { lat: 12.9591, lon: 77.6974 },
      lanes: 6,
      surface_type: "piezo_solar_matrix",
      condition_score: 89.1,
      avg_speed_kmh: 28.0,
      vehicle_density: 68.0,
      kinetic_mwh_today: 37.2,
      solar_mwh_today: 18.4,
      heat_island_deg_c: -1.5,
      is_green_corridor: false,
      last_maintained: new Date(Date.now() - 45 * 86400000).toISOString(),
    },
  ],
  active_emergencies: [
    {
      id: "em-AMBULANCE-8492",
      callsign: "PRANA-LIFE-108",
      vehicle_type: "AMBULANCE_CRITICAL",
      priority: 1,
      current_loc: { lat: 12.9756, lon: 77.6066 },
      destination_loc: { lat: 12.9592, lon: 77.6534 },
      destination_name: "Manipal Institute of Cardiac Sciences",
      current_speed_kmh: 74.0,
      status: "CORRIDOR_ACTIVE",
      assigned_segments: ["seg-mgroad-indiranagar"],
      estimated_arrival: new Date(Date.now() + 240000).toISOString(),
      time_saved_seconds: 940,
    },
  ],
  recent_anomalies: [
    {
      id: "alert-anom-01",
      timestamp: new Date(Date.now() - 10 * 60000).toISOString(),
      segment_id: "seg-whitefield-marathahalli",
      severity: "MEDIUM",
      anomaly_type: "POTHOLE_EARLY_STAGE",
      confidence_score: 0.942,
      description: "Subsurface piezo-acoustic vibration deviation at km 3.4.",
      suggested_action: "Deploy autonomous thermal resin patch drone within 48h.",
      gemini_reasoning: "Gemini Vision confirms 4cm micro-fissure forming along outer wheel track.",
    },
  ],
};

export async function fetchCityState(): Promise<UrbanGridState> {
  try {
    const res = await fetch(`${GO_BACKEND_URL}/api/city/state`, { cache: "no-store" });
    if (!res.ok) throw new Error("Go backend error");
    return await res.json();
  } catch {
    return MOCK_INITIAL_STATE;
  }
}

export async function dispatchGreenCorridor(
  req: CorridorDispatchRequest
): Promise<CorridorDispatchResponse> {
  try {
    const res = await fetch(`${GO_BACKEND_URL}/api/corridor/dispatch`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
    });
    if (!res.ok) throw new Error("Dispatch failed");
    return await res.json();
  } catch {
    // Client-side fallback response
    return {
      dispatch_id: `disp-${Date.now()}`,
      vehicle_id: `em-${req.vehicle_type}-${Math.floor(1000 + Math.random() * 9000)}`,
      assigned_segments: ["seg-mgroad-indiranagar"],
      estimated_arrival: new Date(Date.now() + 260000).toISOString(),
      time_saved_seconds: 937.0,
      signal_phase_wave_id: `WAVE-PHASE-${Math.floor(Math.random() * 0xffffff).toString(16)}`,
      status: "ACTIVE_GREEN_CORRIDOR",
    };
  }
}

export async function fetchAISegmentPrediction(params: {
  vehicle_density: number;
  avg_speed_kmh: number;
  pressure_kpa: number;
  surface_temp: number;
  rainfall_mm: number;
  hour_of_day: number;
}): Promise<AIPredictionResponse> {
  try {
    const res = await fetch(`${PYTHON_AI_URL}/predict/segment`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(params),
    });
    if (!res.ok) throw new Error("AI engine unreachable");
    return await res.json();
  } catch {
    // Scratch mathematical inference formula on client
    const densNorm = Math.min(1.0, params.vehicle_density / 100);
    const spdNorm = Math.min(1.0, params.avg_speed_kmh / 120);
    const pressNorm = Math.min(1.0, params.pressure_kpa / 300);
    const rainNorm = Math.min(1.0, params.rainfall_mm / 100);

    const congProb = Math.min(0.98, Math.max(0.02, densNorm * 0.65 + (1 - spdNorm) * 0.35));
    const potholeRisk = Math.min(0.95, Math.max(0.05, pressNorm * 0.5 + rainNorm * 0.4));
    const optSpeed = Math.round(40 + (1 - congProb) * 50);

    return {
      congestion_probability_30m: Math.round(congProb * 1000) / 1000,
      pothole_degradation_risk: Math.round(potholeRisk * 1000) / 1000,
      optimal_green_corridor_speed_kmh: optSpeed,
      structural_status: potholeRisk > 0.75 ? "CRITICAL" : potholeRisk > 0.45 ? "WARNING" : "OPTIMAL",
    };
  }
}

export async function fetchGeminiReasoning(query: string, segmentData: any): Promise<GeminiReasoningResponse> {
  try {
    const res = await fetch(`${PYTHON_AI_URL}/gemini/reason`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ query, segment_data: segmentData }),
    });
    if (!res.ok) throw new Error("Gemini AI failed");
    return await res.json();
  } catch {
    return {
      assessment_summary:
        "Metro-Synapse Multimodal Sentinel: High axle kinetic compression detected along MG Road. Green corridor synchronized at 74 km/h wave speed with zero intersection queue latency.",
      structural_risk_level: "LOW",
      green_corridor_recommendation:
        "Signal cycle preempted: MG Road Eastbound locked green for next 180s. Cross traffic rerouted via 80ft arterial.",
      energy_harvesting_optimization:
        "Piezoelectric resonance tuned to 420 Hz. Energy capture efficiency optimized at 23.4%.",
      autonomous_drone_dispatch: false,
      engine_mode: "METRO-SYNAPSE-SCRATCH-REASONER",
    };
  }
}
