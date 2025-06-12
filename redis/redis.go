package mencache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/renaldyhidayatt/movie_grpc/dto"
	pb "github.com/renaldyhidayatt/movie_grpc/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type MovieServiceCache interface {
	GetMovie(ctx context.Context, id string) (*pb.Movie, error)
	SetMovie(ctx context.Context, movie *pb.Movie) error
	DeleteMovie(ctx context.Context, id string) error

	GetMovieList(ctx context.Context, key string) (*dto.MovieListResult, error)
	SetMovieList(ctx context.Context, key string, result *dto.MovieListResult) error
	DeleteMovieList(ctx context.Context, key string) error
}

type redisMovieCache struct {
	client     *redis.Client
	expiration time.Duration
}

func NewMovieServiceCache(client *redis.Client, expiration time.Duration) MovieServiceCache {
	return &redisMovieCache{
		client:     client,
		expiration: expiration,
	}
}

func (r *redisMovieCache) key(id string) string {
	return fmt.Sprintf("movie:%s", id)
}

func (r *redisMovieCache) GetMovie(ctx context.Context, id string) (*pb.Movie, error) {
	val, err := r.client.Get(ctx, r.key(id)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var movie pb.Movie
	if err := protojson.Unmarshal([]byte(val), &movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r *redisMovieCache) SetMovie(ctx context.Context, movie *pb.Movie) error {
	jsonData, err := protojson.Marshal(movie)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.key(movie.Id), jsonData, r.expiration).Err()
}

func (r *redisMovieCache) DeleteMovie(ctx context.Context, id string) error {
	return r.client.Del(ctx, r.key(id)).Err()
}

func (r *redisMovieCache) GetMovieList(ctx context.Context, key string) (*dto.MovieListResult, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var result dto.MovieListResult
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *redisMovieCache) SetMovieList(ctx context.Context, key string, result *dto.MovieListResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, r.expiration).Err()
}

func (r *redisMovieCache) DeleteMovieList(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
