package main

import (
	//"context"
	//"fmt"
	"log"

	//c "chat"

	//"github.com/Emilia-Victoria/ChittyChat/chat"

	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	defer conn.Close()

	//c := chat.NewChittyChatClient(conn)
	//JoinChat(c)

	for {
		//update chat
		//c.Sleep(5 * c.Second)
	}
}

func joinChat() {}

func leaveCat() {}

func publishMessage() {}
