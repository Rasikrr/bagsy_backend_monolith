package auth

import "github.com/go-chi/chi/v5"

// RegisterRoutes mounts auth registration endpoints onto the given router.
//
//	POST /api/v1/auth/register        — start registration, send OTP
//	POST /api/v1/auth/register/verify  — confirm OTP, create entities
//	POST /api/v1/auth/register/resend  — resend OTP code
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/register/verify", h.Verify)
		r.Post("/register/resend", h.Resend)
	})
}
