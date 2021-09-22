package main

import (
	"context"
	"github.com/google/uuid"
	pb "golang-minecraft-event-ingress/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"testing"
)

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024*1024)

	server := grpc.NewServer()
	pb.RegisterEventIngressServer(server, &eventIngressServer{})

	go func() {
		if err:=server.Serve(listener); err != nil {
			log.Fatalf("Cannot listen due to %v", err)
		}
	}()

	return func(context.Context, string)(net.Conn, error) {
		return listener.Dial()
	}
}

func TestEventIngressServer_ReportEvent(t *testing.T) {
	tests := []struct {
		name string
		event *pb.Event
		response emptypb.Empty
		err error
	}{
		{
		"Sanity check",
			&pb.Event{
				Meta:    &pb.Event_EventMeta{OriginatorId: uuid.New().String()},
				Payload: []*pb.Event_EventPayload{&pb.Event_EventPayload{SomeValue: uuid.New().String()}},
			},
		emptypb.Empty{},
		nil,
		},
	}
	ctx := context.Background()

	conn,err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewEventIngressClient(conn)

	for _,tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.ReportEvent(ctx, tt.event) // drop response, is pb.empty

			if err != nil {
				t.Errorf("Something went wrong %v", err)
			}
		})
	}


	client = client
	tests = tests
}