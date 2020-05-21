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

func (c *CafeGRPC) GetByOwnerId(ctx context.Context, id int) ([]models.Cafe, error) {
	idProto := &proto.ID{}
	idProto.Id = int64(id)

	cafeList, err := c.client.GetByOwnerID(ctx, idProto)
	return cafeListProtoToModel(cafeList), err

}

func cafeListProtoToModel(listCafe *proto.ListCafe) []models.Cafe {
	if listCafe == nil {
		return []models.Cafe{}
	}
	var resList []models.Cafe
	for _, cafe := range listCafe.Cafe {
		resList = append(resList, cafeProtoToModel(cafe))
	}
	return resList
}

func cafeProtoToModel(cafe *proto.Cafe) models.Cafe {
	if cafe == nil {
		return models.Cafe{}
	}
	return models.Cafe{
		CafeID:      int(cafe.CafeID),
		CafeName:    cafe.CafeName,
		Address:     cafe.Address,
		Description: cafe.Description,
		StaffID:     int(cafe.StaffID),
		OpenTime:    time.Unix(cafe.GetOpenTime().GetSeconds(), 0).UTC(),
		CloseTime:   time.Unix(cafe.GetCloseTime().GetSeconds(), 0).UTC(),
		Photo:       cafe.Photo,
	}
}
