package models

import "orijinplus/utils/faulterr"

// Validate Address
func (r *Address) Validate() *faulterr.FaultErr {
	if r.Tag == "" {
		return faulterr.NewBadRequestError("Tag is required")
	}
	if r.Line1 == "" {
		return faulterr.NewBadRequestError("Line 1 is required")
	}
	if r.Line2 == "" {
		return faulterr.NewBadRequestError("Line 2 is required")
	}
	if r.City == "" {
		return faulterr.NewBadRequestError("City is required")
	}
	if r.State == "" {
		return faulterr.NewBadRequestError("State is required")
	}
	if r.Country == "" {
		return faulterr.NewBadRequestError("Country is required")
	}
	if r.Pincode == "" {
		return faulterr.NewBadRequestError("Pincode is required")
	}
	return nil
}
