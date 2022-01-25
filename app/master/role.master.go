package master

import (
	"context"
	"fmt"
	"orijinplus/app/models"
	"orijinplus/app/store/dbstore"
	"orijinplus/utils/faulterr"

	"github.com/jackc/pgx/v4"
)

type RoleMaster struct {
	dbstore *dbstore.DBStore
}

func NewRoleMaster(s *dbstore.DBStore) *RoleMaster {
	return &RoleMaster{
		s,
	}
}

func (m *RoleMaster) Create(ctx context.Context, tx pgx.Tx, r models.RoleCreateRequest) (*models.Role, *faulterr.FaultErr) {
	if err := m.validate(r); err != nil {
		return nil, err
	}

	// Get last inserted row id
	lastRowID, err := m.dbstore.RoleStore.GetLastInsertedRow(ctx)
	if err != nil {
		return nil, err
	}

	count := lastRowID + 1
	role := models.Role{
		Code:        fmt.Sprintf("ROLE%05d", count),
		Name:        r.Name,
		Permissions: r.Permissions,
		IsOrgAdmin:  r.IsOrgAdmin,
		IsArchived:  false,
	}
	role.OrganizationID = r.OrganizationID

	return m.dbstore.RoleStore.Insert(ctx, tx, role)
}

func (m *RoleMaster) Update(
	ctx context.Context,
	tx pgx.Tx,
	obj *models.Role,
	req models.RoleCreateRequest,
) (*models.Role, *faulterr.FaultErr) {
	// Validate request
	if err := m.validate(req); err != nil {
		return nil, err
	}

	// Update fields
	obj.Name = req.Name

	if err := m.dbstore.RoleStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *RoleMaster) ValidateUpdate(r models.RoleUpdateRequest) *faulterr.FaultErr {
	if r.Name == "" {
		return faulterr.NewBadRequestError("name cannot be empty")
	}
	if len(r.Permissions) <= 0 {
		return faulterr.NewBadRequestError("permissions list cannot be empty")
	}

	return nil
}

func (m *RoleMaster) GrantPermission(ctx context.Context, roleID int64, permission string) *faulterr.FaultErr {
	role, err := m.dbstore.RoleStore.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	if !role.IsOrgAdmin {
		for _, perm := range role.Permissions {
			if perm == permission {
				return nil
			}
		}
		return faulterr.NewUnauthorizedError("User Not Authorized")
	}

	return nil
}

func (m *RoleMaster) validate(r models.RoleCreateRequest) *faulterr.FaultErr {
	if r.Name == "" {
		return faulterr.NewBadRequestError("Role Name is required")
	}
	if !r.IsOrgAdmin && len(r.Permissions) == 0 {
		return faulterr.NewBadRequestError("Permissions cannot be empty")
	}
	if r.OrganizationID == 0 {
		return faulterr.NewBadRequestError("Organization ID is required")
	}
	return nil
}
