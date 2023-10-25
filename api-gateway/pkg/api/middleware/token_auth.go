package middleware

import (
	"api-gateway/pkg/api/handler/response"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderKey = "Authorization"
	AuthorizationType      = "Bearer"
)

var (
	ErrInvalidHeaderValues      = errors.New("header doesn't contain token and authorization type properly")
	ErrInvalidAuthorizationType = errors.New("invalid authorization type")
	ErrRoleNotMatching          = errors.New("request user role not matching")
)

// To authenticate user request with token and role
func (m *middleware) Authenticate(role string) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// get the token from request header
		authorizationValue := ctx.GetHeader(AuthorizationHeaderKey)

		// split token and authorization type (Bearer)
		authFields := strings.Split(authorizationValue, " ")

		// check header value  minimum contain token and authorization type
		if len(authFields) < 2 {
			response := response.ErrorResponse(http.StatusUnauthorized, "failed to authorize request", ErrInvalidHeaderValues, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// check the authorization type
		if authFields[0] != AuthorizationType {
			response := response.ErrorResponse(http.StatusUnauthorized, "failed to authorize request", ErrInvalidAuthorizationType, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// take access token and validate
		accessToken := authFields[1]
		// call the auth verify client to verify token
		payload, err := m.authClient.VerifyAccessToken(ctx, accessToken)

		if err != nil {
			response := response.ErrorResponse(http.StatusUnauthorized, "failed to authorize request", err, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// verify the role
		if payload.Role != role {
			response := response.ErrorResponse(http.StatusUnauthorized, "failed to authorize request", ErrRoleNotMatching, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// set payload values to context for next handler func to use.
		ctx.Set("userId", payload.UserID)
		ctx.Set("role", payload.Role)
		ctx.Set("email", payload.Email)
	}
}
