package service

import (
	"company-service/pkg/models"
	"company-service/pkg/pb"
	"company-service/pkg/usecase/interfaces"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	maxBufferSize = 1024
)

type companyServiceServer struct {
	companyUseCase interfaces.CompanyUseCase

	pb.UnimplementedCompanyServiceServer
}

func NewCompanyServiceServer(companyUseCase interfaces.CompanyUseCase) pb.CompanyServiceServer {

	return &companyServiceServer{
		companyUseCase: companyUseCase,
	}
}

func (c *companyServiceServer) Create(req *pb.CompanyCreateRequest, stream pb.CompanyService_CreateServer) error {

	companyReq := models.CompanyRequest{
		Name:           req.GetName(),
		CEO:            req.GetCeo(),
		TotalEmployees: int(req.GetTotalEmployees()),
	}

	company := c.companyUseCase.Create(companyReq)

	data, err := json.Marshal(company)

	if err != nil {
		return status.Errorf(codes.Internal, "failed to marshal school to json: %v", err)
	}

	var (
		start  = 0
		end    = maxBufferSize
		buffer []byte
	)

	for {

		// check the end is out of bound of data
		if end >= len(data) {
			buffer = data[start:]
			if err := stream.Send(&pb.CompanyCreateResponse{Data: buffer}); err != nil {
				return err
			}

			break
		} else {
			// slice out the buffer from data with max buffer size
			buffer = data[start:end]

			// update the start and end
			start += maxBufferSize
			end += maxBufferSize
		}

		// send the buffer as stream
		if err := stream.Send(&pb.CompanyCreateResponse{Data: buffer}); err != nil {
			return err
		}
	}
	return nil
}
