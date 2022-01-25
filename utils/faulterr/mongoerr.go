package faulterr

// NewMongoError structure
func NewMongoError(err error, message string) *FaultErr {
	switch err.Error() {
	case "mongo: no documents in result":
		return notFoundErr(message, err)
	default:
		message := "Something went wrong, please try again"
		return interalServerErr(message, err)
	}
}
