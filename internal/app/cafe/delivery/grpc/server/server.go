package server

import (
	"2020_1_drop_table/internal/app/cafe"
	proto "2020_1_drop_table/internal/app/cafe/delivery/grpc/protobuff"
	"2020_1_drop_table/internal/app/cafe/models"
	"context"
	googleProtobuf "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	cafeUseCase cafe.Usecase
}

func NewCafeServerGRPC(gServer *grpc.Server, cafeUCase cafe.Usecase) {
	cafeServer := &server{
		cafeUseCase: cafeUCase,
	}
	proto.RegisterCafeGRPCHandlerServer(gServer, cafeServer)
	reflection.Register(gServer)
}

func (s *server) GetByID(ctx context.Context, id *proto.ID) (*proto.Cafe, error) {
	cafeObj, err := s.cafeUseCase.GetByID(ctx, int(id.Id))
	return cafeModelToProto(cafeObj), err
}

func cafeModelToProto(cafe models.Cafe) *proto.Cafe {
	return &proto.Cafe{
		CafeID:      int64(cafe.CafeID),
		CafeName:    cafe.CafeName,
		Address:     cafe.Address,
		Description: cafe.Description,
		StaffID:     int64(cafe.StaffID),
		OpenTime:    &googleProtobuf.Timestamp{Seconds: cafe.OpenTime.Unix()},
		CloseTime:   &googleProtobuf.Timestamp{Seconds: cafe.CloseTime.Unix()},
		Photo:       cafe.Photo,
	}
}
