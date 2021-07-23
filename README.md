# go-grpc
An application with gRPC client and server.
For running it open the application in GoLand and sync dependencies of App then
run the server, after that run client.

# GreetService
 rpc Greet(GreetingRequest) returns (GreetingResponse) {};
 
 rpc GreetManyTimes(GreetingManyTimeRequest) returns (stream GreetingManyTimesResponse) {};
 
 rpc LongGreet(stream LongGreetRequest) returns (LongGreetResponse) {};

# SumService
rpc GetSumResult(SumRequest) returns (SumResponse) {};

rpc GetStreamingSumResult(SumRequest) returns (stream SumResponse) {};

rpc AverageStreamingResult(stream NumberRequest) returns (AverageResultResponse) {};
