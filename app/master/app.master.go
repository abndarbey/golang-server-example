package master

import "orijinplus/app/store/dbstore"

type Master struct {
	OrganizationMaster *OrganizationMaster
	RoleMaster         *RoleMaster
	UserMaster         *UserMaster
	ContainerMaster    *ContainerMaster
	PalletMaster       *PalletMaster
}

func NewMaster(dbStore *dbstore.DBStore) *Master {
	return &Master{
		NewOrganizationMaster(dbStore),
		NewRoleMaster(dbStore),
		NewUserMaster(dbStore),
		NewContainerMaster(dbStore),
		NewPalletMaster(dbStore),
	}
}
