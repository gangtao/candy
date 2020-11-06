package config

import (
	"context"
	"log"

	pb "github.com/gangtao/candy/protobuf"
)

type ConfigEvent struct {
    Content   string
}

type ConfigEventCallback func(*ConfigEvent)

func initEventCallback(stream pb.Configuration_MonitorConfigServer) ConfigEventCallback {
	return func(event *ConfigEvent) {
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
	DeleteConfig(dataId string, group string) error
}

type Server struct {
	pb.UnimplementedConfigurationServer
	Client KVStore
}

func (s *Server) GetConfig(ctx context.Context, in *pb.ConfigRequest) (*pb.ConfigItem, error) {
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

func (s *Server) MonitorConfig(in *pb.ConfigRequest, stream pb.Configuration_MonitorConfigServer) error {
	log.Printf("Received: %v %v %v", in.GetDataId(), in.GetGroup(), in.GetTimeout())

	var callback = initEventCallback(stream)

	err := s.Client.MonitorConfig(in.GetDataId(), in.GetGroup(), callback)
	if err != nil {
		return err
	}

	// TODO: how to cancel the monitor
	return nil
}

func (s *Server) DeletesConfig(ctx context.Context, in *pb.ConfigRequest) (*pb.DeleteConfigResponse, error) {
	log.Printf("Received: %v %v %v", in.GetDataId(), in.GetGroup())

	err := s.Client.DeleteConfig(in.GetDataId(), in.GetGroup())
	if err != nil {
		return &pb.DeleteConfigResponse{Result: false}, nil
	}

	return &pb.DeleteConfigResponse{Result: true}, nil
}