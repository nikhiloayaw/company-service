//go:build wireinject
// +build wireinject

package di

import (
	"auth-service/pkg/api"
	"auth-service/pkg/api/service"
	"auth-service/pkg/config"
	"auth-service/pkg/db"
	"auth-service/pkg/repository"
	"auth-service/pkg/service/token"
	"auth-service/pkg/usecase"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*api.Server, error) {

	wire.Build(
		db.ConnectDatabase,

		token.NewJwtTokenAuth,

		repository.NewAuthRepo,
		usecase.NewAuthUseCase,
		service.NewAuthServiceServer,

		api.NewServerGRPC,
	)

	return &api.Server{}, nil
}
