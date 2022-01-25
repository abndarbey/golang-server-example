package dbstore

import (
	"context"
	"orijinplus/app/models"
	"orijinplus/utils/faulterr"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PalletStore struct {
	conn *pgxpool.Pool
}

var _ PalletStoreInterface = &PalletStore{}

type PalletStoreInterface interface {
	GetLastInsertedRow(ctx context.Context) (int64, *faulterr.FaultErr)
	List(ctx context.Context) ([]models.Pallet, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64) (*models.Pallet, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string) (*models.Pallet, *faulterr.FaultErr)
	Insert(ctx context.Context, tx pgx.Tx, o models.Pallet) (*models.Pallet, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, o models.Pallet) *faulterr.FaultErr
	Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr
}

func NewPalletStore(conn *pgxpool.Pool) *PalletStore {
	return &PalletStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetLastInsertedRow retrives last row from database
func (s *PalletStore) GetLastInsertedRow(ctx context.Context) (int64, *faulterr.FaultErr) {
	queryStmt := `
	SELECT id FROM pallets
	ORDER BY id DESC
	LIMIT 1
	`

	var ID int64
	errMsg := "error getting last inserted row"

	rows, err := s.conn.Query(context.Background(), queryStmt)
	if err != nil {
		return 0, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&ID); err != nil {
			return 0, faulterr.NewPostgresError(err, errMsg)
		}
	}

	return ID, nil
}

// GetMany get all pallets by ids
func (s *PalletStore) GetMany(ctx context.Context, ids []int64) ([]*models.Pallet, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		index := strconv.Itoa(i + 1)
		placeholders[i] = "$" + index
		args[i] = ids[i]
	}

	queryStmt := "SELECT * from pallets WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	// errMsg := "error when trying to get pallets"

	pallets := []models.Pallet{}
	obj := models.Pallet{}
	rows, err := s.conn.Query(ctx, queryStmt, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&obj.ID,
			&obj.UID,
			&obj.Code,
			&obj.Description,
			&obj.ContainerID,
			&obj.IsArchived,
			&obj.OrganizationID,
			&obj.CreatedByID,
			&obj.CreatedAt,
			&obj.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pallets = append(pallets, obj)
	}

	result := []*models.Pallet{}
	for i := 0; i < len(pallets); i++ {
		result = append(result, &pallets[i])
	}

	return result, nil
}

// List retrives all pallets from database
func (s *PalletStore) List(ctx context.Context) ([]models.Pallet, *faulterr.FaultErr) {
	queryStmt := `SELECT * FROM pallets`

	pallets := []models.Pallet{}
	obj := models.Pallet{}
	errMsg := "error when trying to get pallets"

	rows, err := s.conn.Query(context.Background(), queryStmt)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&obj.ID,
			&obj.UID,
			&obj.Code,
			&obj.Description,
			&obj.ContainerID,
			&obj.IsArchived,
			&obj.OrganizationID,
			&obj.CreatedByID,
			&obj.CreatedAt,
			&obj.UpdatedAt,
		); err != nil {
			return nil, faulterr.NewPostgresError(err, errMsg)
		}
		pallets = append(pallets, obj)
	}

	return pallets, nil
}

// ListByOrgID retrives all pallets from database
func (s *PalletStore) ListByOrgID(ctx context.Context, orgID int64) ([]models.Pallet, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM pallets
	WHERE pallets.organization_id = $1
	`

	pallets := []models.Pallet{}
	obj := models.Pallet{}
	errMsg := "error when trying to get pallets"

	rows, err := s.conn.Query(ctx, queryStmt, orgID)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&obj.ID,
			&obj.UID,
			&obj.Code,
			&obj.Description,
			&obj.ContainerID,
			&obj.IsArchived,
			&obj.OrganizationID,
			&obj.CreatedByID,
			&obj.CreatedAt,
			&obj.UpdatedAt,
		); err != nil {
			return nil, faulterr.NewPostgresError(err, errMsg)
		}
		pallets = append(pallets, obj)
	}

	return pallets, nil
}

// GetByID gets pallet by ID from database
func (s *PalletStore) GetByID(ctx context.Context, id int64) (*models.Pallet, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM pallets
	WHERE pallets.id = $1
	`

	obj := models.Pallet{}

	row := s.conn.QueryRow(ctx, queryStmt, id)
	if err := row.Scan(
		&obj.ID,
		&obj.UID,
		&obj.Code,
		&obj.Description,
		&obj.ContainerID,
		&obj.IsArchived,
		&obj.OrganizationID,
		&obj.CreatedByID,
		&obj.CreatedAt,
		&obj.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, "error when trying to get pallet")
	}

	return &obj, nil
}

// GetByID gets pallet by UID from database
func (s *PalletStore) GetByUID(ctx context.Context, uid uuid.UUID) (*models.Pallet, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM pallets
	WHERE pallets.uid = $1
	`

	obj := models.Pallet{}

	row := s.conn.QueryRow(ctx, queryStmt, uid)
	if err := row.Scan(
		&obj.ID,
		&obj.UID,
		&obj.Code,
		&obj.Description,
		&obj.ContainerID,
		&obj.IsArchived,
		&obj.OrganizationID,
		&obj.CreatedByID,
		&obj.CreatedAt,
		&obj.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, "error when trying to get pallet")
	}

	return &obj, nil
}

// GetByCode gets pallet by code from database
func (s *PalletStore) GetByCode(ctx context.Context, code string) (*models.Pallet, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM pallets
	WHERE pallets.code = $1
	`

	obj := models.Pallet{}

	row := s.conn.QueryRow(ctx, queryStmt, code)
	if err := row.Scan(
		&obj.ID,
		&obj.UID,
		&obj.Code,
		&obj.Description,
		&obj.ContainerID,
		&obj.IsArchived,
		&obj.OrganizationID,
		&obj.CreatedByID,
		&obj.CreatedAt,
		&obj.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, "error when trying to get pallet")
	}

	return &obj, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Mutate****///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// Insert inserts a pallet in database
func (s *PalletStore) Insert(ctx context.Context, tx pgx.Tx, obj models.Pallet) (*models.Pallet, *faulterr.FaultErr) {
	queryStmt := `
	INSERT INTO
	pallets(
		uid,
		code,
		description,
		container_id,
		is_archived,
		organization_id,
		created_by_id
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING *
	`

	row := tx.QueryRow(ctx, queryStmt,
		&obj.UID,
		&obj.Code,
		&obj.Description,
		&obj.ContainerID,
		&obj.IsArchived,
		&obj.OrganizationID,
		&obj.CreatedByID,
	)

	if err := row.Scan(
		&obj.ID,
		&obj.UID,
		&obj.Code,
		&obj.Description,
		&obj.ContainerID,
		&obj.IsArchived,
		&obj.OrganizationID,
		&obj.CreatedByID,
		&obj.CreatedAt,
		&obj.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, "error when trying to insert pallet")
	}

	return &obj, nil
}

// Update updates a pallet in database
func (s *PalletStore) Update(ctx context.Context, tx pgx.Tx, obj models.Pallet) *faulterr.FaultErr {
	queryStmt := `
	UPDATE pallets
	SET
		description = $1,
		container_id = $2,
		is_archived = $3
	WHERE id=$4
	`

	_, err := tx.Exec(ctx, queryStmt,
		&obj.Description,
		&obj.ContainerID,
		&obj.IsArchived,
		&obj.ID,
	)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to update pallet")
	}

	return nil
}

// Delete deletes a pallet from database
func (s *PalletStore) Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr {
	queryStmt := `DELETE FROM pallets WHERE id=$1`

	_, err := tx.Exec(ctx, queryStmt, id)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to delete pallet")
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Helpers****//////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
