package main

import (
	feedservice "cid.com/feed-service/src/contracts/feed-content"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	address := ":8089"

	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Cannot create listener on port %s. %s", address, err)
	}

	server := grpc.NewServer()
	service := &myFeedServer{}
	feedservice.RegisterFeedContentServer(server, service)
	err = server.Serve(listener)

	if err != nil {
		log.Fatalf("Something went wrong while trying to start the server: %s", err)
	}
}

type myFeedServer struct {
	feedservice.UnimplementedFeedContentServer
}

func (s myFeedServer) Feed(ctx context.Context, req *feedservice.FeedContentRequest) (*feedservice.FeedContentResponse, error) {
	return &feedservice.FeedContentResponse{
		MediaItems:    make([]*feedservice.MediaItem, 0),
		Notifications: make([]*feedservice.Notification, 0),
	}, nil
}
