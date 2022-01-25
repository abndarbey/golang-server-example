package services

import (
	"context"
	"orijinplus/app/master"
	"orijinplus/app/models"
	"orijinplus/app/store/dbstore"
	"orijinplus/utils/faulterr"
)

type OrganizationService struct {
	dbstore *dbstore.DBStore
	master  *master.Master
}

var _ OrganizationServiceInterface = &OrganizationService{}

type OrganizationServiceInterface interface {
	List(ctx context.Context) ([]models.Organization, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64, auther *models.Auther) (*models.Organization, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string, auther *models.Auther) (*models.Organization, *faulterr.FaultErr)
	Update(ctx context.Context, request models.Organization, auther *models.Auther) (*models.Organization, *faulterr.FaultErr)
	Archive(ctx context.Context, id int64, auther *models.Auther) (*models.Organization, *faulterr.FaultErr)
	Unarchive(ctx context.Context, id int64, auther *models.Auther) (*models.Organization, *faulterr.FaultErr)
	Delete(ctx context.Context, id int64, auther *models.Auther) *faulterr.FaultErr
}

func NewOrganizationService(s *dbstore.DBStore, m *master.Master) *OrganizationService {
	return &OrganizationService{s, m}
}

// List gets all skus
func (s *OrganizationService) List(ctx context.Context) ([]models.Organization, *faulterr.FaultErr) {
	return s.dbstore.OrganizationStore.List(ctx)
}

func (s *OrganizationService) GetByID(ctx context.Context, id int64, auther *models.Auther) (*models.Organization, *faulterr.FaultErr) {
	org, err := s.dbstore.OrganizationStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if auther.IsMember && auther.OrganizationID.Int64 != org.ID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return org, nil
}

func (s *OrganizationService) GetByCode(ctx context.Context, code string, auther *models.Auther) (*models.Organization, *faulterr.FaultErr) {
	org, err := s.dbstore.OrganizationStore.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if auther.IsMember && auther.OrganizationID.Int64 != org.ID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return org, nil
}

func (s *OrganizationService) Update(ctx context.Context, request models.Organization, auther *models.Auther) (*models.Organization, *faulterr.FaultErr) {
	org, err := s.GetByID(ctx, request.ID, auther)
	if err != nil {
		return nil, err
	}
	if err = s.master.OrganizationMaster.ValidateUpdate(request); err != nil {
		return nil, err
	}

	// Compare changes
	if request.Name != "" {
		org.Name = request.Name
	}
	org.Website = request.Website

	// Begin transaction
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.OrganizationStore.Update(ctx, tx, *org); err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *OrganizationService) Archive(ctx context.Context, id int64, auther *models.Auther) (*models.Organization, *faulterr.FaultErr) {
	org, err := s.GetByID(ctx, id, auther)
	if err != nil {
		return nil, err
	}

	if org.IsArchived {
		return nil, faulterr.NewBadRequestError("organization is already archived")
	}

	org.IsArchived = true
	if _, err := s.Update(ctx, *org, auther); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *OrganizationService) Unarchive(ctx context.Context, id int64, auther *models.Auther) (*models.Organization, *faulterr.FaultErr) {
	org, err := s.GetByID(ctx, id, auther)
	if err != nil {
		return nil, err
	}

	if !org.IsArchived {
		return nil, faulterr.NewBadRequestError("organization is already unarchived")
	}

	org.IsArchived = false
	if _, err := s.Update(ctx, *org, auther); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *OrganizationService) Delete(ctx context.Context, id int64, auther *models.Auther) *faulterr.FaultErr {
	if !auther.IsAdmin {
		return faulterr.NewUnauthorizedError("permission not granted")
	}

	_, err := s.dbstore.OrganizationStore.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Start db transaction
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.OrganizationStore.Delete(ctx, tx, id); err != nil {
		return err
	}
	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return err
	}

	return nil
}
