package interfaces

import (
	"company-service/pkg/domain"
	"company-service/pkg/models"
)

type CompanyUseCase interface {
	Create(companyReq models.CompanyRequest) domain.Company
}
