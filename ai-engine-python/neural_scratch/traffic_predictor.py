"""
Urban Neural Network trained from scratch for congestion forecasting and structural degradation detection.
"""

import json
import math
import random
from typing import Dict, List, Tuple
from .layers import Dense, Dropout, Layer, ReLU, Sigmoid
from .matrix_ops import Matrix, mse_loss, shape
from .optimizers import Adam


class UrbanNeuralModel:
    """End-to-End deep neural network built from scratch in pure Python"""
    def __init__(self):
        # Input features (6): [vehicle_density, avg_speed, pressure_kpa, surface_temp, rainfall_mm, hour_norm]
        # Output predictions (3): [congestion_prob_30m, pothole_risk_score, optimal_corridor_speed_norm]
        self.layers: List[Layer] = [
            Dense(6, 32, init_scheme="he"),
            ReLU(),
            Dropout(0.1),
            Dense(32, 16, init_scheme="he"),
            ReLU(),
            Dense(16, 3, init_scheme="xavier"),
            Sigmoid()
        ]
        self.optimizer = Adam(lr=0.005, weight_decay=1e-4)

    def forward(self, x: Matrix, training: bool = False) -> Matrix:
        out = x
        for layer in self.layers:
            out = layer.forward(out, training=training)
        return out

    def backward(self, grad_loss: Matrix):
        grad = grad_loss
        for layer in reversed(self.layers):
            grad = layer.backward(grad)

    def train_step(self, x: Matrix, y_true: Matrix) -> float:
        # Forward pass
        y_pred = self.forward(x, training=True)
        # Compute Loss and analytical gradient
        loss, grad = mse_loss(y_pred, y_true)
        # Backward pass
        self.backward(grad)
        # Gather parameters and apply Adam update
        param_groups = []
        for layer in self.layers:
            param_groups.extend(layer.get_params_and_grads())
        self.optimizer.step(param_groups)
        return loss

    def predict_segment(self, vehicle_density: float, avg_speed: float, pressure_kpa: float,
                        surface_temp: float, rainfall_mm: float, hour_of_day: float) -> Dict[str, float]:
        """Perform normalized forward inference on a single urban segment"""
        # Feature normalization
        x = [[
            min(1.0, max(0.0, vehicle_density / 100.0)),
            min(1.0, max(0.0, avg_speed / 120.0)),
            min(1.0, max(0.0, pressure_kpa / 300.0)),
            min(1.0, max(0.0, (surface_temp - 10.0) / 45.0)),
            min(1.0, max(0.0, rainfall_mm / 100.0)),
            min(1.0, max(0.0, hour_of_day / 24.0))
        ]]

        out = self.forward(x, training=False)[0]

        congestion_prob = round(out[0], 4)
        pothole_risk = round(out[1], 4)
        optimal_speed_kmh = round(40.0 + out[2] * 50.0, 1) # Range 40 - 90 km/h

        return {
            "congestion_probability_30m": congestion_prob,
            "pothole_degradation_risk": pothole_risk,
            "optimal_green_corridor_speed_kmh": optimal_speed_kmh,
            "structural_status": "CRITICAL" if pothole_risk > 0.75 else "WARNING" if pothole_risk > 0.45 else "OPTIMAL"
        }

    def train_synthetic(self, num_samples: int = 500, epochs: int = 40) -> List[float]:
        """Train model on synthetic physical urban dynamics"""
        losses = []
        
        # Generate synthetic training pairs based on ground truth physical laws
        X_train: Matrix = []
        Y_train: Matrix = []

        for _ in range(num_samples):
            v_dens = random.uniform(5.0, 95.0)
            v_spd = random.uniform(10.0, 100.0)
            press = random.uniform(50.0, 280.0)
            temp = random.uniform(15.0, 48.0)
            rain = random.uniform(0.0, 60.0)
            hour = random.uniform(0.0, 24.0)

            # Ground truth targets with non-linear physics
            # High density + low speed + peak hour (8-10 or 17-20) -> High congestion
            peak_factor = 1.0 if (8 <= hour <= 10 or 17 <= hour <= 20) else 0.3
            c_prob = min(0.99, max(0.01, (v_dens / 100.0) * 0.6 + (1.0 - v_spd / 120.0) * 0.3 + peak_factor * 0.2))

            # High pressure + high rain + extreme temp -> High pothole degradation
            p_risk = min(0.99, max(0.01, (press / 300.0) * 0.45 + (rain / 100.0) * 0.40 + (temp / 50.0) * 0.15))

            # Optimal speed wave
            opt_speed_norm = max(0.1, 1.0 - c_prob * 0.8)

            X_train.append([v_dens / 100.0, v_spd / 120.0, press / 300.0, (temp - 10.0) / 45.0, rain / 100.0, hour / 24.0])
            Y_train.append([c_prob, p_risk, opt_speed_norm])

        # Batch training
        batch_size = 32
        for epoch in range(epochs):
            epoch_loss = 0.0
            batches = len(X_train) // batch_size
            for b in range(batches):
                start = b * batch_size
                end = start + batch_size
                batch_x = X_train[start:end]
                batch_y = Y_train[start:end]
                loss = self.train_step(batch_x, batch_y)
                epoch_loss += loss
            avg_loss = epoch_loss / max(1, batches)
            losses.append(avg_loss)

        return losses


# Singleton model instance
_urban_model: UrbanNeuralModel = None

def get_trained_urban_model() -> UrbanNeuralModel:
    global _urban_model
    if _urban_model is None:
        _urban_model = UrbanNeuralModel()
        _urban_model.train_synthetic(num_samples=400, epochs=30)
    return _urban_model
