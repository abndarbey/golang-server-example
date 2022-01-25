package models

import (
	"github.com/volatiletech/null"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Password      string `json:"password"`
	IsAdmin       bool   `json:"IsAdmin"`
	IsMember      bool   `json:"isMember"`
	IsCustomer    bool   `json:"isCustomer"`
	ReferralCode  string `json:"referralCode"`
	PurchaseToken string `json:"purchaseToken"`
	OrgName       string `json:"orgName"`
	Website       string `json:"website"`
}

type OrganizationRequest struct {
	OrgName   string `json:"orgName"`
	Website   string `json:"website"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type RoleCreateRequest struct {
	Name           string   `json:"name"`
	Permissions    []string `json:"permissions"`
	IsOrgAdmin     bool     `json:"isOrgAdmin"`
	OrganizationID int64    `json:"organizationID"`
}

type RoleUpdateRequest struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	IsOrgAdmin  bool     `json:"isOrgAdmin"`
	IsArchived  bool     `json:"isArchived"`
}

type SuperAdminRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type MemberRequest struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Password       string `json:"password"`
	OrganizationID int64  `json:"organizationID"`
	RoleID         int64  `json:"roleID"`
}

type CustomerRequest struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Password      string `json:"password"`
	ReferralCode  string `json:"referralCode"`
	PurchaseToken string `json:"purchaseToken"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
}

type ContainerRequest struct {
	Description    string     `json:"description"`
	IsArchived     bool       `json:"isArchived"`
	OrganizationID null.Int64 `json:"organizationID"`
}

type PalletRequest struct {
	Code           string     `json:"code"`
	Description    string     `json:"description"`
	ContainerID    null.Int64 `json:"containerID"`
	IsArchived     bool       `json:"isArchived"`
	OrganizationID null.Int64 `json:"organizationID"`
}
