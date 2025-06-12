package dto

import (
	pb "github.com/renaldyhidayatt/movie_grpc/proto"
)

type MovieListResult struct {
	Movies       []*pb.Movie
	TotalRecords int64
}
