package handlers

import (
	"orijinplus/app/services"
	"orijinplus/app/store/filestore"
)

type Handlers struct {
	AuthHandler    *AuthHandler
	GraphQLHandler *GraphQLHandler
}

func NewHandlers(s *services.Services, fs *filestore.FileStore) *Handlers {
	return &Handlers{
		NewAuthHandler(s),
		NewGraphQLHandler(s, fs),
	}
}
