package resolvers

import (
	"context"
	"fmt"
	"orijinplus/app/api/dataloaders"
	"orijinplus/app/api/graphql/generated/graph"
	"orijinplus/app/models"

	"github.com/gofrs/uuid"
)

type palletResolver struct{ *Resolver }

// Pallet returns graph.PalletResolver implementation.
func (r *Resolver) Pallet() graph.PalletResolver { return &palletResolver{r} }

func (r *palletResolver) UID(ctx context.Context, obj *models.Pallet) (string, error) {
	return obj.UID.String(), nil
}

func (r *palletResolver) Container(ctx context.Context, obj *models.Pallet) (*models.Container, error) {
	if obj.ContainerID.Valid {
		return dataloaders.ContainerLoaderFromContext(ctx, obj.ContainerID.Int64)
	}
	return nil, nil
}

func (r *palletResolver) Organization(ctx context.Context, obj *models.Pallet) (*models.Organization, error) {
	if obj.OrganizationID.Valid {
		return dataloaders.OrganizationLoaderFromContext(ctx, obj.OrganizationID.Int64)
	}
	return nil, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Pallets(
	ctx context.Context,
	search graph.SearchFilter,
	limit int,
	offset int,
	containerID *int64,
) (*graph.PalletResult, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadPallet, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	if containerID != nil && *containerID != 0 {
		pallets, err := r.services.PalletService.List(ctx, auther)
		if err != nil {
			return nil, fmt.Errorf(err.Message)
		}
		return &graph.PalletResult{Pallets: pallets, Total: len(pallets)}, nil
	}

	pallets, err := r.services.PalletService.List(ctx, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}
	return &graph.PalletResult{Pallets: pallets, Total: len(pallets)}, nil
}

func (r *queryResolver) PalletByID(ctx context.Context, id int64) (*models.Pallet, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadPallet, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	obj, err := r.services.PalletService.GetByID(ctx, id, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}

func (r *queryResolver) PalletByUID(ctx context.Context, uid string) (*models.Pallet, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadPallet, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	objUUID, uuidErr := uuid.FromString(uid)
	if uuidErr != nil {
		return nil, fmt.Errorf("invalid uid")
	}

	obj, err := r.services.PalletService.GetByUID(ctx, objUUID, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}

func (r *queryResolver) PalletByCode(ctx context.Context, code string) (*models.Pallet, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadPallet, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	obj, err := r.services.PalletService.GetByCode(ctx, code, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}

///////////////
// Mutations //
///////////////

func (r *mutationResolver) PalletCreate(ctx context.Context, input graph.UpdatePallet) (*models.Pallet, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.CreatePallet, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	request := models.PalletRequest{}
	if input.Description != nil {
		request.Description = input.Description.String
	}
	if input.ContainerID != nil {
		request.ContainerID = *input.ContainerID
	}
	if input.OrganizationID != nil {
		request.OrganizationID = *input.OrganizationID
	}

	obj, err := r.services.PalletService.Create(ctx, request, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}

func (r *mutationResolver) PalletUpdate(ctx context.Context, id int64, input graph.UpdatePallet) (*models.Pallet, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.UpdatePallet, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	request := models.PalletRequest{}
	if input.Description != nil {
		request.Description = input.Description.String
	}
	if input.ContainerID != nil {
		request.ContainerID = *input.ContainerID
	}

	obj, err := r.services.PalletService.Create(ctx, request, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}

func (r *mutationResolver) PalletArchive(ctx context.Context, id int64) (*models.Pallet, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.UpdatePallet, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	obj, err := r.services.PalletService.Archive(ctx, id, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}

func (r *mutationResolver) PalletUnarchive(ctx context.Context, id int64) (*models.Pallet, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.UpdatePallet, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	obj, err := r.services.PalletService.Unarchive(ctx, id, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}
