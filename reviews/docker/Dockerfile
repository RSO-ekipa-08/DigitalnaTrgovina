FROM rust:latest AS builder

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

FROM debian:bookworm-slim
# Install OpenSSL and PostgreSQL client libraries
RUN apt-get update && \
    apt-get install -y \
    libssl3 \
    ca-certificates \
    libpq5 && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/app/target/release/reviews /usr/local/bin/
EXPOSE 50051 8080

CMD ["reviews"]
