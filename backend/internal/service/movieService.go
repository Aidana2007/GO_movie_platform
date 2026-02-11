package service

import (
	"errors"
	"github.com/Aidana2007/GO_movie_platform/internal/model"
	"github.com/Aidana2007/GO_movie_platform/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieService struct {
	movieRepo *repository.MovieRepository
}

func NewMovieService(movieRepo *repository.MovieRepository) *MovieService {
	return &MovieService{
		movieRepo: movieRepo,
	}
}

func (s *MovieService) CreateMovie(req *model.CreateMovieRequest, userID primitive.ObjectID) (*model.Movie, error) {
	movie := &model.Movie{
		Title:       req.Title,
		Description: req.Description,
		Year:        req.Year,
		Director:    req.Director,
		Cast:        req.Cast,
		Genre:       req.Genre,
		PosterURL:   req.PosterURL,
		TrailerURL:  req.TrailerURL,
		CreatedBy:   userID,
	}

	if err := s.movieRepo.Create(movie); err != nil {
		return nil, err
	}

	return movie, nil
}

func (s *MovieService) GetAllMovies() ([]*model.Movie, error) {
	return s.movieRepo.FindAll()
}

func (s *MovieService) GetMovieByID(id primitive.ObjectID) (*model.Movie, error) {
	return s.movieRepo.FindByID(id)
}

func (s *MovieService) GetMoviesByIDs(ids []primitive.ObjectID) ([]*model.Movie, error) {
	return s.movieRepo.FindByIDs(ids)
}

func (s *MovieService) UpdateMovie(id primitive.ObjectID, req *model.CreateMovieRequest, userID primitive.ObjectID) (*model.Movie, error) {
	movie, err := s.movieRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user is the creator or admin
	if movie.CreatedBy != userID {
		return nil, errors.New("unauthorized to update this movie")
	}

	movie.Title = req.Title
	movie.Description = req.Description
	movie.Year = req.Year
	movie.Director = req.Director
	movie.Cast = req.Cast
	movie.Genre = req.Genre
	movie.PosterURL = req.PosterURL
	movie.TrailerURL = req.TrailerURL

	if err := s.movieRepo.Update(id, movie); err != nil {
		return nil, err
	}

	return movie, nil
}

func (s *MovieService) DeleteMovie(id primitive.ObjectID, userID primitive.ObjectID) error {
	movie, err := s.movieRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check if user is the creator or admin
	if movie.CreatedBy != userID {
		return errors.New("unauthorized to delete this movie")
	}

	return s.movieRepo.Delete(id)
}

func (s *MovieService) GetTopRated(limit int) ([]*model.Movie, error) {
	return s.movieRepo.GetTopRated(limit)
}

func (s *MovieService) SearchMovies(query, genre, sort string, minRating float64) ([]*model.Movie, error) {
	return s.movieRepo.Search(query, genre, sort, minRating)
}
