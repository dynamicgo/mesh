package agent

import (
	"context"
	"net"
	"time"

	"github.com/dynamicgo/slf4go"

	config "github.com/dynamicgo/go-config"
	"github.com/dynamicgo/mesh"
	"github.com/dynamicgo/mesh/proto"
	"google.golang.org/grpc"
)

type serviceImpl struct {
	slf4go.Logger
	server      *grpc.Server
	listener    net.Listener
	agent       *agentImpl
	serviceName string
	timeout     time.Duration
	nodeid      string
}

// NewService .
func (agent *agentImpl) newService(serviceName string, server *grpc.Server, listener net.Listener) mesh.Service {
	return &serviceImpl{
		Logger:      agent.Logger,
		server:      server,
		listener:    listener,
		agent:       agent,
		serviceName: serviceName,
		timeout:     agent.heartbeatTimeout,
		nodeid:      agent.network.ID(),
	}
}

func (service *serviceImpl) runHeartBeat(ctx context.Context) {

	ticker := time.NewTicker(service.timeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			{
				service.InfoF("[%s] service '%s' exit", service.nodeid, service.serviceName)
				return
			}
		case <-ticker.C:
			{
				if err := service.doHeartBeat(ctx); err != nil {
					service.ErrorF("[%s] service '%s' send heartbeat err %s", service.nodeid, service.serviceName, err)
				}
			}
		}
	}

}

func (service *serviceImpl) doHeartBeat(ctx context.Context) error {
	agent := service.agent

	if agent.serviceHub == nil {
		agent.WarnF("[%s] skip register service to service hub", agent.network.ID())
	} else {
		agent.DebugF("[%s] register service to service hub", agent.network.ID())

		_, err := agent.serviceHub.Register(ctx, &proto.RegisterRequest{
			Name:  service.serviceName,
			Addrs: agent.network.Peer().MeshAddrs(),
		})

		if err != nil {
			return err
		}

		agent.DebugF("[%s] register service to service hub -- success", agent.network.ID())
	}

	return nil
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

	if err := main(service.agent, service.server, config); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(service.agent.ctx)

	defer cancel()

	if err := service.doHeartBeat(ctx); err != nil {
		service.ErrorF("[%s] service '%s' send heartbeat err %s", service.nodeid, service.serviceName, err)
	}

	go service.runHeartBeat(ctx)

	return service.server.Serve(service.listener)
}
