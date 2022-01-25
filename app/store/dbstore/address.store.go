package dbstore

import (
	"context"
	"orijinplus/app/models"
	"orijinplus/utils/faulterr"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AddressStore struct {
	conn *pgxpool.Pool
}

var _ AddressStoreInterface = &AddressStore{}

type AddressStoreInterface interface {
	GetMany(ctx context.Context, ids []int64) ([]*models.Address, []error)
	List(ctx context.Context) ([]models.Address, *faulterr.FaultErr)
	ListByUserID(ctx context.Context, userID int64) ([]models.Address, *faulterr.FaultErr)
	GetByID(ctx context.Context, userID int64) (*models.Address, *faulterr.FaultErr)
	Insert(ctx context.Context, tx pgx.Tx, p models.Address) (*models.Address, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, p models.Address) *faulterr.FaultErr
	Delete(ctx context.Context, tx pgx.Tx, userID int64) *faulterr.FaultErr
}

func NewAddressStore(conn *pgxpool.Pool) *AddressStore {
	return &AddressStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetMany get all addresses by ids
func (s *AddressStore) GetMany(ctx context.Context, ids []int64) ([]*models.Address, []error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		placeholders[i] = "?"
		args[i] = i
	}

	qeryStmt := "SELECT * from addresses WHERE id IN (" + strings.Join(placeholders, ",") + ")"

	// errMsg := "error when trying to get addresses"

	rows, err := s.conn.Query(ctx, qeryStmt, args...)
	if err != nil {
		var errs []error
		errs = append(errs, err)
		return nil, errs
	}

	addresses, err := s.scanList(rows)
	if err != nil {
		var errs []error
		errs = append(errs, err)
		return nil, errs
	}

	addressList := []*models.Address{}
	for _, a := range addresses {
		addressList = append(addressList, &a)
	}

	return addressList, nil
}

// List gets all addresses
func (s *AddressStore) List(ctx context.Context) ([]models.Address, *faulterr.FaultErr) {
	qeryStmt := `SELECT * FROM addresses`

	addresses := []models.Address{}
	a := models.Address{}
	errMsg := "error when trying to get addresses"

	rows, err := s.conn.Query(ctx, qeryStmt)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	for rows.Next() {
		if err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Tag,
			&a.Line1,
			&a.Line2,
			&a.Line3,
			&a.City,
			&a.State,
			&a.Country,
			&a.Pincode,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, faulterr.NewPostgresError(err, errMsg)
		}
		addresses = append(addresses, a)
	}

	return addresses, nil
}

// ListByUserID gets all addresses
func (s *AddressStore) ListByUserID(ctx context.Context, userID int64) ([]models.Address, *faulterr.FaultErr) {
	qeryStmt := `SELECT * FROM addresses WHERE user_id=$1`

	addresses := []models.Address{}
	a := models.Address{}
	errMsg := "error when trying to get addresses"

	rows, err := s.conn.Query(ctx, qeryStmt, userID)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	for rows.Next() {
		if err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Tag,
			&a.Line1,
			&a.Line2,
			&a.Line3,
			&a.City,
			&a.State,
			&a.Country,
			&a.Pincode,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, faulterr.NewPostgresError(err, errMsg)
		}
		addresses = append(addresses, a)
	}

	return addresses, nil
}

// GetByID Address
func (s *AddressStore) GetByID(ctx context.Context, id int64) (*models.Address, *faulterr.FaultErr) {
	queryStmt := `SELECT * FROM address WHERE id=$1`

	a := models.Address{}
	errMsg := "error when trying to get address by id"

	row := s.conn.QueryRow(ctx, queryStmt, id)
	if err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.Tag,
		&a.Line1,
		&a.Line2,
		&a.Line3,
		&a.City,
		&a.State,
		&a.Country,
		&a.Pincode,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return &a, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Mutate****///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// Insert address
func (s *AddressStore) Insert(ctx context.Context, tx pgx.Tx, a models.Address) (*models.Address, *faulterr.FaultErr) {
	queryStmt := `
	INSERT INTO
	addresses(user_id, tag, line_1, line_2, line_3, city, state, country, pincode)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING *
	`

	errMsg := "error when trying to insert address"
	row := tx.QueryRow(ctx, queryStmt,
		a.UserID,
		a.Tag,
		a.Line1,
		a.Line2,
		a.Line3,
		a.City,
		a.State,
		a.Country,
		a.Pincode,
	)
	if err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.Tag,
		&a.Line1,
		&a.Line2,
		&a.Line3,
		&a.City,
		&a.State,
		&a.Country,
		&a.Pincode,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return &a, nil
}

// Update Address
func (s *AddressStore) Update(ctx context.Context, tx pgx.Tx, a models.Address) *faulterr.FaultErr {
	queryStmt := `
	UPDATE addresses
	SET 
		tag=$1
		line_1=$2
		line_2=$3
		line_3=$4
		city=$5
		state=$6
		country=$7
		pincode=$8
	WHERE addresses.id=$9
	`

	errMsg := "error when trying to update address"
	_, err := tx.Exec(ctx, queryStmt,
		a.Tag,
		a.Line1,
		a.Line2,
		a.Line3,
		a.City,
		a.State,
		a.Country,
		a.Pincode,
	)
	if err != nil {
		return faulterr.NewPostgresError(err, errMsg)
	}

	return nil
}

// Delete Address
func (s *AddressStore) Delete(ctx context.Context, tx pgx.Tx, userID int64) *faulterr.FaultErr {
	queryStmt := `DELETE FROM addresses WHERE id=$1`
	_, err := tx.Exec(ctx, queryStmt, userID)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to delete address")
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Helpers****//////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

func (s *AddressStore) scanList(rows pgx.Rows) ([]models.Address, error) {
	addresses := []models.Address{}
	a := models.Address{}

	for rows.Next() {
		if err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Tag,
			&a.Line1,
			&a.Line2,
			&a.Line3,
			&a.City,
			&a.State,
			&a.Country,
			&a.Pincode,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}

	return addresses, nil
}

func (s *AddressStore) scanRow(row pgx.Row) (*models.Address, error) {
	a := &models.Address{}
	if err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.Tag,
		&a.Line1,
		&a.Line2,
		&a.Line3,
		&a.City,
		&a.State,
		&a.Country,
		&a.Pincode,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return a, nil
}
