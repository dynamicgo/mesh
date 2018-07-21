package hub

import (
	config "github.com/dynamicgo/go-config"
	"github.com/dynamicgo/mesh/proto"
	"google.golang.org/grpc"
)

// Main servicehub main
func Main(server *grpc.Server, config config.Config) error {

	hub, err := new(config)

	if err != nil {
		return err
	}

	proto.RegisterServiceHubServer(server, hub)

	return nil
}
