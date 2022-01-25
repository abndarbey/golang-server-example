package resolvers

import (
	"context"
	"fmt"
	"orijinplus/app/api/authentication"
	"orijinplus/app/api/graphql/generated/graph"
	"orijinplus/app/models"
	"orijinplus/app/services"
	"orijinplus/app/store/filestore"
)

type Resolver struct {
	services  *services.Services
	filestore *filestore.FileStore
}

func NewResolver(s *services.Services, fs *filestore.FileStore) *Resolver {
	return &Resolver{s, fs}
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func (r *Resolver) GetAuther(ctx context.Context) (*models.Auther, error) {
	auther := authentication.AutherFromContext(ctx)
	if auther == nil {
		return nil, fmt.Errorf("no credentials provided")
	}

	return auther, nil
}
