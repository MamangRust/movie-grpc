package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/renaldyhidayatt/movie_grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Movie struct {
	ID    string `json:"id"`
	Title string `json:"Title"`
	Genre string `json:"genre"`
}

func main() {
	flag.Parse()

	var conn *grpc.ClientConn

	conn, err := grpc.NewClient("server:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewMovieServiceClient(conn)

	r := gin.Default()

	r.GET("/movies", func(ctx *gin.Context) {
		page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		search := ctx.Query("search")

		req := &pb.ReadMoviesRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Search:   search,
		}

		res, err := client.GetMovies(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"movies":       res.Movies,
			"totalRecords": res.TotalRecords,
			"page":         page,
			"pageSize":     pageSize,
		})
	})

	r.GET("/movies/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		res, err := client.GetMovie(ctx, &pb.ReadMovieRequest{Id: id})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": res.Movie,
		})
	})
	r.POST("/movies", func(ctx *gin.Context) {
		var movie Movie

		err := ctx.ShouldBind(&movie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		data := &pb.Movie{
			Title: movie.Title,
			Genre: movie.Genre,
		}
		res, err := client.CreateMovie(ctx, &pb.CreateMovieRequest{
			Movie: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"movie": res.Movie,
		})
	})
	r.PUT("/movies/:id", func(ctx *gin.Context) {
		var movie Movie
		err := ctx.ShouldBind(&movie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.UpdateMovie(ctx, &pb.UpdateMovieRequest{
			Movie: &pb.Movie{
				Id:    movie.ID,
				Title: movie.Title,
				Genre: movie.Genre,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": res.Movie,
		})

	})
	r.DELETE("/movies/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		res, err := client.DeleteMovie(ctx, &pb.DeleteMovieRequest{Id: id})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if res.Success == true {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Movie deleted successfully",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error deleting movie",
			})
			return
		}

	})

	r.Run(":5000")

}
