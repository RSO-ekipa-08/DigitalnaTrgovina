# go-simple-demo

This is a simple demo of a Go server and client, using Connect RPC, following the tutorial at https://connectrpc.com/docs/go/getting-started

## Server

Run the server with:
```sh
go run cmd/server/main.go
```

## Client

Run the client with:
```sh
go run cmd/client/main.go
```

Alternatively, make a request to the server with curl:
```sh
curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Jane"}' \
    http://localhost:8080/greet.v1.GreetService/Greet
```

Or with grpcurl:
```sh
grpcurl \
    -protoset <(buf build -o -) -plaintext \
    -d '{"name": "Jane"}' \
    localhost:8080 greet.v1.GreetService/Greet
```
