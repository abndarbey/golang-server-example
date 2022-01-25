package master

import (
	"context"
	"fmt"
	"orijinplus/app/models"
	"orijinplus/app/store/dbstore"
	"orijinplus/utils/faulterr"

	"github.com/jackc/pgx/v4"
	"github.com/volatiletech/null"
)

type OrganizationMaster struct {
	dbstore *dbstore.DBStore
}

func NewOrganizationMaster(dbstore *dbstore.DBStore) *OrganizationMaster {
	return &OrganizationMaster{dbstore}
}

func (m *OrganizationMaster) Create(ctx context.Context, tx pgx.Tx, r models.OrganizationRequest) (*models.Organization, *faulterr.FaultErr) {
	if err := m.validate(r); err != nil {
		return nil, err
	}

	// Get last inserted row ID
	lastRowID, err := m.dbstore.OrganizationStore.GetLastInsertedRow(ctx)
	if err != nil {
		return nil, err
	}

	count := lastRowID + 1
	o := models.Organization{
		Code:       fmt.Sprintf("ORG%05d", count),
		Name:       r.OrgName,
		Website:    null.StringFrom(r.Website),
		IsArchived: false,
	}

	return m.dbstore.OrganizationStore.Insert(ctx, tx, o)
}

func (m *OrganizationMaster) ValidateUpdate(request models.Organization) *faulterr.FaultErr {
	if request.Name == "" {
		return faulterr.NewBadRequestError("Name cannot be empty")
	}
	return nil
}

func (m *OrganizationMaster) validate(r models.OrganizationRequest) *faulterr.FaultErr {
	if r.OrgName == "" {
		return faulterr.NewBadRequestError("Organization Name is required")
	}
	if r.FirstName == "" {
		return faulterr.NewBadRequestError("First Name is required")
	}
	if r.LastName == "" {
		return faulterr.NewBadRequestError("Last Name is required")
	}
	if r.Email == "" {
		return faulterr.NewBadRequestError("Email is required")
	}
	if r.Password == "" {
		return faulterr.NewBadRequestError("Password is required")
	}
	return nil
}
