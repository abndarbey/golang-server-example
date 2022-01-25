package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/null"
)

type Address struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"userID"`
	Tag       string      `json:"tag"`
	Line1     string      `json:"line1"`
	Line2     string      `json:"line2"`
	Line3     null.String `json:"line3"`
	City      string      `json:"city"`
	State     string      `json:"state"`
	Country   string      `json:"country"`
	Pincode   string      `json:"pincode"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type Container struct {
	ID             int64      `json:"id"`
	UID            uuid.UUID  `json:"uid"`
	Code           string     `json:"code"`
	Description    string     `json:"description"`
	IsArchived     bool       `json:"isArchived"`
	OrganizationID null.Int64 `json:"organizationID"`
	CreatedByID    int64      `json:"createdByID"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type Organization struct {
	ID         int64       `json:"id"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	Website    null.String `json:"website"`
	IsArchived bool        `json:"isArchived"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

type Pallet struct {
	ID             int64      `json:"id"`
	UID            uuid.UUID  `json:"uid"`
	Code           string     `json:"code"`
	Description    string     `json:"description"`
	ContainerID    null.Int64 `json:"containerID"`
	IsArchived     bool       `json:"isArchived"`
	OrganizationID null.Int64 `json:"organizationID"`
	CreatedByID    int64      `json:"createdByID"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type Permission struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Profile struct {
	UserID       int64       `json:"userID"`
	DateOfBirth  null.String `json:"dateOfBirth"`
	ReferralCode null.String `json:"referralCode"`
	WalletPoints int64       `json:"walletPoints"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

type Role struct {
	ID             int64     `json:"id"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Permissions    []string  `json:"permissions"`
	IsOrgAdmin     bool      `json:"isOrgAdmin"`
	IsArchived     bool      `json:"isArchived"`
	OrganizationID int64     `json:"organizationID"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type User struct {
	ID             int64      `json:"id"`
	FirstName      string     `json:"firstName"`
	LastName       string     `json:"lastName"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	IsAdmin        bool       `json:"isAdmin"`
	IsMember       bool       `json:"isMember"`
	IsCustomer     bool       `json:"isCustomer"`
	PasswordHash   string     `json:"passwordHash"`
	OrganizationID null.Int64 `json:"organizationID"`
	RoleID         null.Int64 `json:"roleID"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
