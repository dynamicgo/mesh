package agent

import (
	"net"

	config "github.com/dynamicgo/go-config"
	"github.com/dynamicgo/mesh"
	"github.com/dynamicgo/mesh/proto"
	"google.golang.org/grpc"
)

type serviceImpl struct {
	server      *grpc.Server
	listener    net.Listener
	agent       *agentImpl
	serviceName string
}

// NewService .
func (agent *agentImpl) newService(serviceName string, server *grpc.Server, listener net.Listener) mesh.Service {
	return &serviceImpl{
		server:      server,
		listener:    listener,
		agent:       agent,
		serviceName: serviceName,
	}
}

func (service *serviceImpl) Run(main mesh.ServiceMain, options ...mesh.ServiceOption) error {

	serviceConfig := &mesh.ServiceConfig{
		RemoteConfig: true,
	}

	for _, option := range options {
		option(serviceConfig)
	}

	if serviceConfig.RemoteConfig {
		source, err := service.agent.createConfigSource(service.serviceName)
		if err != nil {
			return err
		}

		serviceConfig.ConfigSources = append(serviceConfig.ConfigSources, source)
	}

	config := config.NewConfig()

	if err := config.Load(serviceConfig.ConfigSources...); err != nil {
		return err
	}

	if err := main(service.server, config); err != nil {
		return err
	}

	agent := service.agent

	if agent.serviceHub == nil {
		agent.WarnF("[%s] skip register service to service hub", agent.network.ID())
	} else {
		agent.DebugF("[%s] register service to service hub", agent.network.ID())
		_, err := agent.serviceHub.Register(agent.ctx, &proto.RegisterRequest{
			Name:  service.serviceName,
			Addrs: agent.network.Addrs(),
		})

		if err != nil {
			return err
		}

		agent.DebugF("[%s] register service to service hub -- success", agent.network.ID())
	}

	return service.server.Serve(service.listener)
}
