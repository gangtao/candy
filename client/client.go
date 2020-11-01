package main

import (
	"context"
	"log"
	"os"
	"time"
	"io"

	"google.golang.org/grpc"
	pb "github.com/gangtao/candy/protobuf"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewConfigurationClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Publish Config
	pr, err := c.PublishConfig(ctx, &pb.PublishConfigRequest{DataId: name, Group:"defaultgroup", Content:"test content"})
	if err != nil {
		log.Fatalf("could not publish config: %v", err)
	}
	log.Printf("Greeting: %s", pr.GetResult())

	// Get Config
	gr, err := c.GetConfig(ctx, &pb.GetConfigRequest{DataId: name, Group:"defaultgroup", Timeout:1000})
	if err != nil {
		log.Fatalf("could not get config: %v", err)
	}
	log.Printf("Greeting: %s", gr.GetContent())

	// Monitor Config
	req := &pb.GetConfigRequest{ DataId: name, Group:"defaultgroup", Timeout:1000 } 
	stream, err := c.MonitorConfig(context.Background(), req)
	if err != nil {
		log.Fatalf("could not monitor config: %v", err)
	}

	for {
		config, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.MonitorConfig(_) = _, %v", c, err)
		}

		log.Println(config.GetContent())
	}
}