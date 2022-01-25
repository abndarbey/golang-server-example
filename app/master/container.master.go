package master

import (
	"context"
	"fmt"
	"orijinplus/app/models"
	"orijinplus/app/store/dbstore"
	"orijinplus/utils/faulterr"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

type ContainerMaster struct {
	dbstore *dbstore.DBStore
}

func NewContainerMaster(s *dbstore.DBStore) *ContainerMaster {
	return &ContainerMaster{s}
}

func (m *ContainerMaster) Create(ctx context.Context, tx pgx.Tx, r models.ContainerRequest, createdByID int64) (*models.Container, *faulterr.FaultErr) {
	if err := m.validate(r); err != nil {
		return nil, err
	}

	// Get last inserted row ID
	lastRowID, err := m.dbstore.ContainerStore.GetLastInsertedRow(ctx)
	if err != nil {
		return nil, err
	}

	uid, uidErr := uuid.NewV4()
	if uidErr != nil {
		return nil, faulterr.NewInternalServerError(uidErr.Error())
	}

	count := lastRowID + 1
	obj := models.Container{
		UID:            uid,
		Code:           fmt.Sprintf("CNT%05d", count),
		Description:    r.Description,
		IsArchived:     false,
		OrganizationID: r.OrganizationID,
		CreatedByID:    createdByID,
	}

	return m.dbstore.ContainerStore.Insert(ctx, tx, obj)
}

func (m *ContainerMaster) Update(
	ctx context.Context,
	tx pgx.Tx,
	obj *models.Container,
	req models.ContainerRequest,
) (*models.Container, *faulterr.FaultErr) {
	// Validate request
	if err := m.validate(req); err != nil {
		return nil, err
	}

	// Update fields
	obj.Description = req.Description

	if err := m.dbstore.ContainerStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *ContainerMaster) validate(r models.ContainerRequest) *faulterr.FaultErr {
	// if !r.Name.Valid {
	// 	return faulterr.NewBadRequestError("contract name is required")
	// }
	return nil
}
