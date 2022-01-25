package dbstore

import (
	"context"
	"fmt"
	"orijinplus/utils/faulterr"
	"orijinplus/utils/logger"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBTX struct {
	conn *pgxpool.Pool
}

func NewDBTX(conn *pgxpool.Pool) *DBTX {
	return &DBTX{conn}
}

func (t *DBTX) BeginTx(ctx context.Context) (pgx.Tx, *faulterr.FaultErr) {
	// txOptions := &sql.TxOptions{
	// 	Isolation: sql.IsolationLevel(1),
	// 	ReadOnly:  false,
	// }

	tx, txErr := t.conn.BeginTx(ctx, pgx.TxOptions{})
	if txErr != nil {
		return nil, faulterr.NewInternalServerError(txErr.Error())
	}
	logger.Info("Beign Transaction")

	return tx, nil

}

func (t *DBTX) CommitTx(ctx context.Context, tx pgx.Tx) *faulterr.FaultErr {
	err := tx.Commit(ctx)
	if err != nil {
		return faulterr.NewInternalServerError(err.Error())
	}
	logger.Info("Commit Transaction")

	return nil
}

func (t *DBTX) RollbackTx(ctx context.Context, tx pgx.Tx) *faulterr.FaultErr {
	err := tx.Rollback(ctx)
	if err != nil {
		// return faulterr.NewInternalServerError(err.Error())
		logger.Info(fmt.Sprintf("Rollback Transaction %s", err.Error()))
		return nil
	}
	logger.Info("Rollback Transaction")

	return nil
}
