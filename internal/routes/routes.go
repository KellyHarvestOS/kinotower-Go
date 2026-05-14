package routes

import (
	"net/http"

	"github.com/KellyHarvestOS/kinotower-Go/internal/handlers"
	"github.com/KellyHarvestOS/kinotower-Go/internal/middleware"
)

func New(h *handlers.Handler, authMiddleware middleware.Middleware) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/kinotower", h.KinotowerHome)
	mux.HandleFunc("/kinotower/", h.KinotowerHome)
	mux.Handle("/web/static/", http.StripPrefix("/web/static/", http.FileServer(http.Dir("web/static"))))
	mux.HandleFunc("/api/v1/films", h.Films)
	mux.HandleFunc("/api/v1/films/", h.FilmRoutes)
	mux.HandleFunc("/api/v1/categories", h.Categories)
	mux.HandleFunc("/api/v1/countries", h.Countries)
	mux.HandleFunc("/api/v1/genders", h.Genders)
	mux.HandleFunc("/api/v1/auth/signup", h.Signup)
	mux.HandleFunc("/api/v1/auth/signin", h.Signin)
	mux.Handle("/api/v1/auth/signout", authMiddleware(http.HandlerFunc(h.Signout)))
	mux.Handle("/api/v1/users", authMiddleware(http.HandlerFunc(h.CurrentUser)))
	mux.Handle("/api/v1/users/", authMiddleware(http.HandlerFunc(h.UserRoutes)))
	return mux
}
