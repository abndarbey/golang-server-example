package faulterr

import (
	"net/http"
	"orijinplus/utils/logger"
)

// badRequestErr structure
func badRequestErr(msg string, err error) *FaultErr {
	logger.Error(err, msg)
	return &FaultErr{
		Message: msg,
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}

// unauthorizedErr structure
func unauthorizedErr(msg string, err error) *FaultErr {
	logger.Error(err, msg)
	return &FaultErr{
		Message: msg,
		Status:  http.StatusUnauthorized,
		Error:   "unauthorized",
	}
}

// notFoundErr structure
func notFoundErr(msg string, err error) *FaultErr {
	// logger.Error(err, msg)
	return &FaultErr{
		Message: msg,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}

// unprocessableEntityErr structure
func unprocessableEntityErr(msg string, err error) *FaultErr {
	logger.Error(err, msg)
	return &FaultErr{
		Message: msg,
		Status:  http.StatusUnprocessableEntity,
		Error:   "unprocessable_entity",
	}
}

// interalServerErr structure
func interalServerErr(msg string, err error) *FaultErr {
	logger.Error(err, msg)
	return &FaultErr{
		Message: msg,
		Status:  http.StatusInternalServerError,
		Error:   "internal_server_error",
	}
}
