package main

import (
	"log"
	"net"
	"net/http"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	pb "authentication/src/gen/proto"
	"authentication/src/platform/authenticator"
	grpcServer "authentication/src/platform/grpc"
	"authentication/src/platform/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	// Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterAuthServiceServer(s, grpcServer.NewServer(auth))

		log.Printf("gRPC server listening on :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server
	rtr := router.New(auth)
	log.Print("HTTP server listening on http://localhost:3000/")
	if err := http.ListenAndServe("0.0.0.0:3000", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
