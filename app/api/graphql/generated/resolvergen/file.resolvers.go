package resolvergen

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"orijinplus/app/api/graphql/generated/graph"
	"orijinplus/app/models"

	"github.com/99designs/gqlgen/graphql"
)

func (r *mutationResolver) FileUpload(ctx context.Context, file graphql.Upload) (*models.File, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) FileUploadMultiple(ctx context.Context, files []graphql.Upload) ([]models.File, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
