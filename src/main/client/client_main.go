package main

import (
	"context"
	"fmt"
	"github.com/DapperBlondie/go-grpc/src/messages/files"
	"google.golang.org/grpc"
	"log"
)

// This file is our unary client that we are using it for our unary Greeting API

func main()  {
	err := runClient()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return
}

func runClient() error {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Error in dialing the rpc server : " + err.Error())
		return err
	}

	log.Println("rpc Client dialed localhost:50051 ...")
	client := files.NewGreetServiceClient(conn)
	//fmt.Printf("The client struct : %f\n", client)

	unaryGreeting(client)

	return nil
}

func unaryGreeting(client files.GreetServiceClient)  {
	req := &files.GreetingRequest{Greeting: &files.Greeting{
		FirstName: "Johnny",
		LastName:  "SilverHand",
	}}
	resp, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Println("Error in getting resp from Greeting rpc : " + err.Error())
		return
	}

	fmt.Println("Response msg : " + resp.GetResult())
	return
}