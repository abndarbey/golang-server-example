package services

import (
	"context"
	"orijinplus/app/master"
	"orijinplus/app/models"
	"orijinplus/app/store/dbstore"
	"orijinplus/utils/faulterr"

	"github.com/volatiletech/null"
)

type RoleService struct {
	dbstore *dbstore.DBStore
	master  *master.Master
}

var _ RoleServiceInterface = &RoleService{}

type RoleServiceInterface interface {
	List(ctx context.Context, orgID null.Int64, auther *models.Auther) ([]models.Role, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64, auther *models.Auther) (*models.Role, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string, auther *models.Auther) (*models.Role, *faulterr.FaultErr)
	Create(ctx context.Context, request models.RoleCreateRequest, auther *models.Auther) (*models.Role, *faulterr.FaultErr)
	Update(ctx context.Context, request models.RoleUpdateRequest, auther *models.Auther) (*models.Role, *faulterr.FaultErr)
	Delete(ctx context.Context, id int64, auther *models.Auther) *faulterr.FaultErr
}

func NewRoleService(s *dbstore.DBStore, m *master.Master) *RoleService {
	return &RoleService{s, m}
}

// List gets all roles for super admin and associated organization roles for members
func (s *RoleService) List(ctx context.Context, orgID null.Int64, auther *models.Auther) ([]models.Role, *faulterr.FaultErr) {
	if auther.IsAdmin {
		if orgID.Valid {
			return s.dbstore.RoleStore.ListByOrgID(ctx, orgID.Int64)
		}
		return s.dbstore.RoleStore.ListAll(ctx)
	}

	return s.dbstore.RoleStore.ListByOrgID(ctx, auther.OrganizationID.Int64)
}

// GetByID gets a role by role id
func (s *RoleService) GetByID(ctx context.Context, id int64, auther *models.Auther) (*models.Role, *faulterr.FaultErr) {
	role, err := s.dbstore.RoleStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if auther.IsMember && auther.OrganizationID.Int64 != role.OrganizationID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return role, nil
}

// GetByID gets a role by role code
func (s *RoleService) GetByCode(ctx context.Context, code string, auther *models.Auther) (*models.Role, *faulterr.FaultErr) {
	role, err := s.dbstore.RoleStore.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if auther.IsMember && auther.OrganizationID.Int64 != role.ID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return role, nil
}

// Create saves a role object in db
func (s *RoleService) Create(ctx context.Context, request models.RoleCreateRequest, auther *models.Auther) (*models.Role, *faulterr.FaultErr) {
	if auther.IsAdmin && request.OrganizationID > 0 {
		// Verify organization
		_, err := s.dbstore.OrganizationStore.GetByID(ctx, request.OrganizationID)
		if err != nil {
			return nil, err
		}
	}
	if auther.IsMember {
		request.OrganizationID = auther.OrganizationID.Int64
	}

	// Begin db transaction
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	role, err := s.master.RoleMaster.Create(ctx, tx, request)
	if err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) Update(ctx context.Context, request models.RoleUpdateRequest, auther *models.Auther) (*models.Role, *faulterr.FaultErr) {
	role, err := s.GetByID(ctx, request.ID, auther)
	if err != nil {
		return nil, err
	}
	if err = s.master.RoleMaster.ValidateUpdate(request); err != nil {
		return nil, err
	}

	// Update fields
	role.Name = request.Name
	role.Permissions = request.Permissions
	role.IsArchived = request.IsArchived

	// Begin db transaction
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.RoleStore.Update(ctx, tx, *role); err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) Delete(ctx context.Context, id int64, auther *models.Auther) *faulterr.FaultErr {
	_, err := s.dbstore.RoleStore.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Start db transaction
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.RoleStore.Delete(ctx, tx, id); err != nil {
		return err
	}
	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (s *RoleService) UniquePermissions(list []string) []string {
	permissions := []string{}

	for _, perm := range list {
		if !s.ContainsPermission(permissions[:], perm) {
			permissions = append(permissions, perm)
		}
	}
	return permissions
}

func (s *RoleService) ContainsPermission(list []string, perm string) bool {
	for _, e := range list {
		if e == perm {
			return true
		}
	}
	return false
}
