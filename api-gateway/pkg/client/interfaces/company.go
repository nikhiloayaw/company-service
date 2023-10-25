package interfaces

import (
	"api-gateway/pkg/api/handler/request"
)

type CompanyServiceClient interface {
	Create(companyReq request.CompanyRequest) ([]byte, error)
}
