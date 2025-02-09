package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/karthikchekuri1208/grpc-weather-service/gen/github.com/karthikchekuri1208/grpc-weather-service/proto"
	"github.com/karthikchekuri1208/grpc-weather-service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// runGRPCServer is to start the GRPC server
func runGRPCServer(apiKey string) {
	// gRPC server listens on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	// Create a new WeatherServiceServer
	weatherService := service.NewWeatherServiceServer(apiKey)
	// Register the service
	pb.RegisterWeatherServiceServer(s, weatherService)

	log.Println("gRPC server is running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// runGateway is to expose grpc as http.
func runGateway() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	// Using insecure credentials for demo purpose.
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register the gRPC service handler
	err := pb.RegisterWeatherServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("HTTP gateway is running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	// Getting API key as environment variable
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENWEATHER_API_KEY environment variable is not set")
	}

	// Start the gRPC server in a goroutine
	go runGRPCServer(apiKey)
	// Start the HTTP gateway
	runGateway()
}
