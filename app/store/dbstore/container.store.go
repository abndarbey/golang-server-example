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

type ContainerStore struct {
	conn *pgxpool.Pool
}

var _ ContainerStoreInterface = &ContainerStore{}

type ContainerStoreInterface interface {
	GetLastInsertedRow(ctx context.Context) (int64, *faulterr.FaultErr)
	List(ctx context.Context) ([]models.Container, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64) (*models.Container, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string) (*models.Container, *faulterr.FaultErr)
	Insert(ctx context.Context, tx pgx.Tx, obj models.Container) (*models.Container, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, obj models.Container) *faulterr.FaultErr
	Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr
}

func NewContainerStore(conn *pgxpool.Pool) *ContainerStore {
	return &ContainerStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetLastInsertedRow retrives last row from database
func (s *ContainerStore) GetLastInsertedRow(ctx context.Context) (int64, *faulterr.FaultErr) {
	queryStmt := `
	SELECT id FROM containers
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

// GetMany get all containers by ids
func (s *ContainerStore) GetMany(ctx context.Context, ids []int64) ([]*models.Container, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		index := strconv.Itoa(i + 1)
		placeholders[i] = "$" + index
		args[i] = ids[i]
	}

	queryStmt := "SELECT * from containers WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	// errMsg := "error when trying to get containers"

	containers := []models.Container{}
	obj := models.Container{}
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
			&obj.IsArchived,
			&obj.OrganizationID,
			&obj.CreatedByID,
			&obj.CreatedAt,
			&obj.UpdatedAt,
		); err != nil {
			return nil, err
		}
		containers = append(containers, obj)
	}

	result := []*models.Container{}
	for i := 0; i < len(containers); i++ {
		result = append(result, &containers[i])
	}

	return result, nil
}

// List retrives all containers from database
func (s *ContainerStore) List(ctx context.Context) ([]models.Container, *faulterr.FaultErr) {
	queryStmt := `SELECT * FROM containers`

	containers := []models.Container{}
	obj := models.Container{}
	errMsg := "error when trying to get containers"

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
			&obj.IsArchived,
			&obj.OrganizationID,
			&obj.CreatedByID,
			&obj.CreatedAt,
			&obj.UpdatedAt,
		); err != nil {
			return nil, faulterr.NewPostgresError(err, errMsg)
		}
		containers = append(containers, obj)
	}

	return containers, nil
}

// ListByOrgID retrives all containers from database
func (s *ContainerStore) ListByOrgID(ctx context.Context, orgID int64) ([]models.Container, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM containers
	WHERE containers.organization_id = $1
	`

	containers := []models.Container{}
	obj := models.Container{}
	errMsg := "error when trying to get containers"

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
			&obj.IsArchived,
			&obj.OrganizationID,
			&obj.CreatedByID,
			&obj.CreatedAt,
			&obj.UpdatedAt,
		); err != nil {
			return nil, faulterr.NewPostgresError(err, errMsg)
		}
		containers = append(containers, obj)
	}

	return containers, nil
}

// GetByID gets container by ID from database
func (s *ContainerStore) GetByID(ctx context.Context, id int64) (*models.Container, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM containers
	WHERE containers.id = $1
	`

	obj := models.Container{}

	row := s.conn.QueryRow(ctx, queryStmt, id)
	if err := row.Scan(
		&obj.ID,
		&obj.UID,
		&obj.Code,
		&obj.Description,
		&obj.IsArchived,
		&obj.OrganizationID,
		&obj.CreatedByID,
		&obj.CreatedAt,
		&obj.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, "error when trying to get container")
	}

	return &obj, nil
}

// GetByID gets container by UID from database
func (s *ContainerStore) GetByUID(ctx context.Context, uid uuid.UUID) (*models.Container, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM containers
	WHERE containers.uid = $1
	`

	obj := models.Container{}

	row := s.conn.QueryRow(ctx, queryStmt, uid)
	if err := row.Scan(
		&obj.ID,
		&obj.UID,
		&obj.Code,
		&obj.Description,
		&obj.IsArchived,
		&obj.OrganizationID,
		&obj.CreatedByID,
		&obj.CreatedAt,
		&obj.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, "error when trying to get container")
	}

	return &obj, nil
}

// GetByCode gets container by code from database
func (s *ContainerStore) GetByCode(ctx context.Context, code string) (*models.Container, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM containers
	WHERE containers.code = $1
	`

	obj := models.Container{}

	row := s.conn.QueryRow(ctx, queryStmt, code)
	if err := row.Scan(
		&obj.ID,
		&obj.UID,
		&obj.Code,
		&obj.Description,
		&obj.IsArchived,
		&obj.OrganizationID,
		&obj.CreatedByID,
		&obj.CreatedAt,
		&obj.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, "error when trying to get container")
	}

	return &obj, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Mutate****///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// Insert inserts a container in database
func (s *ContainerStore) Insert(ctx context.Context, tx pgx.Tx, obj models.Container) (*models.Container, *faulterr.FaultErr) {
	queryStmt := `
	INSERT INTO
	containers(
		uid,
		code,
		description,
		is_archived,
		organization_id,
		created_by_id
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING *
	`

	row := tx.QueryRow(ctx, queryStmt,
		&obj.UID,
		&obj.Code,
		&obj.Description,
		&obj.IsArchived,
		&obj.OrganizationID,
		&obj.CreatedByID,
	)

	if err := row.Scan(
		&obj.ID,
		&obj.UID,
		&obj.Code,
		&obj.Description,
		&obj.IsArchived,
		&obj.OrganizationID,
		&obj.CreatedByID,
		&obj.CreatedAt,
		&obj.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, "error when trying to insert container")
	}

	return &obj, nil
}

// Update updates a container in database
func (s *ContainerStore) Update(ctx context.Context, tx pgx.Tx, obj models.Container) *faulterr.FaultErr {
	queryStmt := `
	UPDATE containers
	SET
		description = $1,
		is_archived = $2
	WHERE id=$3
	`

	_, err := tx.Exec(ctx, queryStmt,
		&obj.Description,
		&obj.IsArchived,
		&obj.ID,
	)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to update container")
	}

	return nil
}

// Delete deletes a container from database
func (s *ContainerStore) Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr {
	queryStmt := `DELETE FROM containers WHERE id=$1`

	_, err := tx.Exec(ctx, queryStmt, id)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to delete container")
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Helpers****//////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
