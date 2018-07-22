package mesh

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/dynamicgo/slf4go"

	"google.golang.org/grpc"
)

// Dialer .
type Dialer interface {
	Dial(ctx context.Context, serviceName string, options ...grpc.DialOption) (*grpc.ClientConn, error)
	Network() Network
}

type dialerWithBalancer struct {
	slf4go.Logger
	balancer DialerBalancer
	network  Network
}

// DialerBalancer .
type DialerBalancer interface {
	NextPeer(serviceName string) (*Peer, error)
}

type defaultBalancer struct {
	peers []*Peer
	index int
}

// DefaultBalancer .
func DefaultBalancer(peers []*Peer) DialerBalancer {
	return &defaultBalancer{
		peers: peers,
	}
}

func (balancer *defaultBalancer) NextPeer(serviceName string) (*Peer, error) {
	if len(balancer.peers) == 0 {
		return nil, nil
	}

	peers := balancer.peers[balancer.index]

	balancer.index++

	if balancer.index >= len(balancer.peers) {
		balancer.index = 0
	}

	return peers, nil
}

// NewDialer .
func NewDialer(network Network, balancer DialerBalancer) Dialer {
	return &dialerWithBalancer{
		Logger:   slf4go.Get("random dialer"),
		network:  network,
		balancer: balancer,
	}
}

func (dialer *dialerWithBalancer) Dial(ctx context.Context, serviceName string, options ...grpc.DialOption) (*grpc.ClientConn, error) {

	dialerOption := grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {

		dialer.DebugF("[%s] try dial to %s", serviceName, addr)

		peer, err := dialer.balancer.NextPeer(serviceName)

		if err != nil {
			return nil, err
		}

		if peer == nil {
			return nil, fmt.Errorf("can't find valid peer for service %s", serviceName)
		}

		return dialer.network.Dial(peer, serviceName, timeout)
	})

	options = append(options, dialerOption)

	return grpc.DialContext(ctx, serviceName, options...)
}

func (dialer *dialerWithBalancer) Network() Network {
	return dialer.network
}
