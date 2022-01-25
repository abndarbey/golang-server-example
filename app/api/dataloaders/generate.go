//go:generate go run github.com/vektah/dataloaden UserLoader int64 *orijinplus/app/models.User
//go:generate go run github.com/vektah/dataloaden ProfileLoader int64 *orijinplus/app/models.Profile
//go:generate go run github.com/vektah/dataloaden OrganizationLoader int64 *orijinplus/app/models.Organization
//go:generate go run github.com/vektah/dataloaden RoleLoader int64 *orijinplus/app/models.Role
//go:generate go run github.com/vektah/dataloaden ContainerLoader int64 *orijinplus/app/models.Container
//go:generate go run github.com/vektah/dataloaden PalletLoader int64 *orijinplus/app/models.Pallet

package dataloaders
