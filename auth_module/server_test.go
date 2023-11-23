package main

import (
	"chatbox/auth_module/pb"
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGrpcServer(t *testing.T) {
	starts := make(chan struct{})
	go testServer(starts)

	select {
	case <-starts:
	case <-time.After(5 * time.Second):
		t.Fatal("Startup timeout")
	}
	conn, err := grpc.Dial(":8008", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("Dial error, %s", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	t.Run("Test login", func(t *testing.T) {
		ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second)
		defer ctxCancel()

		req := pb.AuthenticateRequest{Username: "user1", Password: "pass1"}

		resp, err := client.AuthenticateUser(ctx, &req)
		if err != nil {
			t.Errorf("Response error, %s", err)
		}

		if !resp.Success {
			t.Error(`Expected "Success:true", but received "Sucess:false"`)
		}
	})
}
