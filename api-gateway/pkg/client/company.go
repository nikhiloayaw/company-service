package client

import (
	"api-gateway/pkg/api/handler/request"
	"api-gateway/pkg/client/interfaces"
	"api-gateway/pkg/config"
	"api-gateway/pkg/pb"
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type companyServiceClient struct {
	client pb.CompanyServiceClient
}

// This client is abstraction over the actual client
func NewCompanyServiceClient(cfg config.Config) (interfaces.CompanyServiceClient, error) {

	// create the company service address
	addr := fmt.Sprintf("%s:%s", cfg.CompanyServiceHost, cfg.CompanyServicePort)

	// create a grpc client connection to company service url
	cc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create a grpc client connection for company service : %w", err)
	}

	// create new company service client with the grpc connection
	client := pb.NewCompanyServiceClient(cc)

	// return our abstracted client with grpc company service client
	return &companyServiceClient{
		client: client,
	}, nil
}

func (c *companyServiceClient) Create(companyReq request.CompanyRequest) ([]byte, error) {

	req := &pb.CompanyCreateRequest{
		Name:           companyReq.Name,
		Ceo:            companyReq.CEO,
		TotalEmployees: int32(companyReq.TotalEmployees),
	}

	stream, err := c.client.Create(context.Background(), req)

	if err != nil {
		return nil, fmt.Errorf("failed to call create school: %w", err)
	}

	var companyData []byte

	for {
		res, err := stream.Recv()
		if err != nil {

			if err == io.EOF {
				return companyData, nil
			}
			return nil, fmt.Errorf("failed to receive school on stream: %w", err)
		}
		companyData = append(companyData, res.Data...)
	}

}
