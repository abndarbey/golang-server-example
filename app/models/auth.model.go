package models

import "github.com/volatiletech/null"

type Auther struct {
	ID             int64      `json:"id"`
	IsAdmin        bool       `json:"IsAdmin"`
	IsMember       bool       `json:"isMember"`
	IsCustomer     bool       `json:"isCustomer"`
	OrganizationID null.Int64 `json:"organizationID"`
	RoleID         null.Int64 `json:"roleID"`
}

// ValueToken struct
type ValueToken struct {
	TokenString string `json:"tokenString"`
}

// TokenHeader struct
type TokenHeader struct {
	TYP string `json:"typ"`
	ALG string `json:"alg"`
}

// TokenPayload struct
type TokenPayload struct {
	Auther
	Authorized bool  `json:"authorized"`
	Expiry     int64 `json:"exp"`
}
