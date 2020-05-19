package staff

import (
	"2020_1_drop_table/configs"
	staff2 "2020_1_drop_table/internal/microservices/staff"
	proto "2020_1_drop_table/internal/microservices/staff/delivery/grpc/protobuff"
	"2020_1_drop_table/internal/microservices/staff/models"
	"context"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
	"time"
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

func StartStaffGrpcServer(staffUCase staff2.Usecase, url string) {
	list, err := net.Listen("tcp", url)
	if err != nil {
		log.Err(err)
	}
	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)
	NewStaffServerGRPC(server, staffUCase)
	_ = server.Serve(list)
}

func makeContextFromMetaDataInContext(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)

	userid := md["userid"]
	intUserId, _ := strconv.Atoi(userid[0])

	session := sessions.Session{Values: map[interface{}]interface{}{"userID": intUserId}}
	return context.WithValue(context.Background(), configs.SessionStaffID, &session)
}

func (s *server) GetFromSession(ctx context.Context, _ *proto.Empty) (*proto.SafeStaff, error) {
	ctx = makeContextFromMetaDataInContext(ctx)

	safeStaff, err := s.staffUseCase.GetFromSession(ctx)

	return transformIntoRPC(&safeStaff), err
}

func (s *server) GetById(ctx context.Context, id *proto.Id) (*proto.SafeStaff, error) {
	ctx = makeContextFromMetaDataInContext(ctx)
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
