FROM gcc:13 AS build
RUN apt-get update && apt-get install -y --no-install-recommends cmake ninja-build && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY . .
RUN cmake -S . -B build -G Ninja -DENGINE_BUILD_TESTS=OFF && cmake --build build --target engine_server

FROM debian:bookworm-slim
COPY --from=build /src/build/engine_server /usr/local/bin/engine_server
EXPOSE 50051
CMD ["/usr/local/bin/engine_server"]
