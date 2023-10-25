package utils

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetHTTPStatusCodeFromGRPCError(err error) int {

	var statusCode int

	// convert the error to grpc code
	switch status.Code(err) {

	case codes.AlreadyExists:
		statusCode = http.StatusConflict
	case codes.Internal:
		statusCode = http.StatusInternalServerError
	case codes.Unauthenticated:
		statusCode = http.StatusUnauthorized
	case codes.InvalidArgument:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusServiceUnavailable
	}

	return statusCode
}
