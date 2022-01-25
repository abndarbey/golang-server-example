package seed

import (
	"context"
	"fmt"
	"orijinplus/app/models"
	"orijinplus/app/services"
	"orijinplus/utils/faulterr"
	"orijinplus/utils/logger"
)

func InsertOrganizations(srv *services.Services) *faulterr.FaultErr {
	organizations := []models.OrganizationRequest{
		{
			OrgName:   "ABL",
			Website:   "www.abl.com",
			FirstName: "Abhinandan",
			LastName:  "Darbey",
			Email:     "abhinandan@example.com",
			Phone:     "7757888393",
			Password:  "password@123",
		},
		{
			OrgName:   "Latitude 28",
			Website:   "www.latitude28.com.au",
			FirstName: "Michael",
			LastName:  "Stone",
			Email:     "michael@example.com",
			Phone:     "1234567893",
			Password:  "password@123",
		},
	}

	for _, r := range organizations {
		ctx := context.Background()
		_, err := srv.AuthService.RegisterOrganization(ctx, r)
		if err != nil {
			return err
		}
	}

	logger.Success(fmt.Sprintf("%v organizations added to the database", len(organizations)))

	return nil
}
