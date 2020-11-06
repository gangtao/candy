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
	defaultName = "world.is.here"
	defaultGroup = "default::group"
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
	pr, err := c.PublishConfig(ctx, &pb.PublishConfigRequest{DataId: name, Group:defaultGroup, Content:"test.content"})
	if err != nil {
		log.Fatalf("could not publish config: %v", err)
	}
	log.Printf("PublishConfig: %s", pr.GetResult())

	//Delete Config
	dr, err := c.DeletesConfig(ctx, &pb.ConfigRequest{DataId: name, Group:defaultGroup , Timeout:1000})
	if err != nil {
		log.Fatalf("could not delete config: %v", err)
	}
	log.Printf("DeletesConfig: %s", dr.GetResult())


	// Publish Config again
	pr, err = c.PublishConfig(ctx, &pb.PublishConfigRequest{DataId: name, Group:defaultGroup, Content:"test.content"})
	if err != nil {
		log.Fatalf("could not publish config: %v", err)
	}
	log.Printf("PublishConfig: %s", pr.GetResult())


	// Get Config
	gr, err := c.GetConfig(ctx, &pb.ConfigRequest{DataId: name, Group:defaultGroup, Timeout:1000})
	if err != nil {
		log.Fatalf("could not get config: %v", err)
	}
	log.Printf("Greeting: %s", gr.GetContent())

	// Monitor Config
	req := &pb.ConfigRequest{ DataId: name, Group:defaultGroup, Timeout:1000 } 
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

		log.Printf("find config change %s", config.GetContent())
	}
}