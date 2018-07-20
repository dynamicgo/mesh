package configservice

import (
	"context"

	config "github.com/dynamicgo/go-config"
	configgrpc "github.com/dynamicgo/go-config/source/grpc/proto"
)

type serviceImpl struct {
}

func new(config config.Config) configgrpc.SourceServer {
	return &serviceImpl{}
}
func (service *serviceImpl) Read(context.Context, *configgrpc.ReadRequest) (*configgrpc.ReadResponse, error) {
	return nil, nil
}
func (service *serviceImpl) Watch(*configgrpc.WatchRequest, configgrpc.Source_WatchServer) error {
	return nil
}
