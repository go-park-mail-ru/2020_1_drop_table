package server

import (
	"2020_1_drop_table/internal/app/customer"
	proto "2020_1_drop_table/internal/app/customer/delivery/grpc/protobuff"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	customerUseCase customer.Usecase
}

func NewCustomerServerGRPC(gServer *grpc.Server, customerUCase customer.Usecase) {
	customerServer := &server{
		customerUseCase: customerUCase,
	}
	proto.RegisterCustomerGRPCHandlerServer(gServer, customerServer)
	reflection.Register(gServer)
}

func (s *server) Add(ctx context.Context, customer *proto.Customer) (*proto.Customer, error) {

}
