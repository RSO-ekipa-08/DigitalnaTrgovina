FROM golang:1.19.3-alpine3.17

# Define current working directory
WORKDIR /authentication

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

RUN go build -o out/auth-service ./src


# Expose port 50051 for grpc
EXPOSE 50051

# Expose 3000 for http
EXPOSE 3000

CMD ["/authentication/out/auth-service"]
