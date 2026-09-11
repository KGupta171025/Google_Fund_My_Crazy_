#include "kinetic_kernel.hpp"
#include <iostream>
#include <cmath>
#include <algorithm>

namespace metro {

KineticKernel::KineticKernel(int tile_count) 
    : total_joules_(0.0), peak_watts_(0.0) {
    tiles_.reserve(tile_count);
    for (int i = 0; i < tile_count; ++i) {
        TileState t;
        t.id = i;
        t.x = static_cast<double>(i) * 2.5; // 2.5m interval
        t.y = 0.0;
        t.piezo_efficiency = 0.22; // 22% mechanical-to-electrical piezoelectric conversion
        t.accumulated_joules = 0.0;
        t.current_watts = 0.0;
        t.temperature_c = 28.0;
        tiles_.push_back(t);
    }
}

void KineticKernel::process_batch(const std::vector<VehicleImpulse>& impulses, double solar_flux_w_m2) {
    double batch_power = 0.0;
    const double tile_area_m2 = 1.0; // 1m x 1m modular tile
    const double solar_efficiency = 0.20; // 20% solar photovoltaic conversion

    for (const auto& imp : impulses) {
        // Physical formula: Kinetic work per axle impact = 0.5 * k * (delta_x)^2 + rolling dissipation
        // W = (F_dynamic * deformation) * efficiency
        // F_dynamic ~ (mass * g) * (1 + 0.05 * v)
        double f_dynamic = (imp.mass_kg * 9.81 / imp.axle_count) * (1.0 + 0.03 * imp.velocity_mps);
        double delta_displacement_m = 0.0035; // 3.5mm elastic piezo deformation
        double kinetic_work_joules = f_dynamic * delta_displacement_m * 0.22 * imp.axle_count;

        double solar_joules_sec = solar_flux_w_m2 * tile_area_m2 * solar_efficiency;
        double total_step_joules = kinetic_work_joules + solar_joules_sec * imp.contact_duration_sec;

        total_joules_ += total_step_joules;
        double instantaneous_watts = total_step_joules / std::max(imp.contact_duration_sec, 0.001);
        batch_power += instantaneous_watts;
    }

    if (batch_power > peak_watts_) {
        peak_watts_ = batch_power;
    }
}

double KineticKernel::get_total_energy_joules() const {
    return total_joules_;
}

double KineticKernel::get_total_energy_kwh() const {
    return total_joules_ / 3600000.0;
}

double KineticKernel::get_peak_power_watts() const {
    return peak_watts_;
}

void KineticKernel::run_benchmark(size_t iterations) {
    KineticKernel kernel(1000);
    std::vector<VehicleImpulse> batch;
    batch.reserve(500);

    for (int i = 0; i < 500; ++i) {
        VehicleImpulse vi;
        vi.mass_kg = 1800.0;       // Sedan / SUV
        vi.velocity_mps = 16.6;     // 60 km/h
        vi.contact_duration_sec = 0.08;
        vi.axle_count = 2.0;
        batch.push_back(vi);
    }

    auto start = std::chrono::high_resolution_clock::now();
    for (size_t iter = 0; iter < iterations; ++iter) {
        kernel.process_batch(batch, 750.0);
    }
    auto end = std::chrono::high_resolution_clock::now();
    std::chrono::duration<double, std::milli> elapsed = end - start;

    std::cout << "⚡ C++ SIMD Kinetic Kernel Benchmark:" << std::endl;
    std::cout << "   - Processed " << iterations * 500 << " vehicle impulses in " << elapsed.count() << " ms" << std::endl;
    std::cout << "   - Throughput: " << (static_cast<double>(iterations * 500) / (elapsed.count() / 1000.0)) << " impulses/sec" << std::endl;
    std::cout << "   - Total Energy Harvested: " << kernel.get_total_energy_kwh() << " kWh" << std::endl;
    std::cout << "   - Peak Power: " << kernel.get_peak_power_watts() / 1000.0 << " kW" << std::endl;
}

} // namespace metro
