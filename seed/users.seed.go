package seed

import (
	"context"
	"fmt"
	"orijinplus/app/models"
	"orijinplus/app/services"
	"orijinplus/utils/faulterr"
	"orijinplus/utils/logger"
)

func InsertAdmin(srv *services.Services) (*models.User, *faulterr.FaultErr) {
	admins := []models.SuperAdminRequest{
		{
			FirstName: "Super",
			LastName:  "Admin",
			Email:     "superadmin@example.com",
			Phone:     "1234567890",
			Password:  "password@123",
		},
		{
			FirstName: "James",
			LastName:  "Williams",
			Email:     "james@example.com",
			Phone:     "1234567891",
			Password:  "password@123",
		},
		{
			FirstName: "Rhys",
			LastName:  "Williams",
			Email:     "rhys@example.com",
			Phone:     "1234567892",
			Password:  "password@123",
		},
	}

	superAdmins := []*models.User{}

	for _, r := range admins {
		ctx := context.Background()
		obj, err := srv.AuthService.RegisterAdmin(ctx, r)
		if err != nil {
			return nil, err
		}

		superAdmins = append(superAdmins, obj)
	}

	logger.Success(fmt.Sprintf("%v super admins added to the database", len(admins)))

	return superAdmins[0], nil
}
