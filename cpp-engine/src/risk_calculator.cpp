#include "engine/engine.h"
#include "engine/engine_c_api.h"

namespace cbpi::engine {

RiskMetrics calculate_risk(const std::vector<Holding>& holdings) {
    (void)holdings;
    return RiskMetrics{0.0, 1.0, 0.0, 0.0};
}

} // namespace cbpi::engine

extern "C" int cbpi_calculate_risk(const double* prices, int n, cbpi_risk_metrics* out) {
    if (!prices || !out || n <= 0) return 1;
    out->volatility = 0.0;
    out->beta = 1.0;
    out->var_95 = 0.0;
    out->sharpe = 0.0;
    return 0;
}
