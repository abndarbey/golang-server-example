package resolvers

import (
	"context"
	"fmt"
	"orijinplus/app/api/graphql/generated/graph"
	"orijinplus/app/models"
)

///////////////
//   Query   //
///////////////

func (r *queryResolver) Organizations(ctx context.Context, search graph.SearchFilter, limit int, offset int) (*graph.OrganizationsResult, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if !auther.IsAdmin {
		return nil, fmt.Errorf("unauthorized request")
	}
	orgs, err := r.services.OrganizationService.List(ctx)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	result := &graph.OrganizationsResult{
		Organizations: orgs,
		Total:         len(orgs),
	}

	return result, nil
}

func (r *queryResolver) Organization(ctx context.Context, id *int64, code *string) (*models.Organization, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadOrganization, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	if id != nil {
		result, err := r.services.OrganizationService.GetByID(ctx, *id, auther)
		if err != nil {
			return nil, fmt.Errorf(err.Message)
		}
		return result, nil
	}

	if code != nil {
		result, err := r.services.OrganizationService.GetByCode(ctx, *code, auther)
		if err != nil {
			return nil, fmt.Errorf(err.Message)
		}
		return result, nil
	}

	return nil, fmt.Errorf("no query parameters provided")
}

func (r *queryResolver) OrganizationByID(ctx context.Context, id int64) (*models.Organization, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadOrganization, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	obj, err := r.services.OrganizationService.GetByID(ctx, id, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}
	return obj, nil
}

func (r *queryResolver) OrganizationByCode(ctx context.Context, code string) (*models.Organization, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadOrganization, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	obj, err := r.services.OrganizationService.GetByCode(ctx, code, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}
	return obj, nil
}

///////////////
// Mutations //
///////////////

func (r *mutationResolver) OrganizationUpdate(ctx context.Context, id int64, input graph.UpdateOrganization) (*models.Organization, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.UpdateOrganization, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	update := models.Organization{
		ID:      id,
		Name:    input.Name.String,
		Website: *input.Website,
	}

	result, err := r.services.OrganizationService.Update(ctx, update, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}
	return result, nil
}
