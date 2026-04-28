#include <gtest/gtest.h>
#include "engine/engine.h"

using namespace cbpi::engine;

TEST(RiskCalculator, ReturnsDefaultsForEmptyHoldings) {
    auto r = calculate_risk({});
    EXPECT_DOUBLE_EQ(r.beta, 1.0);
    EXPECT_DOUBLE_EQ(r.volatility, 0.0);
}

TEST(PortfolioSimulator, ReturnsThreePercentiles) {
    SimulationRequest req{};
    req.portfolio_id = "p1";
    req.horizon_days = 10;
    req.monte_carlo_runs = 100;
    auto r = simulate_portfolio(req);
    EXPECT_EQ(r.percentile_paths.size(), 3u);
}
