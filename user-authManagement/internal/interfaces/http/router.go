// Package http wires chi routes onto the handlers and middleware.
package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
	"github.com/lindritprekaj/user-authmanagement/internal/interfaces/http/handlers"
	"github.com/lindritprekaj/user-authmanagement/internal/interfaces/http/middleware"
)

// RouterDeps groups handler and policy dependencies for NewRouter.
type RouterDeps struct {
	Logger              *slog.Logger
	TokenService        ports.TokenService
	AuthHandler         *handlers.AuthHandler
	UserHandler         *handlers.UserHandler
	AdminHandler        *handlers.AdminHandler
	HealthHandler       *handlers.HealthHandler
	CORSAllowedOrigins  []string
	LoginRateLimit      int
	RegisterRateLimit   int
}

// NewRouter assembles the HTTP routes.
func NewRouter(d RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(middleware.RequestID)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.Recovery(d.Logger))
	r.Use(middleware.Logging(d.Logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   d.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", middleware.HeaderRequestID},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/healthz", d.HealthHandler.Liveness)
	r.Get("/readyz", d.HealthHandler.Readiness)

	loginRL := middleware.NewRateLimit(d.LoginRateLimit)
	registerRL := middleware.NewRateLimit(d.RegisterRateLimit)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.With(registerRL.Middleware).Post("/register", d.AuthHandler.Register)
			r.With(loginRL.Middleware).Post("/login", d.AuthHandler.Login)
			r.Post("/refresh", d.AuthHandler.Refresh)
			r.Post("/password/reset", d.AuthHandler.RequestPasswordReset)
			r.Post("/password/confirm", d.AuthHandler.ConfirmPasswordReset)

			r.Group(func(r chi.Router) {
				r.Use(middleware.Authn(d.TokenService))
				r.Post("/logout", d.AuthHandler.Logout)
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.Authn(d.TokenService))
			r.Get("/me", d.UserHandler.Me)
			r.Patch("/me", d.UserHandler.UpdateMe)
			r.Delete("/me", d.UserHandler.DeleteMe)
		})
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.Authn(d.TokenService))
			r.Use(middleware.RequireRole("admin"))
			r.Get("/users", d.AdminHandler.ListUsers)
			r.Get("/users/{id}", d.AdminHandler.GetUser)
			r.Patch("/users/{id}/roles", d.AdminHandler.SetRoles)
			r.Delete("/users/{id}", d.AdminHandler.DeleteUser)
		})
	})

	return r
}
