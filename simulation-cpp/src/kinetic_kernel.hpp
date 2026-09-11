#ifndef METRO_KINETIC_KERNEL_HPP
#define METRO_KINETIC_KERNEL_HPP

#include <vector>
#include <string>
#include <chrono>

namespace metro {

struct TileState {
    int id;
    double x;
    double y;
    double piezo_efficiency; // 0.15 - 0.28 conversion factor
    double accumulated_joules;
    double current_watts;
    double temperature_c;
};

struct VehicleImpulse {
    double mass_kg;
    double velocity_mps;
    double contact_duration_sec;
    double axle_count;
};

class KineticKernel {
public:
    KineticKernel(int tile_count);
    
    // Process vehicle rolling impulses across tiles
    void process_batch(const std::vector<VehicleImpulse>& impulses, double solar_flux_w_m2);
    
    // Total energy generated in Joules and kWh
    double get_total_energy_joules() const;
    double get_total_energy_kwh() const;
    
    // Peak instantaneous power in Watts
    double get_peak_power_watts() const;
    
    // High-performance benchmark
    static void run_benchmark(size_t iterations = 100000);

private:
    std::vector<TileState> tiles_;
    double total_joules_;
    double peak_watts_;
};

} // namespace metro

#endif // METRO_KINETIC_KERNEL_HPP
