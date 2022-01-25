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

func InsertContainers(
	ctx context.Context,
	tx pgx.Tx,
	d *dbstore.DBStore,
	creator *models.User,
) *faulterr.FaultErr {
	count := 1
	for _, obj := range containersDefault {
		uid, uuidErr := uuid.NewV4()
		if uuidErr != nil {
			log.Fatal(uuidErr)
		}

		obj.UID = uid
		obj.Code = fmt.Sprintf("CNT%05d", count)
		obj.CreatedByID = creator.ID
		_, err := d.ContainerStore.Insert(ctx, tx, obj)
		if err != nil {
			return err
		}
		count++
	}

	logger.Success(fmt.Sprintf("%v containers added to the database", len(containersDefault)))

	return nil
}

var (
	containersDefault = []models.Container{
		{
			Description: "Container Description One",
			IsArchived:  false,
		},
		{
			Description: "Container Description Two",
			IsArchived:  false,
		},
	}
)
