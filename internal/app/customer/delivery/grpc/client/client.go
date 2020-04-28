package main

import (
	proto "2020_1_drop_table/internal/app/customer/delivery/grpc/protobuff"
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

func (s *CustomerGRPC) Add(ctx context.Context, newCustomer *proto.Customer) (*proto.Customer, error) {
	return s.client.Add(ctx, newCustomer)
}
