package dbstore

import (
	"context"
	"fmt"
	"orijinplus/app/models"
	"orijinplus/utils/faulterr"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserStore struct {
	conn *pgxpool.Pool
}

var _ UserStoreInterface = &UserStore{}

type UserStoreInterface interface {
	GetMany(ctx context.Context, ids []int64) ([]*models.User, error)
	List(ctx context.Context, isAdmin, isMember, isCustomer bool) ([]models.User, *faulterr.FaultErr)
	ListAdmin(ctx context.Context) ([]models.User, *faulterr.FaultErr)
	ListMembers(ctx context.Context) ([]models.User, *faulterr.FaultErr)
	ListMembersByOrgID(ctx context.Context, orgID int64) ([]models.User, *faulterr.FaultErr)
	ListCustomers(ctx context.Context) ([]models.User, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64) (*models.User, *faulterr.FaultErr)
	GetByEmail(ctx context.Context, email string) (*models.User, *faulterr.FaultErr)
	GetByPhone(ctx context.Context, phone string) (*models.User, *faulterr.FaultErr)
	Insert(ctx context.Context, tx pgx.Tx, u models.User) (*models.User, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, u models.User) *faulterr.FaultErr
	Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr
}

func NewUserStore(conn *pgxpool.Pool) *UserStore {
	return &UserStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetMany get all users by ids
func (s *UserStore) GetMany(ctx context.Context, ids []int64) ([]*models.User, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		index := strconv.Itoa(i + 1)
		placeholders[i] = "$" + index
		args[i] = ids[i]
	}

	qeryStmt := "SELECT * from users WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	// errMsg := "error when trying to get users"

	rows, err := s.conn.Query(ctx, qeryStmt, args...)
	if err != nil {
		return nil, err
	}

	users, err := s.scanList(rows)
	if err != nil {
		return nil, err
	}

	result := []*models.User{}
	for i := 0; i < len(users); i++ {
		result = append(result, &users[i])
	}

	return result, nil
}

// List gets all users
func (s *UserStore) List(ctx context.Context, isAdmin, isMember, isCustomer bool) ([]models.User, *faulterr.FaultErr) {
	qeryStmt := `
	SELECT * FROM users
	WHERE is_admin=$1 OR is_member=$2 OR is_customer=$3
	`

	errMsg := "error when trying to get users"

	rows, err := s.conn.Query(ctx, qeryStmt, isAdmin, isMember, isCustomer)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	users, err := s.scanList(rows)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return users, nil
}

// ListAdmin all users registered as admin
func (s *UserStore) ListAdmin(ctx context.Context) ([]models.User, *faulterr.FaultErr) {
	qeryStmt := `
	SELECT * FROM users
	WHERE is_admin = true
	`

	errMsg := "error when trying to get users"

	rows, err := s.conn.Query(ctx, qeryStmt)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	users, err := s.scanList(rows)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return users, nil
}

// ListMembers get all users registered as member
func (s *UserStore) ListMembers(ctx context.Context) ([]models.User, *faulterr.FaultErr) {
	qeryStmt := `
	SELECT * FROM users
	WHERE users.is_member = true
	`

	errMsg := "error when trying to get members"

	rows, err := s.conn.Query(ctx, qeryStmt)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	users, err := s.scanList(rows)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return users, nil
}

// ListMembers get all users registered as member
func (s *UserStore) ListMembersByOrgID(ctx context.Context, orgID int64) ([]models.User, *faulterr.FaultErr) {
	qeryStmt := `
	SELECT * FROM users
	WHERE users.is_member = true and users.organization_id = $1
	`

	errMsg := "error when trying to get users"

	rows, err := s.conn.Query(ctx, qeryStmt, orgID)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	users, err := s.scanList(rows)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return users, nil
}

// ListCustomers gets all users registered as consumers
func (s *UserStore) ListCustomers(ctx context.Context) ([]models.User, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM users
	WHERE is_customer = true
	`

	errMsg := "error when trying to get consumers"

	rows, err := s.conn.Query(context.Background(), queryStmt)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	users, err := s.scanList(rows)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return users, nil
}

// GetByID User
func (s *UserStore) GetByID(ctx context.Context, id int64) (*models.User, *faulterr.FaultErr) {
	queryStmt := `SELECT * FROM users WHERE id=$1`

	u := models.User{}
	errMsg := "error when trying to get user by id"

	row := s.conn.QueryRow(ctx, queryStmt, id)
	if err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Phone,
		&u.IsAdmin,
		&u.IsMember,
		&u.IsCustomer,
		&u.PasswordHash,
		&u.OrganizationID,
		&u.RoleID,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return &u, nil
}

// GetByEmail User
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*models.User, *faulterr.FaultErr) {
	queryStmt := `SELECT * FROM users WHERE email=$1`

	errMsg := fmt.Sprintf("error when trying to get user by email - %s", email)

	row := s.conn.QueryRow(ctx, queryStmt, email)
	u, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return u, nil
}

// GetByPhone User
func (s *UserStore) GetByPhone(ctx context.Context, phone string) (*models.User, *faulterr.FaultErr) {
	queryStmt := `SELECT * FROM users WHERE phone=$1`

	errMsg := fmt.Sprintf("error when trying to get user by phone - %s", phone)

	row := s.conn.QueryRow(ctx, queryStmt, phone)
	u, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return u, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Mutate****///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// Insert User
func (s *UserStore) Insert(ctx context.Context, tx pgx.Tx, u models.User) (*models.User, *faulterr.FaultErr) {
	queryStmt := `
	INSERT INTO
	users(first_name, last_name, email, phone, is_admin, is_member, is_customer, password_hash, organization_id, role_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING *
	`
	row := tx.QueryRow(ctx, queryStmt,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		u.IsAdmin,
		u.IsMember,
		u.IsCustomer,
		u.PasswordHash,
		u.OrganizationID,
		u.RoleID,
	)

	if err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Phone,
		&u.IsAdmin,
		&u.IsMember,
		&u.IsCustomer,
		&u.PasswordHash,
		&u.OrganizationID,
		&u.RoleID,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, "error when trying to insert user")
	}

	return &u, nil
}

// Update User
func (s *UserStore) Update(ctx context.Context, tx pgx.Tx, u models.User) *faulterr.FaultErr {
	queryStmt := `
	UPDATE users
	SET first_name=$1, last_name=$2, email=$3, phone=$4, role_id=$5
	WHERE id=$6
	`

	errMsg := "error when trying to update user"
	_, err := tx.Exec(ctx, queryStmt,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Phone,
		&u.RoleID,
		&u.ID,
	)
	if err != nil {
		return faulterr.NewPostgresError(err, errMsg)
	}

	return nil
}

// Delete User
func (s *UserStore) Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr {
	queryStmt := `DELETE FROM users WHERE id=$1`
	_, err := tx.Exec(ctx, queryStmt, id)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to delete user")
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Helpers****//////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

func (s *UserStore) scanList(rows pgx.Rows) ([]models.User, error) {
	users := []models.User{}
	u := models.User{}

	for rows.Next() {
		if err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Phone,
			&u.IsAdmin,
			&u.IsMember,
			&u.IsCustomer,
			&u.PasswordHash,
			&u.OrganizationID,
			&u.RoleID,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (s *UserStore) scanRow(row pgx.Row) (*models.User, error) {
	u := &models.User{}
	if err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Phone,
		&u.IsAdmin,
		&u.IsMember,
		&u.IsCustomer,
		&u.PasswordHash,
		&u.OrganizationID,
		&u.RoleID,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return u, nil
}
