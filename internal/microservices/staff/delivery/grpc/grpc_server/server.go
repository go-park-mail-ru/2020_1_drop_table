package staff

import (
	staff2 "2020_1_drop_table/internal/microservices/staff"
	proto "2020_1_drop_table/internal/microservices/staff/delivery/grpc/protobuff"
	"2020_1_drop_table/internal/microservices/staff/models"
	"context"
	"fmt"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
)

type server struct {
	staffUseCase staff2.Usecase
}

func NewArticleServerGrpc(gserver *grpc.Server, staffUCase staff2.Usecase) {
	articleServer := &server{
		staffUseCase: staffUCase,
	}
	proto.RegisterStaffGRPCHandlerServer(gserver, articleServer)
	reflection.Register(gserver)
}

func StartGrpcServer(staffUCase staff2.Usecase) {
	list, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Err(err)
	}
	server := grpc.NewServer()
	NewArticleServerGrpc(server, staffUCase)
	server.Serve(list)
}

func (s *server) GetFromSession(ctx context.Context, in *proto.Empty) (*proto.SafeStaff, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println(md, ok)
	safeStaff, err := s.staffUseCase.GetFromSession(ctx)
	return transformIntoRPC(&safeStaff), err

}

func transformIntoRPC(staff *models.SafeStaff) *proto.SafeStaff {
	if staff == nil {
		return nil
	}

	edited_at := &google_protobuf.Timestamp{
		Seconds: staff.EditedAt.Unix(),
	}

	res := &proto.SafeStaff{
		StaffID:  int64(staff.StaffID),
		Name:     staff.Name,
		Email:    staff.Email,
		EditedAt: edited_at,
		Photo:    staff.Photo,
		IsOwner:  staff.IsOwner,
		CafeId:   int64(staff.CafeId),
		Position: staff.Position,
	}
	return res
}
