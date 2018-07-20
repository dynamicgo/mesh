package configservice

import (
	config "github.com/dynamicgo/go-config"
	configgrpc "github.com/dynamicgo/go-config/source/grpc/proto"
	"google.golang.org/grpc"
)

// Main servicehub main
func Main(server *grpc.Server, config config.Config) error {
	service := new(config)

	configgrpc.RegisterSourceServer(server, service)

	return nil
}
