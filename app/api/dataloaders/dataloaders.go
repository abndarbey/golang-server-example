package dataloaders

import (
	"context"
	"net/http"
	"orijinplus/app/models"
	"orijinplus/app/store/dbstore"
	"time"
)

// ContextKey holds a custom String func for uniqueness
type ContextKey string

func (k ContextKey) String() string {
	return "dataloader_" + string(k)
}

// UserLoaderKey declares a statically typed key for context reference in other packages
const UserLoaderKey ContextKey = "user_loader"

// ProfileLoaderKey declares a statically typed key for context reference in other packages
const ProfileLoaderKey ContextKey = "profile_loader"

// OrganizationLoaderKey declares a statically typed key for context reference in other packages
const OrganizationLoaderKey ContextKey = "organization_loader"

// RoleLoaderKey declares a statically typed key for context reference in other packages
const RoleLoaderKey ContextKey = "role_loader"

// ContainerLoaderKey declares a statically typed key for context reference in other packages
const ContainerLoaderKey ContextKey = "container_loader"

// PalletLoaderKey declares a statically typed key for context reference in other packages
const PalletLoaderKey ContextKey = "pallet_loader"

// UserLoaderFromContext runs the dataloader inside the context
func UserLoaderFromContext(ctx context.Context, id int64) (*models.User, error) {
	return ctx.Value(UserLoaderKey).(*UserLoader).Load(id)
}

// ProfileLoaderFromContext runs the dataloader inside the context
func ProfileLoaderFromContext(ctx context.Context, id int64) (*models.Profile, error) {
	return ctx.Value(ProfileLoaderKey).(*ProfileLoader).Load(id)
}

// OrganizationLoaderFromContext runs the dataloader inside the context
func OrganizationLoaderFromContext(ctx context.Context, id int64) (*models.Organization, error) {
	return ctx.Value(OrganizationLoaderKey).(*OrganizationLoader).Load(id)
}

// RoleLoaderFromContext runs the dataloader inside the context
func RoleLoaderFromContext(ctx context.Context, id int64) (*models.Role, error) {
	return ctx.Value(RoleLoaderKey).(*RoleLoader).Load(id)
}

// ContainerLoaderFromContext runs the dataloader inside the context
func ContainerLoaderFromContext(ctx context.Context, id int64) (*models.Container, error) {
	return ctx.Value(ContainerLoaderKey).(*ContainerLoader).Load(id)
}

// PalletLoaderFromContext runs the dataloader inside the context
func PalletLoaderFromContext(ctx context.Context, id int64) (*models.Pallet, error) {
	return ctx.Value(PalletLoaderKey).(*PalletLoader).Load(id)
}

// WithDataloaders returns a new context that contains dataloaders
func WithDataloaders(
	ctx context.Context,
	dbstore *dbstore.DBStore,
) context.Context {
	userLoader := NewUserLoader(
		UserLoaderConfig{
			Fetch: func(ids []int64) ([]*models.User, []error) {
				data, err := dbstore.UserStore.GetMany(ctx, ids)
				if err != nil {
					return nil, []error{err}
				}

				slice := make(map[interface{}]*models.User, len(data))
				for _, e := range data {
					slice[e.ID] = e
				}

				result := make([]*models.User, len(ids))
				for i, key := range ids {
					result[i] = slice[key]
				}

				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	profileLoader := NewProfileLoader(
		ProfileLoaderConfig{
			Fetch: func(ids []int64) ([]*models.Profile, []error) {
				data, err := dbstore.ProfileStore.GetMany(ctx, ids)
				if err != nil {
					return nil, []error{err}
				}

				// make result and ids of the same order
				slice := make(map[interface{}]*models.Profile, len(data))
				for _, e := range data {
					slice[e.UserID] = e
				}

				result := make([]*models.Profile, len(ids))
				for i, key := range ids {
					result[i] = slice[key]
				}

				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	organizationLoader := NewOrganizationLoader(
		OrganizationLoaderConfig{
			Fetch: func(ids []int64) ([]*models.Organization, []error) {
				data, err := dbstore.OrganizationStore.GetMany(ctx, ids)
				if err != nil {
					return nil, []error{err}
				}

				slice := make(map[interface{}]*models.Organization, len(data))
				for _, e := range data {
					slice[e.ID] = e
				}

				result := make([]*models.Organization, len(ids))
				for i, key := range ids {
					result[i] = slice[key]
				}

				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	roleLoader := NewRoleLoader(
		RoleLoaderConfig{
			Fetch: func(ids []int64) ([]*models.Role, []error) {
				data, err := dbstore.RoleStore.GetMany(ctx, ids)
				if err != nil {
					return nil, []error{err}
				}

				slice := make(map[interface{}]*models.Role, len(data))
				for _, e := range data {
					slice[e.ID] = e
				}

				result := make([]*models.Role, len(ids))
				for i, key := range ids {
					result[i] = slice[key]
				}

				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	containerLoader := NewContainerLoader(
		ContainerLoaderConfig{
			Fetch: func(ids []int64) ([]*models.Container, []error) {
				data, err := dbstore.ContainerStore.GetMany(ctx, ids)
				if err != nil {
					return nil, []error{err}
				}

				// make result and ids of the same order
				slice := make(map[interface{}]*models.Container, len(data))
				for _, e := range data {
					slice[e.ID] = e
				}

				result := make([]*models.Container, len(ids))
				for i, key := range ids {
					result[i] = slice[key]
				}

				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	palletLoader := NewPalletLoader(
		PalletLoaderConfig{
			Fetch: func(ids []int64) ([]*models.Pallet, []error) {
				data, err := dbstore.PalletStore.GetMany(ctx, ids)
				if err != nil {
					return nil, []error{err}
				}

				// make result and ids of the same order
				slice := make(map[interface{}]*models.Pallet, len(data))
				for _, e := range data {
					slice[e.ID] = e
				}

				result := make([]*models.Pallet, len(ids))
				for i, key := range ids {
					result[i] = slice[key]
				}

				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	ctx = context.WithValue(ctx, UserLoaderKey, userLoader)
	ctx = context.WithValue(ctx, ProfileLoaderKey, profileLoader)
	ctx = context.WithValue(ctx, OrganizationLoaderKey, organizationLoader)
	ctx = context.WithValue(ctx, RoleLoaderKey, roleLoader)
	ctx = context.WithValue(ctx, ContainerLoaderKey, containerLoader)
	ctx = context.WithValue(ctx, PalletLoaderKey, palletLoader)
	return ctx
}

// DataloaderMiddleware runs before each API call and loads the dataloaders into context
func DataloaderMiddleware(
	dbstore *dbstore.DBStore,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(WithDataloaders(r.Context(), dbstore))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
