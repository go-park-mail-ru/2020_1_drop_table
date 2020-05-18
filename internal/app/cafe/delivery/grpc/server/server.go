package server

import (
	"2020_1_drop_table/internal/app/cafe"
	proto "2020_1_drop_table/internal/app/cafe/delivery/grpc/protobuff"
	"2020_1_drop_table/internal/app/cafe/models"
	"context"
	googleProtobuf "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type server struct {
	cafeUseCase cafe.Usecase
}

func StartCafeGrpcServer(cafeUseCase cafe.Usecase, url string) {
	list, err := net.Listen("tcp", url)
	if err != nil {
		log.Err(err)
	}
	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)
	NewCafeServerGRPC(server, cafeUseCase)
	server.Serve(list)
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

func (s *server) GetByOwnerID(ctx context.Context, id *proto.ID) (*proto.ListCafe, error) {
	cafeList, err := s.cafeUseCase.GetByOwnerIDWithOwnerID(ctx, int(id.Id))
	return &proto.ListCafe{Cafe: cafeListToProto(cafeList)}, err
}

func cafeListToProto(cafeList []models.Cafe) []*proto.Cafe {
	var resList []*proto.Cafe
	for _, caf := range cafeList {
		resList = append(resList, cafeModelToProto(caf))
	}
	return resList
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
