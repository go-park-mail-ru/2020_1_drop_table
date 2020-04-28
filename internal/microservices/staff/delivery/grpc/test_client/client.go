package staff

import (
	proto "2020_1_drop_table/internal/microservices/staff/delivery/grpc/protobuff"
	"2020_1_drop_table/internal/microservices/staff/models"
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strconv"
	"time"
)

type StaffClient struct {
	client proto.StaffGRPCHandlerClient
}

func NewStaffClient(conn *grpc.ClientConn) *StaffClient {
	c := proto.NewStaffGRPCHandlerClient(conn)
	return &StaffClient{
		client: c,
	}
}

func transformStaffFromRPC(staff *proto.SafeStaff) models.SafeStaff {
	edited_at := time.Unix(staff.GetEditedAt().GetSeconds(), 0).UTC()
	res := models.SafeStaff{
		StaffID:  int(staff.StaffID),
		Name:     staff.Name,
		Email:    staff.Email,
		EditedAt: edited_at,
		Photo:    staff.Photo,
		IsOwner:  staff.IsOwner,
		CafeId:   int(staff.CafeId),
		Position: staff.Position,
	}
	return res
}

func (s StaffClient) GetFromSession(ctx context.Context) (models.SafeStaff, error) {
	ctx = s.AddSessionInMetadata(ctx)
	empt := proto.Empty{}
	r, err := s.client.GetFromSession(ctx, &empt, grpc.EmptyCallOption{})
	if err != nil {
		fmt.Println("Unexpected Error", err)
	}
	return transformStaffFromRPC(r), err
}

func (s StaffClient) AddSessionInMetadata(ctx context.Context) context.Context {
	value := ctx.Value("session").(*sessions.Session)
	el := value.Values["userID"].(int)
	fmt.Println(el)
	return metadata.AppendToOutgoingContext(ctx, "userID", strconv.Itoa(el))
}

//func main() {
//	conn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
//	defer conn.Close()
//	if err != nil {
//		fmt.Println("Unexpected Error", err)
//	}
//	client := NewStaffClient(conn)
//	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 41}}
//	ctx := context.WithValue(context.Background(), "session", &session)
//	client.GetFromSession(ctx)
//
//}
