package server

import (
	"2020_1_drop_table/internal/app/customer"
	proto "2020_1_drop_table/internal/app/customer/delivery/grpc/protobuff"
	"2020_1_drop_table/internal/app/customer/models"
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type server struct {
	customerUseCase customer.Usecase
}

func StartCustomerGrpcServer(customerUseCase customer.Usecase, url string) {
	list, err := net.Listen("tcp", url)
	if err != nil {
		log.Err(err)
	}
	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)
	NewCustomerServerGRPC(server, customerUseCase)
	server.Serve(list)
}

func NewCustomerServerGRPC(gServer *grpc.Server, customerUCase customer.Usecase) {
	customerServer := &server{
		customerUseCase: customerUCase,
	}
	proto.RegisterCustomerGRPCHandlerServer(gServer, customerServer)
	reflection.Register(gServer)
}

func (s *server) Add(ctx context.Context, customer *proto.Customer) (*proto.Customer, error) {
	modelCustomer, err := s.customerUseCase.Add(ctx, CustomerProtoToModel(customer))
	return CustomerModelToProto(modelCustomer), err
}

func CustomerProtoToModel(protoCustomer *proto.Customer) models.Customer {
	return models.Customer{
		CustomerID:   protoCustomer.CustomerID,
		CafeID:       int(protoCustomer.CafeID),
		Type:         protoCustomer.Type,
		Points:       protoCustomer.Points,
		SurveyResult: protoCustomer.SurveyResult,
	}
}

func CustomerModelToProto(modelCustomer models.Customer) *proto.Customer {
	return &proto.Customer{
		CustomerID:   modelCustomer.CustomerID,
		CafeID:       int64(modelCustomer.CafeID),
		Type:         modelCustomer.Type,
		Points:       modelCustomer.Points,
		SurveyResult: modelCustomer.SurveyResult,
	}
}
