package web

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	topav1 "github.com/pranayhere/touringparty/gen/go/v1"
	"github.com/pranayhere/touringparty/internal/config"
	"github.com/pranayhere/touringparty/internal/port/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	cfg     config.Config
	errChan chan error
}

func NewService(cfg config.Config) (*Service, error) {
	return &Service{
		cfg:     cfg,
		errChan: make(chan error),
	}, nil
}

func (svc *Service) GetErrorChannel() chan error {
	return svc.errChan
}

func (svc *Service) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":"+svc.cfg.Server.Grpc.Port)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	topav1.RegisterStatusCheckServiceServer(s, handler.NewStatusCheckService())
	log.Println("Serving gRPC on " + svc.cfg.Server.Grpc.Host + ":" + svc.cfg.Server.Grpc.Port)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		svc.cfg.Server.Grpc.Host+":"+svc.cfg.Server.Grpc.Port,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = topav1.RegisterStatusCheckServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":" + svc.cfg.Server.Http.Port,
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:" + svc.cfg.Server.Http.Port)
	log.Fatalln(gwServer.ListenAndServe())
	return nil
}
