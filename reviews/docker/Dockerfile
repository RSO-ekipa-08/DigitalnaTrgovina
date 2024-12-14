FROM rust:1.75 as builder

WORKDIR /usr/src/app
COPY . .
RUN SQLX_OFFLINE=true cargo build --release

FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y libpq5 ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/app/target/release/reviews /usr/local/bin/
EXPOSE 50051

CMD ["reviews"]
