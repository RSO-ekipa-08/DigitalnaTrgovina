.PHONY: proto
proto:
	mkdir -p src/gen/proto
	protoc --go_out=./src/gen --go_opt=module=authentication/src/gen \
		--go-grpc_out=./src/gen --go-grpc_opt=module=authentication/src/gen \
		proto/auth.proto
