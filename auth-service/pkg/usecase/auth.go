package usecase

import (
	"auth-service/pkg/api/service/response"
	"auth-service/pkg/domain"
	repoInterfaces "auth-service/pkg/repository/interfaces"
	"auth-service/pkg/service/token"
	"auth-service/pkg/usecase/interfaces"
	"auth-service/pkg/utils"
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	accessTokenDuration = 30 * time.Minute
)

var (
	ErrAlreadyExist  = errors.New("resource already exist")
	ErrNotExist      = errors.New("resource not exist")
	ErrWrongPassword = errors.New("wrong password")

	ErrExpired = errors.New("resource got expired")
	ErrInvalid = errors.New("invalid resource")
)

type authUseCase struct {
	authRepo  repoInterfaces.AuthRepo
	tokenAuth token.TokenAuth
}

func NewAuthUseCase(
	authRepo repoInterfaces.AuthRepo,
	tokenAuth token.TokenAuth,
) interfaces.AuthUseCase {

	return &authUseCase{
		authRepo:  authRepo,
		tokenAuth: tokenAuth,
	}
}

func (a *authUseCase) SignUp(user domain.User) (domain.User, error) {

	// check the user already exist or not
	alreadyExist, err := a.authRepo.IsUserExist(user.Email)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to check user already exist in db: %w", err)
	}

	if alreadyExist {
		return domain.User{}, ErrAlreadyExist
	}

	// hash user password
	hashPass, err := utils.GenerateHashFromPassword(user.Password)

	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// update user password to hashed password
	user.Password = hashPass

	// save user
	user, err = a.authRepo.SaveUser(user)

	if err != nil {
		return domain.User{}, fmt.Errorf("failed to save user on in db: %w", err)
	}

	return user, nil
}

func (a *authUseCase) SignIn(user domain.User) (domain.User, error) {

	// get user by email
	dbUser, err := a.authRepo.FindUserByEmail(user.Email)

	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user from db: %w", err)
	}

	// check user exist or not
	if dbUser.ID == 0 {
		return domain.User{}, ErrNotExist
	}

	// check user password with hashed password
	if valid := utils.VerifyHashAndPassword(dbUser.Password, user.Password); !valid {
		return domain.User{}, ErrWrongPassword
	}

	return dbUser, nil
}

func (a *authUseCase) GenerateAccessToken(role string, user domain.User) (response.Token, error) {

	// generate new token id and expire time with access token duration.
	var (
		tokenID  = utils.GenerateUUID()
		expireAt = time.Now().Add(accessTokenDuration)
	)

	// create token payload with user details and role
	payload := token.Payload{
		TokenID:  tokenID,
		UserID:   user.ID26,
		Email:    user.Email,
		Role:     role,
		ExpireAt: expireAt,
	}

	// generate the token
	token, err := a.tokenAuth.GenerateToken(payload)

	if err != nil {
		return response.Token{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	// return token and expire time
	return response.Token{
		AccessToken:         token,
		AccessTokenExpireAt: expireAt,
	}, nil
}

func (a *authUseCase) VerifyAccessToken(tokenString string) (token.Payload, error) {

	payload, err := a.tokenAuth.VerifyToken(tokenString)

	if err != nil {
		log.Println(err)
		// check the error and return new error according the token error
		switch err {
		case token.ErrExpiredToken:
			return token.Payload{}, ErrExpired
		case token.ErrInvalidToken:
			return token.Payload{}, ErrInvalid
		default:
			return token.Payload{}, fmt.Errorf("failed to verify token: %w", err)
		}
	}







	
	return payload, nil
}
