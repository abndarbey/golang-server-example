package routes

import (
	"orijinplus/app/api/authentication"

	"github.com/go-chi/chi"
)

// GraphQL Routes function
func (rt *Routes) GraphQL(r chi.Router) {
	h := rt.Handlers.GraphQLHandler

	r.Route("/gql", func(r chi.Router) {
		r.Use(authentication.Middleware())
		r.Handle("/", h.Playground())
		r.Handle("/query", h.Query())
	})
}
