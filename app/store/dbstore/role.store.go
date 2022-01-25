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

type RoleStore struct {
	conn *pgxpool.Pool
}

var _ RoleStoreInterface = &RoleStore{}

type RoleStoreInterface interface {
	GetLastInsertedRow(ctx context.Context) (int64, *faulterr.FaultErr)
	GetMany(ctx context.Context, ids []int64) ([]*models.Role, error)
	ListAll(ctx context.Context) ([]models.Role, *faulterr.FaultErr)
	ListByOrgID(ctx context.Context, orgID int64) ([]models.Role, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64) (*models.Role, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string) (*models.Role, *faulterr.FaultErr)
	Insert(ctx context.Context, tx pgx.Tx, r models.Role) (*models.Role, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, r models.Role) *faulterr.FaultErr
	Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr
}

func NewRoleStore(conn *pgxpool.Pool) *RoleStore {
	return &RoleStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetLastInsertedRow retrives last row from database
func (s *RoleStore) GetLastInsertedRow(ctx context.Context) (int64, *faulterr.FaultErr) {
	queryStmt := `
	SELECT id FROM roles
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

// GetMany get all roles by ids
func (s *RoleStore) GetMany(ctx context.Context, ids []int64) ([]*models.Role, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		index := strconv.Itoa(i + 1)
		placeholders[i] = "$" + index
		args[i] = ids[i]
	}

	queryStmt := "SELECT * from roles WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	// errMsg := "error when trying to get roles"

	rows, err := s.conn.Query(ctx, queryStmt, args...)
	if err != nil {
		return nil, err
	}

	roles, err := s.scanList(rows)
	if err != nil {
		return nil, err
	}

	result := []*models.Role{}
	for i := 0; i < len(roles); i++ {
		result = append(result, &roles[i])
	}

	return result, nil
}

// List retrives all roles from database
func (s *RoleStore) ListAll(ctx context.Context) ([]models.Role, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM roles
	`

	errMsg := "error when trying to get roles"
	rows, err := s.conn.Query(ctx, queryStmt)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	roles, err := s.scanList(rows)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return roles, nil
}

// ListByOrgID retrives all roles of an organization by orgID from database
func (s *RoleStore) ListByOrgID(ctx context.Context, orgID int64) ([]models.Role, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM roles
	WHERE organization_id = $1
	`

	errMsg := "error when trying to get roles"
	rows, err := s.conn.Query(ctx, queryStmt, orgID)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	roles, err := s.scanList(rows)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return roles, nil
}

// GetByID gets role by ID from database
func (s *RoleStore) GetByID(ctx context.Context, id int64) (*models.Role, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM roles
	WHERE roles.id = $1
	`
	errMsg := "error when trying to get role by id"

	row := s.conn.QueryRow(ctx, queryStmt, id)

	r, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return r, nil
}

// GetByCode gets role by code from database
func (s *RoleStore) GetByCode(ctx context.Context, code string) (*models.Role, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM roles
	WHERE roles.code = $1
	`

	errMsg := "error when trying to get role by code"

	row := s.conn.QueryRow(ctx, queryStmt, code)

	r, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return r, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Mutate****///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// Insert inserts a role in database
func (s *RoleStore) Insert(ctx context.Context, tx pgx.Tx, r models.Role) (*models.Role, *faulterr.FaultErr) {
	queryStmt := `
	INSERT INTO
	roles(
		code,
		name,
		permissions,
		is_org_admin,
		is_archived,
		organization_id
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING *
	`

	errMsg := "error when trying to insert role"

	row := tx.QueryRow(ctx, queryStmt,
		&r.Code,
		&r.Name,
		&r.Permissions,
		&r.IsOrgAdmin,
		&r.IsArchived,
		&r.OrganizationID,
	)

	role, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return role, nil
}

// Update updates a role in database
func (s *RoleStore) Update(ctx context.Context, tx pgx.Tx, r models.Role) *faulterr.FaultErr {
	queryStmt := `
	UPDATE roles
	SET name=$1, permissions=$2, is_archived=$3
	WHERE id=$4
	`

	errMsg := "error when trying to update role"

	_, err := tx.Exec(ctx, queryStmt,
		&r.Name,
		&r.Permissions,
		&r.IsArchived,
		&r.ID,
	)
	if err != nil {
		return faulterr.NewPostgresError(err, errMsg)
	}

	return nil
}

// Delete deletes a role from database
func (s *RoleStore) Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr {
	queryStmt := `DELETE FROM roles WHERE id=$1`

	_, err := tx.Exec(ctx, queryStmt)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to delete role")
	}

	return nil
}

func (s *RoleStore) scanList(rows pgx.Rows) ([]models.Role, error) {
	roles := []models.Role{}
	r := models.Role{}

	for rows.Next() {
		if err := rows.Scan(
			&r.ID,
			&r.Code,
			&r.Name,
			&r.Permissions,
			&r.IsOrgAdmin,
			&r.IsArchived,
			&r.OrganizationID,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}

	return roles, nil
}

func (s *RoleStore) scanRow(row pgx.Row) (*models.Role, error) {
	r := models.Role{}

	if err := row.Scan(
		&r.ID,
		&r.Code,
		&r.Name,
		&r.Permissions,
		&r.IsOrgAdmin,
		&r.IsArchived,
		&r.OrganizationID,
		&r.CreatedAt,
		&r.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &r, nil
}
