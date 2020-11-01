package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/gangtao/candy/numbers"
	"rsc.io/quote"

	pb "github.com/gangtao/candy/protobuf"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedConfigurationServer
}

func (s *server) GetConfig(ctx context.Context, in *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
	log.Printf("Received: %v", in.GetDataId())
	return &pb.GetConfigResponse{Content: "Hello " + in.GetDataId()}, nil
}

func main() {
	fmt.Println(quote.Go())
	fmt.Println(numbers.IsPrime(19))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterConfigurationServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
