package config

import (
	"context"
	"log"

	pb "github.com/gangtao/candy/protobuf"
)

type ConfigEvent struct {
    Content   string
}

type ConfigEventCallback func(ConfigEvent)

func initEventCallback(stream pb.Configuration_MonitorConfigServer) ConfigEventCallback {
	return func(event ConfigEvent) {
        item := pb.ConfigItem{Content: event.Content}
		if err := stream.Send(&item); err != nil {
			log.Printf("Error : something terrible happen -> %s", err)
		}
    }

}

type KVStore interface {
	GetConfig(dataId string, group string) (string, error)
	PublishConfig(dataId string, group string, content string) error
	MonitorConfig(dataId string, group string, callback ConfigEventCallback) error
}

type Server struct {
	pb.UnimplementedConfigurationServer
	Client KVStore
}

func (s *Server) GetConfig(ctx context.Context, in *pb.GetConfigRequest) (*pb.ConfigItem, error) {
	log.Printf("Received: %v %v %v", in.GetDataId(), in.GetGroup(), in.GetTimeout())

	content, err := s.Client.GetConfig(in.GetDataId(), in.GetGroup())
	if err != nil {
		return nil, err
	}

	return &pb.ConfigItem{Content: content}, nil
}

func (s *Server) PublishConfig(ctx context.Context, in *pb.PublishConfigRequest) (*pb.PublishConfigResponse, error) {
	log.Printf("Received: %v %v %v", in.GetDataId(), in.GetGroup(), in.GetContent())

	err := s.Client.PublishConfig(in.GetDataId(), in.GetGroup(), in.GetContent())
	if err != nil {
		return &pb.PublishConfigResponse{Result: false}, nil
	}

	return &pb.PublishConfigResponse{Result: true}, nil
}

func (s *Server) MonitorConfig(in *pb.GetConfigRequest, stream pb.Configuration_MonitorConfigServer) error {
	log.Printf("Received: %v %v %v", in.GetDataId(), in.GetGroup(), in.GetTimeout())


	
	item := pb.ConfigItem{Content: "test content"}
	if err := stream.Send(&item); err != nil {
    	return err
	}
	return nil
}