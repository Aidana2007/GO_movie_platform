package service

import (
	"errors"

	"github.com/Aidana2007/GO_movie_platform/internal/model"
	"github.com/Aidana2007/GO_movie_platform/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	movieRepo  *repository.MovieRepository
}

func NewReviewService(reviewRepo *repository.ReviewRepository, movieRepo *repository.MovieRepository) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		movieRepo:  movieRepo,
	}
}

func (s *ReviewService) CreateReview(movieID primitive.ObjectID, req *model.CreateReviewRequest, userID primitive.ObjectID, username string) (*model.Review, error) {
	if req.Rating < 1 || req.Rating > 10 {
		return nil, errors.New("rating must be between 1 and 10")
	}

	_, err := s.movieRepo.FindByID(movieID)
	if err != nil {
		return nil, errors.New("movie not found")
	}

	// Check if user already reviewed this movie
	existingReview, err := s.reviewRepo.FindByUserAndMovie(userID, movieID)
	if err == nil && existingReview != nil {
		return nil, errors.New("you have already reviewed this movie")
	}

	review := &model.Review{
		User:     userID,
		Movie:    movieID,
		Username: username,
		Rating:   req.Rating,
		Comment:  req.Comment,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	if err := s.movieRepo.AddReview(movieID, review.ID); err != nil {
		return nil, err
	}

	avgRating, err := s.reviewRepo.GetAverageRating(movieID)
	if err == nil {
		s.movieRepo.UpdateRanking(movieID, avgRating)
	}

	return review, nil
}

func (s *ReviewService) GetMovieReviews(movieID primitive.ObjectID) ([]*model.Review, error) {
	return s.reviewRepo.FindByMovieID(movieID)
}

func (s *ReviewService) DeleteReview(reviewID primitive.ObjectID, userID primitive.ObjectID) error {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return err
	}

	if review.User != userID {
		return errors.New("unauthorized to delete this review")
	}

	if err := s.movieRepo.RemoveReview(review.Movie, reviewID); err != nil {
		return err
	}

	if err := s.reviewRepo.Delete(reviewID); err != nil {
		return err
	}

	avgRating, err := s.reviewRepo.GetAverageRating(review.Movie)
	if err == nil {
		s.movieRepo.UpdateRanking(review.Movie, avgRating)
	}

	return nil
}

func (s *ReviewService) UpdateReview(reviewID primitive.ObjectID, req *model.CreateReviewRequest, userID primitive.ObjectID) (*model.Review, error) {
	if req.Rating < 1 || req.Rating > 10 {
		return nil, errors.New("rating must be between 1 and 10")
	}

	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return nil, errors.New("review not found")
	}

	if review.User != userID {
		return nil, errors.New("unauthorized to update this review")
	}

	if err := s.reviewRepo.Update(reviewID, req.Rating, req.Comment); err != nil {
		return nil, err
	}

	// Update movie ranking
	avgRating, err := s.reviewRepo.GetAverageRating(review.Movie)
	if err == nil {
		s.movieRepo.UpdateRanking(review.Movie, avgRating)
	}

	// Get updated review
	updatedReview, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return nil, err
	}

	return updatedReview, nil
}

func (s *ReviewService) DeleteReviewAdmin(reviewID primitive.ObjectID) error {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return err
	}

	if err := s.movieRepo.RemoveReview(review.Movie, reviewID); err != nil {
		return err
	}

	if err := s.reviewRepo.Delete(reviewID); err != nil {
		return err
	}

	avgRating, err := s.reviewRepo.GetAverageRating(review.Movie)
	if err == nil {
		s.movieRepo.UpdateRanking(review.Movie, avgRating)
	}

	return nil
}
