package client

import (
	proto "2020_1_drop_table/internal/app/cafe/delivery/grpc/protobuff"
	"2020_1_drop_table/internal/app/cafe/models"
	"context"
	"google.golang.org/grpc"
	"time"
)

type CafeGRPC struct {
	client proto.CafeGRPCHandlerClient
}

func NewCafeClient(conn *grpc.ClientConn) *CafeGRPC {
	c := proto.NewCafeGRPCHandlerClient(conn)
	return &CafeGRPC{
		client: c,
	}
}

func (c *CafeGRPC) GetByID(ctx context.Context, id int) (models.Cafe, error) {
	idProto := &proto.ID{}
	idProto.Id = int64(id)

	cafeProto, err := c.client.GetByID(ctx, idProto)
	return cafeProtoToModel(cafeProto), err
}

func cafeProtoToModel(cafe *proto.Cafe) models.Cafe {
	return models.Cafe{
		CafeID:      int(cafe.CafeID),
		CafeName:    cafe.CafeName,
		Address:     cafe.Address,
		Description: cafe.Description,
		StaffID:     int(cafe.StaffID),
		OpenTime:    time.Unix(cafe.OpenTime.GetSeconds(), 0).UTC(),
		CloseTime:   time.Unix(cafe.CloseTime.GetSeconds(), 0).UTC(),
		Photo:       cafe.Photo,
	}
}
