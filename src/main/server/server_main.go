package main

import (
	"context"
	"fmt"
	"github.com/DapperBlondie/go-grpc/src/messages/files"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

// GreetService our structure that implemented our rpc
type GreetService struct{}

type SumService struct{}

func (ss *SumService) GetStreamingSumResult(r *files.SumRequest, stream files.SumService_GetStreamingSumResultServer) error {
	lst := r.GetList()
	var number int32 = 0
	for _, value := range lst {
		number += value
		resp := &files.SumResponse{Result: number}
		err := stream.Send(resp)
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 300)
	}

	return nil
}

// GetSumResult use for computing the sum of the repeated values from our request
func (ss *SumService) GetSumResult(ctx context.Context, r *files.SumRequest) (*files.SumResponse, error) {
	input := r.GetList()
	var result int32 = 0
	for _, num := range input {
		result += num
	}
	resp := &files.SumResponse{Result: result}

	return resp, nil
}

// GreetManyTimes use for stream many greet to our client
func (gs *GreetService) GreetManyTimes(r *files.GreetingManyTimeRequest, stream files.GreetService_GreetManyTimesServer) error {
	firstName := r.GetGreeting().GetFirstName()
	lastName := r.GetGreeting().GetLastName()

	result := firstName + " : " + lastName
	for i := 0; i < 10; i += 1 {
		result += " number " + strconv.Itoa(i) + "\n"
		resp := &files.GreetingManyTimesResponse{Result: result}
		err := stream.Send(resp)
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 500)
	}
	return nil
}

//Greet our Greeting API just have on rpc service that we implemented that
func (gs *GreetService) Greet(ctx context.Context, r *files.GreetingRequest) (*files.GreetingResponse, error) {
	firstName := r.GetGreeting().FirstName
	lastName := r.GetGreeting().LastName

	result := fmt.Sprintf("Hello, %s %s\nI am Greeting API", firstName, lastName)

	resp := &files.GreetingResponse{Result: result}

	return resp, nil
}

func main() {
	err := runServer()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	return
}

func runServer() error {
	log.Println("gRPC server is running ...")
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalln("Error in listening to rpc : " + err.Error())
		return err
	}

	srv := grpc.NewServer()
	files.RegisterGreetServiceServer(srv, &GreetService{})
	files.RegisterSumServiceServer(srv, &SumService{})

	log.Println("rpc Server listening on localhost:50051 ...")
	err = srv.Serve(listener)
	if err != nil {
		log.Fatalln("Error in serving our rpc server : " + err.Error())
		return err
	}

	return nil
}
