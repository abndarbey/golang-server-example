package faulterr

// NewBadRequestError structure
func NewBadRequestError(message string) *FaultErr {
	var err error
	return badRequestErr(message, err)
}

// NewUnauthorizedError structure
func NewUnauthorizedError(message string) *FaultErr {
	var err error
	return unauthorizedErr(message, err)
}

// NewNotFoundError structure
func NewNotFoundError(message string) *FaultErr {
	var err error
	return notFoundErr(message, err)
}

// NewUnprocessableEntityError structure
func NewUnprocessableEntityError(message string) *FaultErr {
	var err error
	return unprocessableEntityErr(message, err)
}

// NewInternalServerError structure
func NewInternalServerError(message string) *FaultErr {
	var err error
	return interalServerErr(message, err)
}
