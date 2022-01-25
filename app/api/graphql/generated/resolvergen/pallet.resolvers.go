package resolvergen

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"orijinplus/app/api/graphql/generated/graph"
	"orijinplus/app/models"
)

func (r *mutationResolver) PalletCreate(ctx context.Context, input graph.UpdatePallet) (*models.Pallet, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PalletUpdate(ctx context.Context, id int64, input graph.UpdatePallet) (*models.Pallet, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PalletArchive(ctx context.Context, id int64) (*models.Pallet, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) PalletUnarchive(ctx context.Context, id int64) (*models.Pallet, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *palletResolver) UID(ctx context.Context, obj *models.Pallet) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *palletResolver) Container(ctx context.Context, obj *models.Pallet) (*models.Container, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *palletResolver) Organization(ctx context.Context, obj *models.Pallet) (*models.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Pallets(ctx context.Context, search graph.SearchFilter, limit int, offset int, containerID *int64) (*graph.PalletResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PalletByID(ctx context.Context, id int64) (*models.Pallet, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PalletByUID(ctx context.Context, uid string) (*models.Pallet, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PalletByCode(ctx context.Context, code string) (*models.Pallet, error) {
	panic(fmt.Errorf("not implemented"))
}

// Pallet returns graph.PalletResolver implementation.
func (r *Resolver) Pallet() graph.PalletResolver { return &palletResolver{r} }

type palletResolver struct{ *Resolver }
