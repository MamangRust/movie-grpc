package service

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/renaldyhidayatt/movie_grpc/logger"
	pb "github.com/renaldyhidayatt/movie_grpc/proto"
	mencache "github.com/renaldyhidayatt/movie_grpc/redis"
	"github.com/renaldyhidayatt/movie_grpc/repository"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func generateMovieListKey(page, pageSize int, search string) string {
	if search == "" {
		search = "_"
	}
	return fmt.Sprintf("movie:list:page=%d:size=%d:search=%s", page, pageSize, search)
}

type MovieService struct {
	trace           trace.Tracer
	logger          logger.LoggerInterface
	mencache        mencache.MovieServiceCache
	repo            repository.MovieRepository
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	pb.UnimplementedMovieServiceServer
}

func NewMovieService(repo repository.MovieRepository, trace trace.Tracer, logger logger.LoggerInterface, mencache mencache.MovieServiceCache) *MovieService {
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

	prometheus.MustRegister(requestCounter, requestDuration)

	return &MovieService{
		repo:            repo,
		trace:           trace,
		logger:          logger,
		mencache:        mencache,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *MovieService) recordMetrics(method string, status string, startTime time.Time) {
	duration := time.Since(startTime).Seconds()
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(duration)
}

func (s *MovieService) CreateMovie(ctx context.Context, req *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	var err error
	ctx, end := s.startTracingAndLogging(
		ctx,
		"CreateMovie",
		attribute.String("movie.title", req.GetMovie().Title),
	)
	defer func() { end(err) }()

	movie := req.GetMovie()
	err = s.repo.CreateMovie(ctx, movie)
	if err != nil {
		return nil, err
	}

	return &pb.CreateMovieResponse{
		Movie: movie,
	}, nil
}

func (s *MovieService) GetMovie(ctx context.Context, req *pb.ReadMovieRequest) (*pb.ReadMovieResponse, error) {
	var err error
	ctx, end := s.startTracingAndLogging(
		ctx,
		"GetMovie",
		attribute.String("movie.id", req.GetId()),
	)
	defer func() { end(err) }()

	movie, err := s.repo.GetMovie(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.ReadMovieResponse{
		Movie: movie,
	}, nil
}

func (s *MovieService) GetMovies(ctx context.Context, req *pb.ReadMoviesRequest) (*pb.ReadMoviesResponse, error) {
	var err error
	ctx, end := s.startTracingAndLogging(
		ctx,
		"GetMovies",
	)
	defer func() { end(err) }()

	page := int(req.GetPage())
	if page < 1 {
		page = 1
	}
	pageSize := int(req.GetPageSize())
	if pageSize < 1 {
		pageSize = 10
	}
	search := req.GetSearch()
	cacheKey := generateMovieListKey(page, pageSize, search)

	result, err := s.mencache.GetMovieList(ctx, cacheKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get movie list from cache: %v", err)
	}
	if result != nil {
		return &pb.ReadMoviesResponse{
			Movies:       result.Movies,
			TotalRecords: result.TotalRecords,
		}, nil
	}

	result, err = s.repo.GetMovies(ctx, page, pageSize, search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch movies: %v", err)
	}

	_ = s.mencache.SetMovieList(ctx, cacheKey, result)

	return &pb.ReadMoviesResponse{
		Movies:       result.Movies,
		TotalRecords: result.TotalRecords,
	}, nil
}

func (s *MovieService) UpdateMovie(ctx context.Context, req *pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	var err error
	ctx, end := s.startTracingAndLogging(
		ctx,
		"UpdateMovie",
		attribute.String("movie.id", req.GetMovie().Id),
	)
	defer func() { end(err) }()

	movie := req.GetMovie()
	updatedMovie, err := s.repo.UpdateMovie(ctx, movie)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateMovieResponse{
		Movie: updatedMovie,
	}, nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, req *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	var err error
	ctx, end := s.startTracingAndLogging(
		ctx,
		"DeleteMovie",
		attribute.String("movie.id", req.GetId()),
	)
	defer func() { end(err) }()

	err = s.repo.DeleteMovie(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.DeleteMovieResponse{
		Success: true,
	}, nil
}

func (s *MovieService) startTracingAndLogging(
	ctx context.Context,
	method string,
	attrs ...attribute.KeyValue,
) (context.Context, func(error)) {
	start := time.Now()
	ctx, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}
	span.AddEvent("Start: " + method)

	s.logger.Debug("Start: " + method)

	end := func(err error) {
		duration := time.Since(start)

		span.SetAttributes(attribute.Float64("execution_duration_ms", float64(duration.Milliseconds())))

		if err != nil {
			span.RecordError(err)
			span.SetStatus(otelcodes.Error, err.Error())
			s.logger.Error("Error in "+method,
				zap.Error(err),
				zap.Duration("duration", duration),
			)
			s.recordMetrics(method, "error", start)
		} else {
			span.SetStatus(otelcodes.Ok, "success")
			s.logger.Info("Success: "+method,
				zap.Duration("duration", duration),
			)
			s.recordMetrics(method, "success", start)
		}

		span.End()
	}

	return ctx, end
}
