package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	pb "github.com/renaldyhidayatt/movie_grpc/proto"
	"github.com/renaldyhidayatt/movie_grpc/repository"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type MovieService struct {
	trace           trace.Tracer
	repo            repository.MovieRepository
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	pb.UnimplementedMovieServiceServer
}

func NewMovieService(repo repository.MovieRepository, trace trace.Tracer) *MovieService {
	// Register Prometheus metrics
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "movie_service_requests_total",
			Help: "Total number of requests to the MovieService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "movie_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MovieService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// Register metrics with Prometheus
	prometheus.MustRegister(requestCounter, requestDuration)

	return &MovieService{
		repo:            repo,
		trace:           trace,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

// Helper to observe metrics
func (s *MovieService) recordMetrics(method string, status string, startTime time.Time) {
	duration := time.Since(startTime).Seconds()
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(duration)
}

func (s *MovieService) CreateMovie(ctx context.Context, req *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	startTime := time.Now()
	method := "CreateMovie"

	ctx, span := s.trace.Start(ctx, method)
	defer span.End()

	log.Println("Create Movie")
	span.SetAttributes(attribute.String("method", method))
	span.SetAttributes(attribute.String("movie_title", req.GetMovie().Title))

	movie := req.GetMovie()
	err := s.repo.CreateMovie(ctx, movie)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Int("http.status_code", http.StatusInternalServerError))
		span.SetStatus(codes.Error, "Failed to create movie")
		s.recordMetrics(method, "error", startTime)
		return nil, err
	}

	span.SetAttributes(attribute.Int("http.status_code", http.StatusOK))
	span.SetStatus(codes.Ok, "Movie created successfully")
	s.recordMetrics(method, "success", startTime)

	return &pb.CreateMovieResponse{
		Movie: movie,
	}, nil
}

func (s *MovieService) GetMovie(ctx context.Context, req *pb.ReadMovieRequest) (*pb.ReadMovieResponse, error) {
	startTime := time.Now()
	method := "GetMovie"

	ctx, span := s.trace.Start(ctx, method)
	defer span.End()

	log.Println("Read Movie", req.GetId())
	span.SetAttributes(attribute.String("method", method))
	span.SetAttributes(attribute.String("movie_id", req.GetId()))

	movie, err := s.repo.GetMovie(ctx, req.GetId())
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Int("http.status_code", http.StatusNotFound))
		span.SetStatus(codes.Error, "Movie not found")
		s.recordMetrics(method, "error", startTime)
		return nil, err
	}

	span.SetAttributes(attribute.Int("http.status_code", http.StatusOK))
	span.SetStatus(codes.Ok, "Movie fetched successfully")
	s.recordMetrics(method, "success", startTime)

	return &pb.ReadMovieResponse{
		Movie: movie,
	}, nil
}

func (s *MovieService) GetMovies(ctx context.Context, req *pb.ReadMoviesRequest) (*pb.ReadMoviesResponse, error) {
	startTime := time.Now()
	method := "GetMovies"

	ctx, span := s.trace.Start(ctx, method)
	defer span.End()

	log.Println("Read Movies")
	span.SetAttributes(attribute.String("method", method))

	movies, err := s.repo.GetMovies(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Int("http.status_code", http.StatusInternalServerError))
		span.SetStatus(codes.Error, "Failed to fetch movies")
		s.recordMetrics(method, "error", startTime)
		return nil, err
	}

	span.SetAttributes(attribute.Int("http.status_code", http.StatusOK))
	span.SetStatus(codes.Ok, "Movies fetched successfully")
	s.recordMetrics(method, "success", startTime)

	return &pb.ReadMoviesResponse{
		Movies: movies,
	}, nil
}

func (s *MovieService) UpdateMovie(ctx context.Context, req *pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	startTime := time.Now()
	method := "UpdateMovie"

	ctx, span := s.trace.Start(ctx, method)
	defer span.End()

	log.Println("Update Movie")
	span.SetAttributes(attribute.String("method", method))
	span.SetAttributes(attribute.String("movie_id", req.GetMovie().Id))

	movie := req.GetMovie()
	updatedMovie, err := s.repo.UpdateMovie(ctx, movie)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Int("http.status_code", http.StatusInternalServerError))
		span.SetStatus(codes.Error, "Failed to update movie")
		s.recordMetrics(method, "error", startTime)
		return nil, err
	}

	span.SetAttributes(attribute.Int("http.status_code", http.StatusOK))
	span.SetStatus(codes.Ok, "Movie updated successfully")
	s.recordMetrics(method, "success", startTime)

	return &pb.UpdateMovieResponse{
		Movie: updatedMovie,
	}, nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, req *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	startTime := time.Now()
	method := "DeleteMovie"

	ctx, span := s.trace.Start(ctx, method)
	defer span.End()

	log.Println("Delete Movie")
	span.SetAttributes(attribute.String("method", method))
	span.SetAttributes(attribute.String("movie_id", req.GetId()))

	err := s.repo.DeleteMovie(ctx, req.GetId())
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Int("http.status_code", http.StatusInternalServerError))
		span.SetStatus(codes.Error, "Failed to delete movie")
		s.recordMetrics(method, "error", startTime)
		return nil, err
	}

	span.SetAttributes(attribute.Int("http.status_code", http.StatusOK))
	span.SetStatus(codes.Ok, "Movie deleted successfully")
	s.recordMetrics(method, "success", startTime)

	return &pb.DeleteMovieResponse{
		Success: true,
	}, nil
}
