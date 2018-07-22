package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/dynamicgo/mesh"
	"github.com/dynamicgo/mesh/proto"
)

type serviceBalancer struct {
	sync.Mutex
	serviceHub proto.ServiceHubClient
	network    mesh.Network
	peers      []*mesh.Peer
	next       int
	ctx        context.Context
	name       string
	indexer    int
}

func (agent *agentImpl) newServiceBalancer(serviceName string) *serviceBalancer {
	return &serviceBalancer{
		serviceHub: agent.serviceHub,
		network:    agent.network,
		ctx:        agent.ctx,
		name:       serviceName,
	}
}

func (balancer *serviceBalancer) NextPeer() (*mesh.Peer, error) {

	for {
		peer, ok := balancer.selectOne()

		if ok {
			return peer, nil
		}

		ok, err := balancer.fetchPeers()

		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, nil
		}
	}
}

func (balancer *serviceBalancer) selectOne() (*mesh.Peer, bool) {
	balancer.Lock()
	defer balancer.Unlock()

	if len(balancer.peers) == 0 {
		return nil, false
	}

	peers := balancer.peers[balancer.indexer]

	balancer.indexer++

	if balancer.indexer >= len(balancer.peers) {
		return nil, false
	}

	return peers, true
}

func (balancer *serviceBalancer) fetchPeers() (bool, error) {
	if balancer.serviceHub == nil {
		return false, fmt.Errorf("[%s] mesh.hub.peers is zero", balancer.network.ID())
	}

	response, err := balancer.serviceHub.Lookup(balancer.ctx, &proto.LockupRequest{
		Name: balancer.name,
	})

	if err != nil {
		return false, err
	}

	peers, err := mesh.AddrsToPeers(response.Addrs)

	if err != nil {
		return false, err
	}

	balancer.Lock()
	defer balancer.Unlock()

	balancer.peers = peers
	balancer.indexer = 0

	return len(balancer.peers) > 0, nil
}
