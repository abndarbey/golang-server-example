package resolvers

import (
	"context"
	"fmt"
	"orijinplus/app/api/dataloaders"
	"orijinplus/app/api/graphql/generated/graph"
	"orijinplus/app/models"

	"github.com/gofrs/uuid"
)

type containerResolver struct{ *Resolver }

// Container returns graph.ContainerResolver implementation.
func (r *Resolver) Container() graph.ContainerResolver { return &containerResolver{r} }

func (r *containerResolver) UID(ctx context.Context, obj *models.Container) (string, error) {
	return obj.UID.String(), nil
}

func (r *containerResolver) Organization(ctx context.Context, obj *models.Container) (*models.Organization, error) {
	if obj.OrganizationID.Valid {
		return dataloaders.OrganizationLoaderFromContext(ctx, obj.OrganizationID.Int64)
	}
	return nil, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Containers(ctx context.Context, search graph.SearchFilter, limit int, offset int) (*graph.ContainerResult, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadContainer, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	containers, err := r.services.ContainerService.List(ctx, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}
	return &graph.ContainerResult{Containers: containers, Total: len(containers)}, nil
}

func (r *queryResolver) ContainerByID(ctx context.Context, id int64) (*models.Container, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadContainer, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	obj, err := r.services.ContainerService.GetByID(ctx, id, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}

func (r *queryResolver) ContainerByUID(ctx context.Context, uid string) (*models.Container, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadContainer, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	objUUID, uuidErr := uuid.FromString(uid)
	if uuidErr != nil {
		return nil, fmt.Errorf("invalid uid")
	}

	obj, err := r.services.ContainerService.GetByUID(ctx, objUUID, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}

func (r *queryResolver) ContainerByCode(ctx context.Context, code string) (*models.Container, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadContainer, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	obj, err := r.services.ContainerService.GetByCode(ctx, code, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}

///////////////
// Mutations //
///////////////

func (r *mutationResolver) ContainerCreate(ctx context.Context, input graph.UpdateContainer) (*models.Container, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.CreateContainer, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	request := models.ContainerRequest{}
	if input.Description != nil {
		request.Description = input.Description.String
	}
	if input.OrganizationID != nil {
		request.OrganizationID = *input.OrganizationID
	}

	obj, err := r.services.ContainerService.Create(ctx, request, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}

func (r *mutationResolver) ContainerUpdate(ctx context.Context, id int64, input graph.UpdateContainer) (*models.Container, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.UpdateContainer, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	request := models.ContainerRequest{}
	if input.Description != nil {
		request.Description = input.Description.String
	}

	obj, err := r.services.ContainerService.Update(ctx, id, request, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}

func (r *mutationResolver) ContainerArchive(ctx context.Context, id int64) (*models.Container, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.UpdateContainer, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	obj, err := r.services.ContainerService.Archive(ctx, id, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}

func (r *mutationResolver) ContainerUnarchive(ctx context.Context, id int64) (*models.Container, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.UpdateContainer, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	obj, err := r.services.ContainerService.Unarchive(ctx, id, auther)
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	return obj, nil
}
