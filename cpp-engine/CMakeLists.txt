cmake_minimum_required(VERSION 3.20)
project(cbpi_engine LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_POSITION_INDEPENDENT_CODE ON)

if(NOT CMAKE_BUILD_TYPE)
  set(CMAKE_BUILD_TYPE Release)
endif()

add_compile_options(-Wall -Wextra -Wpedantic)

# --- core engine library ----------------------------------------------------
add_library(engine_core
  src/portfolio_simulator.cpp
  src/risk_calculator.cpp
  src/allocation_optimizer.cpp
)
target_include_directories(engine_core PUBLIC include)

# --- engine binary (gRPC server entry point — wired in follow-up) -----------
add_executable(engine_server src/main.cpp)
target_link_libraries(engine_server PRIVATE engine_core)

# --- tests ------------------------------------------------------------------
option(ENGINE_BUILD_TESTS "Build cpp-engine tests" ON)
if(ENGINE_BUILD_TESTS)
  enable_testing()
  add_subdirectory(tests)
endif()
