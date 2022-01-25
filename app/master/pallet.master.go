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

type PalletMaster struct {
	dbstore *dbstore.DBStore
}

func NewPalletMaster(s *dbstore.DBStore) *PalletMaster {
	return &PalletMaster{s}
}

func (m *PalletMaster) Create(ctx context.Context, tx pgx.Tx, r models.PalletRequest, createdByID int64) (*models.Pallet, *faulterr.FaultErr) {
	if err := m.validate(r); err != nil {
		return nil, err
	}

	// Get last inserted row ID
	lastRowID, err := m.dbstore.PalletStore.GetLastInsertedRow(ctx)
	if err != nil {
		return nil, err
	}

	uid, uidErr := uuid.NewV4()
	if uidErr != nil {
		return nil, faulterr.NewInternalServerError(uidErr.Error())
	}

	count := lastRowID + 1
	obj := models.Pallet{
		UID:            uid,
		Code:           fmt.Sprintf("PLT%05d", count),
		Description:    r.Description,
		ContainerID:    r.ContainerID,
		IsArchived:     false,
		OrganizationID: r.OrganizationID,
		CreatedByID:    createdByID,
	}

	return m.dbstore.PalletStore.Insert(ctx, tx, obj)
}

func (m *PalletMaster) Update(
	ctx context.Context,
	tx pgx.Tx,
	obj *models.Pallet,
	req models.PalletRequest,
) (*models.Pallet, *faulterr.FaultErr) {
	// Validate request
	if err := m.validate(req); err != nil {
		return nil, err
	}

	// Update fields
	obj.Description = req.Description

	if err := m.dbstore.PalletStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *PalletMaster) validate(r models.PalletRequest) *faulterr.FaultErr {
	// if !r.Name.Valid {
	// 	return faulterr.NewBadRequestError("contract name is required")
	// }
	return nil
}
