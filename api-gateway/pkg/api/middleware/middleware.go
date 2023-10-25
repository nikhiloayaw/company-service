package middleware

import (
	"api-gateway/pkg/api/middleware/interfaces"
	clientInterfaces "api-gateway/pkg/client/interfaces"
)

// To hold all middleware and it's dependencies
type middleware struct {
	authClient clientInterfaces.AuthServiceClient
}

func NewMiddleware(authClient clientInterfaces.AuthServiceClient) interfaces.Middleware {

	return &middleware{
		authClient: authClient,
	}
}
