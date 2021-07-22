package main

import (
	"context"
	"fmt"
	"github.com/DapperBlondie/go-grpc/src/messages/files"
	"google.golang.org/grpc"
	"io"
	"log"
	"strconv"
)

// This file is our unary client that we are using it for our unary Greeting API

func main() {
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
	sumClient := files.NewSumServiceClient(conn)

	unaryGreeting(client)
	unarySum(sumClient)
	serverStreamingGreeting(client)
	serverStreamingSum(sumClient)

	return nil
}

func unarySum(client files.SumServiceClient) {
	req := &files.SumRequest{List: []int32{12, 34, 56, 78, 90, 89, 668}}
	resp, err := client.GetSumResult(context.Background(), req)
	if err != nil {
		log.Println("Error in getting the response from Sum Service : " + err.Error())
		return
	}
	fmt.Println(strconv.Itoa(int(resp.Result)) + ", result of our sum request.")
	return
}

func serverStreamingSum(client files.SumServiceClient) {
	req := &files.SumRequest{List: []int32{12, 56, 89, 90, 772, 8761}}
	stream, err := client.GetStreamingSumResult(context.Background(), req)
	if err != nil {
		log.Println("Error in streaming : " + err.Error())
		return
	}

	fmt.Println("Result of sum streaming : ")
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("Error occurred in streaming : " + err.Error())
			return
		} else if err != nil {
			log.Println("Error occurred in streaming : " + err.Error())
			return
		} else {
			fmt.Println(resp.GetResult())
		}
	}
}

func serverStreamingGreeting(client files.GreetServiceClient) {
	req := &files.GreetingManyTimeRequest{Greeting: &files.Greeting{
		FirstName: "Johnny",
		LastName:  "SilverHand",
	}}

	stream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Println("Error in streaming in server-side : " + err.Error())
		return
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF occurred, reached the end of the stream")
			break
		} else if err != nil {
			log.Fatalln("Error occurred in streaming : " + err.Error())
			return
		} else {
			fmt.Println(msg.Result)
		}
	}
	return
}

func unaryGreeting(client files.GreetServiceClient) {
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
