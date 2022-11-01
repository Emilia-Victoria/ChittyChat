package main

import (
	//"io"
	//"context"

	"fmt"
	"io"
	"log"
	"net"

	chat "github.com/Emilia-Victoria/ChittyChat/chat"
	"google.golang.org/grpc"
)

type Server struct {
	chat.UnimplementedChittyChatServer
	messageChannel []chan *chat.Message
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
	fmt.Printf("Server started successfully")
}

func (s *Server) PublishMessage(msgStream chat.ChittyChat_PublishMessageServer) error {

	msg, err := msgStream.Recv()

	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	ack := chat.MessageAck{IsSent: true}
	msgStream.SendAndClose(&ack)

	go s.sendMessage(msg)

	return nil
}

func (s *Server) JoinChat(ch *chat.JoinRequest, msgStream chat.ChittyChat_JoinChatServer) error {

	msgChannel := make(chan *chat.Message)
	s.messageChannel = append(s.messageChannel, msgChannel)

	announcementMsg := chat.Message{Sender: "Server", Message: ch.User + " joined the chat! ( ･_･)♡"}
	s.sendMessage(&announcementMsg)

	for {
		select {
		case <-msgStream.Context().Done():
			s.sendMessage(&chat.Message{
				Message: ch.User + " left the chat",
				Sender:  "Server",
			})
			return nil
		case msg := <-msgChannel:
			fmt.Printf("GO ROUTINE (got message): %v \n", msg)
			msgStream.Send(msg)
		}
	}
}

func (s *Server) sendMessage(msg *chat.Message) {
	for _, msgChan := range s.messageChannel {
		select {
		case msgChan <- msg:
		default:
		}
	}
}

func (s *Server) LeaveChat(*chat.LeaveRequest, chat.ChittyChat_LeaveChatServer) error {
	return nil
}
