include(FetchContent)
FetchContent_Declare(
  googletest
  GIT_REPOSITORY https://github.com/google/googletest.git
  GIT_TAG        v1.14.0
)
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

add_executable(engine_tests
  risk_test.cpp
)
target_link_libraries(engine_tests PRIVATE engine_core GTest::gtest_main)

include(GoogleTest)
gtest_discover_tests(engine_tests)
