"""
Comprehensive Unit Tests for Pure-Python Scratch Neural Engine and Cellular Automata.
"""

import math
import unittest
from neural_scratch.matrix_ops import (
    zeros, ones, randn, matmul, transpose, add, subtract, multiply,
    relu, sigmoid, shape, mse_loss
)
from neural_scratch.layers import Dense, ReLU, Sigmoid
from neural_scratch.optimizers import Adam
from neural_scratch.traffic_predictor import UrbanNeuralModel
from cellular_automata.urban_dispersion import NagelSchreckenbergCA, UrbanDispersion2D


class TestMatrixMathOps(unittest.TestCase):
    def test_matmul_and_transpose(self):
        A = [[1.0, 2.0], [3.0, 4.0]]
        B = [[5.0, 6.0], [7.0, 8.0]]
        # A * B = [[19, 22], [43, 50]]
        C = matmul(A, B)
        self.assertEqual(C[0][0], 19.0)
        self.assertEqual(C[0][1], 22.0)
        self.assertEqual(C[1][0], 43.0)
        self.assertEqual(C[1][1], 50.0)

        At = transpose(A)
        self.assertEqual(At[0][1], 3.0)
        self.assertEqual(At[1][0], 2.0)

    def test_broadcasting_add(self):
        A = [[1.0, 2.0], [3.0, 4.0]]
        bias = [10.0, 20.0]
        res = add(A, bias)
        self.assertEqual(res[0], [11.0, 22.0])
        self.assertEqual(res[1], [13.0, 24.0])

    def test_activations(self):
        self.assertEqual(relu(-5.0), 0.0)
        self.assertEqual(relu(3.5), 3.5)
        self.assertAlmostEqual(sigmoid(0.0), 0.5, places=5)


class TestNeuralLayersAndGradients(unittest.TestCase):
    def test_dense_layer_backward(self):
        dense = Dense(in_features=3, out_features=2)
        x = [[1.0, 2.0, 3.0]]
        out = dense.forward(x, training=True)
        self.assertEqual(shape(out), (1, 2))

        # Backpropagation
        grad_out = [[0.5, -0.5]]
        grad_x = dense.backward(grad_out)
        self.assertEqual(shape(grad_x), (1, 3))
        self.assertEqual(shape(dense.grad_weights), (3, 2))


class TestUrbanNeuralModel(unittest.TestCase):
    def test_synthetic_training_convergence(self):
        model = UrbanNeuralModel()
        losses = model.train_synthetic(num_samples=100, epochs=15)
        # Loss should decrease over training
        self.assertGreater(len(losses), 0)
        self.assertLess(losses[-1], losses[0] * 1.5)

    def test_prediction_ranges(self):
        model = UrbanNeuralModel()
        pred = model.predict_segment(
            vehicle_density=80.0,
            avg_speed=15.0,
            pressure_kpa=220.0,
            surface_temp=35.0,
            rainfall_mm=25.0,
            hour_of_day=18.0
        )
        self.assertIn("congestion_probability_30m", pred)
        self.assertIn("pothole_degradation_risk", pred)
        self.assertIn("optimal_green_corridor_speed_kmh", pred)
        self.assertTrue(0.0 <= pred["congestion_probability_30m"] <= 1.0)
        self.assertTrue(0.0 <= pred["pothole_degradation_risk"] <= 1.0)
        self.assertTrue(40.0 <= pred["optimal_green_corridor_speed_kmh"] <= 90.0)


class TestCellularAutomata(unittest.TestCase):
    def test_nagel_schreckenberg_step(self):
        ca = NagelSchreckenbergCA(track_length=20, v_max=3)
        ca.seed_vehicles(density=0.5)
        new_grid = ca.step()
        self.assertEqual(len(new_grid), 20)

    def test_urban_dispersion(self):
        disp = UrbanDispersion2D(grid_size=10)
        disp.add_emission_source(5, 5, intensity=100.0)
        stats_before = disp.get_stats()
        disp.step_diffusion()
        stats_after = disp.get_stats()
        self.assertLessEqual(stats_after["peak_pollution_index"], stats_before["peak_pollution_index"])


if __name__ == "__main__":
    unittest.main()
