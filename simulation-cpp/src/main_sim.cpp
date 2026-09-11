#include "kinetic_kernel.hpp"
#include "green_corridor.hpp"
#include <iostream>
#include <iomanip>

int main() {
    std::cout << "==================================================================" << std::endl;
    std::cout << "⚡ METRO-SYNAPSE HIGH-PERFORMANCE C++ PHYSICAL SIMULATION KERNEL" << std::endl;
    std::cout << "🚀 Google Gemini 'Fund My Crazy' Moonshot Engine" << std::endl;
    std::cout << "==================================================================" << std::endl;

    // 1. Run SIMD Kinetic Kernel Benchmark
    metro::KineticKernel::run_benchmark(50000);

    std::cout << "\n------------------------------------------------------------------" << std::endl;
    std::cout << "🚑 Testing Emergency Green-Corridor Wave-Clearing Optimization:" << std::endl;

    metro::GreenCorridorOptimizer opt(12);
    // Add realistic urban road network edges
    opt.add_edge(0, 1, 1200.0, 0.45); // MG Road -> Trinity
    opt.add_edge(1, 2, 850.0, 0.60);  // Trinity -> Halasuru
    opt.add_edge(2, 3, 1400.0, 0.70); // Halasuru -> Indiranagar
    opt.add_edge(3, 4, 1800.0, 0.50); // Indiranagar -> Domlur Flyover
    opt.add_edge(4, 5, 2100.0, 0.85); // Domlur -> Manipal Hospital (Cardiac Centre)
    opt.add_edge(0, 6, 2500.0, 0.90); // Alternative congested route
    opt.add_edge(6, 5, 3000.0, 0.75);

    auto plan = opt.compute_optimal_corridor(0, 5, 75.0); // 75 km/h emergency speed

    std::cout << "   - Route Nodes Cleared: ";
    for (size_t i = 0; i < plan.node_sequence.size(); ++i) {
        std::cout << "Junction-" << plan.node_sequence[i] + 1;
        if (i + 1 < plan.node_sequence.size()) std::cout << " -> ";
    }
    std::cout << std::endl;
    std::cout << "   - Total Distance: " << plan.total_distance_m / 1000.0 << " km" << std::endl;
    std::cout << "   - Standard Gridlock ETA: " << plan.standard_eta_sec / 60.0 << " mins" << std::endl;
    std::cout << "   - Metro-Synapse Green Corridor ETA: " << plan.green_corridor_eta_sec / 60.0 << " mins" << std::endl;
    std::cout << "   - Time Saved: " << std::fixed << std::setprecision(1) << plan.time_saved_sec / 60.0 << " minutes (" << plan.time_saved_sec << " seconds)" << std::endl;
    std::cout << "   - Patient Trauma Survival Gain: +" << plan.survival_probability_gain_percent << "%" << std::endl;
    std::cout << "==================================================================" << std::endl;

    return 0;
}
