package response

import (
	"strings"
)

type Response struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Error      any    `json:"error,omitempty"`
	Data       any    `json:"data,omitempty"`
}

// To create response structure for success response
func SuccessResponse(statusCode int, message string, data any) Response {
	return Response{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

// To create response structure for error response
func ErrorResponse(statusCode int, message string, err error, data any) Response {
	// separate error string by newline
	errArr := strings.Split(err.Error(), "\n")
	return Response{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Error:      errArr,
		Data:       data,
	}
}
