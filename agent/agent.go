package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	config "github.com/dynamicgo/go-config"
	"github.com/dynamicgo/go-config/source"
	configgrpc "github.com/dynamicgo/go-config/source/grpc/proto"
	"github.com/dynamicgo/mesh"
	"github.com/dynamicgo/mesh/proto"
	"github.com/dynamicgo/slf4go"
	"google.golang.org/grpc"
)

type agentImpl struct {
	sync.Mutex
	slf4go.Logger
	network          mesh.Network
	serviceHub       proto.ServiceHubClient
	ctx              context.Context
	cancel           context.CancelFunc
	configservices   []string
	dialer           mesh.Dialer
	serviceBalancers map[string]*serviceBalancer
}

// New create new mesh node with config
func New(config config.Config) (mesh.Agent, error) {

	plugin, err := mesh.GetNetworkPlugin()

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	agent := &agentImpl{
		Logger:           slf4go.Get("micro-agent"),
		ctx:              ctx,
		cancel:           cancel,
		serviceBalancers: make(map[string]*serviceBalancer),
	}

	network, err := plugin(ctx, config)

	agent.network = network
	agent.dialer = mesh.NewDialer(network, agent)

	if err := agent.connectAdmin(config); err != nil {
		return nil, err
	}

	if err := agent.connectConfigServer(config); err != nil {
		return nil, err
	}

	return agent, nil
}

func (agent *agentImpl) Stop() error {
	agent.cancel()
	return agent.network.Stop()
}

func (agent *agentImpl) Network() mesh.Network {
	return agent.network
}

func (agent *agentImpl) NextPeer(serviceName string) (*mesh.Peer, error) {
	agent.Lock()
	balancer, ok := agent.serviceBalancers[serviceName]

	if !ok {
		balancer = agent.newServiceBalancer(serviceName)
		agent.serviceBalancers[serviceName] = balancer
	}

	agent.Unlock()

	return balancer.NextPeer()
}

func (agent *agentImpl) RegisterService(name string, options ...grpc.ServerOption) (mesh.Service, error) {

	listener, err := agent.network.Listen(name)

	if err != nil {
		return nil, err
	}

	server := grpc.NewServer(options...)

	return agent.newService(name, server, listener), nil
}

func (agent *agentImpl) connectAdmin(config config.Config) (err error) {

	var addrs []string

	err = config.Get("mesh", "hub", "peers").Scan(&addrs)

	if err != nil {
		return
	}

	if len(addrs) == 0 {
		agent.InfoF("[%s] no need connect to mesh.hub.peers", agent.network.ID())
		return nil
	}

	conn, err := agent.dialWithDefaultBalancer(mesh.ServiceHub, addrs)

	if err != nil {
		return err
	}

	agent.serviceHub = proto.NewServiceHubClient(conn)

	return
}

func (agent *agentImpl) createConfigSource(service string) (source.Source, error) {
	if len(agent.configservices) == 0 {
		return nil, fmt.Errorf("[%s] mesh.config.peers is zero", agent.network.ID())
	}

	agent.DebugF("[%s] grpc config source dial to %s", agent.network.ID(), strings.Join(agent.configservices, ","))

	conn, err := agent.dialWithDefaultBalancer(mesh.ConfigService, agent.configservices)

	if err != nil {
		return nil, err
	}

	agent.DebugF("[%s] grpc config source dial to %s -- success", agent.network.ID(), strings.Join(agent.configservices, ","))

	sourceClient := configgrpc.NewSourceClient(conn)

	return agent.newGrpcSource(service, sourceClient), nil
}

func (agent *agentImpl) connectConfigServer(config config.Config) (err error) {

	var addrs []string

	err = config.Get("mesh", "config", "peers").Scan(&addrs)

	if err != nil {
		return
	}

	agent.configservices = addrs

	return
}

func (agent *agentImpl) dialWithDefaultBalancer(serviceName string, addrs []string, options ...grpc.DialOption) (*grpc.ClientConn, error) {

	peers, err := mesh.AddrsToPeers(addrs)

	if err != nil {
		return nil, err
	}

	dialer := mesh.NewDialer(agent.network, mesh.DefaultBalancer(peers))

	timeout := config.Get("mesh", "connect", "timeout").Duration(time.Second * 10)

	options = append(options, grpc.WithInsecure())
	options = append(options, grpc.WithBlock())
	options = append(options, grpc.WithTimeout(timeout))

	return dialer.Dial(agent.ctx, serviceName, options...)
}

func (agent *agentImpl) FindService(name string, options ...mesh.FindOption) (*grpc.ClientConn, error) {
	return agent.dialer.Dial(agent.ctx, name)
}
