package main

import (
	"context"
	"fmt"
	"grpc/calculator/calculatorpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Printf("Client func")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection not established: %v", err)
	}
	c := calculatorpb.NewCalculatorServiceClient(cc)

	// dounary(c)
	// doServerSide(c)
	doClientSide(c)
}
func dounary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.NumberRequest{
		Num1: 10,
		Num2: 3,
	}
	res, err := c.Add(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Add Func: %v", err)
	}
	fmt.Printf("Result is :%v", res)
}
func doServerSide(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberRequest{
		Num: 120,
	}
	resstream, err := c.PrimeNumber(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumber func: %v", err)
	}
	for {
		msg, err := resstream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while calling receving value: %v", err)

		}
		fmt.Printf("Message from server: %v \n", msg)
	}
}
func doClientSide(c calculatorpb.CalculatorServiceClient) {
	request := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Num: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 4,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 5,
		},
	}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling computeAverage: %v", err)
	}
	for _, req := range request {
		stream.Send(req)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while Receving Value: %v", err)
	}
	fmt.Printf("Value Returned: %v", res.GetResult())
}
