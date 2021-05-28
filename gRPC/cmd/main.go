package main

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/tm8619/gRPC-vs-REST/gRPC/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":8081"
)

var validate = validator.New()

type server struct {
	service.UnimplementedGRPCServer
}

func (s *server) GetNumbers(ctx context.Context, in *service.GetNumbersInput) (out *service.GetNumbersOutput, e error) {
	if in == nil {
		e = errors.New("invalid input")
		return
	}

	if in.To < in.From || in.To-in.From > 10000000 {
		e = errors.New("invalid input")
		return
	}

	out.Numbers = []int64{}
	for i := in.From; i <= in.To; i++ {
		out.Numbers = append(out.Numbers, i)
	}

	return
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service.RegisterGRPCServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
