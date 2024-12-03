package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/renaldyhidayatt/movie_grpc/models"
	pb "github.com/renaldyhidayatt/movie_grpc/proto"
	"gorm.io/gorm"
)

type MovieRepository interface {
	CreateMovie(ctx context.Context, movie *pb.Movie) error
	GetMovie(ctx context.Context, id string) (*pb.Movie, error)
	GetMovies(ctx context.Context) ([]*pb.Movie, error)
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

func (r *movieRepository) GetMovies(ctx context.Context) ([]*pb.Movie, error) {
	var movies []*models.Movie
	res := r.db.Find(&movies)
	if res.RowsAffected == 0 {
		return nil, errors.New("movies not found")
	}

	// Convert []*Movie to []*pb.Movie
	pbMovies := make([]*pb.Movie, len(movies))
	for i, m := range movies {
		pbMovies[i] = &pb.Movie{
			Id:    m.ID,
			Title: m.Title,
			Genre: m.Genre,
		}
	}

	return pbMovies, nil
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
