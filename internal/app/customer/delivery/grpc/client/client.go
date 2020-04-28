package main

import (
	proto "2020_1_drop_table/internal/app/customer/delivery/grpc/protobuff"
	"2020_1_drop_table/internal/app/customer/delivery/grpc/server"
	"2020_1_drop_table/internal/app/customer/models"
	"context"
	"google.golang.org/grpc"
)

type CustomerGRPC struct {
	client proto.CustomerGRPCHandlerClient
}

func NewCustomerClient(conn *grpc.ClientConn) *CustomerGRPC {
	c := proto.NewCustomerGRPCHandlerClient(conn)
	return &CustomerGRPC{
		client: c,
	}
}

func (s *CustomerGRPC) Add(ctx context.Context, newCustomer models.Customer) (models.Customer, error) {
	customerProto := server.CustomerModelToProto(newCustomer)
	customerProto, err := s.client.Add(ctx, customerProto)
	return server.CustomerProtoToModel(customerProto), err
}
