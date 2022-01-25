package resolvergen

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"orijinplus/app/api/graphql/generated/graph"
	"orijinplus/app/models"
)

func (r *mutationResolver) RoleCreate(ctx context.Context, input graph.NewRole) (*models.Role, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RoleUpdate(ctx context.Context, id int64, input graph.UpdateRole) (*models.Role, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Roles(ctx context.Context, search graph.SearchFilter, limit int, offset int, organizationID *int64) (*graph.RolesResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Role(ctx context.Context, id *int64, code *string) (*models.Role, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *roleResolver) Organization(ctx context.Context, obj *models.Role) (*models.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

// Role returns graph.RoleResolver implementation.
func (r *Resolver) Role() graph.RoleResolver { return &roleResolver{r} }

type roleResolver struct{ *Resolver }
