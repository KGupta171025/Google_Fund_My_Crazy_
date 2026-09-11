#ifndef METRO_GREEN_CORRIDOR_HPP
#define METRO_GREEN_CORRIDOR_HPP

#include <vector>
#include <string>

namespace metro {

struct IntersectionNode {
    int id;
    std::string name;
    double x;
    double y;
    int current_signal_state; // 0: RED, 1: YELLOW, 2: GREEN
    int queue_length_vehicles;
    double clearance_time_sec;
};

struct CorridorPlan {
    std::vector<int> node_sequence;
    double total_distance_m;
    double standard_eta_sec;
    double green_corridor_eta_sec;
    double time_saved_sec;
    double survival_probability_gain_percent;
};

class GreenCorridorOptimizer {
public:
    GreenCorridorOptimizer(int node_count);
    
    void add_edge(int from_node, int to_node, double length_meters, double initial_traffic_density);
    
    CorridorPlan compute_optimal_corridor(int start_node, int target_hospital_node, double ambulance_speed_kmh);

private:
    struct Edge {
        int to;
        double length;
        double density;
    };
    
    int node_count_;
    std::vector<IntersectionNode> nodes_;
    std::vector<std::vector<Edge>> adj_;
};

} // namespace metro

#endif // METRO_GREEN_CORRIDOR_HPP
