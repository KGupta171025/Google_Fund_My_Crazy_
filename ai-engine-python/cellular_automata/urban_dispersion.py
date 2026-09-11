"""
Urban Cellular Automata & Fluid Pollutant Dispersion Engine.
Implements the Nagel-Schreckenberg traffic model and 2D Fickian finite-difference gas dispersion.
"""

import math
import random
from typing import Dict, List, Tuple


class NagelSchreckenbergCA:
    """1D/2D Vehicular Cellular Automata"""
    def __init__(self, track_length: int = 50, v_max: int = 5, p_slowdown: float = 0.25):
        self.length = track_length
        self.v_max = v_max
        self.p_slowdown = p_slowdown
        # Grid: -1 indicates empty cell, >=0 indicates vehicle with speed v
        self.grid = [-1] * track_length

    def seed_vehicles(self, density: float = 0.3):
        for i in range(self.length):
            if random.random() < density:
                self.grid[i] = random.randint(0, self.v_max)
            else:
                self.grid[i] = -1

    def step(self) -> List[int]:
        new_grid = [-1] * self.length
        
        # Rule 1: Acceleration (v -> min(v+1, v_max))
        # Rule 2: Deceleration (distance to next car)
        # Rule 3: Randomization (with probability p, v -> max(v-1, 0))
        # Rule 4: Movement (car moves v cells forward)
        
        for i in range(self.length):
            v = self.grid[i]
            if v == -1:
                continue

            # 1. Accelerate
            v = min(v + 1, self.v_max)

            # 2. Check distance to next car ahead
            dist = 1
            while dist <= v:
                idx = (i + dist) % self.length
                if self.grid[idx] != -1:
                    break
                dist += 1
            v = min(v, dist - 1)

            # 3. Random slowdown (driver fluctuation)
            if v > 0 and random.random() < self.p_slowdown:
                v -= 1

            # 4. Move
            new_pos = (i + v) % self.length
            new_grid[new_pos] = v

        self.grid = new_grid
        return self.grid


class UrbanDispersion2D:
    """2D Finite-Difference Grid simulating PM2.5 and thermal heat dissipation"""
    def __init__(self, grid_size: int = 20, diffusion_coeff: float = 0.18):
        self.size = grid_size
        self.diff = diffusion_coeff
        # 2D concentration matrix
        self.grid = [[0.0 for _ in range(grid_size)] for _ in range(grid_size)]
        # Heat island reduction factor from smart kinetic tiles (green corridors)
        self.smart_corridors = set()

    def add_emission_source(self, x: int, y: int, intensity: float = 50.0):
        if 0 <= x < self.size and 0 <= y < self.size:
            self.grid[x][y] += intensity

    def set_green_corridor(self, corridor_x: int):
        for y in range(self.size):
            self.smart_corridors.add((corridor_x, y))

    def step_diffusion(self, decay: float = 0.05) -> List[List[float]]:
        """Compute one time step of 2D diffusion: C_new = C + D * laplacian(C) - decay"""
        new_grid = [[0.0 for _ in range(self.size)] for _ in range(self.size)]

        for x in range(self.size):
            for y in range(self.size):
                # 5-point discrete Laplacian stencil
                center = self.grid[x][y]
                left = self.grid[x - 1][y] if x > 0 else center
                right = self.grid[x + 1][y] if x < self.size - 1 else center
                top = self.grid[x][y - 1] if y > 0 else center
                bottom = self.grid[x][y + 1] if y < self.size - 1 else center

                laplacian = (left + right + top + bottom) - 4.0 * center
                absorption = 0.25 if (x, y) in self.smart_corridors else 0.0

                new_val = center + self.diff * laplacian - (decay + absorption) * center
                new_grid[x][y] = max(0.0, round(new_val, 3))

        self.grid = new_grid
        return self.grid

    def get_stats(self) -> Dict[str, float]:
        total_conc = sum(sum(row) for row in self.grid)
        max_conc = max(max(row) for row in self.grid)
        avg_aqi = round(min(500.0, total_conc / (self.size * self.size) * 12.0), 1)
        return {
            "total_particulate_mass": round(total_conc, 2),
            "peak_pollution_index": round(max_conc, 2),
            "estimated_city_aqi": avg_aqi
        }
