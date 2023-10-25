package api

import (
	"auth-service/pkg/config"
	"auth-service/pkg/pb"
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

func NewServerGRPC(cfg config.Config, srv pb.AuthServiceServer) (*Server, error) {

	addr := fmt.Sprintf("%s:%s", cfg.AuthServiceHost, cfg.AuthServicePort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	gsr := grpc.NewServer()

	pb.RegisterAuthServiceServer(gsr, srv)

	return &Server{
		lis:  lis,
		gsr:  gsr,
		port: cfg.AuthServicePort,
	}, err
}

func (c *Server) Start() error {
	log.Println("Auth service listening on port: ", c.port)
	return c.gsr.Serve(c.lis)
}
