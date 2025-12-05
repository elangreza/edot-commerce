package rest

import (
	"github.com/go-chi/chi/v5"
)

func NewHandlerWithMiddleware(
	publicRoute chi.Router,
	authService AuthService,
) {

	authMiddleware := AuthMiddleware{
		svc: authService,
	}

	publicRoute.Group(func(r chi.Router) {
		r.Use(authMiddleware.MustAuthMiddleware())
		// r.Get("/profile", profileHandler.ProfileUserHandler)
	})
}
