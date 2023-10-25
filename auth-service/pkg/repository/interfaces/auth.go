package interfaces

import (
	"auth-service/pkg/domain"
)

type AuthRepo interface {
	IsUserExist(email string) (bool, error)
	FindUserByEmail(email string) (domain.User, error)
	SaveUser(user domain.User) (domain.User, error)
}
