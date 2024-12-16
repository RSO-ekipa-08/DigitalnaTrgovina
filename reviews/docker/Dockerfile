FROM rust:latest as builder

WORKDIR /usr/src/app

# Install protobuf compiler
RUN apt-get update && \
    apt-get install -y protobuf-compiler && \
    rm -rf /var/lib/apt/lists/*

# Copy Cargo.toml and Cargo.lock first to cache dependencies
COPY . .

# Build the application
ENV SQLX_OFFLINE=true
RUN cargo build --release



FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y libpq5 ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/app/target/release/reviews /usr/local/bin/
EXPOSE 50051

CMD ["reviews"]
