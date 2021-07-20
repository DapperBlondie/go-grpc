package main

import (
	"context"
	"fmt"
	"github.com/DapperBlondie/go-grpc/src/messages/files"
	"google.golang.org/grpc"
	"log"
	"net"
)

// GreetService our structure that implemented our rpc
type GreetService struct{}

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

	log.Println("rpc Server listening on localhost:50051 ...")
	err = srv.Serve(listener)
	if err != nil {
		log.Fatalln("Error in serving our rpc server : " + err.Error())
		return err
	}

	return nil
}
