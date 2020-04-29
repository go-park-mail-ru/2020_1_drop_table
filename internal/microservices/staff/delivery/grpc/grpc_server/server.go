package staff

import (
	"2020_1_drop_table/configs"
	staff2 "2020_1_drop_table/internal/microservices/staff"
	proto "2020_1_drop_table/internal/microservices/staff/delivery/grpc/protobuff"
	"2020_1_drop_table/internal/microservices/staff/models"
	"context"
	"fmt"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
)

type server struct {
	staffUseCase staff2.Usecase
}

func NewStaffServerGRPC(gserver *grpc.Server, staffUCase staff2.Usecase) {
	articleServer := &server{
		staffUseCase: staffUCase,
	}
	proto.RegisterStaffGRPCHandlerServer(gserver, articleServer)
	reflection.Register(gserver)
}

func StartStaffGrpcServer(staffUCase staff2.Usecase) {
	list, err := net.Listen("tcp", configs.GRPCStaffUrl)
	if err != nil {
		log.Err(err)
	}
	server := grpc.NewServer()
	NewStaffServerGRPC(server, staffUCase)
	server.Serve(list)
}

func (s *server) GetFromSession(ctx context.Context, in *proto.Empty) (*proto.SafeStaff, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Println(md)
	userid, _ := md["userid"]
	intUserId, _ := strconv.Atoi(userid[0])

	session := sessions.Session{Values: map[interface{}]interface{}{"userID": intUserId}}
	ctx = context.WithValue(context.Background(), "session", &session)
	safeStaff, err := s.staffUseCase.GetFromSession(ctx)
	fmt.Println(safeStaff, err)
	return transformIntoRPC(&safeStaff), err
}

func (s *server) GetById(ctx context.Context, id *proto.Id) (*proto.SafeStaff, error) {
	safeStaff, err := s.staffUseCase.GetByID(ctx, int(id.GetId()))
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