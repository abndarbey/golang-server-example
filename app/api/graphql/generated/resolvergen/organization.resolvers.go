package resolvergen

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"orijinplus/app/api/graphql/generated/graph"
	"orijinplus/app/models"
)

func (r *mutationResolver) OrganizationUpdate(ctx context.Context, id int64, input graph.UpdateOrganization) (*models.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Organizations(ctx context.Context, search graph.SearchFilter, limit int, offset int) (*graph.OrganizationsResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Organization(ctx context.Context, id *int64, code *string) (*models.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) OrganizationByID(ctx context.Context, id int64) (*models.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) OrganizationByCode(ctx context.Context, code string) (*models.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}
