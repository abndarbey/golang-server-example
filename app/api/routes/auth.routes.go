package routes

import (
	"github.com/go-chi/chi"
)

// AuthRoutes Routes function
func (rt *Routes) AuthRoutes(r chi.Router) {
	h := rt.Handlers.AuthHandler

	r.Route("/auth", func(r chi.Router) {
		r.Get("/permissions", h.ListPermissions)
		r.Post("/admin/login", h.LoginAdmin)
		r.Post("/member/login", h.LoginMember)
		r.Post("/customer/login", h.LoginCustomer)
		r.Post("/admin/register", h.RegisterAdmin)
		r.Post("/member/register", h.RegisterMember)
		r.Post("/customer/register", h.RegisterCustomer)
		r.Post("/organization/register", h.RegisterOrganization)
		r.Get("/logout", h.Logout)
	})
}
