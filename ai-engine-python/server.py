"""
Metro-Synapse Python AI Microservice.
Built entirely with standard Python libraries and zero-dependency scratch neural engines.
"""

import json
import os
import sys
from http.server import BaseHTTPRequestHandler, HTTPServer
from typing import Any, Dict

from neural_scratch.traffic_predictor import get_trained_urban_model
from cellular_automata.urban_dispersion import NagelSchreckenbergCA, UrbanDispersion2D
from gemini_agent.sentinel import GeminiUrbanSentinel


# Initialize AI models
print("🧠 Initializing Metro-Synapse Custom Neural Models from scratch...")
urban_model = get_trained_urban_model()
gemini_sentinel = GeminiUrbanSentinel()
ca_sim = NagelSchreckenbergCA(track_length=40, v_max=5)
ca_sim.seed_vehicles(density=0.35)
dispersion_sim = UrbanDispersion2D(grid_size=20)
dispersion_sim.add_emission_source(10, 10, intensity=80.0)
dispersion_sim.set_green_corridor(10)
print("✅ Custom Neural & Cellular Automata Models Ready.")


class AIRequestHandler(BaseHTTPRequestHandler):
    def _set_cors_headers(self):
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type, Authorization")

    def do_OPTIONS(self):
        self.send_response(200)
        self._set_cors_headers()
        self.end_headers()

    def do_GET(self):
        if self.path == "/health":
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self._set_cors_headers()
            self.end_headers()
            resp = {
                "status": "healthy",
                "service": "metro-synapse-python-ai",
                "engine": "custom-scratch-neural-net",
                "version": "1.0.0"
            }
            self.wfile.write(json.dumps(resp).encode("utf-8"))
        else:
            self.send_error(404, "Endpoint not found")

    def do_POST(self):
        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length)
        data = json.loads(body.decode("utf-8")) if body else {}

        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self._set_cors_headers()
        self.end_headers()

        response_data: Dict[str, Any] = {}

        if self.path == "/predict/segment":
            v_dens = float(data.get("vehicle_density", 35.0))
            v_spd = float(data.get("avg_speed_kmh", 45.0))
            press = float(data.get("pressure_kpa", 120.0))
            temp = float(data.get("surface_temp", 28.0))
            rain = float(data.get("rainfall_mm", 0.0))
            hour = float(data.get("hour_of_day", 14.0))

            prediction = urban_model.predict_segment(v_dens, v_spd, press, temp, rain, hour)
            response_data = prediction

        elif self.path == "/simulate/cellular":
            # Step cellular automata and pollutant dispersion
            track_state = ca_sim.step()
            dispersion_grid = dispersion_sim.step_diffusion()
            stats = dispersion_sim.get_stats()
            response_data = {
                "track_state": track_state,
                "dispersion_grid": dispersion_grid,
                "air_quality_stats": stats
            }

        elif self.path == "/gemini/reason":
            # Gemini Multimodal Sentinel reasoning
            query = data.get("query", "")
            segment_data = data.get("segment_data", {
                "vehicle_density": 45.0,
                "avg_speed_kmh": 32.0,
                "pothole_degradation_risk": 0.38,
                "is_green_corridor": False
            })
            drone_image = data.get("drone_image_base64", None)
            reasoning = gemini_sentinel.analyze_urban_anomaly(segment_data, drone_image)
            response_data = reasoning

        else:
            response_data = {"error": "Invalid endpoint"}

        self.wfile.write(json.dumps(response_data).encode("utf-8"))

    def log_message(self, format, *args):
        # Concise logging
        sys.stderr.write(f"[AI-Engine] {self.address_string()} - {format % args}\n")


def run(port: int = 8000):
    server_address = ("", port)
    httpd = HTTPServer(server_address, AIRequestHandler)
    print(f"🚀 Metro-Synapse Python AI Microservice running on http://localhost:{port}")
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\n🛑 Shutting down Python AI Microservice...")
        httpd.server_close()


if __name__ == "__main__":
    port = int(os.environ.get("PORT", 8000))
    run(port)
