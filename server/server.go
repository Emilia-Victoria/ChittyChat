package main

import (
	//"io"
	"log"
	"net"

	chat "github.com/Emilia-Victoria/ChittyChat/chat"
	"google.golang.org/grpc"
)

type Server struct {
	chat.UnimplementedChittyChatServer
}

func main() {
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()
	//chat.RegisterGetCurrentTimeServer(grpcServer, &Server{})

	chat.RegisterChittyChatServer(grpcServer, &Server{}) //Registers the server to the gRPC server.

	//log.Printf("Server %s: Listening at %v\n", *serverName, list.Addr())

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
