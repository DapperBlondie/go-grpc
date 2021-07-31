package main

import (
	"context"
	"fmt"
	"github.com/DapperBlondie/go-grpc/src/messages/files"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//sumClient := files.NewSumServiceClient(conn)

	//unaryGreeting(client)
	//unarySum(sumClient)
	//serverStreamingGreeting(client)
	//serverStreamingSum(sumClient)
	//clientStreamingGreeting(client)
	//clientStreamingAverage(sumClient)
	//biDirectionalGreetStreaming(client)
	//evenOrOddStreaming(sumClient)
	//errorInvalidArgumentInUnaryAPI(sumClient, 12)
	//errorInvalidArgumentInUnaryAPI(sumClient, -10)
	withDeadlineGreetUnary(client)

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
			time.Sleep(time.Millisecond * 200)
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
			if err == io.EOF {
				log.Println(err.Error())
				waitC <- 0
				return
			} else if err != nil {
				log.Println("Error in receiving data : " + err.Error())
				waitC <- 0
				return
			} else {
				fmt.Println(resp.GetResult())
			}
		}
	}()

	for result := range waitC {
		if result == 0 {
			break
		}
	}
	return
}

func evenOrOddStreaming(client files.SumServiceClient) {
	stream, err := client.EvenOrOdd(context.Background())
	if err != nil {
		log.Fatalln(err.Error() + "Occurred in getting stream for client")
		return
	}
	var requests = []int32{12, 345, 567, 69, 42, 556, 245, 6789, 123}
	reqS := &files.NumReq{ReqNum: -1}
	waitC := make(chan int)

	go func() {
		for _, req := range requests {
			reqS.ReqNum = req
			err := stream.Send(reqS)
			if err != nil {
				log.Println(err.Error() + "occurred in sending requests")
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Println(err.Error() + " ;Error in closing the stream")
			return
		}
		return
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println(err)
				waitC <- 0
				return
			} else if err != nil {
				log.Println(err.Error())
				waitC <- 0
				return
			}
			fmt.Println(resp.RespNum)
		}
	}()

	for signal := range waitC {
		if signal == 0 {
			break
		}
	}

	return
}

func errorInvalidArgumentInUnaryAPI(client files.SumServiceClient, reqNumber int32) {
	req := &files.SquareRootRequest{Number: reqNumber}
	resp, err := client.SquareRoot(context.Background(), req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			log.Println(respErr.Code().String() + " : " + respErr.Message())
			return
		} else {
			log.Fatalln("Unexpected error occurred ! " + err.Error())
			return
		}
	}
	fmt.Println(fmt.Sprintf("The result of sqrt : %f\n", resp.GetRootNumber()))
	return
}

func withDeadlineGreetUnary(client files.GreetServiceClient) {
	req := &files.GreetWithDeadlineRequest{Greet: &files.Greeting{
		FirstName: "Johnny",
		LastName:  "SilverHand",
	}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := client.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Println(statusErr.Err().Error() + " : " + statusErr.Code().String())
				return
			} else {
				log.Println(statusErr.Err().Error())
				return
			}
		} else {
			log.Fatalln(err.Error())
			return
		}
		return
	}

	log.Println("The result of greeting to the server : " + resp.GetResult())
	return
}
