//go:build wireinject
// +build wireinject

package di

import (
	"company-service/pkg/api"
	"company-service/pkg/api/service"
	"company-service/pkg/config"
	"company-service/pkg/service/random"
	"company-service/pkg/usecase"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*api.Server, error) {

	wire.Build(
		random.NewRandomGenerator,
		usecase.NewCompanyUseCase,
		service.NewCompanyServiceServer,
		api.NewServerGRPC,
	)

	return &api.Server{}, nil
}
