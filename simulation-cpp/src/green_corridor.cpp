#include "green_corridor.hpp"
#include <queue>
#include <cmath>
#include <limits>
#include <algorithm>

namespace metro {

GreenCorridorOptimizer::GreenCorridorOptimizer(int node_count)
    : node_count_(node_count), adj_(node_count) {
    nodes_.resize(node_count);
    for (int i = 0; i < node_count; ++i) {
        nodes_[i].id = i;
        nodes_[i].name = "Junction-" + std::to_string(i + 1);
        nodes_[i].current_signal_state = 0; // RED
        nodes_[i].queue_length_vehicles = 12 + (i % 10) * 3;
        nodes_[i].clearance_time_sec = static_cast<double>(nodes_[i].queue_length_vehicles) * 1.8;
    }
}

void GreenCorridorOptimizer::add_edge(int from_node, int to_node, double length_meters, double initial_traffic_density) {
    if (from_node >= 0 && from_node < node_count_ && to_node >= 0 && to_node < node_count_) {
        adj_[from_node].push_back({to_node, length_meters, initial_traffic_density});
        adj_[to_node].push_back({from_node, length_meters, initial_traffic_density}); // Bidirectional
    }
}

CorridorPlan GreenCorridorOptimizer::compute_optimal_corridor(int start_node, int target_hospital_node, double ambulance_speed_kmh) {
    // Dijkstra algorithm with dynamic traffic resistance weight
    std::vector<double> dist(node_count_, std::numeric_limits<double>::infinity());
    std::vector<int> prev(node_count_, -1);
    
    // Min-priority queue: pair<cost, node>
    typedef std::pair<double, int> PQElement;
    std::priority_queue<PQElement, std::vector<PQElement>, std::greater<PQElement>> pq;

    dist[start_node] = 0.0;
    pq.push({0.0, start_node});

    while (!pq.empty()) {
        PQElement topElem = pq.top();
        pq.pop();
        double current_cost = topElem.first;
        int u = topElem.second;

        if (current_cost > dist[u]) continue;
        if (u == target_hospital_node) break;

        for (const auto& edge : adj_[u]) {
            int v = edge.to;
            // Cost = Physical Distance * (1 + 0.8 * Density)
            double edge_cost = edge.length * (1.0 + 0.8 * edge.density);
            if (dist[u] + edge_cost < dist[v]) {
                dist[v] = dist[u] + edge_cost;
                prev[v] = u;
                pq.push({dist[v], v});
            }
        }
    }

    // Reconstruct path
    std::vector<int> path;
    int curr = target_hospital_node;
    while (curr != -1) {
        path.push_back(curr);
        curr = prev[curr];
    }
    std::reverse(path.begin(), path.end());

    double total_length = 0.0;
    for (size_t i = 0; i + 1 < path.size(); ++i) {
        int u = path[i];
        int v = path[i + 1];
        for (const auto& e : adj_[u]) {
            if (e.to == v) {
                total_length += e.length;
                break;
            }
        }
    }

    double amb_mps = ambulance_speed_kmh * (1000.0 / 3600.0);
    double gridlock_mps = 18.0 * (1000.0 / 3600.0); // 18 km/h standard traffic

    double standard_eta = (total_length / gridlock_mps) + path.size() * 35.0; // signal delays
    double green_eta = (total_length / amb_mps) + 4.0; // zero signal stop, pure wave clearance
    double time_saved = std::max(0.0, standard_eta - green_eta);

    // Medical trauma literature: each 1 min saved in Golden Hour increases survival chance by ~6.8%
    double survival_gain = (time_saved / 60.0) * 6.8;

    CorridorPlan plan;
    plan.node_sequence = path;
    plan.total_distance_m = total_length;
    plan.standard_eta_sec = standard_eta;
    plan.green_corridor_eta_sec = green_eta;
    plan.time_saved_sec = time_saved;
    plan.survival_probability_gain_percent = std::min(45.0, survival_gain);

    return plan;
}

} // namespace metro
