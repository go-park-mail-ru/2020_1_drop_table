package main

import (
	proto "2020_1_drop_table/internal/microservices/staff/delivery/grpc/protobuff"
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure())

	if err != nil {
		fmt.Println("Unexpected Error", err)
	}

	defer conn.Close()
	// CALL the NewArticleHandlerClient from generated File
	c := proto.NewStaffGRPCHandlerClient(conn)
	// SingleRequest is a generated Struct from proto file
	empt := proto.Empty{}
	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 2}}
	ctx := context.WithValue(context.Background(), "session", &session)
	md, ok := metadata.FromOutgoingContext(ctx)
	fmt.Println(md, ok)
	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs("key1", "val1", "key2", "val2"),
	)
	md, ok = metadata.FromOutgoingContext(ctx)
	fmt.Println(md, ok)
	r, err := c.GetFromSession(ctx, &empt, grpc.EmptyCallOption{})
	if err != nil {
		fmt.Println("Unexpected Error", err)
	}
	fmt.Println("Article : ", r)
}
