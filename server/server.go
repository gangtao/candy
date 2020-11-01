package main

import (
	"fmt"
	"log"
	"net"

	"github.com/gangtao/candy/numbers"
	"rsc.io/quote"

	pb "github.com/gangtao/candy/protobuf"
	"google.golang.org/grpc"

	"github.com/gangtao/candy/config"
)

const (
	port = ":50051"
)

func main() {
	fmt.Println(quote.Go())
	fmt.Println(numbers.IsPrime(19))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterConfigurationServer(s, &config.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
