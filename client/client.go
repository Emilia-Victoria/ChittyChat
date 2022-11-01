package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	chat "github.com/Emilia-Victoria/ChittyChat/chat"

	"google.golang.org/grpc"
)

var channel = flag.String("Channel", "default", "Chitty-Chat")
var username = flag.String("user", "default", "username")
var lamportTime = flag.Int64("time", 0, "lamportTimeStamp")

func main() {
	fmt.Println(".❀.❀.❀. Welcome to Chitty Chat .❀.❀.❀.")
	fmt.Println("Please enter your name: ")
	fmt.Scanf("%s", username)

	fmt.Println("Connecting... (° ͜ʖ °)")

	file, err := os.OpenFile("logger.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer file.Close()

	log.SetOutput(file)

	conn, err := grpc.Dial(":9080", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %s", err)
		return
	}

	defer conn.Close()

	ctx := context.Background()
	client := chat.NewChittyChatClient(conn)

	go joinChat(ctx, client)
	// exitChannel(ctx, client)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		go publishMessage(ctx, client, scanner.Text())

	}

}

func joinChat(ctx context.Context, client chat.ChittyChatClient) {

	joinreq := chat.JoinRequest{User: *username, Channel: *channel}
	stream, err := client.JoinChat(ctx, &joinreq)

	if err != nil {
		log.Fatalf("client.JoinChannel(ctx, &channel) throws: %v", err)
	}
	welcome := *username + " has joined the chat: " + *channel + "! ( ･_･)♡"
	publishMessage(ctx, client, welcome)

	waitc := make(chan struct{}) //go never stops with this

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive message from channel joining. \nErr: %v", err)
			}
			if *username != in.Sender {
				if in.LampTime > *lamportTime {
					*lamportTime = in.LampTime + 1
				} else {
					*lamportTime++
				}
				fmt.Printf("(%v) %v: %v \n", *lamportTime, in.Sender, in.Message)
				log.Printf("(%v) %v: %v \n", *lamportTime, in.Sender, in.Message)

			}
		}
	}()
	<-waitc
}

func leaveChat(ctx context.Context, client chat.ChittyChatClient) {

	*lamportTime++
	stream, err := client.PublishMessage(ctx)
	if err != nil {
		log.Printf("Unable to send leave message: error: %v", err)
	}
	msg := chat.Message{
		Sender:   *username,
		LampTime: *lamportTime,
	}
	log.Printf("This user has left the chat: %v, (%v)", *username, *lamportTime)

	stream.Send(&msg)

	ack, _ := stream.CloseAndRecv()
	fmt.Printf("Message has been sent: %v \n", ack)
}

//( ͡ಥ ͜ʖ ͡ಥ) ⊙︿⊙ (つ◉益◉)つ ୧( ಠ Д ಠ )୨ ᕕ(˵•̀෴•́˵)ᕗ (＞﹏＜)
//（╯‵□′）╯︵┴─┴   ( ･_･)♡

func publishMessage(ctx context.Context, client chat.ChittyChatClient, message string) {

	*lamportTime++
	stream, err := client.PublishMessage(ctx)
	if err != nil {
		log.Printf("Cannot send message: error: %v", err)
	}

	msg := chat.Message{Sender: *username, Message: message, LampTime: *lamportTime}

	stream.Send(&msg)
	ack, _ := stream.CloseAndRecv()
	fmt.Printf(" %v \n", ack)

}
