package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/renaldyhidayatt/movie_grpc/dto"
	"github.com/renaldyhidayatt/movie_grpc/models"
	pb "github.com/renaldyhidayatt/movie_grpc/proto"
	"gorm.io/gorm"
)

type MovieRepository interface {
	CreateMovie(ctx context.Context, movie *pb.Movie) error
	GetMovie(ctx context.Context, id string) (*pb.Movie, error)
	GetMovies(ctx context.Context, page, pageSize int, search string) (*dto.MovieListResult, error)
	UpdateMovie(ctx context.Context, movie *pb.Movie) (*pb.Movie, error)
	DeleteMovie(ctx context.Context, id string) error
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{
		db: db,
	}
}

func (r *movieRepository) CreateMovie(ctx context.Context, movie *pb.Movie) error {
	movie.Id = uuid.New().String()

	data := models.Movie{
		ID:    movie.GetId(),
		Title: movie.GetTitle(),
		Genre: movie.GetGenre(),
	}

	res := r.db.Create(&data)
	if res.RowsAffected == 0 {
		return errors.New("movie creation unsuccessful")
	}
	return nil
}

func (r *movieRepository) GetMovie(ctx context.Context, id string) (*pb.Movie, error) {
	var movie models.Movie
	res := r.db.Find(&movie, "id = ?", id)
	if res.RowsAffected == 0 {
		return nil, errors.New("movie not found")
	}
	return &pb.Movie{
		Id:    movie.ID,
		Title: movie.Title,
		Genre: movie.Genre,
	}, nil
}

func (r *movieRepository) GetMovies(ctx context.Context, page, pageSize int, search string) (*dto.MovieListResult, error) {
	var (
		movies       []*models.Movie
		totalRecords int64
	)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	query := r.db.WithContext(ctx).Model(&models.Movie{})

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(genre) LIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, fmt.Errorf("failed to count movies: %w", err)
	}

	if err := query.Limit(pageSize).Offset(offset).Find(&movies).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}
	if len(movies) == 0 {
		return nil, errors.New("movies not found")
	}

	pbMovies := make([]*pb.Movie, len(movies))
	for i, m := range movies {
		pbMovies[i] = &pb.Movie{
			Id:    m.ID,
			Title: m.Title,
			Genre: m.Genre,
		}
	}

	return &dto.MovieListResult{
		Movies:       pbMovies,
		TotalRecords: totalRecords,
	}, nil
}

func (r *movieRepository) UpdateMovie(ctx context.Context, movie *pb.Movie) (*pb.Movie, error) {
	var m models.Movie
	res := r.db.Model(&m).Where("id=?", movie.Id).Updates(models.Movie{Title: movie.Title, Genre: movie.Genre})
	if res.RowsAffected == 0 {
		return nil, errors.New("movie not found")
	}
	return &pb.Movie{
		Id:    m.ID,
		Title: m.Title,
		Genre: m.Genre,
	}, nil
}

func (r *movieRepository) DeleteMovie(ctx context.Context, id string) error {
	var movie models.Movie
	res := r.db.Where("id = ?", id).Delete(&movie)
	if res.RowsAffected == 0 {
		return errors.New("movie not found")
	}
	return nil
}
