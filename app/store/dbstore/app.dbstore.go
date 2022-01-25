package dbstore

import "github.com/jackc/pgx/v4/pgxpool"

type DBStore struct {
	DBTX              *DBTX
	OrganizationStore *OrganizationStore
	RoleStore         *RoleStore
	UserStore         *UserStore
	ProfileStore      *ProfileStore
	AddressStore      *AddressStore
	ContainerStore    *ContainerStore
	PalletStore       *PalletStore
}

func NewDBStore(conn *pgxpool.Pool) *DBStore {
	return &DBStore{
		NewDBTX(conn),
		NewOrganizationStore(conn),
		NewRoleStore(conn),
		NewUserStore(conn),
		NewProfileStore(conn),
		NewAddressStore(conn),
		NewContainerStore(conn),
		NewPalletStore(conn),
	}
}
