package services

import (
	"context"
	"orijinplus/app/master"
	"orijinplus/app/models"
	"orijinplus/app/store/dbstore"
	"orijinplus/utils/faulterr"

	"github.com/gofrs/uuid"
)

type PalletService struct {
	dbstore *dbstore.DBStore
	master  *master.Master
}

var _ PalletServiceInterface = &PalletService{}

type PalletServiceInterface interface {
	List(ctx context.Context, auther *models.Auther) ([]models.Pallet, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64, auther *models.Auther) (*models.Pallet, *faulterr.FaultErr)
	GetByUID(ctx context.Context, uid uuid.UUID, auther *models.Auther) (*models.Pallet, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string, auther *models.Auther) (*models.Pallet, *faulterr.FaultErr)
	Create(ctx context.Context, request models.PalletRequest, auther *models.Auther) (*models.Pallet, *faulterr.FaultErr)
	Update(ctx context.Context, id int64, request models.PalletRequest, auther *models.Auther) (*models.Pallet, *faulterr.FaultErr)
	Delete(ctx context.Context, id int64, auther *models.Auther) *faulterr.FaultErr
}

func NewPalletService(s *dbstore.DBStore, m *master.Master) *PalletService {
	return &PalletService{s, m}
}

// List gets all skus
func (s *PalletService) List(ctx context.Context, auther *models.Auther) ([]models.Pallet, *faulterr.FaultErr) {
	if auther.IsAdmin {
		return s.dbstore.PalletStore.List(ctx)
	}
	return s.dbstore.PalletStore.ListByOrgID(ctx, auther.OrganizationID.Int64)
}

func (s *PalletService) GetByID(ctx context.Context, id int64, auther *models.Auther) (*models.Pallet, *faulterr.FaultErr) {
	pallet, err := s.dbstore.PalletStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !auther.IsAdmin && auther.OrganizationID.Int64 != pallet.OrganizationID.Int64 {
		return nil, faulterr.NewNotFoundError("no pallet found")
	}
	return pallet, nil
}

func (s *PalletService) GetByUID(ctx context.Context, uid uuid.UUID, auther *models.Auther) (*models.Pallet, *faulterr.FaultErr) {
	pallet, err := s.dbstore.PalletStore.GetByUID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if !auther.IsAdmin && auther.OrganizationID.Int64 != pallet.OrganizationID.Int64 {
		return nil, faulterr.NewNotFoundError("no pallet found")
	}
	return pallet, nil
}

func (s *PalletService) GetByCode(ctx context.Context, code string, auther *models.Auther) (*models.Pallet, *faulterr.FaultErr) {
	pallet, err := s.dbstore.PalletStore.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if !auther.IsAdmin && auther.OrganizationID.Int64 != pallet.OrganizationID.Int64 {
		return nil, faulterr.NewNotFoundError("no pallet found")
	}
	return pallet, nil
}

// Create gets all skus
func (s *PalletService) Create(ctx context.Context, req models.PalletRequest, auther *models.Auther) (*models.Pallet, *faulterr.FaultErr) {
	if auther.IsAdmin && !req.OrganizationID.Valid {
		return nil, faulterr.NewBadRequestError("organization id is required")
	}
	// Reassign organization ID to the request
	if !auther.IsAdmin {
		req.OrganizationID = auther.OrganizationID
	}
	createdByID := auther.ID

	// Verify container and verify organization
	if req.ContainerID.Valid {
		container, err := s.dbstore.ContainerStore.GetByID(ctx, req.ContainerID.Int64)
		if err != nil {
			return nil, err
		}
		if !auther.IsAdmin && container.OrganizationID != auther.OrganizationID {
			return nil, faulterr.NewNotFoundError("no container found with given container id")
		}
	}

	// Start transactions
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	pallet, err := s.master.PalletMaster.Create(ctx, tx, req, createdByID)
	if err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return pallet, nil
}

func (s *PalletService) Update(ctx context.Context, id int64, request models.PalletRequest, auther *models.Auther) (*models.Pallet, *faulterr.FaultErr) {
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

	pallet, err := s.master.PalletMaster.Update(ctx, tx, current, request)
	if err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return pallet, nil
}

func (s *PalletService) Archive(ctx context.Context, id int64, auther *models.Auther) (*models.Pallet, *faulterr.FaultErr) {
	pallet, err := s.GetByID(ctx, id, auther)
	if err != nil {
		return nil, err
	}
	pallet.IsArchived = true

	// Start transactions
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.PalletStore.Update(ctx, tx, *pallet); err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return pallet, nil
}

func (s *PalletService) Unarchive(ctx context.Context, id int64, auther *models.Auther) (*models.Pallet, *faulterr.FaultErr) {
	pallet, err := s.GetByID(ctx, id, auther)
	if err != nil {
		return nil, err
	}
	pallet.IsArchived = false

	// Start transactions
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.PalletStore.Update(ctx, tx, *pallet); err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return pallet, nil
}

func (s *PalletService) Delete(ctx context.Context, id int64, auther *models.Auther) *faulterr.FaultErr {
	if !auther.IsAdmin {
		return faulterr.NewUnauthorizedError("Permission not granted")
	}
	_, err := s.dbstore.PalletStore.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Start db transaction
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.PalletStore.Delete(ctx, tx, id); err != nil {
		return err
	}
	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return err
	}

	return nil
}
