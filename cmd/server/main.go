// main.go

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/renaldyhidayatt/movie_grpc/config"
	"github.com/renaldyhidayatt/movie_grpc/models"
	pb "github.com/renaldyhidayatt/movie_grpc/proto"
	"github.com/renaldyhidayatt/movie_grpc/repository"
	"github.com/renaldyhidayatt/movie_grpc/service"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	DatabaseConnection()
}

var DB *gorm.DB
var err error

var serviceName = semconv.ServiceNameKey.String("test-service")

func DatabaseConnection() {
	dbPath := "movie_grpc.db" // Path file SQLite
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}

	// Migrasi model ke database SQLite
	if err := DB.AutoMigrate(&models.Movie{}); err != nil {
		log.Fatalf("Error during migration: %v", err)
	}

	fmt.Println("Database connection successful using SQLite...")
}

func main() {
	fmt.Println("gRPC server running ...")
	ctx := context.Background()
	flag.Parse()

	// Listener untuk gRPC
	grpcLis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	// Listener untuk metrics
	metricsLis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen for metrics: %v", err)
	}

	conn, err := config.InitConn()
	if err != nil {
		log.Fatal(err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(serviceName),
	)
	if err != nil {
		log.Fatal(err)
	}

	shutdownTracerProvider, err := config.InitTracerProvider(ctx, res, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdownTracerProvider(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %s", err)
		}
	}()

	tracer := otel.Tracer("hello")
	movieRepo := repository.NewMovieRepository(DB)
	movieService := service.NewMovieService(movieRepo, tracer)

	// Create the gRPC server
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(),
		),
	)

	// Create metrics server
	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.Handler())

	// Register the movie service with the gRPC server
	pb.RegisterMovieServiceServer(grpcServer, movieService)

	// Use WaitGroup to ensure both servers run
	var wg sync.WaitGroup
	wg.Add(2)

	// Start metrics server
	go func() {
		defer wg.Done()
		log.Println("Metrics server listening on :8080")
		if err := http.Serve(metricsLis, metricsServer); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	// Start gRPC server
	go func() {
		defer wg.Done()
		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Wait for both servers to finish
	wg.Wait()
}
