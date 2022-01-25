package services

import (
	"orijinplus/app/master"
	"orijinplus/app/store/blockchain"
	"orijinplus/app/store/dbstore"
)

type Services struct {
	AuthService         *AuthService
	OrganizationService *OrganizationService
	RoleService         *RoleService
	UserService         *UserService
	ContainerService    *ContainerService
	PalletService       *PalletService
}

func NewService(
	dbstore *dbstore.DBStore,
	blk *blockchain.BlockchainStore,
	master *master.Master,
) *Services {
	return &Services{
		NewAuthService(dbstore, master),
		NewOrganizationService(dbstore, master),
		NewRoleService(dbstore, master),
		NewUserService(dbstore, master),
		NewContainerService(dbstore, master),
		NewPalletService(dbstore, master),
	}
}
