package main

import (
	"context"
	"fmt"
	"github.com/DapperBlondie/go-grpc/src/messages/files"
	"google.golang.org/grpc"
	"io"
	"log"
	"strconv"
	"time"
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
	clientStreamingGreeting(client)
	clientStreamingAverage(sumClient)
	biDirectionalGreetStreaming(client)

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

func clientStreamingAverage(client files.SumServiceClient) {
	stream, err := client.AverageStreamingResult(context.Background())
	if err != nil {
		log.Println("Error in getting the Average Stream : " + err.Error())
		return
	}

	req := &files.NumberRequest{Num: 10}

	for j := 0; j < 10; j += 1 {
		err = stream.Send(req)
		if err != nil {
			log.Println("Error in client sending average : " + err.Error())
			return
		}
		time.Sleep(time.Millisecond * 200)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("Error in receiving response in average : " + err.Error())
		return
	}

	fmt.Printf("Average of streaming : %f\n", resp.GetAverage())
	return
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

func clientStreamingGreeting(client files.GreetServiceClient) {
	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Println("Error in greeting client streaming : " + err.Error())
		return
	}

	req := &files.LongGreetRequest{Greeting: &files.Greeting{
		FirstName: "Alireza",
		LastName:  "Gharib",
	}}

	for j := 0; j < 10; j += 1 {
		err := stream.Send(req)
		if err != nil {
			log.Println("Error in send msg : " + err.Error())
			break
		}
		time.Sleep(time.Millisecond * 200)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("Error in receiving the response : " + err.Error())
		return
	}
	fmt.Println(resp.GetResult())
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

func biDirectionalGreetStreaming(client files.GreetServiceClient) {
	stream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Println("Error in getting the stream from client : " + err.Error())
		return
	}

	var requests []*files.GreetEveryoneRequest
	requests = append(requests, &files.GreetEveryoneRequest{Greeting: &files.Greeting{
		FirstName: "Dexter",
		LastName:  "Deshawn",
	}})
	requests = append(requests, &files.GreetEveryoneRequest{Greeting: &files.Greeting{
		FirstName: "Tony",
		LastName:  "McKay",
	}})
	requests = append(requests, &files.GreetEveryoneRequest{Greeting: &files.Greeting{
		FirstName: "Johnny",
		LastName:  "SilverHand",
	}})

	waitC := make(chan int)

	go func() {
		for _, req := range requests {
			err := stream.Send(req)
			if err != nil {
				log.Println("Error in sending request : " + err.Error())
				continue
			}
		}
		err := stream.CloseSend()
		if err != nil {
			log.Print("Error in closing the stream : " + err.Error())
			return
		}
		return
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				log.Println("Error in receiving data : " + err.Error())
				continue
			} else if err == io.EOF {
				log.Println("Error in receiving data : " + err.Error())
				close(waitC)
				return
			} else {
				fmt.Println(resp.GetResult())
			}
		}
	}()

	<-waitC
}
