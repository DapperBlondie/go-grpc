First of all I want to mention that all of my applications may be are not
appropriate for this environment ((github)) because, my apps are not reusable or a pkg
that you could be able to expand it of use in other projects but I am trying and hardworking
everyday to reach that point.
This repository just a way or sign for showing my abilities in a specific language.
***

# go-grpc
An application with gRPC client and server.
For running it open the application in GoLand and sync dependencies of App then
run the server, after that run client.
Also you can use its data for your starter template.
I modify it everyday with new features if gRPC such as gRPC errors.

# GreetService
 rpc Greet(GreetingRequest) returns (GreetingResponse) {};
 
 rpc GreetManyTimes(GreetingManyTimeRequest) returns (stream GreetingManyTimesResponse) {};
 
 rpc LongGreet(stream LongGreetRequest) returns (LongGreetResponse) {};

# SumService
rpc GetSumResult(SumRequest) returns (SumResponse) {};

rpc GetStreamingSumResult(SumRequest) returns (stream SumResponse) {};

rpc AverageStreamingResult(stream NumberRequest) returns (AverageResultResponse) {};

rpc SquareRoot(SquareRootRequest) returns (SquareRootResponse) {};
