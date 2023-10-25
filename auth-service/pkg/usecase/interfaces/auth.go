package interfaces

import (
	"auth-service/pkg/api/service/response"
	"auth-service/pkg/domain"
	"auth-service/pkg/service/token"
)

type AuthUseCase interface {
	SignUp(user domain.User) (domain.User, error)
	SignIn(user domain.User) (domain.User, error)
	GenerateAccessToken(role string, user domain.User) (response.Token, error)
	VerifyAccessToken(token string) (token.Payload, error)
}
