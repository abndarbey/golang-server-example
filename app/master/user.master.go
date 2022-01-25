package master

import (
	"context"
	"orijinplus/app/models"
	"orijinplus/app/store/dbstore"
	"orijinplus/utils/encrypt"
	"orijinplus/utils/faulterr"

	"github.com/jackc/pgx/v4"
	"github.com/volatiletech/null"
)

type UserMaster struct {
	dbstore *dbstore.DBStore
}

func NewUserMaster(dbstore *dbstore.DBStore) *UserMaster {
	return &UserMaster{dbstore}
}

// CreateUser creates and saves user in the db
func (m *UserMaster) CreateUser(ctx context.Context, tx pgx.Tx, r models.RegisterRequest) (*models.User, *faulterr.FaultErr) {
	if err := m.validateRegisterRequest(r); err != nil {
		return nil, err
	}

	u := models.User{
		FirstName:    r.FirstName,
		LastName:     r.LastName,
		Email:        r.Email,
		Phone:        r.Phone,
		PasswordHash: encrypt.GetMd5(r.Password),
		IsAdmin:      r.IsAdmin,
		IsMember:     r.IsMember,
		IsCustomer:   r.IsCustomer,
	}

	err := m.verifyUniqueFields(ctx, u)
	if err != nil {
		return nil, err
	}

	return m.dbstore.UserStore.Insert(ctx, tx, u)
}

// CreateMember creates and saves user in the db
func (m *UserMaster) CreateMember(ctx context.Context, tx pgx.Tx, r models.MemberRequest) (*models.User, *faulterr.FaultErr) {
	if err := m.validateMemberRequest(r); err != nil {
		return nil, err
	}

	u := models.User{
		FirstName:    r.FirstName,
		LastName:     r.LastName,
		Email:        r.Email,
		Phone:        r.Phone,
		IsAdmin:      false,
		IsMember:     true,
		IsCustomer:   false,
		PasswordHash: encrypt.GetMd5(r.Password),
	}
	u.OrganizationID = null.Int64From(r.OrganizationID)
	u.RoleID = null.Int64From(r.RoleID)

	err := m.verifyUniqueFields(ctx, u)
	if err != nil {
		return nil, err
	}

	return m.dbstore.UserStore.Insert(ctx, tx, u)
}

// CreateSuperAdmin creates and saves user in the db
func (m *UserMaster) CreateSuperAdmin(ctx context.Context, tx pgx.Tx, r models.SuperAdminRequest) (*models.User, *faulterr.FaultErr) {
	if err := m.validateSuperAdminRequest(r); err != nil {
		return nil, err
	}

	u := models.User{
		FirstName:    r.FirstName,
		LastName:     r.LastName,
		Email:        r.Email,
		Phone:        r.Phone,
		IsAdmin:      true,
		IsMember:     false,
		IsCustomer:   false,
		PasswordHash: encrypt.GetMd5(r.Password),
	}

	err := m.verifyUniqueFields(ctx, u)
	if err != nil {
		return nil, err
	}

	return m.dbstore.UserStore.Insert(ctx, tx, u)
}

// CreateCustomer creates and saves user in the db
func (m *UserMaster) CreateCustomer(ctx context.Context, tx pgx.Tx, r models.CustomerRequest) (*models.User, *faulterr.FaultErr) {
	if err := m.validateCustomerRequest(r); err != nil {
		return nil, err
	}

	u := models.User{
		FirstName:    r.FirstName,
		LastName:     r.LastName,
		Email:        r.Email,
		Phone:        r.Phone,
		IsAdmin:      false,
		IsMember:     false,
		IsCustomer:   true,
		PasswordHash: encrypt.GetMd5(r.Password),
	}

	err := m.verifyUniqueFields(ctx, u)
	if err != nil {
		return nil, err
	}

	return m.dbstore.UserStore.Insert(ctx, tx, u)
}

// CreateProfile creates and saves user parofile in the db
func (m *UserMaster) CreateProfile(ctx context.Context, tx pgx.Tx, u *models.User) (*models.Profile, *faulterr.FaultErr) {
	p := models.Profile{
		UserID:       u.ID,
		WalletPoints: 0,
		ReferralCode: null.StringFrom(encrypt.GenerateRandomString(10)),
	}

	return m.dbstore.ProfileStore.Insert(ctx, tx, p)
}

// UpdatePassword updates the user password hash in db
func (m *UserMaster) UpdatePassword(ctx context.Context, r models.UpdatePasswordRequest, auther *models.Auther) (*models.User, *faulterr.FaultErr) {
	if err := m.validatePasswordRequest(r); err != nil {
		return nil, err
	}

	// Get user user from db
	u, err := m.dbstore.UserStore.GetByID(ctx, auther.ID)
	if err != nil {
		return nil, err
	}
	passwordHash := encrypt.GetMd5(r.Password)
	if passwordHash == u.PasswordHash {
		return nil, faulterr.NewBadRequestError("This is your current password, please enter new password")
	}
	u.PasswordHash = passwordHash

	// Start db transaction
	tx, err := m.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer m.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := m.dbstore.UserStore.Update(ctx, tx, *u); err != nil {
		return nil, err
	}

	if err := m.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}
	return u, nil
}

// Helpers

// verifyUniqueFields verifies the uniqueness of user
func (m *UserMaster) verifyUniqueFields(ctx context.Context, u models.User) *faulterr.FaultErr {
	// Verify unique email
	_, err := m.dbstore.UserStore.GetByEmail(ctx, u.Email)
	if err == nil {
		return faulterr.NewBadRequestError("email already registered")
	}

	// Verify unique phone
	_, err = m.dbstore.UserStore.GetByPhone(ctx, u.Phone)
	if err == nil {
		return faulterr.NewBadRequestError("phone already registered")
	}

	return nil
}

// Validation

// verifyUniqueFields verifies the uniqueness of user
func (m *UserMaster) validateRegisterRequest(r models.RegisterRequest) *faulterr.FaultErr {
	if r.FirstName == "" {
		return faulterr.NewBadRequestError("First Name is required")
	}
	if r.LastName == "" {
		return faulterr.NewBadRequestError("Last Name is required")
	}
	if r.Email == "" {
		return faulterr.NewBadRequestError("Email is required")
	}
	if r.Password == "" {
		return faulterr.NewBadRequestError("Password is required")
	}

	return nil
}

// ValidateSuperAdminRequest validates the super admin request
func (m *UserMaster) validateSuperAdminRequest(r models.SuperAdminRequest) *faulterr.FaultErr {
	if r.FirstName == "" {
		return faulterr.NewBadRequestError("First Name is required")
	}
	if r.LastName == "" {
		return faulterr.NewBadRequestError("Last Name is required")
	}
	if r.Email == "" {
		return faulterr.NewBadRequestError("Email is required")
	}
	if r.Password == "" {
		return faulterr.NewBadRequestError("Password is required")
	}
	return nil
}

// ValidateMemberRequest validates the member request
func (m *UserMaster) validateMemberRequest(r models.MemberRequest) *faulterr.FaultErr {
	if r.FirstName == "" {
		return faulterr.NewBadRequestError("First Name is required")
	}
	if r.LastName == "" {
		return faulterr.NewBadRequestError("Last Name is required")
	}
	if r.Email == "" {
		return faulterr.NewBadRequestError("Email is required")
	}
	if r.OrganizationID <= 0 {
		return faulterr.NewBadRequestError("Organization ID is required")
	}
	if r.RoleID <= 0 {
		return faulterr.NewBadRequestError("Role ID is required")
	}
	if r.Password == "" {
		return faulterr.NewBadRequestError("Password is required")
	}
	return nil
}

// ValidateCustomerRequest validates the consumer request
func (m *UserMaster) validateCustomerRequest(r models.CustomerRequest) *faulterr.FaultErr {
	if r.FirstName == "" {
		return faulterr.NewBadRequestError("First Name is required")
	}
	if r.LastName == "" {
		return faulterr.NewBadRequestError("Last Name is required")
	}
	if r.Email == "" {
		return faulterr.NewBadRequestError("Email is required")
	}
	if r.Password == "" {
		return faulterr.NewBadRequestError("Password is required")
	}
	return nil
}

// verifyUniqueFields verifies the uniqueness of user
func (m *UserMaster) validatePasswordRequest(r models.UpdatePasswordRequest) *faulterr.FaultErr {
	if r.Password == "" {
		return faulterr.NewBadRequestError("Password is required")
	}

	return nil
}
