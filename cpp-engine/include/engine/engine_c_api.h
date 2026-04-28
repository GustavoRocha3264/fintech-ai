// C-compatible facade over the engine, suitable for cgo / ctypes / JNI.
// Kept intentionally minimal — full gRPC service is the primary integration.
#ifndef CBPI_ENGINE_C_API_H
#define CBPI_ENGINE_C_API_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct cbpi_risk_metrics {
    double volatility;
    double beta;
    double var_95;
    double sharpe;
} cbpi_risk_metrics;

// Returns 0 on success, non-zero on error. Inputs are flat C arrays so the
// header stays ABI-stable.
int cbpi_calculate_risk(const double* prices, int n, cbpi_risk_metrics* out);

#ifdef __cplusplus
} // extern "C"
#endif

#endif // CBPI_ENGINE_C_API_H
