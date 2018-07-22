package mesh

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/dynamicgo/slf4go"

	config "github.com/dynamicgo/go-config"
	"google.golang.org/grpc"
)

// Builtin errors
var (
	ErrServiceNotFound = errors.New("service not found")
	ErrNetworkClosed   = errors.New("underground network closed")
)

// Admin .
type Admin interface {
	Stop() error
}

// Agent microsrevice mesh agent node
type Agent interface {
	Network() Network
	Stop() error
	RegisterService(name string, options ...grpc.ServerOption) (Service, error)
	FindService(name string, options ...FindOption) (*grpc.ClientConn, error)
	// Dial(serviceName string, addrs []string, options ...grpc.DialOption) (*grpc.ClientConn, error)
}

// Network .
type Network interface {
	Stop() error
	Listen(serviceName string) (net.Listener, error)
	Dial(peer *Peer, serviceName string, timeout time.Duration) (net.Conn, error)
	ID() string
	Addrs() []string
}

// NetworkPlugin .
type NetworkPlugin func(ctx context.Context, config config.Config) (Network, error)

type pluginRegister struct {
	slf4go.Logger
	sync.RWMutex
	network NetworkPlugin
}

var register = &pluginRegister{
	Logger: slf4go.Get("mciro-register"),
}

// RegisterNetwork register network plugin
func RegisterNetwork(name string, network NetworkPlugin) {
	register.Lock()
	defer register.Unlock()

	register.DebugF("register network %s", name)

	register.network = network
}

// GetNetworkPlugin .
func GetNetworkPlugin() (NetworkPlugin, error) {
	if register.network == nil {
		return nil, fmt.Errorf("no valid network plugin")
	}

	return register.network, nil
}
