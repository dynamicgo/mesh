package mesh

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	config "github.com/dynamicgo/go-config"
	"github.com/dynamicgo/slf4go"
)

// Network .
type Network interface {
	Stop() error
	Listen(serviceName string) (net.Listener, error)
	Dial(peer *Peer, serviceName string, timeout time.Duration) (net.Conn, error)
	ID() string
	Peer() *Peer
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
