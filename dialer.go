package mesh

import (
	"context"
	"math/rand"
	"net"
	"time"

	"github.com/dynamicgo/slf4go"

	"google.golang.org/grpc"
)

// Dialer .
type Dialer interface {
	Dial(ctx context.Context, options ...grpc.DialOption) (*grpc.ClientConn, error)
	Network() Network
}

type randomDialer struct {
	slf4go.Logger
	rand        *rand.Rand
	addrs       []string
	network     Network
	serviceName string
}

// NewDialer .
func NewDialer(serviceName string, addrs []string, network Network) Dialer {
	return &randomDialer{
		Logger:      slf4go.Get("random dialer"),
		addrs:       addrs,
		network:     network,
		rand:        rand.New(rand.NewSource(time.Now().Unix())),
		serviceName: serviceName,
	}
}

func (dialer *randomDialer) randomSelect() string {
	id := dialer.rand.Intn(len(dialer.addrs))
	return dialer.addrs[id]
}

func (dialer *randomDialer) Dial(ctx context.Context, options ...grpc.DialOption) (*grpc.ClientConn, error) {

	dialerOption := grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		addr = dialer.randomSelect()
		dialer.DebugF("[%s] try dial to %s", dialer.serviceName, addr)
		return dialer.network.Dial(addr, dialer.serviceName, timeout)
	})

	options = append(options, dialerOption)

	return grpc.DialContext(ctx, dialer.randomSelect(), options...)
}

func (dialer *randomDialer) Network() Network {
	return dialer.network
}
