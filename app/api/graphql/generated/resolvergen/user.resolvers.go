package resolvergen

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"orijinplus/app/api/graphql/generated/graph"
	"orijinplus/app/models"

	"github.com/volatiletech/null"
)

func (r *mutationResolver) ChangePassword(ctx context.Context, oldPassword string, password string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChangeDetails(ctx context.Context, input graph.UpdateUser) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UserUpdate(ctx context.Context, id int64, input graph.UpdateUser) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ForgotPassword(ctx context.Context, email string, viaSms *bool) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ResetPassword(ctx context.Context, token string, password string, email *null.String) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ResendEmailVerification(ctx context.Context, email string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context, search graph.SearchFilter, limit int, offset int, isAdmin bool, isMember bool, isCustomer bool, organizationID *int64) (*graph.UserResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, id *int64, email *string, phone *string) (*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) UserType(ctx context.Context, obj *models.User) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Organization(ctx context.Context, obj *models.User) (*models.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Role(ctx context.Context, obj *models.User) (*models.Role, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Profile(ctx context.Context, obj *models.User) (*models.Profile, error) {
	panic(fmt.Errorf("not implemented"))
}

// User returns graph.UserResolver implementation.
func (r *Resolver) User() graph.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
