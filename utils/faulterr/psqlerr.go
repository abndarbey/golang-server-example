package faulterr

// NewPostgresError structure
func NewPostgresError(err error, message string) *FaultErr {
	switch err.Error() {
	case "no rows in result set":
		return notFoundErr(message, err)
	default:
		message := "Something went wrong, please try again"
		return interalServerErr(message, err)
	}
}
