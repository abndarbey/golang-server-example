package resolvers

import (
	"context"
	"fmt"
	"orijinplus/app/api/dataloaders"
	"orijinplus/app/api/graphql/generated/graph"
	"orijinplus/app/models"

	"github.com/volatiletech/null"
)

// User returns graph.UserResolver implementation.
func (r *Resolver) User() graph.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

func (r *userResolver) UserType(ctx context.Context, obj *models.User) (string, error) {
	if obj.IsAdmin {
		return "Super Admin", nil
	}
	if obj.IsMember {
		return "Member", nil
	}
	if obj.IsCustomer {
		return "Customer", nil
	}
	return "", nil
}

func (r *userResolver) Organization(ctx context.Context, obj *models.User) (*models.Organization, error) {
	return dataloaders.OrganizationLoaderFromContext(ctx, obj.OrganizationID.Int64)
}

func (r *userResolver) Role(ctx context.Context, obj *models.User) (*models.Role, error) {
	return dataloaders.RoleLoaderFromContext(ctx, obj.RoleID.Int64)
}

func (r *userResolver) Profile(ctx context.Context, obj *models.User) (*models.Profile, error) {
	return dataloaders.ProfileLoaderFromContext(ctx, obj.ID)
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Users(
	ctx context.Context,
	search graph.SearchFilter,
	limit int,
	offset int,
	isAdmin, isMember, isCustomer bool,
	organizationID *int64,
) (*graph.UserResult, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadUser, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	var orgID null.Int64
	if organizationID != nil {
		orgID.Valid = true
		orgID.Int64 = *organizationID
	}

	users, err := r.services.UserService.List(ctx, isAdmin, isMember, isCustomer, orgID, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	userResult := &graph.UserResult{
		Users: users,
		Total: len(users),
	}

	return userResult, nil
}

func (r *queryResolver) User(ctx context.Context, id *int64, email, phone *string) (*models.User, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadUser, true, true); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	if id != nil {
		result, err := r.services.UserService.GetByID(ctx, *id, auther)
		if err != nil {
			return nil, fmt.Errorf(err.Message)
		}
		return result, nil
	}

	if email != nil {
		result, err := r.services.UserService.GetByEmail(ctx, *email, auther)
		if err != nil {
			return nil, fmt.Errorf(err.Message)
		}
		return result, nil
	}

	if phone != nil {
		result, err := r.services.UserService.GetByPhone(ctx, *phone, auther)
		if err != nil {
			return nil, fmt.Errorf(err.Message)
		}
		return result, nil
	}

	return nil, fmt.Errorf("no query parameters provided")
}

///////////////
// Mutations //
///////////////

func (r *mutationResolver) ChangePassword(ctx context.Context, oldPassword string, password string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) ChangeDetails(ctx context.Context, input graph.UpdateUser) (*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) UserUpdate(ctx context.Context, id int64, input graph.UpdateUser) (*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) ForgotPassword(ctx context.Context, email string, viaSms *bool) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) ResetPassword(ctx context.Context, token string, password string, email *null.String) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) ResendEmailVerification(ctx context.Context, email string) (bool, error) {
	panic("not implemented")
}
