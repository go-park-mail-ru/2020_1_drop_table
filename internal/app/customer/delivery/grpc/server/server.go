package server

import (
	"2020_1_drop_table/internal/app/customer"
	proto "2020_1_drop_table/internal/app/customer/delivery/grpc/protobuff"
	"2020_1_drop_table/internal/app/customer/models"
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
	modelCustomer, err := s.customerUseCase.Add(ctx, customerProtoToModel(customer))
	return customerModelToProto(modelCustomer), err
}

func customerProtoToModel(protoCustomer *proto.Customer) models.Customer {
	return models.Customer{
		CustomerID:   protoCustomer.CustomerID,
		CafeID:       int(protoCustomer.CafeID),
		Type:         protoCustomer.Type,
		Points:       protoCustomer.Points,
		SurveyResult: protoCustomer.SurveyResult,
	}
}

func customerModelToProto(modelCustomer models.Customer) *proto.Customer {
	return &proto.Customer{
		CustomerID:   modelCustomer.CustomerID,
		CafeID:       int64(modelCustomer.CafeID),
		Type:         modelCustomer.Type,
		Points:       modelCustomer.Points,
		SurveyResult: modelCustomer.SurveyResult,
	}
}
