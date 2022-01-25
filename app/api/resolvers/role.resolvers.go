package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"orijinplus/app/api/dataloaders"
	"orijinplus/app/api/graphql/generated/graph"
	"orijinplus/app/models"

	"github.com/volatiletech/null"
)

type roleResolver struct{ *Resolver }

// Role returns graph.RoleResolver implementation.
func (r *Resolver) Role() graph.RoleResolver { return &roleResolver{r} }

func (r *roleResolver) Organization(ctx context.Context, obj *models.Role) (*models.Organization, error) {
	return dataloaders.OrganizationLoaderFromContext(ctx, obj.OrganizationID)
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Roles(ctx context.Context, search graph.SearchFilter, limit int, offset int, organizationID *int64) (*graph.RolesResult, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadRole, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	orgID := null.Int64{}
	if organizationID != nil {
		orgID.Valid = true
		orgID.Int64 = *organizationID
	}

	roles, err := r.services.RoleService.List(ctx, orgID, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	obj := &graph.RolesResult{
		Roles: roles,
		Total: len(roles),
	}

	return obj, nil
}

func (r *queryResolver) Role(ctx context.Context, id *int64, code *string) (*models.Role, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadRole, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	if id != nil {
		obj, err := r.services.RoleService.GetByID(ctx, *id, auther)
		if err != nil {
			return nil, fmt.Errorf(err.Message)
		}
		return obj, nil
	}

	if code != nil {
		obj, err := r.services.RoleService.GetByCode(ctx, *code, auther)
		if err != nil {
			return nil, fmt.Errorf(err.Message)
		}
		return obj, nil
	}

	return nil, fmt.Errorf("no query parameters provided")
}

///////////////
// Mutations //
///////////////

func (r *mutationResolver) RoleCreate(ctx context.Context, input graph.NewRole) (*models.Role, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.CreateRole, true, false); err != nil {
		return nil, fmt.Errorf(err.Error)
	}

	var orgID int64
	rolePermissions := []string{}

	if input.OrganizationID != nil {
		orgID = *input.OrganizationID
	}

	if !input.IsOrgAdmin && len(input.Permissions) > 0 {
		permissions := models.ListPermissions()
		requestPermissions := r.services.RoleService.UniquePermissions(input.Permissions)
		for _, perm := range requestPermissions {
			ok := r.services.RoleService.ContainsPermission(permissions, perm)
			if ok {
				rolePermissions = append(rolePermissions, perm)
			}
		}
	}

	request := models.RoleCreateRequest{
		Name:           input.Name,
		IsOrgAdmin:     input.IsOrgAdmin,
		OrganizationID: orgID,
		Permissions:    rolePermissions,
	}

	role, err := r.services.RoleService.Create(ctx, request, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return role, nil
}

func (r *mutationResolver) RoleUpdate(ctx context.Context, id int64, input graph.UpdateRole) (*models.Role, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.CreateRole, true, false); err != nil {
		return nil, fmt.Errorf(err.Error)
	}

	rolePermissions := []string{}

	if len(input.Permissions) > 0 {
		permissions := models.ListPermissions()
		requestPermissions := r.services.RoleService.UniquePermissions(input.Permissions)
		for _, perm := range requestPermissions {
			ok := r.services.RoleService.ContainsPermission(permissions, perm)
			if ok {
				rolePermissions = append(rolePermissions, perm)
			}
		}
	}

	request := models.RoleUpdateRequest{ID: id}

	if input.Name != nil {
		request.Name = input.Name.String
	}
	if input.IsArchived != nil {
		request.IsArchived = input.IsArchived.Bool
	}
	request.Permissions = rolePermissions

	role, err := r.services.RoleService.Update(ctx, request, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return role, nil
}
