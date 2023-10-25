package api

import (
	"company-service/pkg/config"
	"company-service/pkg/pb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	lis  net.Listener
	gsr  *grpc.Server
	port string
}

func NewServerGRPC(cfg config.Config, srv pb.CompanyServiceServer) (*Server, error) {

	addr := fmt.Sprintf("%s:%s", cfg.CompanyServiceHost, cfg.CompanyServicePort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	gsr := grpc.NewServer()

	pb.RegisterCompanyServiceServer(gsr, srv)

	return &Server{
		lis:  lis,
		gsr:  gsr,
		port: cfg.CompanyServicePort,
	}, err
}

func (c *Server) Start() error {
	log.Println("Company service listening on port: ", c.port)
	return c.gsr.Serve(c.lis)
}
