FROM golang:1.23-alpine AS builder

RUN apk add --no-cache make protoc protobuf-dev

# Define current working directory
WORKDIR /apigateway

# Install protoc-gen-go and protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Download modules to local cache so we can skip re-
# downloading on consecutive docker build commands
COPY go.mod go.sum ./
RUN go mod download

# Add sources
COPY . .

# Generate proto files
RUN make proto

RUN CGO_ENABLED=0 GOOS=linux go build -o out/api-gateway ./src

# Final stage
FROM alpine:latest

RUN adduser -D -g '' appuser

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /apigateway/out/api-gateway /app/api-gateway

USER appuser

# Expose port 50051 for grpc
EXPOSE 50051

# Expose 3000 for http
EXPOSE 3000

CMD ["/app/api-gateway"]
