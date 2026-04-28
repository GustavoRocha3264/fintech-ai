// Entry point for the engine binary. The gRPC server (using stubs generated
// from /proto/engine.proto) will be wired up here. For the scaffold we just
// keep a no-op main so the executable target builds clean.
#include <iostream>

int main() {
    std::cout << "cbpi engine — gRPC server scaffold (not yet wired)\n";
    return 0;
}
