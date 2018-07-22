package configservice

import (
	config "github.com/dynamicgo/go-config"
	configgrpc "github.com/dynamicgo/go-config/source/grpc/proto"
	"github.com/dynamicgo/mesh"
	"google.golang.org/grpc"
)

// Main servicehub main
func Main(agent mesh.Agent, server *grpc.Server, config config.Config) error {
	service, err := new(config)

	if err != nil {
		return err
	}

	configgrpc.RegisterSourceServer(server, service)

	return nil
}
