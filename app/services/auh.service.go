package services

import (
	"context"
	"orijinplus/app/master"
	"orijinplus/app/models"
	"orijinplus/app/store/dbstore"
	"orijinplus/utils/encrypt"
	"orijinplus/utils/faulterr"
)

type AuthService struct {
	dbstore *dbstore.DBStore
	master  *master.Master
}

func NewAuthService(dbstore *dbstore.DBStore, master *master.Master) *AuthService {
	return &AuthService{dbstore, master}
}

// Login validates password and returns user
func (s *AuthService) Login(ctx context.Context, r *models.LoginRequest) (*models.Auther, *faulterr.FaultErr) {
	u, err := s.dbstore.UserStore.GetByEmail(ctx, r.Email)
	if err != nil {
		return nil, err
	}
	if u.PasswordHash != encrypt.GetMd5(r.Password) {
		return nil, faulterr.NewBadRequestError("wrong password")
	}

	auther := &models.Auther{
		ID:             u.ID,
		IsAdmin:        u.IsAdmin,
		IsMember:       u.IsMember,
		IsCustomer:     u.IsCustomer,
		OrganizationID: u.OrganizationID,
		RoleID:         u.RoleID,
	}

	return auther, nil
}

// RegisterAdmin creates a user as an admin
func (s *AuthService) RegisterAdmin(ctx context.Context, r models.SuperAdminRequest) (*models.User, *faulterr.FaultErr) {
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	u, err := s.master.UserMaster.CreateSuperAdmin(ctx, tx, r)
	if err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return u, nil
}

// RegisterMember creates a user as a member
func (s *AuthService) RegisterMember(ctx context.Context, r models.MemberRequest) (*models.User, *faulterr.FaultErr) {
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	org, err := s.dbstore.OrganizationStore.GetByID(ctx, r.OrganizationID)
	if err != nil {
		return nil, err
	}

	role, err := s.dbstore.RoleStore.GetByID(ctx, r.RoleID)
	if err != nil {
		return nil, err
	}

	member := models.MemberRequest{
		FirstName:      r.FirstName,
		LastName:       r.LastName,
		Email:          r.Email,
		Phone:          r.Phone,
		Password:       r.Password,
		OrganizationID: org.ID,
		RoleID:         role.ID,
	}

	u, err := s.master.UserMaster.CreateMember(ctx, tx, member)
	if err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return u, nil
}

// RegisterCustomer creates a user as a consumer
func (s *AuthService) RegisterCustomer(ctx context.Context, r models.CustomerRequest) (*models.User, *faulterr.FaultErr) {
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	u, err := s.master.UserMaster.CreateCustomer(ctx, tx, r)
	if err != nil {
		return nil, err
	}

	_, err = s.master.UserMaster.CreateProfile(ctx, tx, u)
	if err != nil {
		return nil, err
	}

	// if r.ReferralCode != "" {
	// 	if err := s.master.ReferralMaster.Create(ctx, tx, u, r.ReferralCode); err != nil {
	// 		return nil, err
	// 	}
	// }

	// s.master.PurchaseRecordMaster.RedeemByEmail(ctx, tx, u)
	// s.master.PurchaseRecordMaster.RedeemByPhone(ctx, tx, u)
	// if r.PurchaseToken != "" {
	// 	s.master.PurchaseRecordMaster.RedeemByToken(ctx, tx, u.ID, r.PurchaseToken)
	// }

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return u, nil
}

// RegisterOrganization creates a user as a consumer
func (s *AuthService) RegisterOrganization(ctx context.Context, r models.OrganizationRequest) (*models.User, *faulterr.FaultErr) {
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	org, err := s.master.OrganizationMaster.Create(ctx, tx, r)
	if err != nil {
		return nil, err
	}

	rr := models.RoleCreateRequest{
		Name:           models.OrgAdmin,
		Permissions:    nil,
		IsOrgAdmin:     true,
		OrganizationID: org.ID,
	}
	role, err := s.master.RoleMaster.Create(ctx, tx, rr)
	if err != nil {
		return nil, err
	}

	member := models.MemberRequest{
		FirstName:      r.FirstName,
		LastName:       r.LastName,
		Email:          r.Email,
		Phone:          r.Phone,
		Password:       r.Password,
		OrganizationID: org.ID,
		RoleID:         role.ID,
	}

	u, err := s.master.UserMaster.CreateMember(ctx, tx, member)
	if err != nil {
		return nil, err
	}

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *AuthService) UpdatePassword(ctx context.Context, r models.UpdatePasswordRequest, auther *models.Auther) (*models.User, *faulterr.FaultErr) {
	return s.master.UserMaster.UpdatePassword(ctx, r, auther)
}

// GetUserByID returns a user
func (s *AuthService) GetUserByID(ctx context.Context, id int64) (*models.User, *faulterr.FaultErr) {
	return s.dbstore.UserStore.GetByID(ctx, id)
}

// GrantPermission verifies the member's permission and returns unauthorized error if not permitted
func (s *AuthService) GrantPermission(
	ctx context.Context,
	auther *models.Auther,
	perm string,
	memberView bool,
	consumerView bool,
) *faulterr.FaultErr {
	errMsg := "user not authorized"

	if !auther.IsAdmin && !memberView && !consumerView {
		return faulterr.NewUnauthorizedError(errMsg)
	}

	if !consumerView && auther.IsCustomer {
		return faulterr.NewUnauthorizedError(errMsg)
	}

	if !memberView && auther.IsMember {
		return faulterr.NewUnauthorizedError(errMsg)
	}

	if memberView && auther.IsMember {
		role, err := s.dbstore.RoleStore.GetByID(ctx, auther.RoleID.Int64)
		if err != nil {
			return err
		}

		if !role.IsOrgAdmin {
			for _, permission := range role.Permissions {
				if permission == perm {
					return nil
				}
			}
			return faulterr.NewUnauthorizedError(errMsg)
		}

		return nil
	}

	return nil
}

// Helpers
func (s *AuthService) ValidateLoginRequest(r models.LoginRequest) *faulterr.FaultErr {
	if r.Email == "" && r.Phone == "" {
		return faulterr.NewBadRequestError("Email or Phone is required")
	}
	if r.Password == "" {
		return faulterr.NewBadRequestError("Password is required")
	}
	return nil
}
