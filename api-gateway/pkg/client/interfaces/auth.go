package interfaces

import (
	"api-gateway/pkg/api/handler/request"
	"api-gateway/pkg/api/handler/response"
	"api-gateway/pkg/models"
	"context"
)

type AuthServiceClient interface {
	SignUp(ctx context.Context, signUpDetails request.SignUp) (response.SignUp, error)
	SignIn(ctx context.Context, signInDetails request.SignIn) (response.SignIn, error)
	VerifyAccessToken(ctx context.Context, accessToken string) (models.TokenPayload, error)
}
