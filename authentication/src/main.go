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
)

func main() {
	godotenv.Load()

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	// Start HTTP server for health check
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		log.Printf("HTTP (/health) server listening on :3000")
		if err := http.ListenAndServe(":3000", nil); err != nil {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	// Start gRPC server
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
}
