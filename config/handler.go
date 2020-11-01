package config

import (
	"context"
	"log"

	pb "github.com/gangtao/candy/protobuf"
)

type Server struct {
	pb.UnimplementedConfigurationServer
}

func (s *Server) GetConfig(ctx context.Context, in *pb.GetConfigRequest) (*pb.ConfigItem, error) {
	log.Printf("Received: %v %v %v", in.GetDataId(), in.GetGroup(), in.GetTimeout())
	return &pb.ConfigItem{Content: "Hello " + in.GetDataId()}, nil
}

func (c *Server) PublishConfig(ctx context.Context, in *pb.PublishConfigRequest) (*pb.PublishConfigResponse, error) {
	log.Printf("Received: %v %v %v", in.GetDataId(), in.GetGroup(), in.GetContent())
	return &pb.PublishConfigResponse{Result: true}, nil
}

func (c *Server) MonitorConfig(in *pb.GetConfigRequest, stream pb.Configuration_MonitorConfigServer) error {
	item := pb.ConfigItem{Content: "test content"}
	if err := stream.Send(&item); err != nil {
    	return err
	}
	return nil
}