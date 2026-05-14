package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/KellyHarvestOS/kinotower-Go/internal/auth"
	"github.com/KellyHarvestOS/kinotower-Go/internal/models"
	"github.com/KellyHarvestOS/kinotower-Go/internal/repositories"
	"github.com/KellyHarvestOS/kinotower-Go/internal/validator"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrForbidden          = errors.New("forbidden")
	ErrScoreExists        = errors.New("score exist")
)

type Service struct {
	repo      *repositories.Repository
	jwtSecret string
	jwtTTL    time.Duration
}

func New(repo *repositories.Repository, jwtSecret string, jwtTTL time.Duration) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret, jwtTTL: jwtTTL}
}

func (s *Service) ListFilms(ctx context.Context, filters models.FilmFilters) ([]models.Film, int, error) {
	filters.Search = strings.TrimSpace(filters.Search)
	if filters.SortBy == "" {
		filters.SortBy = "name"
	}
	if filters.SortDir != "desc" {
		filters.SortDir = "asc"
	}
	return s.repo.ListFilms(ctx, filters)
}

func (s *Service) GetFilm(ctx context.Context, id int) (models.Film, error) {
	return s.repo.GetFilm(ctx, id)
}
func (s *Service) Categories(ctx context.Context) ([]models.Category, error) {
	return s.repo.ListCategories(ctx)
}
func (s *Service) Countries(ctx context.Context) ([]models.Country, error) {
	return s.repo.ListCountries(ctx)
}
func (s *Service) Genders(ctx context.Context) ([]models.Gender, error) {
	return s.repo.ListGenders(ctx)
}

func (s *Service) Signup(ctx context.Context, req models.SignupRequest) (models.AuthResponse, error) {
	if err := validator.User(req.FIO, req.Email, req.Birthday, req.GenderID, &req.Password); err != nil {
		return models.AuthResponse{}, err
	}
	birthday, _ := time.Parse("2006-01-02", req.Birthday)
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return models.AuthResponse{}, err
	}
	user, err := s.repo.CreateUser(ctx, req, hash, birthday)
	if err != nil {
		return models.AuthResponse{}, err
	}
	token, err := auth.NewToken(user.ID, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return models.AuthResponse{}, err
	}
	return models.AuthResponse{Status: "success", Token: token, ID: user.ID, FIO: user.FIO}, nil
}

func (s *Service) Signin(ctx context.Context, req models.SigninRequest) (models.AuthResponse, error) {
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil || !auth.CheckPassword(user.Password, req.Password) {
		return models.AuthResponse{}, ErrInvalidCredentials
	}
	token, err := auth.NewToken(user.ID, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return models.AuthResponse{}, err
	}
	return models.AuthResponse{Status: "success", Token: token, ID: user.ID, FIO: user.FIO}, nil
}

func (s *Service) User(ctx context.Context, id int) (models.User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *Service) UpdateUser(ctx context.Context, authUserID int, req models.UpdateUserRequest) error {
	if err := validator.User(req.FIO, req.Email, req.Birthday, req.GenderID, nil); err != nil {
		return err
	}
	birthday, _ := time.Parse("2006-01-02", req.Birthday)
	return s.repo.UpdateUser(ctx, authUserID, req, birthday)
}

func (s *Service) DeleteUser(ctx context.Context, authUserID int) error {
	return s.repo.DeleteUser(ctx, authUserID)
}

func (s *Service) FilmReviews(ctx context.Context, filmID int) ([]models.Review, error) {
	if !s.repo.Exists(ctx, `SELECT EXISTS(SELECT 1 FROM films WHERE id=$1 AND deleted_at IS NULL)`, filmID) {
		return nil, repositories.ErrNotFound
	}
	return s.repo.FilmReviews(ctx, filmID)
}

func (s *Service) CreateReview(ctx context.Context, authUserID, pathUserID int, req models.CreateReviewRequest) (models.Review, error) {
	if authUserID != pathUserID {
		return models.Review{}, ErrForbidden
	}
	if req.FilmID < 1 || len(req.Message) < 4 || len(req.Message) > 1024 {
		return models.Review{}, errors.New("invalid review data")
	}
	if !s.repo.Exists(ctx, `SELECT EXISTS(SELECT 1 FROM films WHERE id=$1 AND deleted_at IS NULL)`, req.FilmID) {
		return models.Review{}, errors.New("invalid review data")
	}
	return s.repo.CreateReview(ctx, authUserID, req)
}

func (s *Service) UserReviews(ctx context.Context, userID int) ([]models.Review, error) {
	return s.repo.UserReviews(ctx, userID)
}

func (s *Service) DeleteReview(ctx context.Context, authUserID, pathUserID, id int) error {
	if authUserID != pathUserID {
		return ErrForbidden
	}
	return s.repo.DeleteReview(ctx, authUserID, id)
}

func (s *Service) CreateRating(ctx context.Context, authUserID, pathUserID int, req models.CreateRatingRequest) (models.Rating, error) {
	if authUserID != pathUserID {
		return models.Rating{}, ErrForbidden
	}
	if req.FilmID < 1 || req.Ball < 1 || req.Ball > 5 {
		return models.Rating{}, errors.New("invalid rating data")
	}
	if !s.repo.Exists(ctx, `SELECT EXISTS(SELECT 1 FROM films WHERE id=$1 AND deleted_at IS NULL)`, req.FilmID) {
		return models.Rating{}, errors.New("invalid rating data")
	}
	rating, err := s.repo.CreateRating(ctx, authUserID, req)
	if err != nil {
		return models.Rating{}, ErrScoreExists
	}
	return rating, nil
}

func (s *Service) UserRatings(ctx context.Context, userID int) ([]models.Rating, error) {
	return s.repo.UserRatings(ctx, userID)
}

func (s *Service) DeleteRating(ctx context.Context, authUserID, pathUserID, id int) error {
	if authUserID != pathUserID {
		return ErrForbidden
	}
	return s.repo.DeleteRating(ctx, authUserID, id)
}
