package main

import (
	proto "2020_1_drop_table/internal/microservices/staff/delivery/grpc/protobuff"
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strconv"
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

func (s StaffClient) GetFromSession(ctx context.Context) {
	ctx = s.AddSessionInMetadata(ctx)
	empt := proto.Empty{}
	r, err := s.client.GetFromSession(ctx, &empt, grpc.EmptyCallOption{})
	if err != nil {
		fmt.Println("Unexpected Error", err)
	}
	fmt.Println("Article : ", r)
}

func (s StaffClient) AddSessionInMetadata(ctx context.Context) context.Context {
	value := ctx.Value("session").(*sessions.Session)
	el := value.Values["userID"].(int)
	fmt.Println(el)
	return metadata.AppendToOutgoingContext(ctx, "userID", strconv.Itoa(el))
}

func main() {
	conn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		fmt.Println("Unexpected Error", err)
	}
	client := NewStaffClient(conn)
	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 41}}
	ctx := context.WithValue(context.Background(), "session", &session)
	client.GetFromSession(ctx)

}
