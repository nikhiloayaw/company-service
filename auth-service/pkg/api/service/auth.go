package service

import (
	"auth-service/pkg/domain"
	"auth-service/pkg/pb"
	"auth-service/pkg/usecase"
	"auth-service/pkg/usecase/interfaces"
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	RoleUser = "user"
)

type authServiceServer struct {
	authUseCase interfaces.AuthUseCase

	pb.UnimplementedAuthServiceServer
}

func NewAuthServiceServer(authUseCase interfaces.AuthUseCase) pb.AuthServiceServer {

	return &authServiceServer{
		authUseCase: authUseCase,
	}
}

func (a *authServiceServer) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {

	// create user with request use details.
	user := domain.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	user, err := a.authUseCase.SignUp(user)

	if err != nil {
		// log the error
		log.Println(err)

		var (
			statusCode codes.Code
			message    string
		)
		// check the error, and according to the error set status code and message
		switch {
		case errors.Is(err, usecase.ErrAlreadyExist):
			statusCode = codes.AlreadyExists
			message = "user already exist with given details"
		default:
			statusCode = codes.Internal
			message = "internal server error"
		}
		return nil, status.Error(statusCode, message)
	}

	return &pb.SignUpResponse{
		UserId: user.ID26,
	}, nil
}

func (a *authServiceServer) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {

	user := domain.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	user, err := a.authUseCase.SignIn(user)

	if err != nil {
		// log the error
		log.Println(err)

		var (
			statusCode codes.Code
			message    string
		)
		// check the error, and according to the error set status code and message
		switch {
		case errors.Is(err, usecase.ErrNotExist):
			statusCode = codes.NotFound
			message = "user not exist with given details"
		case errors.Is(err, usecase.ErrWrongPassword):
			statusCode = codes.Unauthenticated
			message = "wrong user password"
		default:
			statusCode = codes.Internal
			message = "internal server error"
		}
		return nil, status.Error(statusCode, message)
	}

	// generate access token for the user
	tokenRes, err := a.authUseCase.GenerateAccessToken(RoleUser, user)

	if err != nil {
		// log the error
		log.Println(err)

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.SignInResponse{
		AccessToken: tokenRes.AccessToken,
		ExpireAt:    timestamppb.New(tokenRes.AccessTokenExpireAt),
	}, nil

}

func (a *authServiceServer) VerifyAccessToken(ctx context.Context, req *pb.VerifyAccessTokenRequest) (*pb.VerifyAccessTokenResponse, error) {

	payload, err := a.authUseCase.VerifyAccessToken(req.GetAccessToken())

	if err != nil {

		var (
			statusCode codes.Code
			message    string
		)
		// check error and according to the error set status code and message
		switch {
		case errors.Is(err, usecase.ErrExpired):
			statusCode = codes.Unauthenticated
			message = "token expired"
		case errors.Is(err, usecase.ErrInvalid):
			statusCode = codes.Unauthenticated
			message = "invalid token"
		default:
			statusCode = codes.Internal
			message = "internal server error"
		}
		return nil, status.Error(statusCode, message)
	}

	return &pb.VerifyAccessTokenResponse{
		TokenId:  payload.TokenID,
		UserId:   payload.UserID,
		Email:    payload.Email,
		Role:     payload.Role,
		ExpireAt: timestamppb.New(payload.ExpireAt),
	}, nil
}
