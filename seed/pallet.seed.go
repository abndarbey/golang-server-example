package seed

import (
	"context"
	"fmt"
	"log"
	"orijinplus/app/models"
	"orijinplus/app/store/dbstore"
	"orijinplus/utils/faulterr"
	"orijinplus/utils/logger"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

func InsertPallets(
	ctx context.Context,
	tx pgx.Tx,
	d *dbstore.DBStore,
	creator *models.User,
) *faulterr.FaultErr {
	count := 1
	for _, obj := range palletsDefault {
		uid, uuidErr := uuid.NewV4()
		if uuidErr != nil {
			log.Fatal(uuidErr)
		}

		obj.UID = uid
		obj.Code = fmt.Sprintf("PLT%05d", count)
		obj.CreatedByID = creator.ID
		_, err := d.PalletStore.Insert(ctx, tx, obj)
		if err != nil {
			return err
		}
		count++
	}

	logger.Success(fmt.Sprintf("%v pallets added to the database", len(palletsDefault)))

	return nil
}

var (
	palletsDefault = []models.Pallet{
		{
			Description: "Pallet Description One",
			IsArchived:  false,
		},
		{
			Description: "Pallet Description One",
			IsArchived:  false,
		},
	}
)
