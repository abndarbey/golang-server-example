package dbstore

import (
	"context"
	"orijinplus/app/models"
	"orijinplus/utils/faulterr"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrganizationStore struct {
	conn *pgxpool.Pool
}

var _ OrganizationStoreInterface = &OrganizationStore{}

type OrganizationStoreInterface interface {
	GetLastInsertedRow(ctx context.Context) (int64, *faulterr.FaultErr)
	GetMany(ctx context.Context, ids []int64) ([]*models.Organization, error)
	List(ctx context.Context) ([]models.Organization, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64) (*models.Organization, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string) (*models.Organization, *faulterr.FaultErr)
	Insert(ctx context.Context, tx pgx.Tx, o models.Organization) (*models.Organization, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, o models.Organization) *faulterr.FaultErr
	Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr
}

func NewOrganizationStore(conn *pgxpool.Pool) *OrganizationStore {
	return &OrganizationStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetLastInsertedRow retrives last row from database
func (s *OrganizationStore) GetLastInsertedRow(ctx context.Context) (int64, *faulterr.FaultErr) {
	queryStmt := `
	SELECT id FROM organizations
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

// GetMany get all organizations by ids
func (s *OrganizationStore) GetMany(ctx context.Context, ids []int64) ([]*models.Organization, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		index := strconv.Itoa(i + 1)
		placeholders[i] = "$" + index
		args[i] = ids[i]
	}

	queryStmt := "SELECT * from organizations WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	// errMsg := "error when trying to get organizations"

	rows, err := s.conn.Query(ctx, queryStmt, args...)
	if err != nil {
		return nil, err
	}

	organizations, err := s.scanList(rows)
	if err != nil {
		return nil, err
	}

	result := []*models.Organization{}
	for i := 0; i < len(organizations); i++ {
		result = append(result, &organizations[i])
	}

	return result, nil
}

// List retrives all organizations from database
func (s *OrganizationStore) List(ctx context.Context) ([]models.Organization, *faulterr.FaultErr) {
	queryStmt := `SELECT * FROM organizations`

	organizations := []models.Organization{}
	o := models.Organization{}
	errMsg := "error when trying to get organizations"

	rows, err := s.conn.Query(ctx, queryStmt)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&o.ID,
			&o.Code,
			&o.Name,
			&o.Website,
			&o.IsArchived,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, faulterr.NewPostgresError(err, errMsg)
		}
		organizations = append(organizations, o)
	}

	return organizations, nil
}

// GetByID gets organization by ID from database
func (s *OrganizationStore) GetByID(ctx context.Context, id int64) (*models.Organization, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM organizations
	WHERE organizations.id = $1
	`

	o := models.Organization{}
	row := s.conn.QueryRow(ctx, queryStmt, id)
	if err := row.Scan(
		&o.ID,
		&o.Code,
		&o.Name,
		&o.Website,
		&o.IsArchived,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, "error when trying to get organization by id")
	}

	return &o, nil
}

// GetByCode gets organization by code from database
func (s *OrganizationStore) GetByCode(ctx context.Context, code string) (*models.Organization, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM organizations
	WHERE organizations.code = $1
	`

	o := models.Organization{}

	row := s.conn.QueryRow(ctx, queryStmt, code)
	if err := row.Scan(
		&o.ID,
		&o.Code,
		&o.Name,
		&o.Website,
		&o.IsArchived,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, "error when trying to get organization by code")
	}

	return &o, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Mutate****///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// Insert inserts a organization in database
func (s *OrganizationStore) Insert(ctx context.Context, tx pgx.Tx, o models.Organization) (*models.Organization, *faulterr.FaultErr) {
	queryStmt := `
	INSERT INTO
	organizations(
		code,
		name,
		website,
		is_archived
	)
	VALUES ($1, $2, $3, $4)
	RETURNING *
	`

	errMsg := "error when trying to insert organization"
	row := tx.QueryRow(ctx, queryStmt,
		&o.Code,
		&o.Name,
		&o.Website,
		&o.IsArchived,
	)

	if err := row.Scan(
		&o.ID,
		&o.Code,
		&o.Name,
		&o.Website,
		&o.IsArchived,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return &o, nil
}

// Update updates a organization in database
func (s *OrganizationStore) Update(ctx context.Context, tx pgx.Tx, o models.Organization) *faulterr.FaultErr {
	queryStmt := `
	UPDATE organizations
	SET name=$1, website=$2, is_archived=$3
	WHERE id=$4
	`

	errMsg := "error when trying to update organization"
	_, err := tx.Exec(ctx, queryStmt,
		&o.Name,
		&o.Website,
		&o.IsArchived,
		&o.ID,
	)
	if err != nil {
		return faulterr.NewPostgresError(err, errMsg)
	}

	return nil
}

// Delete deletes a organization from database
func (s *OrganizationStore) Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr {
	queryStmt := `DELETE FROM organizations WHERE id=$1`

	_, err := tx.Exec(ctx, queryStmt)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to delete organization")
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Helpers****//////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

func (s *OrganizationStore) scanList(rows pgx.Rows) ([]models.Organization, error) {
	orgs := []models.Organization{}
	o := models.Organization{}

	for rows.Next() {
		if err := rows.Scan(
			&o.ID,
			&o.Code,
			&o.Name,
			&o.Website,
			&o.IsArchived,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}

	return orgs, nil
}

func (s *OrganizationStore) scanRow(row pgx.Row) (*models.Organization, error) {
	o := &models.Organization{}
	if err := row.Scan(
		&o.ID,
		&o.Code,
		&o.Name,
		&o.Website,
		&o.IsArchived,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return o, nil
}
