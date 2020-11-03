package main

import (
	"fmt"
	"log"
	"net"

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

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var configServer config.Server
	configServer.Client = config.NewZKStore()

	s := grpc.NewServer()
	pb.RegisterConfigurationServer(s, &configServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
