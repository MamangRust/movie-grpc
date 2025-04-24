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
	"google.golang.org/grpc"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	DatabaseConnection()
}

var DB *gorm.DB
var err error

func DatabaseConnection() {
	dbPath := "movie_grpc.db"
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}

	if err := DB.AutoMigrate(&models.Movie{}); err != nil {
		log.Fatalf("Error during migration: %v", err)
	}

	fmt.Println("Database connection successful using SQLite...")
}

func main() {
	fmt.Println("gRPC server running ...")
	ctx := context.Background()
	flag.Parse()

	grpcLis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	metricsLis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen for metrics: %v", err)
	}

	shutdownTracerProvider, err := config.InitTracerProvider(ctx)
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

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
				otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
			),
		),
	)

	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.Handler())

	pb.RegisterMovieServiceServer(grpcServer, movieService)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Println("Metrics server listening on :8080")
		if err := http.Serve(metricsLis, metricsServer); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	wg.Wait()
}
