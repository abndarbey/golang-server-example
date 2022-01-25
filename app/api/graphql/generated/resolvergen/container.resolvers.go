package resolvergen

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"orijinplus/app/api/graphql/generated/graph"
	"orijinplus/app/models"
)

func (r *containerResolver) UID(ctx context.Context, obj *models.Container) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *containerResolver) Organization(ctx context.Context, obj *models.Container) (*models.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ContainerCreate(ctx context.Context, input graph.UpdateContainer) (*models.Container, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ContainerUpdate(ctx context.Context, id int64, input graph.UpdateContainer) (*models.Container, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ContainerArchive(ctx context.Context, id int64) (*models.Container, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ContainerUnarchive(ctx context.Context, id int64) (*models.Container, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Containers(ctx context.Context, search graph.SearchFilter, limit int, offset int) (*graph.ContainerResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ContainerByID(ctx context.Context, id int64) (*models.Container, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ContainerByUID(ctx context.Context, uid string) (*models.Container, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ContainerByCode(ctx context.Context, code string) (*models.Container, error) {
	panic(fmt.Errorf("not implemented"))
}

// Container returns graph.ContainerResolver implementation.
func (r *Resolver) Container() graph.ContainerResolver { return &containerResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type containerResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
