package main

import (
	"context"
	"fmt"
	"grpc/calculator/calculatorpb"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Add(ctx context.Context, req *calculatorpb.NumberRequest) (*calculatorpb.NumberResponse, error) {
	num1 := req.GetNum1()
	num2 := req.GetNum2()
	res := &calculatorpb.NumberResponse{
		Result: num1 + num2,
	}
	return res, nil
}

func (*server) PrimeNumber(req *calculatorpb.PrimeNumberRequest, stream calculatorpb.CalculatorService_PrimeNumberServer) error {
	num := req.GetNum()
	var k int32
	k = 2
	for num > 1 {
		if num%k == 0 {
			result := &calculatorpb.PrimeNumberResponse{
				Result: k,
			}
			stream.Send(result)
			num = num / k // divide num by k so that we have the rest of the number left.
		} else {
			k = k + 1
		}
	}
	return nil
}
func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	result := int32(0)
	i := int32(0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			result = int32(result / i)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while calling receving value: %v", err)

		}
		fmt.Printf("Message from Client: %v \n", msg)
		result += msg.GetNum()
		i++
	}
}
func main() {
	fmt.Println("Hello server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error While listening: %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while listing :%v", err)
	}
}
