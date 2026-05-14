package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/KellyHarvestOS/kinotower-Go/internal/middleware"
	"github.com/KellyHarvestOS/kinotower-Go/internal/models"
	"github.com/KellyHarvestOS/kinotower-Go/internal/response"
	"github.com/KellyHarvestOS/kinotower-Go/internal/services"
	"github.com/KellyHarvestOS/kinotower-Go/internal/utils"
)

type Handler struct {
	service *services.Service
}

func New(service *services.Service) *Handler { return &Handler{service: service} }

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "web/templates/index.html")
}

func (h *Handler) KinotowerHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/kinotower" && r.URL.Path != "/kinotower/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "web/pages/kinotower/index.html")
}

func (h *Handler) Films(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	p := utils.NewPagination(r)
	category, _ := strconv.Atoi(r.URL.Query().Get("category"))
	country, _ := strconv.Atoi(r.URL.Query().Get("country"))
	films, total, err := h.service.ListFilms(r.Context(), models.FilmFilters{
		Search: r.URL.Query().Get("search"), SortBy: r.URL.Query().Get("sortBy"), SortDir: r.URL.Query().Get("sortDir"), Category: category, Country: country, Limit: p.Limit, Offset: p.Offset,
	})
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"page": p.Page, "size": len(films), "total": total, "films": films})
}

func (h *Handler) FilmRoutes(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/films/"), "/")
	filmID, ok := parseID(first(parts))
	if !ok {
		response.ErrorJSON(w, http.StatusNotFound, "Film not found")
		return
	}
	if len(parts) == 1 && r.Method == http.MethodGet {
		h.FilmByID(w, r, filmID)
		return
	}
	if len(parts) == 2 && parts[1] == "reviews" && r.Method == http.MethodGet {
		h.FilmReviews(w, r, filmID)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) FilmByID(w http.ResponseWriter, r *http.Request, id int) {
	film, err := h.service.GetFilm(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		response.ErrorJSON(w, http.StatusNotFound, "Film not found")
		return
	}
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, film)
}

func (h *Handler) Categories(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.Categories(r.Context())
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"categories": items})
}

func (h *Handler) Countries(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.Countries(r.Context())
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"countries": items})
}

func (h *Handler) Genders(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.Genders(r.Context())
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"genders": items})
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req models.SignupRequest
	if err := response.Decode(r, &req); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	res, err := h.service.Signup(r.Context(), req)
	if err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, res)
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req models.SigninRequest
	if err := response.Decode(r, &req); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	res, err := h.service.Signin(r.Context(), req)
	if errors.Is(err, services.ErrInvalidCredentials) {
		response.JSON(w, http.StatusUnauthorized, map[string]string{"status": "invalid", "message": "Wrong email or password"})
		return
	}
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Signout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *Handler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		var req models.UpdateUserRequest
		if err := response.Decode(r, &req); err != nil {
			response.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		if err := h.service.UpdateUser(r.Context(), middleware.UserID(r), req); err != nil {
			response.ErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"status": "success"})
	case http.MethodDelete:
		if err := h.service.DeleteUser(r.Context(), middleware.UserID(r)); err != nil {
			response.ErrorJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) UserRoutes(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/users/"), "/")
	userID, ok := parseID(first(parts))
	if !ok {
		response.ErrorJSON(w, http.StatusNotFound, "User not found")
		return
	}
	if _, err := h.service.User(r.Context(), userID); err != nil {
		response.ErrorJSON(w, http.StatusNotFound, "User not found")
		return
	}

	if len(parts) == 1 && r.Method == http.MethodGet {
		h.User(w, r, userID)
		return
	}
	if len(parts) == 2 && parts[1] == "reviews" && r.Method == http.MethodPost {
		h.CreateReview(w, r, userID)
		return
	}
	if len(parts) == 2 && parts[1] == "reviews" && r.Method == http.MethodGet {
		h.UserReviews(w, r, userID)
		return
	}
	if len(parts) == 3 && parts[1] == "reviews" && r.Method == http.MethodDelete {
		id, ok := parseID(parts[2])
		if !ok {
			response.ErrorJSON(w, http.StatusNotFound, "Review not found")
			return
		}
		h.DeleteReview(w, r, userID, id)
		return
	}
	if len(parts) == 2 && parts[1] == "ratings" && r.Method == http.MethodPost {
		h.CreateRating(w, r, userID)
		return
	}
	if len(parts) == 2 && parts[1] == "ratings" && r.Method == http.MethodGet {
		h.UserRatings(w, r, userID)
		return
	}
	if len(parts) == 3 && parts[1] == "ratings" && r.Method == http.MethodDelete {
		id, ok := parseID(parts[2])
		if !ok {
			response.ErrorJSON(w, http.StatusNotFound, "Rating not found")
			return
		}
		h.DeleteRating(w, r, userID, id)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handler) User(w http.ResponseWriter, r *http.Request, userID int) {
	user, err := h.service.User(r.Context(), userID)
	if err != nil {
		response.ErrorJSON(w, http.StatusNotFound, "User not found")
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"id": user.ID, "fio": user.FIO, "email": user.Email, "birthday": utils.DateISO(user.Birthday), "gender": map[string]any{"id": user.GenderID, "name": user.GenderName}, "reviewCount": user.ReviewCount, "ratingCount": user.RatingCount})
}

func (h *Handler) FilmReviews(w http.ResponseWriter, r *http.Request, filmID int) {
	items, err := h.service.FilmReviews(r.Context(), filmID)
	if err != nil {
		response.ErrorJSON(w, http.StatusNotFound, "Film not found")
		return
	}
	reviews := []map[string]any{}
	for _, item := range items {
		reviews = append(reviews, map[string]any{"id": item.ID, "user": map[string]any{"id": item.UserID, "fio": item.UserFIO}, "message": item.Message, "created_at": utils.TimeISO(item.CreatedAt)})
	}
	response.JSON(w, http.StatusOK, map[string]any{"reviews": reviews})
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request, userID int) {
	var req models.CreateReviewRequest
	if err := response.Decode(r, &req); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	item, err := h.service.CreateReview(r.Context(), middleware.UserID(r), userID, req)
	if errors.Is(err, services.ErrForbidden) {
		response.ErrorJSON(w, http.StatusForbidden, "Forbidden")
		return
	}
	if err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"id": item.ID, "film": map[string]any{"id": item.FilmID, "name": item.FilmName}, "message": item.Message, "is_approved": 0, "created_at": utils.TimeISO(item.CreatedAt)})
}

func (h *Handler) UserReviews(w http.ResponseWriter, r *http.Request, userID int) {
	items, err := h.service.UserReviews(r.Context(), userID)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	reviews := []map[string]any{}
	for _, item := range items {
		reviews = append(reviews, map[string]any{"id": item.ID, "film": map[string]any{"id": item.FilmID, "name": item.FilmName}, "message": item.Message, "is_approved": utils.BoolInt(item.IsApproved), "created_at": utils.TimeISO(item.CreatedAt)})
	}
	response.JSON(w, http.StatusOK, map[string]any{"reviews": reviews})
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request, userID, id int) {
	err := h.service.DeleteReview(r.Context(), middleware.UserID(r), userID, id)
	if errors.Is(err, services.ErrForbidden) {
		response.ErrorJSON(w, http.StatusForbidden, "Forbidden")
		return
	}
	if err != nil {
		response.ErrorJSON(w, http.StatusNotFound, "Review not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateRating(w http.ResponseWriter, r *http.Request, userID int) {
	var req models.CreateRatingRequest
	if err := response.Decode(r, &req); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	item, err := h.service.CreateRating(r.Context(), middleware.UserID(r), userID, req)
	if errors.Is(err, services.ErrForbidden) {
		response.ErrorJSON(w, http.StatusForbidden, "Forbidden")
		return
	}
	if errors.Is(err, services.ErrScoreExists) {
		response.JSON(w, http.StatusUnauthorized, map[string]string{"status": "invalid", "message": "Score exist"})
		return
	}
	if err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"id": item.ID, "film": map[string]any{"id": item.FilmID, "name": item.FilmName}, "score": item.Score, "created_at": utils.TimeISO(item.CreatedAt)})
}

func (h *Handler) UserRatings(w http.ResponseWriter, r *http.Request, userID int) {
	items, err := h.service.UserRatings(r.Context(), userID)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	ratings := []map[string]any{}
	for _, item := range items {
		ratings = append(ratings, map[string]any{"id": item.ID, "film": map[string]any{"id": item.FilmID, "name": item.FilmName}, "score": item.Score, "created_at": utils.TimeISO(item.CreatedAt)})
	}
	response.JSON(w, http.StatusOK, map[string]any{"ratings": ratings})
}

func (h *Handler) DeleteRating(w http.ResponseWriter, r *http.Request, userID, id int) {
	err := h.service.DeleteRating(r.Context(), middleware.UserID(r), userID, id)
	if errors.Is(err, services.ErrForbidden) {
		response.ErrorJSON(w, http.StatusForbidden, "Forbidden")
		return
	}
	if err != nil {
		response.ErrorJSON(w, http.StatusNotFound, "Rating not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(value string) (int, bool) {
	id, err := strconv.Atoi(value)
	return id, err == nil && id > 0
}
func first(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}
