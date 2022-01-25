package services

import (
	"context"
	"orijinplus/app/master"
	"orijinplus/app/models"
	"orijinplus/app/store/dbstore"
	"orijinplus/utils/faulterr"

	"github.com/gofrs/uuid"
)

type ContainerService struct {
	dbstore *dbstore.DBStore
	master  *master.Master
}

var _ ContainerServiceInterface = &ContainerService{}

type ContainerServiceInterface interface {
	List(ctx context.Context, auther *models.Auther) ([]models.Container, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64, auther *models.Auther) (*models.Container, *faulterr.FaultErr)
	GetByUID(ctx context.Context, uid uuid.UUID, auther *models.Auther) (*models.Container, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string, auther *models.Auther) (*models.Container, *faulterr.FaultErr)
	Create(ctx context.Context, request models.ContainerRequest, auther *models.Auther) (*models.Container, *faulterr.FaultErr)
	Update(ctx context.Context, id int64, request models.ContainerRequest, auther *models.Auther) (*models.Container, *faulterr.FaultErr)
	Delete(ctx context.Context, id int64, auther *models.Auther) *faulterr.FaultErr
}

func NewContainerService(s *dbstore.DBStore, m *master.Master) *ContainerService {
	return &ContainerService{s, m}
}

// List gets all skus
func (s *ContainerService) List(ctx context.Context, auther *models.Auther) ([]models.Container, *faulterr.FaultErr) {
	if auther.IsAdmin {
		return s.dbstore.ContainerStore.List(ctx)
	}
	return s.dbstore.ContainerStore.ListByOrgID(ctx, auther.OrganizationID.Int64)
}

func (s *ContainerService) GetByID(ctx context.Context, id int64, auther *models.Auther) (*models.Container, *faulterr.FaultErr) {
	obj, err := s.dbstore.ContainerStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !auther.IsAdmin && auther.OrganizationID.Int64 != obj.OrganizationID.Int64 {
		return nil, faulterr.NewNotFoundError("no object found")
	}
	return obj, nil
}

func (s *ContainerService) GetByUID(ctx context.Context, uid uuid.UUID, auther *models.Auther) (*models.Container, *faulterr.FaultErr) {
	obj, err := s.dbstore.ContainerStore.GetByUID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if !auther.IsAdmin && auther.OrganizationID.Int64 != obj.OrganizationID.Int64 {
		return nil, faulterr.NewNotFoundError("no object found")
	}
	return obj, nil
}

func (s *ContainerService) GetByCode(ctx context.Context, code string, auther *models.Auther) (*models.Container, *faulterr.FaultErr) {
	obj, err := s.dbstore.ContainerStore.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if !auther.IsAdmin && auther.OrganizationID.Int64 != obj.OrganizationID.Int64 {
		return nil, faulterr.NewNotFoundError("no object found")
	}
	return obj, nil
}

// Create gets all skus
func (s *ContainerService) Create(ctx context.Context, r models.ContainerRequest, auther *models.Auther) (*models.Container, *faulterr.FaultErr) {
	if auther.IsAdmin && !r.OrganizationID.Valid {
		return nil, faulterr.NewBadRequestError("organization id is required")
	}
	// Reassign organization ID to the request
	if !auther.IsAdmin {
		r.OrganizationID = auther.OrganizationID
	}
	createdByID := auther.ID

	// Start transactions
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	obj, err := s.master.ContainerMaster.Create(ctx, tx, r, createdByID)
	if err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *ContainerService) Update(ctx context.Context, id int64, request models.ContainerRequest, auther *models.Auther) (*models.Container, *faulterr.FaultErr) {
	current, err := s.GetByID(ctx, id, auther)
	if err != nil {
		return nil, err
	}

	// Start transactions
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	container, err := s.master.ContainerMaster.Update(ctx, tx, current, request)
	if err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return container, nil
}

func (s *ContainerService) Archive(ctx context.Context, id int64, auther *models.Auther) (*models.Container, *faulterr.FaultErr) {
	container, err := s.GetByID(ctx, id, auther)
	if err != nil {
		return nil, err
	}
	container.IsArchived = true

	// Start transactions
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.ContainerStore.Update(ctx, tx, *container); err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return container, nil
}

func (s *ContainerService) Unarchive(ctx context.Context, id int64, auther *models.Auther) (*models.Container, *faulterr.FaultErr) {
	container, err := s.GetByID(ctx, id, auther)
	if err != nil {
		return nil, err
	}
	container.IsArchived = false

	// Start transactions
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.ContainerStore.Update(ctx, tx, *container); err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return container, nil
}

func (s *ContainerService) Delete(ctx context.Context, id int64, auther *models.Auther) *faulterr.FaultErr {
	if !auther.IsAdmin {
		return faulterr.NewUnauthorizedError("Permission not granted")
	}
	_, err := s.dbstore.ContainerStore.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Start db transaction
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.ContainerStore.Delete(ctx, tx, id); err != nil {
		return err
	}
	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return err
	}

	return nil
}
