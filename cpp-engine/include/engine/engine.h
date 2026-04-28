// CBPI high-performance engine — public C++ API.
//
// Three pure-compute entry points used by the Go backend via gRPC:
//   - simulate_portfolio: Monte Carlo path simulation
//   - calculate_risk    : volatility / beta / VaR / Sharpe
//   - optimize_allocation: mean-variance optimizer
//
// A C-compatible facade is exposed in engine_c_api.h for FFI consumers.
#ifndef CBPI_ENGINE_H
#define CBPI_ENGINE_H

#include <string>
#include <vector>

namespace cbpi::engine {

struct Holding {
    std::string asset_id;
    double quantity;
    double price;
    std::string currency;
};

struct SimulationRequest {
    std::string portfolio_id;
    std::vector<Holding> holdings;
    int horizon_days = 30;
    int monte_carlo_runs = 10000;
};

struct SimulationResult {
    double expected_return;
    double std_dev;
    std::vector<double> percentile_paths; // p5, p50, p95
};

struct RiskMetrics {
    double volatility;
    double beta;
    double var_95;
    double sharpe;
};

struct AllocationSuggestion {
    std::string asset_id;
    double current_weight;
    double target_weight;
    std::string rationale;
};

SimulationResult simulate_portfolio(const SimulationRequest& req);
RiskMetrics      calculate_risk(const std::vector<Holding>& holdings);
std::vector<AllocationSuggestion> optimize_allocation(
    const std::vector<Holding>& holdings, double target_return);

} // namespace cbpi::engine

#endif // CBPI_ENGINE_H
