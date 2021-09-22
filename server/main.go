package main

import (
	"errors"
	pb "golang-minecraft-event-ingress/server/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type eventIngressServer struct {
	pb.UnimplementedEventIngressServer
}

func (srv *eventIngressServer) ReportEvent(ctx context.Context, event *pb.Event) (*emptypb.Empty, error) {
	log.Printf("Received event %s", event)

	if event.Meta == nil || event.Meta.OriginatorId == "" {
		errMsg := "Could not extract originatorId from meta"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	if event.Payload == nil {
		errMsg := "Could not extract payload"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}

	return &emptypb.Empty{}, nil
}

func main() {
	defaultPort := ":8000"
	listener, err := net.Listen("tcp", defaultPort)

	if err != nil {
		log.Fatalf("Couldn't listen on %v because of %v", defaultPort, err)
	}

	grpcServer := grpc.NewServer()

	log.Println("Here we go")
	pb.RegisterEventIngressServer(grpcServer, &eventIngressServer{})
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Can't listen on %s", defaultPort)
	}
}
