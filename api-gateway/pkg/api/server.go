package api

import (
	"api-gateway/pkg/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	server *http.Server
	port   string

	graceFulTimeout time.Duration
}

// To create a new server with HTTP
// @title Go Microservices BoilerPlate
// @description BoilerPlate for micro services with GRPC
//
// @BasePath					/api/v1
// @SecurityDefinitions.apikey	BearerAuth
// @Name						Authorization
// @In							header
// @Description				Add prefix of Bearer before  token Ex: "Bearer token"
// @Query.collection.format	multi
func NewServerHTTP(cfg config.Config, router http.Handler) *Server {

	// once found how to read value of time from env, this will be remove
	{
		cfg.ReadTimeout = time.Second * 15
		cfg.WriteTimeout = time.Second * 15
		cfg.GraceFulTimeout = time.Second * 10
	}

	addr := fmt.Sprintf("localhost:%s", cfg.ApiPort)
	// create an http server with handler as gin
	server := &http.Server{
		Handler: router,
		Addr:    addr,

		WriteTimeout: cfg.ReadTimeout,
		ReadTimeout:  cfg.WriteTimeout,
	}

	return &Server{
		server:          server,
		port:            cfg.ApiPort,
		graceFulTimeout: cfg.GraceFulTimeout,
	}
}

// To start the server.
func (s *Server) Start() {

	// run the listen and serve in separate goroutines to do graceful shutdown.
	go func() {
		log.Printf("server listening at port: %s", s.port)
		s.server.ListenAndServe()
	}()

	// create an os.Signal channel to receive os signal.
	ch := make(chan os.Signal, 1)
	// register our channel with os.Interrupt
	signal.Notify(ch, os.Interrupt)
	// block until the receive the signal
	<-ch
	// create a timeout context with graceful shout down time duration
	ctx, cancel := context.WithTimeout(context.Background(), s.graceFulTimeout)
	defer cancel()

	log.Println("wait for serving request to complete")
	// call server shutdown to stop receiving new request, and wait for the current request to complete before timeout
	if err := s.server.Shutdown(ctx); err != nil {
		log.Println("failed to shutdown server: ", err)
	}

	log.Println("server shuting down...")
	os.Exit(0)
}
