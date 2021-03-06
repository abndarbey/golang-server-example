// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graph

import (
	"fmt"
	"io"
	"orijinplus/app/models"
	"strconv"

	"github.com/volatiletech/null"
)

type ContainerResult struct {
	Containers []models.Container `json:"containers"`
	Total      int                `json:"total"`
}

type FileInput struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type NewCustomer struct {
	FirstName     *null.String `json:"firstName"`
	LastName      *null.String `json:"lastName"`
	Email         *null.String `json:"email"`
	Phone         *null.String `json:"phone"`
	Password      *null.String `json:"password"`
	ReferralCode  *null.String `json:"ReferralCode"`
	PurchaseToken *null.String `json:"purchaseToken"`
}

type NewMember struct {
	FirstName      *null.String `json:"firstName"`
	LastName       *null.String `json:"lastName"`
	Email          *null.String `json:"email"`
	Phone          *null.String `json:"phone"`
	Password       *null.String `json:"password"`
	OrganizationID *null.String `json:"organizationID"`
	RoleID         *null.String `json:"roleID"`
}

type NewRole struct {
	Name           string   `json:"name"`
	IsOrgAdmin     bool     `json:"isOrgAdmin"`
	OrganizationID *int64   `json:"organizationID"`
	Permissions    []string `json:"permissions"`
}

type NewSuperAdmin struct {
	FirstName *null.String `json:"firstName"`
	LastName  *null.String `json:"lastName"`
	Email     *null.String `json:"email"`
	Phone     *null.String `json:"phone"`
	Password  *null.String `json:"password"`
}

type OrganizationsResult struct {
	Organizations []models.Organization `json:"organizations"`
	Total         int                   `json:"total"`
}

type PageInfo struct {
	StartCursor int64 `json:"startCursor"`
	EndCursor   int64 `json:"endCursor"`
}

type PalletResult struct {
	Pallets []models.Pallet `json:"pallets"`
	Total   int             `json:"total"`
}

type RolesResult struct {
	Roles []models.Role `json:"roles"`
	Total int           `json:"total"`
}

type SearchFilter struct {
	Search  *null.String  `json:"search"`
	Filter  *FilterOption `json:"filter"`
	SortBy  *SortByOption `json:"sortBy"`
	SortDir *SortDir      `json:"sortDir"`
}

type UpdateContainer struct {
	Description    *null.String `json:"description"`
	OrganizationID *null.Int64  `json:"organizationID"`
}

type UpdateOrganization struct {
	Name       *null.String `json:"name"`
	Website    *null.String `json:"website"`
	IsArchived *null.Bool   `json:"isArchived"`
}

type UpdatePallet struct {
	Description    *null.String `json:"description"`
	ContainerID    *null.Int64  `json:"containerID"`
	OrganizationID *null.Int64  `json:"organizationID"`
}

type UpdateRole struct {
	Name        *null.String `json:"name"`
	Permissions []string     `json:"permissions"`
	IsArchived  *null.Bool   `json:"isArchived"`
}

type UpdateUser struct {
	FirstName *null.String `json:"firstName"`
	LastName  *null.String `json:"lastName"`
	Email     *null.String `json:"email"`
	Phone     *null.String `json:"phone"`
	RoleID    *null.String `json:"roleID"`
	Password  *null.String `json:"password"`
}

type UserResult struct {
	Users []models.User `json:"users"`
	Total int           `json:"total"`
}

type FilterOption string

const (
	FilterOptionAll                    FilterOption = "All"
	FilterOptionActive                 FilterOption = "Active"
	FilterOptionArchived               FilterOption = "Archived"
	FilterOptionProductWithoutOrder    FilterOption = "ProductWithoutOrder"
	FilterOptionProductWithoutCarton   FilterOption = "ProductWithoutCarton"
	FilterOptionProductWithoutSku      FilterOption = "ProductWithoutSKU"
	FilterOptionCartonWithoutPallet    FilterOption = "CartonWithoutPallet"
	FilterOptionPalletWithoutContainer FilterOption = "PalletWithoutContainer"
	FilterOptionSystem                 FilterOption = "System"
	FilterOptionBlockchain             FilterOption = "Blockchain"
	FilterOptionPending                FilterOption = "Pending"
)

var AllFilterOption = []FilterOption{
	FilterOptionAll,
	FilterOptionActive,
	FilterOptionArchived,
	FilterOptionProductWithoutOrder,
	FilterOptionProductWithoutCarton,
	FilterOptionProductWithoutSku,
	FilterOptionCartonWithoutPallet,
	FilterOptionPalletWithoutContainer,
	FilterOptionSystem,
	FilterOptionBlockchain,
	FilterOptionPending,
}

func (e FilterOption) IsValid() bool {
	switch e {
	case FilterOptionAll, FilterOptionActive, FilterOptionArchived, FilterOptionProductWithoutOrder, FilterOptionProductWithoutCarton, FilterOptionProductWithoutSku, FilterOptionCartonWithoutPallet, FilterOptionPalletWithoutContainer, FilterOptionSystem, FilterOptionBlockchain, FilterOptionPending:
		return true
	}
	return false
}

func (e FilterOption) String() string {
	return string(e)
}

func (e *FilterOption) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FilterOption(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FilterOption", str)
	}
	return nil
}

func (e FilterOption) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortByOption string

const (
	SortByOptionDateCreated  SortByOption = "DateCreated"
	SortByOptionDateUpdated  SortByOption = "DateUpdated"
	SortByOptionAlphabetical SortByOption = "Alphabetical"
)

var AllSortByOption = []SortByOption{
	SortByOptionDateCreated,
	SortByOptionDateUpdated,
	SortByOptionAlphabetical,
}

func (e SortByOption) IsValid() bool {
	switch e {
	case SortByOptionDateCreated, SortByOptionDateUpdated, SortByOptionAlphabetical:
		return true
	}
	return false
}

func (e SortByOption) String() string {
	return string(e)
}

func (e *SortByOption) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortByOption(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortByOption", str)
	}
	return nil
}

func (e SortByOption) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortDir string

const (
	SortDirAscending  SortDir = "Ascending"
	SortDirDescending SortDir = "Descending"
)

var AllSortDir = []SortDir{
	SortDirAscending,
	SortDirDescending,
}

func (e SortDir) IsValid() bool {
	switch e {
	case SortDirAscending, SortDirDescending:
		return true
	}
	return false
}

func (e SortDir) String() string {
	return string(e)
}

func (e *SortDir) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortDir(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortDir", str)
	}
	return nil
}

func (e SortDir) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
