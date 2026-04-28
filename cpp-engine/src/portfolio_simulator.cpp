#include "engine/engine.h"

namespace cbpi::engine {

SimulationResult simulate_portfolio(const SimulationRequest& req) {
    SimulationResult r{};
    r.expected_return = 0.0;
    r.std_dev = 0.0;
    r.percentile_paths = {0.0, 0.0, 0.0};
    (void)req;
    return r;
}

} // namespace cbpi::engine
