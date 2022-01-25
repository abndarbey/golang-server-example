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

type ProfileStore struct {
	conn *pgxpool.Pool
}

var _ ProfileStoreInterface = &ProfileStore{}

type ProfileStoreInterface interface {
	GetMany(ctx context.Context, userIDs []int64) ([]*models.Profile, error)
	List(ctx context.Context) ([]models.Profile, *faulterr.FaultErr)
	GetByUserID(ctx context.Context, userID int64) (*models.Profile, *faulterr.FaultErr)
	GetByReferralCode(ctx context.Context, refCode string) (*models.Profile, *faulterr.FaultErr)
	Insert(ctx context.Context, tx pgx.Tx, p models.Profile) (*models.Profile, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, p models.Profile) *faulterr.FaultErr
	Delete(ctx context.Context, tx pgx.Tx, userID int64) *faulterr.FaultErr
}

func NewProfileStore(conn *pgxpool.Pool) *ProfileStore {
	return &ProfileStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetMany get all profiles by userIDs
func (s *ProfileStore) GetMany(ctx context.Context, userIDs []int64) ([]*models.Profile, error) {
	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, len(userIDs))
	for i := 0; i < len(userIDs); i++ {
		index := strconv.Itoa(i + 1)
		placeholders[i] = "$" + index
		args[i] = userIDs[i]
	}

	qeryStmt := "SELECT * from profiles WHERE user_id IN (" + strings.Join(placeholders, ",") + ")"

	// errMsg := "error when trying to get profiles"

	rows, err := s.conn.Query(ctx, qeryStmt, args...)
	if err != nil {
		return nil, err
	}

	profiles, err := s.scanList(rows)
	if err != nil {
		return nil, err
	}

	result := []*models.Profile{}
	for i := 0; i < len(profiles); i++ {
		result = append(result, &profiles[i])
	}

	return result, nil
}

// List gets all profiles
func (s *ProfileStore) List(ctx context.Context) ([]models.Profile, *faulterr.FaultErr) {
	qeryStmt := `SELECT * FROM profiles`

	profiles := []models.Profile{}
	p := models.Profile{}
	errMsg := "error when trying to get profiles"

	rows, err := s.conn.Query(ctx, qeryStmt)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	for rows.Next() {
		if err := rows.Scan(
			&p.UserID,
			&p.DateOfBirth,
			&p.ReferralCode,
			&p.WalletPoints,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, faulterr.NewPostgresError(err, errMsg)
		}
		profiles = append(profiles, p)
	}

	return profiles, nil
}

// GetByUserID Profile
func (s *ProfileStore) GetByUserID(ctx context.Context, userID int64) (*models.Profile, *faulterr.FaultErr) {
	queryStmt := `SELECT * FROM profiles WHERE user_id=$1`

	p := models.Profile{}
	errMsg := "error when trying to get profile by user_id"

	row := s.conn.QueryRow(ctx, queryStmt, userID)
	if err := row.Scan(
		&p.UserID,
		&p.DateOfBirth,
		&p.ReferralCode,
		&p.WalletPoints,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return &p, nil
}

// GetByReferralCode Profile
func (s *ProfileStore) GetByReferralCode(ctx context.Context, refCode string) (*models.Profile, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM profiles
	WHERE profiles.referral_code=$1
	`

	p := models.Profile{}
	errMsg := "error when trying to get user with referral code"

	row := s.conn.QueryRow(ctx, queryStmt, refCode)
	if err := row.Scan(
		&p.UserID,
		&p.DateOfBirth,
		&p.ReferralCode,
		&p.WalletPoints,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return &p, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Mutate****///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// Insert profile
func (s *ProfileStore) Insert(ctx context.Context, tx pgx.Tx, p models.Profile) (*models.Profile, *faulterr.FaultErr) {
	queryStmt := `
	INSERT INTO
	profiles(user_id, date_of_birth, wallet_points, referral_code)
	VALUES ($1, $2, $3, $4)
	RETURNING *
	`

	errMsg := "error when trying to insert profile"
	_, err := tx.Exec(ctx, queryStmt,
		p.UserID,
		p.DateOfBirth,
		p.WalletPoints,
		p.ReferralCode,
	)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	return &p, nil
}

// Update Profile
func (s *ProfileStore) Update(ctx context.Context, tx pgx.Tx, p models.Profile) *faulterr.FaultErr {
	queryStmt := `
	UPDATE profiles
	SET 
		date_of_birth=$1
		wallet_points=$2
	WHERE profiles.user_id=$3
	`

	errMsg := "error when trying to update profile"
	_, err := tx.Exec(ctx, queryStmt,
		&p.DateOfBirth,
		&p.WalletPoints,
		&p.UserID,
	)
	if err != nil {
		return faulterr.NewPostgresError(err, errMsg)
	}

	return nil
}

// Delete Profile
func (s *ProfileStore) Delete(ctx context.Context, tx pgx.Tx, userID int64) *faulterr.FaultErr {
	queryStmt := `DELETE FROM profiles WHERE user_id=$1`
	_, err := tx.Exec(ctx, queryStmt, userID)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to delete profile")
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Helpers****//////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

func (s *ProfileStore) scanList(rows pgx.Rows) ([]models.Profile, error) {
	profiles := []models.Profile{}
	p := models.Profile{}

	for rows.Next() {
		if err := rows.Scan(
			&p.UserID,
			&p.DateOfBirth,
			&p.ReferralCode,
			&p.WalletPoints,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}

	return profiles, nil
}

func (s *ProfileStore) scanRow(row pgx.Row) (*models.Profile, error) {
	p := &models.Profile{}
	if err := row.Scan(
		&p.UserID,
		&p.DateOfBirth,
		&p.ReferralCode,
		&p.WalletPoints,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return p, nil
}
