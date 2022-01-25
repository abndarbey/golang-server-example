package services

import (
	"context"
	"orijinplus/app/master"
	"orijinplus/app/models"
	"orijinplus/app/store/dbstore"
	"orijinplus/utils/faulterr"

	"github.com/volatiletech/null"
)

type UserService struct {
	dbstore *dbstore.DBStore
	master  *master.Master
}

var _ UserServiceInterface = &UserService{}

type UserServiceInterface interface {
	List(
		ctx context.Context,
		isAdmin, isMember, isCustomer bool,
		organizationID null.Int64,
		auther *models.Auther,
	) ([]models.User, *faulterr.FaultErr)
	ListAdmin(ctx context.Context) ([]models.User, *faulterr.FaultErr)
	ListMembers(ctx context.Context, auther *models.Auther) ([]models.User, *faulterr.FaultErr)
	ListCustomers(ctx context.Context) ([]models.User, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64, auther *models.Auther) (*models.User, *faulterr.FaultErr)
	GetProfile(ctx context.Context, auther *models.Auther) (*models.Profile, *faulterr.FaultErr)
	Update(ctx context.Context, request models.User, auther *models.Auther) (*models.User, *faulterr.FaultErr)
	Delete(ctx context.Context, id int64, auther *models.Auther) *faulterr.FaultErr
}

func NewUserService(s *dbstore.DBStore, m *master.Master) *UserService {
	return &UserService{s, m}
}

// List gets all admin, members, and consumers
func (s *UserService) List(
	ctx context.Context,
	isAdmin, isMember, isCustomer bool,
	organizationID null.Int64,
	auther *models.Auther,
) ([]models.User, *faulterr.FaultErr) {
	if organizationID.Valid {
		if auther.IsAdmin {
			return s.dbstore.UserStore.ListMembersByOrgID(ctx, organizationID.Int64)
		}
		if auther.IsMember {
			return s.dbstore.UserStore.ListMembersByOrgID(ctx, auther.OrganizationID.Int64)
		}
	}
	return s.dbstore.UserStore.List(ctx, isAdmin, isMember, isCustomer)
}

// ListAdmin gets all admins
func (s *UserService) ListAdmin(ctx context.Context) ([]models.User, *faulterr.FaultErr) {
	return s.dbstore.UserStore.ListAdmin(ctx)
}

// ListMembers gets all members
func (s *UserService) ListMembers(ctx context.Context, auther *models.Auther) ([]models.User, *faulterr.FaultErr) {
	if auther.IsAdmin {
		return s.dbstore.UserStore.ListMembers(ctx)
	}

	return s.dbstore.UserStore.ListMembersByOrgID(ctx, auther.OrganizationID.Int64)
}

// ListCustomers gets all consumers
func (s *UserService) ListCustomers(ctx context.Context) ([]models.User, *faulterr.FaultErr) {
	return s.dbstore.UserStore.ListCustomers(ctx)
}

// ListCustomers gets aa user by id
func (s *UserService) GetByID(ctx context.Context, id int64, auther *models.Auther) (*models.User, *faulterr.FaultErr) {
	if !auther.IsAdmin {
		return nil, faulterr.NewUnauthorizedError("permission not granted")
	}

	user, err := s.dbstore.UserStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ListCustomers gets a user by user profile
func (s *UserService) GetProfile(ctx context.Context, auther *models.Auther) (*models.Profile, *faulterr.FaultErr) {
	if auther.IsCustomer {
		return s.dbstore.ProfileStore.GetByUserID(ctx, auther.ID)
	}
	return nil, faulterr.NewNotFoundError("no profile found")
}

// GetByEmail gets a user by user profile
func (s *UserService) GetByEmail(ctx context.Context, email string, auther *models.Auther) (*models.User, *faulterr.FaultErr) {
	user, err := s.dbstore.UserStore.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if auther.IsMember && user.OrganizationID != auther.OrganizationID {
		return nil, faulterr.NewNotFoundError("user not found")
	}
	if auther.IsCustomer && user.ID != auther.ID {
		return nil, faulterr.NewNotFoundError("user not found")
	}

	return nil, faulterr.NewNotFoundError("no profile found")
}

// GetByPhone gets a user by user profile
func (s *UserService) GetByPhone(ctx context.Context, phone string, auther *models.Auther) (*models.User, *faulterr.FaultErr) {
	user, err := s.dbstore.UserStore.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	if auther.IsMember && user.OrganizationID != auther.OrganizationID {
		return nil, faulterr.NewNotFoundError("user not found")
	}
	if auther.IsCustomer && user.ID != auther.ID {
		return nil, faulterr.NewNotFoundError("user not found")
	}

	return nil, faulterr.NewNotFoundError("no profile found")
}

func (s *UserService) AddAddress(ctx context.Context, a models.Address, auther *models.Auther) (*models.User, *faulterr.FaultErr) {
	if !auther.IsCustomer {
		return nil, faulterr.NewUnauthorizedError("only consumers can add address")
	}
	if err := a.Validate(); err != nil {
		return nil, err
	}
	a.UserID = auther.ID

	// Start db transaction
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}
	_, err = s.dbstore.AddressStore.Insert(ctx, tx, a)
	if err != nil {
		return nil, err
	}

	user, err := s.dbstore.UserStore.GetByID(ctx, auther.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, request models.User, auther *models.Auther) (*models.User, *faulterr.FaultErr) {
	user, err := s.GetByID(ctx, request.ID, auther)
	if err != nil {
		return nil, err
	}

	// Verify request fields
	if request.FirstName != "" {
		user.FirstName = request.FirstName
	}

	if request.LastName != "" {
		user.LastName = request.LastName
	}

	if request.Email != "" && request.Email != user.Email {
		// verify unique email
		_, err = s.dbstore.UserStore.GetByEmail(ctx, request.Email)
		if err == nil {
			return nil, faulterr.NewBadRequestError("email already registered")
		}
		user.Email = request.Email
	}

	if request.Phone != "" && request.Phone != user.Phone {
		// verify unique phone
		_, err = s.dbstore.UserStore.GetByPhone(ctx, request.Phone)
		if err == nil {
			return nil, faulterr.NewBadRequestError("phone already registered")
		}
		user.Phone = request.Phone
	}

	// Start db transaction
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.UserStore.Update(ctx, tx, *user); err != nil {
		return nil, err
	}
	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id int64, auther *models.Auther) *faulterr.FaultErr {
	if !auther.IsAdmin {
		return faulterr.NewUnauthorizedError("permission not granted")
	}
	_, err := s.GetByID(ctx, id, auther)
	if err != nil {
		return err
	}

	// Start db transaction
	tx, err := s.dbstore.DBTX.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer s.dbstore.DBTX.RollbackTx(ctx, tx)

	if err := s.dbstore.UserStore.Delete(ctx, tx, id); err != nil {
		return err
	}
	if err := s.dbstore.DBTX.CommitTx(ctx, tx); err != nil {
		return err
	}

	return nil
}
